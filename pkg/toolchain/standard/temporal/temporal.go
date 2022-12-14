package temporal

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/chartutils"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// componentName is the name of the component.
	componentName = "temporal"
	// temporalVersion is the version of the temporal application.
	temporalVersion = "1.19.0"
)

//go:embed temporal-0.19.0.tgz
var chartArchive []byte

type Temporal struct {
	profile profile.Profile
}

type values struct {
	Server           valuesServer        `yaml:"server"`
	Web              valuesWeb           `yaml:"web"`
	Prometheus       valuesPrometheus    `yaml:"prometheus"`
	Grafana          valuesGrafana       `yaml:"grafana"`
	Elasticsearch    valuesElasticsearch `yaml:"elasticsearch"`
	Cassandra        valuesCassandra     `yaml:"cassandra"`
	MySQL            valuesMySQL         `yaml:"mysql"`
	Schema           valuesSchema        `yaml:"schema"`
	FullnameOverride string              `yaml:"fullnameOverride"`
}

type valuesServer struct {
	ReplicaCount uint8              `yaml:"replicaCount"`
	Config       valuesServerConfig `yaml:"config"`
}

type valuesServerConfig struct {
	Persistence valuesServerConfigPersistence `yaml:"persistence"`
}

type valuesServerConfigPersistence struct {
	Default    valuesServerConfigPersistenceDriver `yaml:"default"`
	Visibility valuesServerConfigPersistenceDriver `yaml:"visibility"`
}

type valuesServerConfigPersistenceDriver struct {
	Driver string                                 `yaml:"driver"`
	SQL    valuesServerConfigPersistenceDriverSQL `yaml:"sql"`
}

type valuesServerConfigPersistenceDriverSQL struct {
	User           string `yaml:"user"`
	Host           string `yaml:"host"`
	Database       string `yaml:"database"`
	ExistingSecret string `yaml:"existingSecret"`
}

type valuesWeb struct {
	Ingress           valuesWebIngress          `yaml:"ingress"`
	Env               []valuesWebEnv            `yaml:"env"`
	SidecarContainers []valuesWebSideContainers `yaml:"sidecarContainers,omitempty"`
}

type valuesWebIngress struct {
	Enabled     bool                  `yaml:"enabled"`
	Annotations map[string]string     `yaml:"annotations,omitempty"`
	TLS         []valuesWebIngressTLS `yaml:"tls"`
	Hosts       []string              `yaml:"hosts"`
}

type valuesWebIngressTLS struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type valuesWebEnv struct {
	Name      string                 `yaml:"name"`
	Value     string                 `yaml:"value,omitempty"`
	ValueFrom valuesWebEnvValuesFrom `yaml:"valueFrom,omitempty"`
}

type valuesWebEnvValuesFrom struct {
	SecretKeyRef valuesWebEnvValuesFromSecretKeyRef `yaml:"secretKeyRef"`
}

type valuesWebEnvValuesFromSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type valuesWebSideContainers struct {
	Name  string                       `yaml:"name"`
	Image string                       `yaml:"image"`
	Env   []valuesWebSideContainersEnv `yaml:"env"`
}

type valuesWebSideContainersEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type valuesPrometheus struct {
	Enabled bool `yaml:"enabled"`
}

type valuesGrafana struct {
	Enabled bool `yaml:"enabled"`
}

type valuesElasticsearch struct {
	Enabled bool `yaml:"enabled"`
}

type valuesCassandra struct {
	Enabled bool `yaml:"enabled"`
}

type valuesMySQL struct {
	Enabled bool               `yaml:"enabled"`
	Auth    valuesMySQLAuth    `yaml:"auth"`
	Primary valuesMySQLPrimary `yaml:"primary"`
}

type valuesMySQLPrimary struct {
	ExtraVolumes      []valuesMySQLPrimaryExtraVolumes      `yaml:"extraVolumes"`
	ExtraVolumeMounts []valuesMySQLPrimaryExtraVolumeMounts `yaml:"extraVolumeMounts"`
	Sidecars          []valuesMySQLPrimarySidecars          `yaml:"sidecars"`
}

// valuesMySQLPrimaryExtraVolumeMounts contains volume mount
// parameters.
type valuesMySQLPrimaryExtraVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

// valuesMySQLPrimaryExtraVolumes contains empty dir volumes.
type valuesMySQLPrimaryExtraVolumes struct {
	Name     string                                 `yaml:"name"`
	EmptyDir valuesMySQLPrimaryExtraVolumesEmptyDir `yaml:"emptyDir"`
}

