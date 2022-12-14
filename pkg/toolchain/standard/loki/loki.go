package loki

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/chartutils"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// componentName is the name of the component.
	componentName = "loki"
)

//go:embed loki-3.6.1.tgz
var chartArchive []byte

type Loki struct {
	profile profile.Profile
}

type values struct {
	Loki         valuesLoki         `yaml:"loki"`
	Monitoring   valuesMonitoring   `yaml:"monitoring"`
	Test         valuesTest         `yaml:"test"`
	SingleBinary valuesSingleBinary `yaml:"singleBinary"`
	Grafana      valuesGrafana      `yaml:"grafana"`
}

type valuesMonitoring struct {
	SelfMonitoring valuesSelfMonitoring `yaml:"selfMonitoring"`
	ServiceMonitor valuesServiceMonitor `yaml:"serviceMonitor"`
}

type valuesLoki struct {
	CommonConfig valuesCommonConfig `yaml:"commonConfig"`
	Storage      valuesStorage      `yaml:"storage"`
	AuthEnabled  bool               `yaml:"auth_enabled"`
}

type valuesCommonConfig struct {
	ReplicationFactor uint8 `yaml:"replication_factor"`
}

type valuesStorage struct {
	Type string `yaml:"type"`
}

type valuesRead struct {
	Replicas int `yaml:"replicas"`
}

type valuesWrite struct {
	Replicas int `yaml:"replicas"`
}

type valuesSelfMonitoring struct {
	Enabled      bool               `yaml:"enabled"`
	GrafanaAgent valuesGrafanaAgent `yaml:"grafanaAgent"`
	LokiCanary   valuesLokiCanary   `yaml:"lokiCanary"`
}

type valuesGrafanaAgent struct {
	InstallOperator bool `yaml:"installOperator"`
}

type valuesLokiCanary struct {
	Enabled bool `yaml:"enabled"`
}

type valuesServiceMonitor struct {
	Enabled bool `yaml:"enabled"`
}

type valuesTest struct {
	Enabled bool `yaml:"enabled"`
}

type valuesSingleBinary struct {
	ExtraContainers   []valuesSingleBinaryExtraContainers   `yaml:"extraContainers"`
	ExtraVolumes      []valuesSingleBinaryExtraVolumes      `yaml:"extraVolumes"`
	ExtraVolumeMounts []valuesSingleBinaryExtraVolumeMounts `yaml:"extraVolumeMounts"`
}

type valuesSingleBinaryExtraContainers struct {
	Name         string                                          `yaml:"name"`
	Image        string                                          `yaml:"image"`
	Command      []string                                        `yaml:"command"`
	Args         []string                                        `yaml:"args"`
	Env          []valuesSingleBinaryExtraContainersEnv          `yaml:"env"`
	VolumeMounts []valuesSingleBinaryExtraContainersVolumeMounts `yaml:"volumeMounts"`
}

type valuesSingleBinaryExtraContainersEnv struct {
	Name      string                                        `yaml:"name"`
	ValueFrom valuesSingleBinaryExtraContainersEnvValueFrom `yaml:"valueFrom,omitempty"`
}

type valuesSingleBinaryExtraContainersEnvValueFrom struct {
	SecretKeyRef valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef `yaml:"secretKeyRef"`
}

type valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type valuesSingleBinaryExtraContainersVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type valuesSingleBinaryExtraVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type valuesSingleBinaryExtraVolumes struct {
	Name     string                                 `yaml:"name"`
	EmptyDir valuesSingleBinaryExtraVolumesEmptyDir `yaml:"emptyDir"`
}

type valuesSingleBinaryExtraVolumesEmptyDir struct {
	SizeLimit string `yaml:"sizeLimit"`
}

type valuesGrafana struct {
	GrafanaIni        map[string]interface{}           `yaml:"grafana.ini"`
	Ingress           valuesGrafanaIngress             `yaml:"ingress"`
	ExtraContainers   string                           `yaml:"extraContainers,omitempty"`
	ExtraSecretMounts []valuesGrafanaExtraSecretMounts `yaml:"extraSecretMounts"`
	Datasources       valuesGrafanaDatasources         `yaml:"datasources"`
}

