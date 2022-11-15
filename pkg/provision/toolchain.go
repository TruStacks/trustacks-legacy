package toolchain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"filippo.io/age"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/trustacks/trustacks/pkg"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	// toolchainRoot is where software toolchain metadata is stored.
	toolchainRoot = filepath.Join(pkg.RootDir, "toolchains")
)

type Config struct {
}

// component represents a toolchain component.
type component struct {
	Repo             string `json:"repository"`
	Chart            string `json:"chart"`
	Version          string `json:"version"`
	Values           string `json:"values"`
	Hooks            string `json:"hooks"`
	ApplicationHooks string `json:"applicationHooks,omitempty"`
}

// componentCatalogConfigParameters .
type componentCatalogConfigParameters struct {
	Name    string `json:"name"`
	Default string `json:"default"`
}

// componentCatalogConfig .
type componentCatalogConfig struct {
	Parameters []componentCatalogConfigParameters `json:"parameters"`
}

// componentCatalog contains the component manifests.
type componentCatalog struct {
	HookSource string                  `json:"hookSource"`
	Version    string                  `json:"version"`
	Components map[string]component    `json:"components"`
	Config     *componentCatalogConfig `json:"config"`
}

// getToolchainCatalog gets the component catalog.
func getToolchainCatalog(url string) (*componentCatalog, error) {
	resp, err := http.Get(fmt.Sprintf("%s/.well-known/catalog-manifest", url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var catalog *componentCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}
	return catalog, nil
}

// toolchainDependencies contains the catalog and required components.
type toolchainDependencies struct {
	Catalog    string   `yaml:"catalog"`
	Components []string `yaml:"components"`
}

// toolchain represents a toolchain helm chart.
type toolchain struct {
	name         string
	Dependencies []toolchainDependencies `yaml:"dependencies"`
}

// addComponents downloads the component charts and adds them to the
// resource components.
func (tc *toolchain) addComponents(components []string, catalog *componentCatalog) error {
	for _, name := range components {
		// Check if the chart already exists.
		if _, err := os.Stat(path.Join(tc.componentsPath(), name)); !os.IsNotExist(err) {
			continue
		}
		component := catalog.Components[name]
		pull := action.NewPullWithOpts(action.WithConfig(&action.Configuration{}))
		pull.Settings = cli.New()
		pull.UntarDir = tc.componentsPath()
		pull.Untar = true
		url := fmt.Sprintf("%s/%s-%s.tgz", component.Repo, component.Chart, component.Version)
		_, err := pull.Run(url)
		if err != nil {
			return err
		}
		if err := os.Remove(path.Join(tc.componentsPath(), fmt.Sprintf("%s-%s.tgz", component.Chart, component.Version))); err != nil {
			return err
		}
	}
	return nil
}

// addHooks creates the hook template file in the chart.
func (tc *toolchain) addHooks(components []string, catalog *componentCatalog, params map[string]interface{}) error {
	for _, name := range components {
		params["image"] = catalog.HookSource
		component := catalog.Components[name]

		var buf bytes.Buffer
		t := template.Must(template.New("hook").Parse(component.Hooks))
		if err := t.Execute(&buf, params); err != nil {
			return err
		}
		path := filepath.Join(tc.componentsPath(), name, "templates", "trustacks-hooks.yaml")
		if err := os.WriteFile(path, buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}

// addSubchartValues adds the subchart values to the helm values
// file.
func (tc *toolchain) addSubChartValues(components []string, catalog *componentCatalog, parameters map[string]interface{}) error {
	for _, name := range components {
		values, err := os.Create(path.Join(tc.componentsPath(), name, "override-values.yaml"))
		if err != nil {
			return err
		}
		component := catalog.Components[name]
		t := template.Must(template.New("values").Funcs(sprig.FuncMap()).Parse(component.Values))
		var buf bytes.Buffer
		if err := t.Execute(&buf, parameters); err != nil {
			return err
		}
		if _, err := values.Write(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// createAgeKeySecret creates the age private and public keys secret.
func (tc *toolchain) createAgeKeySecret() error {
	privateKey, err := age.GenerateX25519Identity()
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v", err)
	}
	secret := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]interface{}{
			"name": "sops-age",
		},
		"stringData": map[string]string{
			"age.agekey": privateKey.String(),
			"age.agepub": privateKey.Recipient().String(),
		},
	}
	yml, err := yaml.Marshal(secret)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(tc.path(), "chart", "templates"), 0755); err != nil {
		return err
	}
	return os.WriteFile(path.Join(tc.path(), "chart", "templates", "sops-age-secret.yaml"), yml, 0644)
}

// install installs the toolchain helm chart.
func (tc *toolchain) install() error {
	slug := fmt.Sprintf("trustacks-toolchain-%s", tc.name)
	chartSpec := helmclient.ChartSpec{
		ReleaseName:     slug,
		ChartName:       filepath.Join(tc.path(), "chart"),
		Namespace:       slug,
		UpgradeCRDs:     true,
		CreateNamespace: true,
		CleanupOnFail:   true,
	}
	kubeconfig, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		return err
	}
	helmClient, err := helmclient.NewClientFromKubeConf(&helmclient.KubeConfClientOptions{
		Options:    &helmclient.Options{Namespace: slug},
		KubeConfig: kubeconfig,
	})
	if err != nil {
		return err
	}
	_, err = helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil)
	return err
}

