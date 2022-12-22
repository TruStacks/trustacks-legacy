package authentik

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/chartutils"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// componentName is the name of the component.
	componentName = "authentik"
	// serviceURL is the authentik kubernetes service name.
	serviceURL = "http://authentik"
)

//go:embed authentik-2022.11.0.tgz
var chartArchive []byte

type Authentik struct {
	profile profile.Profile
}

type values struct {
	FullnameOverride string             `yaml:"fullnameOverride"`
	Authentik        valuesAuthentik    `yaml:"authentik"`
	Ingress          valuesIngress      `yaml:"ingress"`
	Postgresql       valuesPostgresql   `yaml:"postgresql"`
	Redis            valuesRedis        `yaml:"redis"`
	EnvValueFrom     valuesEnvValueFrom `yaml:"envValueFrom"`
}

type valuesAuthentik struct {
	SecretKey  string                    `yaml:"secret_key"`
	Postgresql valuesAuthentikPostgresql `yaml:"postgresql"`
}

type valuesAuthentikPostgresql struct {
	Host string `yaml:"host"`
}

type valuesIngress struct {
	Enabled     bool                 `yaml:"enabled"`
	Hosts       []valuesIngressHosts `yaml:"hosts"`
	Annotations map[string]string    `yaml:"annotations,omitempty"`
	TLS         []valuesIngressTLS   `yaml:"tls,omitempty"`
}

type valuesIngressHosts struct {
	Host  string                  `yaml:"host"`
	Paths []valuesIngressHostPath `yaml:"paths"`
}

type valuesIngressHostPath struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
}

type valuesIngressTLS struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type valuesPostgresql struct {
	Enabled          bool                    `yaml:"enabled"`
	FullnameOverride string                  `yaml:"fullnameOverride"`
	ExistingSecret   string                  `yaml:"existingSecret"`
	Primary          valuesPostgresqlPrimary `yaml:"primary"`
}

type valuesPostgresqlPrimary struct {
	ExtraVolumes      []valuesPostgresqlPrimaryExtraVolumes      `yaml:"extraVolumes"`
	ExtraVolumeMounts []valuesPostgresqlPrimaryExtraVolumeMounts `yaml:"extraVolumeMounts"`
	Sidecars          []valuesPostgresqlPrimarySidecars          `yaml:"sidecars"`
}

type valuesPostgresqlPrimaryExtraVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type valuesPostgresqlPrimaryExtraVolumes struct {
	Name     string                                      `yaml:"name"`
	EmptyDir valuesPostgresqlPrimaryExtraVolumesEmptyDir `yaml:"emptyDir"`
}

type valuesPostgresqlPrimaryExtraVolumesEmptyDir struct {
	SizeLimit string `yaml:"sizeLimit"`
}

type valuesPostgresqlPrimarySidecars struct {
	Name         string                                        `yaml:"name"`
	Image        string                                        `yaml:"image"`
	Command      []string                                      `yaml:"command"`
	Args         []string                                      `yaml:"args"`
	VolumeMounts []valuesPostgresqlPrimarySidecarsVolumeMounts `yaml:"volumeMounts"`
	Env          []valuesPostgresqlPrimarySidecarsEnv          `yaml:"env"`
}

type valuesPostgresqlPrimarySidecarsVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type valuesPostgresqlPrimarySidecarsEnv struct {
	Name      string                                      `yaml:"name"`
	ValueFrom valuesPostgresqlPrimarySidecarsEnvValueFrom `yaml:"valueFrom,omitempty"`
}

type valuesPostgresqlPrimarySidecarsEnvValueFrom struct {
	SecretKeyRef valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef `yaml:"secretKeyRef"`
}

type valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type valuesRedis struct {
	Enabled bool `yaml:"enabled"`
}

type valuesEnvValueFrom struct {
	AuthentikPostgresqlPassword valuesEnvValueFromAuthentikPostgresqlPassword `yaml:"AUTHENTIK_POSTGRESQL__PASSWORD"`
	AuthentikBootstrapToken     valuesEnvValueFromAuthentikBootstrapToken     `yaml:"AUTHENTIK_BOOTSTRAP_TOKEN"`
}