type valuesGrafanaExtraContainers struct {
	Name  string                            `yaml:"name"`
	Image string                            `yaml:"image"`
	Env   []valuesGrafanaExtraContainersEnv `yaml:"env"`
}

type valuesGrafanaExtraContainersEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type valuesGrafanaIngress struct {
	Enabled     bool                      `yaml:"enabled"`
	Hosts       []string                  `yaml:"hosts"`
	Annotations map[string]string         `yaml:"annotations"`
	TLS         []valuesGrafanaIngressTLS `yaml:"tls"`
}

type valuesGrafanaIngressTLS struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type valuesGrafanaExtraSecretMounts struct {
	Name       string `yaml:"name"`
	SecretName string `yaml:"secretName"`
	MountPath  string `yaml:"mountPath"`
	ReadOnly   bool   `yaml:"readOnly"`
}

type valuesGrafanaDatasources struct {
	DatasourcesYaml valuesGrafanaDatasourcesYaml `yaml:"datasources.yaml"`
}

type valuesGrafanaDatasourcesYaml struct {
	ApiVersion  int                                       `yaml:"apiVersion"`
	Datasources []valuesGrafanaDatasourcesYamlDatasources `yaml:"datasources"`
}

type valuesGrafanaDatasourcesYamlDatasources struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Access string `yaml:"access"`
	URL    string `yaml:"url"`
}

