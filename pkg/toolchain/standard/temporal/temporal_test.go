package temporal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend/storagetest"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/chartutils"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetChart(t *testing.T) {
	path, err := (&Temporal{}).GetChart()
	if err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	expectedFiles := []string{"Chart.yaml", "charts"}
	for _, f := range files {
		assert.Contains(t, expectedFiles, f.Name())
	}
	charts, err := os.ReadDir(filepath.Join(path, "charts"))
	if err != nil {
		t.Fatal(err)
	}
	// check that the chart version matches
	fd, err := os.Open(fmt.Sprintf("%s/Chart.yaml", path))
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(fd)
	if err != nil {
		t.Fatal(err)
	}
	chartYaml := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &chartYaml); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, chartVersion, chartYaml["version"].(string))
	assert.Equal(t, "temporal.tgz", charts[0].Name())
}

func TestGetValues(t *testing.T) {
	uid, err := chartutils.UniqueID("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("insecure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "local.gd", Port: 8081, Insecure: true}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, uint8(1), v.Server.ReplicaCount)
		assert.Equal(t, "sql", v.Server.Config.Persistence.Default.Driver)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Default.SQL.Database)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Default.SQL.User)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Default.SQL.Host)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Default.SQL.ExistingSecret)
		assert.Equal(t, "sql", v.Server.Config.Persistence.Default.Driver)
		assert.Equal(t, "temporal_visibility", v.Server.Config.Persistence.Visibility.SQL.Database)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Visibility.SQL.User)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Visibility.SQL.Host)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Visibility.SQL.ExistingSecret)
		assert.True(t, v.Web.Ingress.Enabled)
		assert.Empty(t, v.Web.Ingress.Annotations)
		assert.Empty(t, v.Web.Ingress.TLS)
		assert.Equal(t, "temporal.local.gd", v.Web.Ingress.Hosts[0])
		assert.Equal(t, "TEMPORAL_AUTH_ENABLED", v.Web.Env[0].Name)
		assert.Equal(t, "true", v.Web.Env[0].Value)
		assert.Equal(t, "TEMPORAL_AUTH_PROVIDER_URL", v.Web.Env[1].Name)
		assert.Equal(t, "http://authentik.local.gd:8081/application/o/temporal/", v.Web.Env[1].Value)
		assert.Equal(t, "TEMPORAL_AUTH_CLIENT_ID", v.Web.Env[2].Name)
		assert.Equal(t, "temporal-oidc-client", v.Web.Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "client-id", v.Web.Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "TEMPORAL_AUTH_CLIENT_SECRET", v.Web.Env[3].Name)
		assert.Equal(t, "temporal-oidc-client", v.Web.Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "client-secret", v.Web.Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "TEMPORAL_AUTH_CALLBACK_URL", v.Web.Env[4].Name)
		assert.Equal(t, "http://temporal.local.gd:8081/auth/sso/callback", v.Web.Env[4].Value)
		assert.Equal(t, "oidc-auth-proxy", v.Web.SidecarContainers[0].Name)
		assert.Equal(t, "quay.io/trustacks/local-gd-proxy", v.Web.SidecarContainers[0].Image)
		assert.Equal(t, "UPSTREAM", v.Web.SidecarContainers[0].Env[0].Name)
		assert.Equal(t, fmt.Sprintf("authentik-%s", uid), v.Web.SidecarContainers[0].Env[0].Value)
		assert.Equal(t, "LISTEN_PORT", v.Web.SidecarContainers[0].Env[1].Name)
		assert.Equal(t, "8081", v.Web.SidecarContainers[0].Env[1].Value)
		assert.Equal(t, "SERVICE", v.Web.SidecarContainers[0].Env[2].Name)
		assert.Equal(t, "authentik", v.Web.SidecarContainers[0].Env[2].Value)
		assert.False(t, v.Prometheus.Enabled)
		assert.False(t, v.Grafana.Enabled)
		assert.False(t, v.Elasticsearch.Enabled)
		assert.False(t, v.Cassandra.Enabled)
		assert.True(t, v.MySQL.Enabled)
		assert.Equal(t, "temporal", v.MySQL.Auth.Username)
		assert.Equal(t, "temporal-mysql", v.MySQL.Auth.ExistingSecret)
		assert.Equal(t, v.MySQL.InitdbScripts["setup.sh"], `
export MYSQL_PWD=$MYSQL_ROOT_PASSWORD 
mysql -u root -e "create database temporal"
mysql -u root -e "grant all privileges on temporal.* to 'temporal'@'%'"
mysql -u root -e "create database temporal_visibility"
mysql -u root -e "grant all privileges on temporal_visibility.* to 'temporal'@'%'"
`)
		assert.False(t, v.Schema.Setup.Enabled)
		assert.False(t, v.Schema.Update.Enabled)
	})

	t.Run("secure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "test.trustacks.io", Port: 443}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, uint8(1), v.Server.ReplicaCount)
		assert.Equal(t, "sql", v.Server.Config.Persistence.Default.Driver)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Default.SQL.Database)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Default.SQL.User)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Default.SQL.Host)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Default.SQL.ExistingSecret)
		assert.Equal(t, "sql", v.Server.Config.Persistence.Default.Driver)
		assert.Equal(t, "temporal_visibility", v.Server.Config.Persistence.Visibility.SQL.Database)
		assert.Equal(t, "temporal", v.Server.Config.Persistence.Visibility.SQL.User)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Visibility.SQL.Host)
		assert.Equal(t, "temporal-mysql", v.Server.Config.Persistence.Visibility.SQL.ExistingSecret)
		assert.True(t, v.Web.Ingress.Enabled)
		assert.Equal(t, "ts-system", v.Web.Ingress.Annotations["cert-manager.io/cluster-issuer"])
		assert.Equal(t, "ts-system", v.Web.Ingress.Annotations["kubernetes.io/ingress.class"])
		assert.Contains(t, v.Web.Ingress.TLS[0].Hosts, "temporal.test.trustacks.io")
		assert.Equal(t, "temporal-ingress-tls-cert", v.Web.Ingress.TLS[0].SecretName)
		assert.Equal(t, "temporal.test.trustacks.io", v.Web.Ingress.Hosts[0])
		assert.Equal(t, "TEMPORAL_AUTH_ENABLED", v.Web.Env[0].Name)
		assert.Equal(t, "true", v.Web.Env[0].Value)
		assert.Equal(t, "TEMPORAL_AUTH_PROVIDER_URL", v.Web.Env[1].Name)
		assert.Equal(t, "https://authentik.test.trustacks.io:443/application/o/temporal/", v.Web.Env[1].Value)
		assert.Equal(t, "TEMPORAL_AUTH_CLIENT_ID", v.Web.Env[2].Name)
		assert.Equal(t, "temporal-oidc-client", v.Web.Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "client-id", v.Web.Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "TEMPORAL_AUTH_CLIENT_SECRET", v.Web.Env[3].Name)
		assert.Equal(t, "temporal-oidc-client", v.Web.Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "client-secret", v.Web.Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "TEMPORAL_AUTH_CALLBACK_URL", v.Web.Env[4].Name)
		assert.Equal(t, "https://temporal.test.trustacks.io:443/auth/sso/callback", v.Web.Env[4].Value)
		assert.Empty(t, v.Web.SidecarContainers)
		assert.False(t, v.Prometheus.Enabled)
		assert.False(t, v.Grafana.Enabled)
		assert.False(t, v.Elasticsearch.Enabled)
		assert.False(t, v.Cassandra.Enabled)
		assert.True(t, v.MySQL.Enabled)
		assert.Equal(t, "temporal", v.MySQL.Auth.Username)
		assert.Equal(t, "temporal-mysql", v.MySQL.Auth.ExistingSecret)
		assert.False(t, v.Schema.Setup.Enabled)
		assert.False(t, v.Schema.Update.Enabled)
	})
}