type valuesEnvValueFromAuthentikPostgresqlPassword struct {
	SecretKeyRef valuesEnvValueFromAuthentikPostgresqlPasswordSecretKeyRef `yaml:"secretKeyRef"`
}

type valuesEnvValueFromAuthentikPostgresqlPasswordSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type valuesEnvValueFromAuthentikBootstrapToken struct {
	SecretKeyRef valuesEnvValueFromAuthentikBootstrapTokenSecretKeyRef `yaml:"secretKeyRef"`
}

type valuesEnvValueFromAuthentikBootstrapTokenSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

// GetValues creates the helm values object.
func (c *Authentik) GetValues(namespace string) (interface{}, error) {
	secretKey, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return "", err
	}
	uid, err := chartutils.UniqueID(namespace)
	if err != nil {
		return "", err
	}
	v := &values{
		FullnameOverride: fmt.Sprintf("authentik-%s", uid),
		Authentik: valuesAuthentik{
			SecretKey: secretKey,
			Postgresql: valuesAuthentikPostgresql{
				Host: "authentik-postgresql",
			},
		},
		Ingress: valuesIngress{
			Enabled: true,
			Hosts: []valuesIngressHosts{
				{
					Host: fmt.Sprintf("authentik.%s", c.profile.Domain),
					Paths: []valuesIngressHostPath{
						{
							Path:     "/",
							PathType: "Prefix",
						},
					},
				},
			},
		},
		Postgresql: valuesPostgresql{
			Enabled:        true,
			ExistingSecret: "authentik-postgresql",
			Primary: valuesPostgresqlPrimary{
				Sidecars: []valuesPostgresqlPrimarySidecars{
					{
						Name:    "restic",
						Image:   "restic/restic",
						Command: []string{"/bin/sh"},
						Args:    []string{"-c", "sleep infinity"},
						Env: []valuesPostgresqlPrimarySidecarsEnv{
							{
								Name: "RESTIC_REPOSITORY",
								ValueFrom: valuesPostgresqlPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "url",
									},
								},
							},
							{
								Name: "RESTIC_PASSWORD",
								ValueFrom: valuesPostgresqlPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "password",
									},
								},
							},
							{
								Name: "AWS_ACCESS_KEY_ID",
								ValueFrom: valuesPostgresqlPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "access-key-id",
									},
								},
							},
							{
								Name: "AWS_SECRET_ACCESS_KEY",
								ValueFrom: valuesPostgresqlPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesPostgresqlPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "secret-access-key",
									},
								},
							},
						},
						VolumeMounts: []valuesPostgresqlPrimarySidecarsVolumeMounts{
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
				ExtraVolumes: []valuesPostgresqlPrimaryExtraVolumes{
					{
						Name: "backup",
						EmptyDir: valuesPostgresqlPrimaryExtraVolumesEmptyDir{
							SizeLimit: "1Gi",
						},
					},
					{
						Name: "restore",
						EmptyDir: valuesPostgresqlPrimaryExtraVolumesEmptyDir{
							SizeLimit: "1Gi",
						},
					},
				},
				ExtraVolumeMounts: []valuesPostgresqlPrimaryExtraVolumeMounts{
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
		Redis: valuesRedis{
			Enabled: true,
		},
		EnvValueFrom: valuesEnvValueFrom{
			AuthentikBootstrapToken: valuesEnvValueFromAuthentikBootstrapToken{
				SecretKeyRef: valuesEnvValueFromAuthentikBootstrapTokenSecretKeyRef{
					Name: "authentik-bootstrap",
					Key:  "api-token",
				},
			},
			AuthentikPostgresqlPassword: valuesEnvValueFromAuthentikPostgresqlPassword{
				SecretKeyRef: valuesEnvValueFromAuthentikPostgresqlPasswordSecretKeyRef{
					Name: "authentik-postgresql",
					Key:  "postgresql-postgres-password",
				},
			},
		},
	}
	if !c.profile.Insecure {
		v.Ingress.Annotations = map[string]string{
			"cert-manager.io/cluster-issuer": "ts-system",
			"kubernetes.io/ingress.class":    "ts-system",
		}
		v.Ingress.TLS = []valuesIngressTLS{
			{
				Hosts:      []string{fmt.Sprintf("authentik.%s", c.profile.Domain)},
				SecretName: "authentik-ingress-tls-cert",
			},
		}
	}
	return v, nil
}