// valuesMySQLPrimaryExtraVolumesEmptyDir contains the empty dir
// size limit.
type valuesMySQLPrimaryExtraVolumesEmptyDir struct {
	SizeLimit string `yaml:"sizeLimit"`
}

// valuesMySQLPrimarySidecars contains sidecar parameters.
type valuesMySQLPrimarySidecars struct {
	Name         string                                   `yaml:"name"`
	Image        string                                   `yaml:"image"`
	Command      []string                                 `yaml:"command"`
	Args         []string                                 `yaml:"args"`
	Env          []valuesMySQLPrimarySidecarsEnv          `yaml:"env"`
	VolumeMounts []valuesMySQLPrimarySidecarsVolumeMounts `yaml:"volumeMounts"`
}

// valuesMySQLPrimarySidecarsVolumeMounts contains sidecar volume mounts.
type valuesMySQLPrimarySidecarsVolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

// valuesMySQLPrimarySidecarsEnv contains sidecar environment variables.
type valuesMySQLPrimarySidecarsEnv struct {
	Name      string                                 `yaml:"name"`
	ValueFrom valuesMySQLPrimarySidecarsEnvValueFrom `yaml:"valueFrom,omitempty"`
}

// valuesMySQLPrimarySidecarsEnvValueFrom contains environemt variable
// secret key references.
type valuesMySQLPrimarySidecarsEnvValueFrom struct {
	SecretKeyRef valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef `yaml:"secretKeyRef"`
}

// valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef contains
// the secret key reference name and key.
type valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

type valuesMySQLAuth struct {
	CreateDatabase bool   `yaml:"createDatabase"`
	Database       string `yaml:"database"`
	Username       string `yaml:"username"`
	ExistingSecret string `yaml:"existingSecret"`
}

type valuesSchema struct {
	Setup  valuesSchemaConfig `yaml:"setup"`
	Update valuesSchemaConfig `yaml:"update"`
}

type valuesSchemaConfig struct {
	Enabled bool `yaml:"enabled"`
}