// installComponents installs the component helm charts.
func (tc *toolchain) installComponents() error {
	slug := fmt.Sprintf("trustacks-toolchain-%s", tc.name)
	components, err := os.ReadDir(tc.componentsPath())
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, component := range components {
		wg.Add(1)
		go func(name string, wg *sync.WaitGroup) {
			values, err := os.ReadFile(filepath.Join(tc.componentsPath(), name, "override-values.yaml"))
			if err != nil {
				log.Fatalf("error reading override values: %s", err)
			}
			chartSpec := helmclient.ChartSpec{
				ReleaseName:     name,
				ChartName:       filepath.Join(tc.componentsPath(), name),
				Namespace:       slug,
				UpgradeCRDs:     true,
				CreateNamespace: true,
				CleanupOnFail:   true,
				ValuesYaml:      string(values),
			}
			kubeconfig, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".kube", "config"))
			if err != nil {
				log.Fatalf("error reading kubeconfig: %s", err)
			}
			helmClient, err := helmclient.NewClientFromKubeConf(&helmclient.KubeConfClientOptions{
				Options:    &helmclient.Options{Namespace: slug},
				KubeConfig: kubeconfig,
			})
			if err != nil {
				log.Fatalf("error creating helm client: %s", err)
			}
			_, err = helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil)
			if err != nil {
				log.Fatalf("error deploying '%s': %s", name, err)
			}
			wg.Done()
		}(component.Name(), &wg)
	}
	wg.Wait()
	return nil
}

// join combines the toolchain configuration parameters with the
// component parameters
//
// Parameter defaults are set if required.
func (tc *toolchain) join(parameters map[string]interface{}, catalogParameters []componentCatalogConfigParameters) map[string]interface{} {
	joined := make(map[string]interface{})
	for _, param := range catalogParameters {
		if _, ok := parameters[param.Name]; !ok {
			if param.Default != "" {
				joined[param.Name] = param.Default
			}
		} else {
			joined[param.Name] = parameters[param.Name]
		}
	}
	return joined
}

// path returns the filesystem path of the toolchain metadata.
func (tc *toolchain) path() string {
	return filepath.Join(toolchainRoot, tc.name)
}

// componentsPath returns the filesystem path of the toolchain
// components.
func (tc *toolchain) componentsPath() string {
	return filepath.Join(toolchainRoot, tc.name, "components")
}

// applicationsPath returns the filesystem path of the applications.
func (tc *toolchain) applicationsPath() string {
	return filepath.Join(toolchainRoot, tc.name, "applications")
}

// newToolchain creates a new toolchain chart instance.
func newToolchain(name, source, version string, force bool, cloneFunc func(string, bool, *git.CloneOptions) (*git.Repository, error)) (*toolchain, error) {
	tc := &toolchain{name: name}
	if _, err := os.Stat(tc.path()); !os.IsNotExist(err) && !force {
		return nil, fmt.Errorf("error: toolchain '%s' already exists", name)
	}
	if force {
		if err := os.RemoveAll(tc.path()); err != nil {
			return nil, err
		}
	}
	if _, err := cloneFunc(tc.path(), false, &git.CloneOptions{
		URL:           source,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
	}); err != nil {
		return nil, err
	}
	if err := tc.createAgeKeySecret(); err != nil {
		return nil, err
	}
	manifest, err := os.ReadFile(filepath.Join(tc.path(), "config.yaml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(manifest, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func newToolchainFromConfig(name string) (*toolchain, error) {
	tc := &toolchain{name: name}
	manifest, err := os.ReadFile(filepath.Join(tc.path(), "config.yaml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(manifest, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

// toolchainConfig contains the toolchain configuration parameters.
type toolchainConfig struct {
	Name         string                 `json:"name"`
	Source       string                 `json:"source"`
	Version      string                 `json:"version"`
	Parameters   map[string]interface{} `json:"parameters"`
	Applications []applicationConfig    `json:"applications"`
}

// loadToolchainConfig loads the config file at the provided path.
func loadToolchainConfig(path string) (*toolchainConfig, error) {
	rawConfig, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config *toolchainConfig
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// Install installs the toolchain.
func Install(configPath string, force bool, cloneFunc func(string, bool, *git.CloneOptions) (*git.Repository, error)) error {
	config, err := loadToolchainConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading the toolchain config: %s", err)
	}
	tc, err := newToolchain(config.Name, config.Source, config.Version, force, cloneFunc)
	if err != nil {
		return fmt.Errorf("error creating the toolchian: %s", err)
	}
	for _, dep := range tc.Dependencies {
		catalog, err := getToolchainCatalog(dep.Catalog)
		if err != nil {
			return fmt.Errorf("error fetching catalog: %s", err)
		}
		parameters := tc.join(config.Parameters, catalog.Config.Parameters)
		if err := tc.addComponents(dep.Components, catalog); err != nil {
			return fmt.Errorf("error adding subcharts: %s", err)
		}
		if err := tc.addHooks(dep.Components, catalog, parameters); err != nil {
			return fmt.Errorf("error adding hook templates: %s", err)
		}
		if err := tc.addSubChartValues(dep.Components, catalog, parameters); err != nil {
			return fmt.Errorf("error adding subchart values: %s", err)
		}
	}
	if err := tc.install(); err != nil {
		return fmt.Errorf("error installing the toolchain chart: %s", err)
	}
	if err := tc.installComponents(); err != nil {
		return fmt.Errorf("error installing the toolchain components: %s", err)
	}
	return nil
}

// Destroy removes the software factory kubernetes resources and the
// toolchain helm assets.
func Destory(name string, clientset kubernetes.Interface) error {
	tc := &toolchain{}
	if _, err := os.Stat(tc.path()); os.IsNotExist(err) {
		return fmt.Errorf("error: toolchain '%s' could not be found", name)
	}
	if err := clientset.CoreV1().Namespaces().Delete(context.TODO(), fmt.Sprintf("trustacks-toolchain-%s", name), metav1.DeleteOptions{}); err != nil {
		return err
	}
	return os.RemoveAll(tc.path())
}