// GetChart returns the authentik helm chart archive.
func (c *Authentik) GetChart() (string, error) {
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
func (c *Authentik) Install(dispatcher client.Dispatcher, namespace string) error {
	if err := c.preInstall(dispatcher.Clientset(), namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	authentikValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"authentik": authentikValues,
	}
	if err := dispatcher.InstallChart("authentik", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return c.postInstall(dispatcher.Clientset(), namespace)
}

// Upgrade backs up and upgrades the component.
func (c *Authentik) Upgrade(dispatcher client.Dispatcher, namespace string) error {
	if err := c.backup(dispatcher, namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	authentikValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"authentik": authentikValues,
	}
	return dispatcher.UpgradeChart("authentik", values, time.Minute*5, chartPath)
}

// Rollback deploys the previous state of the component and restores
// the component's data.
func (c *Authentik) Rollback(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	authentikValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"authentik": authentikValues,
	}
	if err := dispatcher.RollbackRelease("authentik", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return c.restore(dispatcher, namespace)
}

// Uninstall removes the component.
func (c *Authentik) Uninstall(dispatcher client.Dispatcher, namespace string) error {
	if err := dispatcher.UninstallChart("authentik"); err != nil {
		return err
	}
	return nil
}

// preInstall creates the authentik admin api token.
func (c *Authentik) preInstall(clientset kubernetes.Interface, namespace string) error {
	apiToken, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return err
	}
	if err := createAPIToken(namespace, apiToken, clientset); err != nil {
		return err
	}
	postgresqlPassword, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		return err
	}
	if err := createPostgresqlPassword(namespace, postgresqlPassword, clientset); err != nil {
		return err
	}
	return nil
}

// postInstall creates the authentik user groups.
func (c *Authentik) postInstall(clientset kubernetes.Interface, namespace string) error {
	token, err := getAPIToken(namespace, clientset)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	url, err := c.getServiceURL(namespace)
	if err != nil {
		return err
	}
	if err := healthCheckService(url, 2, ctx); err != nil {
		return err
	}
	if err := createGroups(url, token); err != nil {
		return err
	}
	return nil
}

// backup creates a postgresql database dump and writes the backup
// to the s3 backend.
func (c *Authentik) backup(dispatcher client.Dispatcher, namespace string) error {
	cmd := `PGPASSWORD=$POSTGRES_PASSWORD pg_dump -U $POSTGRES_USER -F c -b -v -f /tmp/backup/authentik-postgresql $POSTGRES_DB`
	if err := dispatcher.ExecCommand("authentik-postgresql-0", "authentik-postgresql", cmd, namespace); err != nil {
		return err
	}
	cmd = `restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/authentik-postgresql`
	return dispatcher.ExecCommand("authentik-postgresql-0", "restic", cmd, namespace)
}

// restore retrieves the s3 backup and restores the postgresql
// database backup.
func (c *Authentik) restore(dispatcher client.Dispatcher, namespace string) error {
	cmd := `restic restore latest --target /tmp/restore --include /tmp/backup/authentik-postgresql`
	if err := dispatcher.ExecCommand("authentik-postgresql-0", "restic", cmd, namespace); err != nil {
		return err
	}
	cmd = `PGPASWORD=$POSTGRES_PASSWORD pg_restore -U $POSTGRES_USER /tmp/restore/tmp/backup/authentik-postgresql`
	return dispatcher.ExecCommand("authentik-postgresql-0", "authentik-postgresql", cmd, namespace)
}

// getServiceURL gets the authentik k8s service url.
func (c *Authentik) getServiceURL(namespace string) (string, error) {
	uid, err := chartutils.UniqueID(namespace)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s.%s.svc.cluster.local", serviceURL, uid, namespace), nil
}

// New creates a new authentik instance.
func New(prof profile.Profile) *Authentik {
	return &Authentik{prof}
}