func (c *Temporal) GetValues(namespace string) (interface{}, error) {
	urlScheme := "https"
	if c.profile.Insecure {
		urlScheme = "http"
	}
	uid, err := chartutils.UniqueID(namespace)
	if err != nil {
		return nil, err
	}
	v := &values{
		Server: valuesServer{
			ReplicaCount: 1,
			Config: valuesServerConfig{
				Persistence: valuesServerConfigPersistence{
					Default: valuesServerConfigPersistenceDriver{
						Driver: "sql",
						SQL: valuesServerConfigPersistenceDriverSQL{
							Database:       "temporal",
							User:           "temporal",
							Host:           "temporal-mysql",
							ExistingSecret: "temporal-mysql",
						},
					},
					Visibility: valuesServerConfigPersistenceDriver{
						Driver: "sql",
						SQL: valuesServerConfigPersistenceDriverSQL{
							Database:       "temporal",
							User:           "temporal",
							Host:           "temporal-mysql",
							ExistingSecret: "temporal-mysql",
						},
					},
				},
			},
		},
		Web: valuesWeb{
			Ingress: valuesWebIngress{
				Enabled: true,
				Hosts: []string{
					fmt.Sprintf("%s.%s", componentName, c.profile.Domain),
				},
			},
			Env: []valuesWebEnv{
				{
					Name:  "TEMPORAL_AUTH_ENABLED",
					Value: "true",
				},
				{
					Name:  "TEMPORAL_AUTH_PROVIDER_URL",
					Value: authentik.GetOIDCDiscoveryURL(c.profile.Domain, componentName, c.profile.Port, c.profile.Insecure),
				},
				{
					Name: "TEMPORAL_AUTH_CLIENT_ID",
					ValueFrom: valuesWebEnvValuesFrom{
						SecretKeyRef: valuesWebEnvValuesFromSecretKeyRef{
							Name: "temporal-oidc-client",
							Key:  "client-id",
						},
					},
				},
				{
					Name: "TEMPORAL_AUTH_CLIENT_SECRET",
					ValueFrom: valuesWebEnvValuesFrom{
						SecretKeyRef: valuesWebEnvValuesFromSecretKeyRef{
							Name: "temporal-oidc-client",
							Key:  "client-secret",
						},
					},
				},
				{
					Name:  "TEMPORAL_AUTH_CALLBACK_URL",
					Value: fmt.Sprintf("%s://%s.%s:%d/auth/sso/callback", urlScheme, componentName, c.profile.Domain, c.profile.Port),
				},
			},
			SidecarContainers: []valuesWebSideContainers{
				{
					Name:  "oidc-auth-proxy",
					Image: "quay.io/trustacks/local-gd-proxy",
					Env: []valuesWebSideContainersEnv{
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
			},
		},
		Prometheus: valuesPrometheus{
			Enabled: false,
		},
		Grafana: valuesGrafana{
			Enabled: false,
		},
		Elasticsearch: valuesElasticsearch{
			Enabled: false,
		},
		Cassandra: valuesCassandra{
			Enabled: false,
		},
		MySQL: valuesMySQL{
			Enabled: true,
			Auth: valuesMySQLAuth{
				CreateDatabase: true,
				Database:       "temporal",
				Username:       "temporal",
				ExistingSecret: "temporal-mysql",
			},
			Primary: valuesMySQLPrimary{
				Sidecars: []valuesMySQLPrimarySidecars{
					{
						Name:    "restic",
						Image:   "restic/restic",
						Command: []string{"/bin/sh"},
						Args:    []string{"-c", "sleep infinity"},
						Env: []valuesMySQLPrimarySidecarsEnv{
							{
								Name: "RESTIC_REPOSITORY",
								ValueFrom: valuesMySQLPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "url",
									},
								},
							},
							{
								Name: "RESTIC_PASSWORD",
								ValueFrom: valuesMySQLPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "password",
									},
								},
							},
							{
								Name: "AWS_ACCESS_KEY_ID",
								ValueFrom: valuesMySQLPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "access-key-id",
									},
								},
							},
							{
								Name: "AWS_SECRET_ACCESS_KEY",
								ValueFrom: valuesMySQLPrimarySidecarsEnvValueFrom{
									SecretKeyRef: valuesMySQLPrimarySidecarsEnvValueFromSecretKeyRef{
										Name: "ts-storage-config",
										Key:  "secret-access-key",
									},
								},
							},
						},
						VolumeMounts: []valuesMySQLPrimarySidecarsVolumeMounts{
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
				ExtraVolumes: []valuesMySQLPrimaryExtraVolumes{
					{
						Name: "backup",
						EmptyDir: valuesMySQLPrimaryExtraVolumesEmptyDir{
							SizeLimit: "1Gi",
						},
					},
					{
						Name: "restore",
						EmptyDir: valuesMySQLPrimaryExtraVolumesEmptyDir{
							SizeLimit: "1Gi",
						},
					},
				},
				ExtraVolumeMounts: []valuesMySQLPrimaryExtraVolumeMounts{
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
		Schema: valuesSchema{
			Setup: valuesSchemaConfig{
				Enabled: false,
			},
			Update: valuesSchemaConfig{
				Enabled: false,
			},
		},
	}
	if !c.profile.Insecure {
		v.Web.Ingress.Annotations = map[string]string{
			"cert-manager.io/cluster-issuer": "ts-system",
			"kubernetes.io/ingress.class":    "ts-system",
		}
		v.Web.Ingress.TLS = []valuesWebIngressTLS{
			{
				Hosts: []string{
					fmt.Sprintf("%s.%s", componentName, c.profile.Domain),
				},
				SecretName: "temporal-ingress-tls-cert",
			},
		}
		v.Web.SidecarContainers = nil
	}
	return v, nil
}

// GetChart returns the temporal helm chart archive.
func (c *Temporal) GetChart() (string, error) {
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

func (c *Temporal) Install(dispatcher client.Dispatcher, namespace string) error {
	if err := c.preInstall(dispatcher.Clientset(), namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	temporalValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"temporal": temporalValues,
	}
	if err := dispatcher.InstallChart("temporal", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return nil
}

// Upgrade .
func (c *Temporal) Upgrade(dispatcher client.Dispatcher, namespace string) error {
	if err := c.backup(dispatcher, namespace); err != nil {
		return err
	}
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	temporalValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"temporal": temporalValues,
	}
	return dispatcher.UpgradeChart("temporal", values, time.Minute*5, chartPath)
}

// Rollback.
func (c *Temporal) Rollback(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	temporalValues, err := c.GetValues(namespace)
	if err != nil {
		return err
	}
	values := map[string]interface{}{
		"temporal": temporalValues,
	}
	if err := dispatcher.RollbackRelease("temporal", values, time.Minute*5, chartPath); err != nil {
		return err
	}
	return c.restore(dispatcher, namespace)
}

// Uninstall .
func (c *Temporal) Uninstall(dispatcher client.Dispatcher, namespace string) error {
	if err := dispatcher.UninstallChart("temporal"); err != nil {
		return err
	}
	return nil
}

// backup .
func (c *Temporal) backup(dispatcher client.Dispatcher, namespace string) error {
	cmd := `MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysqldump -u root temporal > /tmp/backup/temporal-mysql`
	if err := dispatcher.ExecCommand("temporal-mysql-0", "mysql", cmd, namespace); err != nil {
		return err
	}
	cmd = `restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/temporal-mysql`
	return dispatcher.ExecCommand("temporal-mysql-0", "restic", cmd, namespace)
}

// restore .
func (c *Temporal) restore(dispatcher client.Dispatcher, namespace string) error {
	cmd := `restic restore latest --target /tmp/restore --include /tmp/backup/temporal-mysql`
	if err := dispatcher.ExecCommand("temporal-mysql-0", "restic", cmd, namespace); err != nil {
		return err
	}
	cmd = `MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysql -u root temporal < /tmp/backup/temporal-mysql`
	return dispatcher.ExecCommand("temporal-mysql-0", "mysql", cmd, namespace)
}

func (c *Temporal) preInstall(clientset kubernetes.Interface, namespace string) error {
	clientID, clientSecret, err := authentik.CreateOIDCClient(componentName, namespace)
	if err != nil {
		return err
	}
	if err := createOIDCClientSecret(clientID, clientSecret, namespace, clientset); err != nil {
		return err
	}
	pwd, err := createMySQLPasswordSecret(namespace, clientset)
	if err != nil {
		return err
	}
	if err := createMySQLSchemaUpdateJob(namespace, pwd, clientset); err != nil {
		return err
	}
	return nil
}

func createOIDCClientSecret(clientID, clientSecret, namespace string, clientset kubernetes.Interface) error {
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "temporal-oidc-client", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "temporal-oidc-client",
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

func createMySQLPasswordSecret(namespace string, clientset kubernetes.Interface) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "temporal-mysql", metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			rootPwd, err := password.Generate(32, 10, 0, false, true)
			if err != nil {
				return "", err
			}
			pwd, err := password.Generate(32, 10, 0, false, true)
			if err != nil {
				return "", err
			}
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "temporal-mysql",
				},
				Data: map[string][]byte{
					"mysql-root-password":        []byte(rootPwd),
					"mysql-replication-password": []byte(rootPwd),
					"mysql-password":             []byte(pwd),
					"password":                   []byte(pwd),
				},
			}
			_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				return "", err
			}
		} else if !strings.Contains(err.Error(), "already exists") {
			return "", err
		}
	}
	return string(secret.Data["password"]), nil
}