// GetValues creates the helm values file.
func (c *Loki) GetValues(namespace string) (interface{}, error) {
	urlScheme := "https"
	if c.profile.Insecure {
		urlScheme = "http"
	}
	oidcEndpoint := authentik.GetOIDCEndpoint(c.profile.Domain, c.profile.Port, c.profile.Insecure)
	v := &values{
		Loki: valuesLoki{
			CommonConfig: valuesCommonConfig{
				ReplicationFactor: 1,
			},
			Storage: valuesStorage{
				Type: "filesystem",
			},
			AuthEnabled: false,
		},
		Monitoring: valuesMonitoring{
			SelfMonitoring: valuesSelfMonitoring{
				Enabled: false,
				GrafanaAgent: valuesGrafanaAgent{
					InstallOperator: false,
				},
				LokiCanary: valuesLokiCanary{
					Enabled: false,
				},
			},
			ServiceMonitor: valuesServiceMonitor{
				Enabled: false,
			},
		},
		Test: valuesTest{
			Enabled: false,
		},
		SingleBinary: valuesSingleBinary{
			ExtraContainers: []valuesSingleBinaryExtraContainers{
				{
					Name:    "restic",
					Image:   "restic/restic",
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", "sleep infinity"},
					Env: []valuesSingleBinaryExtraContainersEnv{
						{
							Name: "RESTIC_REPOSITORY",
							ValueFrom: valuesSingleBinaryExtraContainersEnvValueFrom{
								SecretKeyRef: valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef{
									Name: "ts-storage-config",
									Key:  "url",
								},
							},
						},
						{
							Name: "RESTIC_PASSWORD",
							ValueFrom: valuesSingleBinaryExtraContainersEnvValueFrom{
								SecretKeyRef: valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef{
									Name: "ts-storage-config",
									Key:  "password",
								},
							},
						},
						{
							Name: "AWS_ACCESS_KEY_ID",
							ValueFrom: valuesSingleBinaryExtraContainersEnvValueFrom{
								SecretKeyRef: valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef{
									Name: "ts-storage-config",
									Key:  "access-key-id",
								},
							},
						},
						{
							Name: "AWS_SECRET_ACCESS_KEY",
							ValueFrom: valuesSingleBinaryExtraContainersEnvValueFrom{
								SecretKeyRef: valuesSingleBinaryExtraContainersEnvValueFromSecretKeyRef{
									Name: "ts-storage-config",
									Key:  "secret-access-key",
								},
							},
						},
					},
					VolumeMounts: []valuesSingleBinaryExtraContainersVolumeMounts{
						{
							Name:      "backup",
							MountPath: "/tmp/backup",
						},
						{
							Name:      "restore",
							MountPath: "/tmp/restore",
						},
					},
				},
			},
			ExtraVolumes: []valuesSingleBinaryExtraVolumes{
				{
					Name: "backup",
					EmptyDir: valuesSingleBinaryExtraVolumesEmptyDir{
						SizeLimit: "1Gi",
					},
				},
				{
					Name: "restore",
					EmptyDir: valuesSingleBinaryExtraVolumesEmptyDir{
						SizeLimit: "1Gi",
					},
				},
			},
			ExtraVolumeMounts: []valuesSingleBinaryExtraVolumeMounts{
				{
					Name:      "backup",
					MountPath: "/tmp/backup",
				},
				{
					Name:      "restore",
					MountPath: "/tmp/restore",
				},
			},
		},
		Grafana: valuesGrafana{
			GrafanaIni: map[string]interface{}{
				"server": map[string]interface{}{
					"domain":   fmt.Sprintf("grafana.%s", c.profile.Domain),
					"root_url": fmt.Sprintf("%s://grafana.%s:%d", urlScheme, c.profile.Domain, c.profile.Port),
				},
				"auth.generic_oauth": map[string]interface{}{
					"enabled":                    true,
					"name":                       "SSO",
					"scopes":                     "openid profile email",
					"auth_url":                   fmt.Sprintf("%sapplication/o/authorize/", oidcEndpoint),
					"token_url":                  fmt.Sprintf("%sapplication/o/token/", oidcEndpoint),
					"api_url":                    fmt.Sprintf("%sapplication/o/userinfo/", oidcEndpoint),
					"client_id":                  "$__file{/etc/secrets/auth_generic_oauth/client-id}",
					"client_secret":              "$__file{/etc/secrets/auth_generic_oauth/client-secret}",
					"allow_assign_grafana_admin": true,
					"role_attribute_path":        "contains(groups[*], 'admins') && 'GrafanaAdmin' || contains(groups[*], 'editors') && 'Editor' || 'Viewer'",
				},
			},
			Ingress: valuesGrafanaIngress{
				Enabled: true,
				Hosts:   []string{fmt.Sprintf("grafana.%s", c.profile.Domain)},
			},
			ExtraSecretMounts: []valuesGrafanaExtraSecretMounts{
				{
					Name:       "oidc-client-credentials",
					SecretName: "loki-oidc-client",
					MountPath:  "/etc/secrets/auth_generic_oauth",
					ReadOnly:   true,
				},
			},
			Datasources: valuesGrafanaDatasources{
				DatasourcesYaml: valuesGrafanaDatasourcesYaml{
					ApiVersion: 1,
					Datasources: []valuesGrafanaDatasourcesYamlDatasources{
						{
							Name:   "Loki",
							Type:   "loki",
							Access: "direct",
							URL:    "http://loki:3100",
						},
					},
				},
			},
		},
	}
	if !c.profile.Insecure {
		v.Grafana.Ingress.Annotations = map[string]string{
			"cert-manager.io/cluster-issuer": "ts-system",
			"kubernetes.io/ingress.class":    "ts-system",
		}
		v.Grafana.Ingress.TLS = []valuesGrafanaIngressTLS{
			{
				Hosts: []string{
					fmt.Sprintf("grafana.%s", c.profile.Domain),
				},
				SecretName: "grafana-ingress-tls-cert",
			},
		}
		v.Grafana.ExtraContainers = ""
	} else {
		uid, err := chartutils.UniqueID(namespace)
		if err != nil {
			return nil, err
		}
		grafanaExtraContainers := []valuesGrafanaExtraContainers{
			{
				Name:  "oidc-auth-proxy",
				Image: "quay.io/trustacks/local-gd-proxy",
				Env: []valuesGrafanaExtraContainersEnv{
					{
						Name:  "UPSTREAM",
						Value: fmt.Sprintf("authentik-%s", uid),
					},
					{
						Name:  "LISTEN_PORT",
						Value: strconv.Itoa(int(c.profile.Port)),
					},
					{
						Name:  "SERVICE",
						Value: "authentik",
					},
				},
			},
		}
		grafanaExtraContainersData, err := yaml.Marshal(grafanaExtraContainers)
		if err != nil {
			return nil, err
		}
		v.Grafana.ExtraContainers = string(grafanaExtraContainersData)
	}
	return v, nil
}