func TestBackup(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Temporal{}).backup(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"temporal-mysql-0",
			"mysql",
			`MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysqldump -u root temporal > /tmp/backup/temporal-mysql`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"temporal-mysql-0",
			"restic",
			`restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/temporal-mysql`,
			"test",
		},
	)
}

func TestRestore(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Temporal{}).restore(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"temporal-mysql-0",
			"restic",
			`restic restore latest --target /tmp/restore --include /tmp/backup/temporal-mysql`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"temporal-mysql-0",
			"mysql",
			`MYSQL_PWD=$MYSQL_ROOT_PASSWORD mysql -u root temporal < /tmp/backup/temporal-mysql`,
			"test",
		},
	)
}

func TestCreateOIDCClientSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if err := createOIDCClientSecret("client-id", "client-secret", "test", clientset); err != nil {
		t.Fatal(err)
	}
	secret, err := clientset.CoreV1().Secrets("test").Get(context.TODO(), "temporal-oidc-client", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "client-id", string(secret.Data["client-id"]))
	assert.Equal(t, "client-secret", string(secret.Data["client-secret"]))
}

func TestCreateMySQLPasswordSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	pwd, err := createMySQLPasswordSecret("test", clientset)
	if err != nil {
		t.Fatal(err)
	}
	secret, err := clientset.CoreV1().Secrets("test").Get(context.TODO(), "temporal-mysql", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, secret.Data["mysql-root-password"], 32)
	assert.Len(t, secret.Data["mysql-replication-password"], 32)
	assert.Equal(t, pwd, string(secret.Data["mysql-password"]))
	assert.Equal(t, pwd, string(secret.Data["password"]))
}

func TestGetMySQLPasswordSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	pwd, err := createMySQLPasswordSecret("test", clientset)
	if err != nil {
		t.Fatal(err)
	}
	mysqlPwd, err := getMySQLPasswordSecret("test", clientset)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mysqlPwd, pwd)
}

func TestCreateMySQLSchemaUpdateJob(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if err := createMySQLSchemaUpdateJob("test", "password123!", clientset); err != nil {
		t.Fatal(err)
	}
	job, err := clientset.BatchV1().Jobs("test").Get(context.TODO(), "temporal-mysql-schema-update", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	ttlSeconds := int32(0)
	assert.Equal(t, "temporal-mysql-tool", job.Spec.Template.Spec.Containers[0].Name)
	assert.Equal(t, "temporalio/admin-tools:1.19.0", job.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, "/bin/sh", job.Spec.Template.Spec.Containers[0].Command[0])
	assert.Equal(t, "SQL_PLUGIN", job.Spec.Template.Spec.Containers[0].Env[0].Name)
	assert.Equal(t, "mysql", job.Spec.Template.Spec.Containers[0].Env[0].Value)
	assert.Equal(t, "SQL_HOST", job.Spec.Template.Spec.Containers[0].Env[1].Name)
	assert.Equal(t, "temporal-mysql", job.Spec.Template.Spec.Containers[0].Env[1].Value)
	assert.Equal(t, "SQL_PORT", job.Spec.Template.Spec.Containers[0].Env[2].Name)
	assert.Equal(t, "3306", job.Spec.Template.Spec.Containers[0].Env[2].Value)
	assert.Equal(t, "SQL_USER", job.Spec.Template.Spec.Containers[0].Env[3].Name)
	assert.Equal(t, "temporal", job.Spec.Template.Spec.Containers[0].Env[3].Value)
	assert.Equal(t, "SQL_PASSWORD", job.Spec.Template.Spec.Containers[0].Env[4].Name)
	assert.Equal(t, "password123!", job.Spec.Template.Spec.Containers[0].Env[4].Value)
	assert.Equal(t, "-c", job.Spec.Template.Spec.Containers[0].Args[0])
	assert.Equal(t, `
while true; do
	curl $SQL_HOST:$SQL_PORT > /dev/null
	if [ "$?" == "1" ]; then
		temporal-sql-tool --db temporal setup-schema -v 0.0
		temporal-sql-tool --db temporal update-schema -d ./schema/mysql/v57/temporal/versioned
		temporal-sql-tool --db temporal_visibility setup-schema -v 0.0
		temporal-sql-tool --db temporal_visibility update-schema -d ./schema/mysql/v57/visibility/versioned
		break
	fi
	sleep 1
done
`,
		job.Spec.Template.Spec.Containers[0].Args[1])
	assert.Equal(t, corev1.RestartPolicyOnFailure, job.Spec.Template.Spec.RestartPolicy)
	assert.Equal(t, &ttlSeconds, job.Spec.TTLSecondsAfterFinished)
}

func TestTemporalLifecycleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()
	namespace := "temporal-integration-test"
	dispatcher, err := client.NewDispatcher(namespace)
	if err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.CreateNamespace(); err != nil {
		t.Fatal(err)
	}
	storageConfig, err := storagetest.NewTestStorageConfig("temporal")
	if err != nil {
		t.Fatal(err)
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		t.Fatal(err)
	}
	defer storagetest.PurgeBucket("temporal")
	pfl := profile.Profile{Domain: "temporal-integration-test.local.gd", Port: 8081, Insecure: true}
	sso := authentik.New(pfl)
	if err := sso.Install(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	c := New(pfl)
	if err := c.Install(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Upgrade(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Rollback(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := c.Uninstall(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := sso.Uninstall(dispatcher, namespace); err != nil {
		t.Fatal(err)
	}
	if err := dispatcher.DeleteNamespace(); err != nil {
		t.Fatal(err)
	}
}