func createMySQLSchemaUpdateJob(namespace, pwd string, clientset kubernetes.Interface) error {
	ttlSeconds := int32(0)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "temporal-mysql-schema-update",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "temporal-mysql-tool",
							Image: fmt.Sprintf("temporalio/admin-tools:%s", temporalVersion),
							Command: []string{
								"/bin/sh",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SQL_PLUGIN",
									Value: "mysql",
								},
								{
									Name:  "SQL_HOST",
									Value: "temporal-mysql",
								},
								{
									Name:  "SQL_PORT",
									Value: "3306",
								},
								{
									Name:  "SQL_DATABASE",
									Value: "temporal",
								},
								{
									Name:  "SQL_USER",
									Value: "temporal",
								},
								{
									Name:  "SQL_PASSWORD",
									Value: pwd,
								},
							},
							Args: []string{
								"-c",
								`
								run() {
									temporal-sql-tool setup-schema -v 0.0
									temporal-sql-tool update -schema-dir schema/mysql/v57/temporal/versioned
									temporal-sql-tool setup-schema -v 0.0
									temporal-sql-tool update -schema-dir schema/mysql/v57/visibility/versioned
								}
								until run; do :; done
								`,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
			TTLSecondsAfterFinished: &ttlSeconds,
		},
	}
	if _, err := clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func New(prof profile.Profile) *Temporal {
	return &Temporal{prof}
}
