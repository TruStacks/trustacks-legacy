package loki

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetChart(t *testing.T) {
	path, err := (&Loki{}).GetChart()
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
	assert.Equal(t, "loki.tgz", charts[0].Name())
}

func TestGetValues(t *testing.T) {
	t.Run("insecure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "local.gd", Port: 8081, Insecure: true}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, uint8(1), v.Loki.CommonConfig.ReplicationFactor)
		assert.Equal(t, "filesystem", v.Loki.Storage.Type)
		assert.False(t, v.Loki.AuthEnabled)
		assert.False(t, v.Monitoring.SelfMonitoring.Enabled)
		assert.False(t, v.Monitoring.SelfMonitoring.GrafanaAgent.InstallOperator)
		assert.False(t, v.Monitoring.SelfMonitoring.LokiCanary.Enabled)
		assert.False(t, v.Monitoring.ServiceMonitor.Enabled)
		assert.False(t, v.Test.Enabled)
		assert.Equal(t, "restic", v.SingleBinary.ExtraContainers[0].Name)
		assert.Equal(t, "restic/restic", v.SingleBinary.ExtraContainers[0].Image)
		assert.Equal(t, []string{"/bin/sh"}, v.SingleBinary.ExtraContainers[0].Command)
		assert.Equal(t, []string{"-c", "sleep infinity"}, v.SingleBinary.ExtraContainers[0].Args)
		assert.Equal(t, "RESTIC_REPOSITORY", v.SingleBinary.ExtraContainers[0].Env[0].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[0].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "url", v.SingleBinary.ExtraContainers[0].Env[0].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "RESTIC_PASSWORD", v.SingleBinary.ExtraContainers[0].Env[1].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[1].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "password", v.SingleBinary.ExtraContainers[0].Env[1].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_ACCESS_KEY_ID", v.SingleBinary.ExtraContainers[0].Env[2].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "access-key-id", v.SingleBinary.ExtraContainers[0].Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_SECRET_ACCESS_KEY", v.SingleBinary.ExtraContainers[0].Env[3].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "secret-access-key", v.SingleBinary.ExtraContainers[0].Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "backup", v.SingleBinary.ExtraContainers[0].VolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.SingleBinary.ExtraContainers[0].VolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.SingleBinary.ExtraContainers[0].VolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.SingleBinary.ExtraContainers[0].VolumeMounts[1].MountPath)
		assert.Equal(t, "backup", v.SingleBinary.ExtraVolumes[0].Name)
		assert.Equal(t, "1Gi", v.SingleBinary.ExtraVolumes[0].EmptyDir.SizeLimit)
		assert.Equal(t, "restore", v.SingleBinary.ExtraVolumes[1].Name)
		assert.Equal(t, "1Gi", v.SingleBinary.ExtraVolumes[1].EmptyDir.SizeLimit)
		assert.Equal(t, "backup", v.SingleBinary.ExtraVolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.SingleBinary.ExtraVolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.SingleBinary.ExtraVolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.SingleBinary.ExtraVolumeMounts[1].MountPath)
		assert.Equal(t, "grafana.local.gd", v.Grafana.GrafanaIni["server"].(map[string]interface{})["domain"])
		assert.Equal(t, "http://grafana.local.gd:8081", v.Grafana.GrafanaIni["server"].(map[string]interface{})["root_url"])
		assert.True(t, v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["enabled"].(bool))
		assert.Equal(t, "SSO", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["name"])
		assert.Equal(t, "http://authentik.local.gd:8081/application/o/authorize/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["auth_url"])
		assert.Equal(t, "http://authentik.local.gd:8081/application/o/token/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["token_url"])
		assert.Equal(t, "http://authentik.local.gd:8081/application/o/userinfo/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["api_url"])
		assert.Equal(t, "$__file{/etc/secrets/auth_generic_oauth/client-id}", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["client_id"])
		assert.Equal(t, "$__file{/etc/secrets/auth_generic_oauth/client-secret}", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["client_secret"])
		assert.True(t, v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["allow_assign_grafana_admin"].(bool))
		assert.Equal(
			t,
			"contains(groups[*], 'admins') && 'GrafanaAdmin' || contains(groups[*], 'editors') && 'Editor' || 'Viewer'",
			v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["role_attribute_path"],
		)
		assert.True(t, v.Grafana.Ingress.Enabled)
		assert.Equal(t, "grafana.local.gd", v.Grafana.Ingress.Hosts[0])
		assert.Empty(t, v.Grafana.Ingress.Annotations)
		assert.Empty(t, v.Grafana.Ingress.TLS)
		assert.Equal(t, "oidc-client-credentials", v.Grafana.ExtraSecretMounts[0].Name)
		assert.Equal(t, "loki-oidc-client", v.Grafana.ExtraSecretMounts[0].SecretName)
		assert.Equal(t, "/etc/secrets/auth_generic_oauth", v.Grafana.ExtraSecretMounts[0].MountPath)
		assert.True(t, v.Grafana.ExtraSecretMounts[0].ReadOnly)
		assert.Equal(t, 1, v.Grafana.Datasources.DatasourcesYaml.ApiVersion)
		assert.Equal(t, "Loki", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Name)
		assert.Equal(t, "loki", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Type)
		assert.Equal(t, "direct", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Access)
		assert.Equal(t, "http://loki:3100", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].URL)
		var grafanaExtraContainers []valuesGrafanaExtraContainers
		if err := yaml.Unmarshal([]byte(v.Grafana.ExtraContainers), &grafanaExtraContainers); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "oidc-auth-proxy", grafanaExtraContainers[0].Name)
		assert.Equal(t, "quay.io/trustacks/local-gd-proxy", grafanaExtraContainers[0].Image)
		assert.Equal(t, "UPSTREAM", grafanaExtraContainers[0].Env[0].Name)
		assert.Equal(t, "authentik-9f86d08", grafanaExtraContainers[0].Env[0].Value)
		assert.Equal(t, "LISTEN_PORT", grafanaExtraContainers[0].Env[1].Name)
		assert.Equal(t, "8081", grafanaExtraContainers[0].Env[1].Value)
		assert.Equal(t, "SERVICE", grafanaExtraContainers[0].Env[2].Name)
		assert.Equal(t, "authentik", grafanaExtraContainers[0].Env[2].Value)
	})

	t.Run("secure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "test.trustacks.io", Port: 443}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, uint8(1), v.Loki.CommonConfig.ReplicationFactor)
		assert.Equal(t, "filesystem", v.Loki.Storage.Type)
		assert.False(t, v.Loki.AuthEnabled)
		assert.False(t, v.Monitoring.SelfMonitoring.Enabled)
		assert.False(t, v.Monitoring.SelfMonitoring.GrafanaAgent.InstallOperator)
		assert.False(t, v.Monitoring.SelfMonitoring.LokiCanary.Enabled)
		assert.False(t, v.Monitoring.ServiceMonitor.Enabled)
		assert.False(t, v.Test.Enabled)
		assert.Equal(t, "restic", v.SingleBinary.ExtraContainers[0].Name)
		assert.Equal(t, "restic/restic", v.SingleBinary.ExtraContainers[0].Image)
		assert.Equal(t, []string{"/bin/sh"}, v.SingleBinary.ExtraContainers[0].Command)
		assert.Equal(t, []string{"-c", "sleep infinity"}, v.SingleBinary.ExtraContainers[0].Args)
		assert.Equal(t, "RESTIC_REPOSITORY", v.SingleBinary.ExtraContainers[0].Env[0].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[0].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "url", v.SingleBinary.ExtraContainers[0].Env[0].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "RESTIC_PASSWORD", v.SingleBinary.ExtraContainers[0].Env[1].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[1].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "password", v.SingleBinary.ExtraContainers[0].Env[1].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_ACCESS_KEY_ID", v.SingleBinary.ExtraContainers[0].Env[2].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "access-key-id", v.SingleBinary.ExtraContainers[0].Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_SECRET_ACCESS_KEY", v.SingleBinary.ExtraContainers[0].Env[3].Name)
		assert.Equal(t, "ts-storage-config", v.SingleBinary.ExtraContainers[0].Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "secret-access-key", v.SingleBinary.ExtraContainers[0].Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "backup", v.SingleBinary.ExtraContainers[0].VolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.SingleBinary.ExtraContainers[0].VolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.SingleBinary.ExtraContainers[0].VolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.SingleBinary.ExtraContainers[0].VolumeMounts[1].MountPath)
		assert.Equal(t, "backup", v.SingleBinary.ExtraVolumes[0].Name)
		assert.Equal(t, "1Gi", v.SingleBinary.ExtraVolumes[0].EmptyDir.SizeLimit)
		assert.Equal(t, "restore", v.SingleBinary.ExtraVolumes[1].Name)
		assert.Equal(t, "1Gi", v.SingleBinary.ExtraVolumes[1].EmptyDir.SizeLimit)
		assert.Equal(t, "backup", v.SingleBinary.ExtraVolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.SingleBinary.ExtraVolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.SingleBinary.ExtraVolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.SingleBinary.ExtraVolumeMounts[1].MountPath)
		assert.Equal(t, "grafana.test.trustacks.io", v.Grafana.GrafanaIni["server"].(map[string]interface{})["domain"])
		assert.Equal(t, "https://grafana.test.trustacks.io:443", v.Grafana.GrafanaIni["server"].(map[string]interface{})["root_url"])
		assert.True(t, v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["enabled"].(bool))
		assert.Equal(t, "SSO", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["name"])
		assert.Equal(t, "https://authentik.test.trustacks.io:443/application/o/authorize/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["auth_url"])
		assert.Equal(t, "https://authentik.test.trustacks.io:443/application/o/token/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["token_url"])
		assert.Equal(t, "https://authentik.test.trustacks.io:443/application/o/userinfo/", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["api_url"])
		assert.Equal(t, "$__file{/etc/secrets/auth_generic_oauth/client-id}", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["client_id"])
		assert.Equal(t, "$__file{/etc/secrets/auth_generic_oauth/client-secret}", v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["client_secret"])
		assert.True(t, v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["allow_assign_grafana_admin"].(bool))
		assert.Equal(
			t,
			"contains(groups[*], 'admins') && 'GrafanaAdmin' || contains(groups[*], 'editors') && 'Editor' || 'Viewer'",
			v.Grafana.GrafanaIni["auth.generic_oauth"].(map[string]interface{})["role_attribute_path"],
		)
		assert.True(t, v.Grafana.Ingress.Enabled)
		assert.Equal(t, "grafana.test.trustacks.io", v.Grafana.Ingress.Hosts[0])
		assert.Equal(t, "ts-system", v.Grafana.Ingress.Annotations["cert-manager.io/cluster-issuer"])
		assert.Equal(t, "ts-system", v.Grafana.Ingress.Annotations["kubernetes.io/ingress.class"])
		assert.Equal(t, "grafana.test.trustacks.io", v.Grafana.Ingress.TLS[0].Hosts[0])
		assert.Equal(t, "grafana-ingress-tls-cert", v.Grafana.Ingress.TLS[0].SecretName)
		assert.Equal(t, "oidc-client-credentials", v.Grafana.ExtraSecretMounts[0].Name)
		assert.Equal(t, "loki-oidc-client", v.Grafana.ExtraSecretMounts[0].SecretName)
		assert.Equal(t, "/etc/secrets/auth_generic_oauth", v.Grafana.ExtraSecretMounts[0].MountPath)
		assert.True(t, v.Grafana.ExtraSecretMounts[0].ReadOnly)
		assert.Equal(t, 1, v.Grafana.Datasources.DatasourcesYaml.ApiVersion)
		assert.Equal(t, "Loki", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Name)
		assert.Equal(t, "loki", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Type)
		assert.Equal(t, "direct", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].Access)
		assert.Equal(t, "http://loki:3100", v.Grafana.Datasources.DatasourcesYaml.Datasources[0].URL)
		assert.Empty(t, v.Grafana.ExtraContainers)
	})
}

func TestBackup(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Loki{}).backup(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"loki-0",
			"single-binary",
			`cd /tmp/backup && tar czf loki-data -C /var loki`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"loki-0",
			"restic",
			`restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/loki-data`,
			"test",
		},
	)
}

func TestRestore(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Loki{}).restore(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"loki-0",
			"restic",
			`restic restore latest --target /tmp/restore --include /tmp/backup/loki-data`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"loki-0",
			"single-binary",
			`cd /tmp/restore && tar xf /tmp/restore/tmp/backup/loki-data && cp -R ./loki /var/loki`,
			"test",
		},
	)
}

func TestCreateOIDCClientSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if err := createOIDCClientSecret("client-id", "client-secret", "test", clientset); err != nil {
		t.Fatal(err)
	}
	secret, err := clientset.CoreV1().Secrets("test").Get(context.TODO(), "loki-oidc-client", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "client-id", string(secret.Data["client-id"]))
	assert.Equal(t, "client-secret", string(secret.Data["client-secret"]))
}