// createAPIToken creates the api token secret.
func createAPIToken(namespace, token string, clientset kubernetes.Interface) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "authentik-bootstrap",
		},
		Data: map[string][]byte{
			"api-token": []byte(token),
		},
	}
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "authentik-bootstrap", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
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

// createPostgresqlPassword creates the authentik postgresql
// database passwords.
func createPostgresqlPassword(namespace string, password string, clientset kubernetes.Interface) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "authentik-postgresql",
		},
		Data: map[string][]byte{
			"postgresql-password":          []byte(password),
			"postgresql-postgres-password": []byte(password),
		},
	}
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "authentik-postgresql", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
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

// getAPIToken gets the api token secret value.
func getAPIToken(namespace string, clientset kubernetes.Interface) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "authentik-bootstrap", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret.Data["api-token"])), nil
}

type group struct {
	Name        string `json:"name"`
	Users       []int  `json:"users"`
	IsSuperuser bool   `json:"is_superuser"`
	Parent      *int   `json:"parent"`
}

// createGroups creates the user groups.
func createGroups(url, token string) error {
	groups := []group{
		{"admins", []int{1}, true, nil},
		{"editors", []int{}, false, nil},
		{"viewers", []int{}, false, nil},
	}
	for _, g := range groups {
		data, err := json.Marshal(g)
		if err != nil {
			return err
		}
		// check if the group already exists.
		resp, err := getAPIResource(url, "core/groups", token, fmt.Sprintf("name=%s", g.Name))
		if err != nil {
			return err
		}
		results := make(map[string]interface{})
		if err := json.Unmarshal(resp, &results); err != nil {
			return err
		}
		if len(results["results"].([]interface{})) > 0 {
			continue
		}
		_, err = postAPIResource(url, "core/groups", token, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// getAPIResource gets the API resource at the provided path.
func getAPIResource(url, resource, token string, search string) ([]byte, error) {
	uri := fmt.Sprintf("%s/api/v3/%s/", url, resource)
	if search != "" {
		uri = fmt.Sprintf("%s?%s", uri, search)
	}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("'%s' get error: %s", resource, body)
	}
	return body, nil
}

// postAPIResource posts the API resource at the provided path.
func postAPIResource(url, resource, token string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v3/%s/", url, resource), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("'%s' post error: %s", resource, body)
	}
	return body, nil
}

// healthCheckService checks the health of the authentik service.
func healthCheckService(url string, interval int, ctx context.Context) error {
	for {
		select {
		case <-time.After(time.Second * time.Duration(interval)):
			if _, err := http.Get(url); err != nil {
				log.Println(err)
				continue
			}
		case <-ctx.Done():
			return errors.New("service health check timeout")
		}
		break
	}
	return nil
}

type propertyMapping struct {
	PK      string `json:"pk"`
	Managed string `json:"managed"`
}

type propertyMappings struct {
	Results []propertyMapping `json:"results"`
}

// getPropertyMappings gets the ids of the oauth2 scope mappings.
func getPropertyMappings(url, token string) ([]string, error) {
	scopes := []string{
		"goauthentik.io/providers/oauth2/scope-email",
		"goauthentik.io/providers/oauth2/scope-openid",
		"goauthentik.io/providers/oauth2/scope-profile",
	}
	pks := make([]string, len(scopes))
	resp, err := getAPIResource(url, "propertymappings/all", token, "")
	if err != nil {
		return nil, err
	}
	pm := &propertyMappings{}
	if err := json.Unmarshal(resp, &pm); err != nil {
		return nil, err
	}
	for _, p := range pm.Results {
		for i, scope := range scopes {
			if p.Managed == scope {
				pks[i] = p.PK
			}
		}
	}
	return pks, nil
}

type flow struct {
	PK   string `json:"pk"`
	Slug string `json:"slug"`
}

type flows struct {
	Results []flow `json:"results"`
}

// getAuthorizationFlow gets the id of the default authorization
// flow.
func getAuthorizationFlow(url, token string) (string, error) {
	resp, err := getAPIResource(url, "flows/instances", token, "")
	if err != nil {
		return "", err
	}
	f := &flows{}
	if err := json.Unmarshal(resp, &f); err != nil {
		return "", err
	}
	for _, flow := range f.Results {
		if flow.Slug == "default-provider-authorization-explicit-consent" {
			return flow.PK, nil
		}
	}
	return "", errors.New("authorization flow not found")
}

type certificateKeypair struct {
	PK   string `json:"pk"`
	Name string `json:"name"`
}

type certificateKeypairs struct {
	Results []certificateKeypair `json:"results"`
}

// getCertificateKeypair gets the certificate key pair.
func getCertificateKeypair(url, token string) (string, error) {
	resp, err := getAPIResource(url, "crypto/certificatekeypairs", token, "")
	if err != nil {
		return "", err
	}
	keyPairs := &certificateKeypairs{}
	if err := json.Unmarshal(resp, &keyPairs); err != nil {
		return "", err
	}
	for _, key := range keyPairs.Results {
		if key.Name == "authentik Self-signed Certificate" {
			return key.PK, err
		}
	}
	return "", errors.New("certificate keypair not found")
}

// createOIDCProvier creates a new openid connection auth provider.
func createOIDCProvider(name, url, token, flow, signingKey string, mappings []string) (int, string, string, error) {
	client_id, err := password.Generate(40, 30, 0, false, true)
	if err != nil {
		return -1, "", "", err
	}
	client_secret, err := password.Generate(128, 96, 0, false, true)
	if err != nil {
		return -1, "", "", err
	}
	body := map[string]interface{}{
		"name":               name,
		"authorization_flow": flow,
		"client_type":        "confidential",
		"client_id":          client_id,
		"client_secret":      client_secret,
		"property_mappings":  mappings,
		"signing_key":        signingKey,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return -1, "", "", err
	}
	resp, err := postAPIResource(url, "providers/oauth2", token, data)
	if err != nil {
		return -1, "", "", err
	}
	provider := map[string]interface{}{}
	if err := json.Unmarshal(resp, &provider); err != nil {
		return -1, "", "", err
	}
	return int(provider["pk"].(float64)), client_id, client_secret, nil
}

// createApplication creates a new application.
func createApplication(provider int, name, url, token string) error {
	body := map[string]interface{}{
		"name":     name,
		"slug":     name,
		"provider": provider,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = postAPIResource(url, "core/applications", token, data)
	return err
}

// CreateOIDCClient creates a consumable end to end oidc client.
func CreateOIDCClient(name, namespace string) (string, string, error) {
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		return "", "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	uid, err := chartutils.UniqueID(namespace)
	if err != nil {
		return "", "", err
	}
	serviceURL := fmt.Sprintf("http://authentik-%s.%s.svc.cluster.local", uid, namespace)
	if err := healthCheckService(serviceURL, 2, ctx); err != nil {
		return "", "", err
	}
	token, err := getAPIToken(namespace, dispatcher.Clientset())
	if err != nil {
		return "", "", err
	}
	mappings, err := getPropertyMappings(serviceURL, token)
	if err != nil {
		return "", "", err
	}
	signingKey, err := getCertificateKeypair(serviceURL, token)
	if err != nil {
		return "", "", err
	}
	flow, err := getAuthorizationFlow(serviceURL, token)
	if err != nil {
		return "", "", err
	}
	pk, id, secret, err := createOIDCProvider(name, serviceURL, token, flow, signingKey, mappings)
	if err != nil {
		return "", "", err
	}
	if err := createApplication(pk, name, serviceURL, token); err != nil {
		return "", "", err
	}
	return id, secret, nil
}

// GetOIDCDiscoveryURL returns the service's oidc discovery url.
func GetOIDCDiscoveryURL(domain, service string, port uint16, insecure bool) string {
	scheme := "https"
	if insecure {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s.%s:%d/application/o/%s/", scheme, componentName, domain, port, service)
}

// GetOIDCEndpoint returns the service's root endpoint.
func GetOIDCEndpoint(domain string, port uint16, insecure bool) string {
	scheme := "https"
	if insecure {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s.%s:%d/", scheme, componentName, domain, port)
}