// GetChart returns the loki helm chart archive.
func (c *Loki) GetChart() (string, error) {
	chart, err := chartutils.NewChart(componentName, chartArchive)
	if err != nil {
		return "", err
	}
	path, err := chart.Save()
	if err != nil {
		return "", err
	}
	return path, nil
}

// Install runs pre installation tasks and install the component.
func (c *Loki) Install(dispatcher client.Dispatcher, namespace string) error {
	if err := c.preInstall(dispatcher.Clientset(), namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	lokiValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"loki": lokiValues,
	}
	if err := dispatcher.InstallChart("loki", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return nil
}

// Upgrade backs up and upgrades the component.
func (c *Loki) Upgrade(dispatcher client.Dispatcher, namespace string) error {
	if err := c.backup(dispatcher, namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	lokiValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"loki": lokiValues,
	}
	return dispatcher.UpgradeChart("loki", values, time.Minute*5, chartPath)
}

// Rollback deploys the previous state of the component and restores
// the component's data.
func (c *Loki) Rollback(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	lokiValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"loki": lokiValues,
	}
	if err := dispatcher.RollbackRelease("loki", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return c.restore(dispatcher, namespace)
}

// Uninstall removes the component.
func (c *Loki) Uninstall(dispatcher client.Dispatcher, namespace string) error {
	if err := dispatcher.UninstallChart("loki"); err != nil {
		return err
	}
	return nil
}

// preInstall creates the oidc client and secret.
func (c *Loki) preInstall(clientset kubernetes.Interface, namespace string) error {
	clientID, clientSecret, err := authentik.CreateOIDCClient(componentName, namespace)
	if err != nil {
		return err
	}
	if err := createOIDCClientSecret(clientID, clientSecret, namespace, clientset); err != nil {
		return err
	}
	return nil
}

// backup tars the loki data and writes the backup toe the s3
// backend.
func (c *Loki) backup(dispatcher client.Dispatcher, namespace string) error {
	cmd := `cd /tmp/backup && tar czf loki-data -C /var loki`
	if err := dispatcher.ExecCommand("loki-0", "single-binary", cmd, namespace); err != nil {
		return err
	}
	cmd = `restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/loki-data`
	return dispatcher.ExecCommand("loki-0", "restic", cmd, namespace)
}

// restore retrieves the s3 backup and untars the loki data to the
// data directory.
func (c *Loki) restore(dispatcher client.Dispatcher, namespace string) error {
	cmd := `restic restore latest --target /tmp/restore --include /tmp/backup/loki-data`
	if err := dispatcher.ExecCommand("loki-0", "restic", cmd, namespace); err != nil {
		return err
	}
	cmd = `cd /tmp/restore && tar xf /tmp/restore/tmp/backup/loki-data && cp -R ./loki /var/loki`
	return dispatcher.ExecCommand("loki-0", "single-binary", cmd, namespace)
}

// New creates a new loki instance.
func New(prof profile.Profile) *Loki {
	return &Loki{prof}
}

// createOIDCClientSecret creates the oidc client secret secret.
func createOIDCClientSecret(clientID, clientSecret, namespace string, clientset kubernetes.Interface) error {
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "loki-oidc-client", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "loki-oidc-client",
				},
				Data: map[string][]byte{
					"client-id":     []byte(clientID),
					"client-secret": []byte(clientSecret),
				},
			}
			_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			return err
		}
		if !strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}
