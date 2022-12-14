package authentik

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetChart(t *testing.T) {
	path, err := (&Authentik{}).GetChart()
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
	assert.Equal(t, "authentik.tgz", charts[0].Name())
}

func TestGetValues(t *testing.T) {
	t.Run("insecure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "local.gd", Insecure: true}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, "authentik-9f86d08", v.FullnameOverride)
		assert.Len(t, v.Authentik.SecretKey, 32)
		assert.Equal(t, "authentik-postgresql", v.Authentik.Postgresql.Host)
		assert.True(t, v.Ingress.Enabled)
		assert.Empty(t, v.Ingress.Annotations)
		assert.Empty(t, v.Ingress.TLS)
		assert.Equal(t, "authentik.local.gd", v.Ingress.Hosts[0].Host)
		assert.Equal(t, "/", v.Ingress.Hosts[0].Paths[0].Path)
		assert.Equal(t, "Prefix", v.Ingress.Hosts[0].Paths[0].PathType)
		assert.True(t, v.Postgresql.Enabled)
		assert.Equal(t, "authentik-postgresql", v.Postgresql.ExistingSecret)
		assert.Equal(t, "restic", v.Postgresql.Primary.Sidecars[0].Name)
		assert.Equal(t, "restic/restic", v.Postgresql.Primary.Sidecars[0].Image)
		assert.ElementsMatch(t, v.Postgresql.Primary.Sidecars[0].Command, []string{"/bin/sh"})
		assert.ElementsMatch(t, v.Postgresql.Primary.Sidecars[0].Args, []string{"-c", "sleep infinity"})
		assert.Equal(t, "RESTIC_REPOSITORY", v.Postgresql.Primary.Sidecars[0].Env[0].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[0].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "url", v.Postgresql.Primary.Sidecars[0].Env[0].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "RESTIC_PASSWORD", v.Postgresql.Primary.Sidecars[0].Env[1].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[1].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "password", v.Postgresql.Primary.Sidecars[0].Env[1].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_ACCESS_KEY_ID", v.Postgresql.Primary.Sidecars[0].Env[2].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "access-key-id", v.Postgresql.Primary.Sidecars[0].Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_SECRET_ACCESS_KEY", v.Postgresql.Primary.Sidecars[0].Env[3].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "secret-access-key", v.Postgresql.Primary.Sidecars[0].Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "backup", v.Postgresql.Primary.Sidecars[0].VolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.Postgresql.Primary.Sidecars[0].VolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.Postgresql.Primary.Sidecars[0].VolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.Postgresql.Primary.Sidecars[0].VolumeMounts[1].MountPath)
		assert.Equal(t, "backup", v.Postgresql.Primary.ExtraVolumes[0].Name)
		assert.Equal(t, "1Gi", v.Postgresql.Primary.ExtraVolumes[0].EmptyDir.SizeLimit)
		assert.Equal(t, "backup", v.Postgresql.Primary.ExtraVolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.Postgresql.Primary.ExtraVolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.Postgresql.Primary.ExtraVolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.Postgresql.Primary.ExtraVolumeMounts[1].MountPath)
		assert.Equal(t, "1Gi", v.Postgresql.Primary.ExtraVolumes[1].EmptyDir.SizeLimit)
		assert.True(t, v.Redis.Enabled)
		assert.Equal(t, "authentik-postgresql", v.EnvValueFrom.AuthentikPostgresqlPassword.SecretKeyRef.Name)
		assert.Equal(t, "postgresql-postgres-password", v.EnvValueFrom.AuthentikPostgresqlPassword.SecretKeyRef.Key)
		assert.Equal(t, "authentik-bootstrap", v.EnvValueFrom.AuthentikBootstrapToken.SecretKeyRef.Name)
		assert.Equal(t, "api-token", v.EnvValueFrom.AuthentikBootstrapToken.SecretKeyRef.Key)
	})

	t.Run("secure-endpoint", func(t *testing.T) {
		c, err := New(profile.Profile{Domain: "test.trustacks.io"}).GetValues("test")
		if err != nil {
			t.Fatal(err)
		}
		v := c.(*values)
		assert.Equal(t, "authentik-9f86d08", v.FullnameOverride)
		assert.Len(t, v.Authentik.SecretKey, 32)
		assert.Equal(t, "authentik-postgresql", v.Authentik.Postgresql.Host)
		assert.True(t, v.Ingress.Enabled)
		assert.Equal(t, "authentik.test.trustacks.io", v.Ingress.Hosts[0].Host)
		assert.Equal(t, "/", v.Ingress.Hosts[0].Paths[0].Path)
		assert.Equal(t, "Prefix", v.Ingress.Hosts[0].Paths[0].PathType)
		assert.Equal(t, "ts-system", v.Ingress.Annotations["cert-manager.io/cluster-issuer"])
		assert.Equal(t, "ts-system", v.Ingress.Annotations["kubernetes.io/ingress.class"])
		assert.Contains(t, v.Ingress.TLS[0].Hosts, "authentik.test.trustacks.io")
		assert.Equal(t, "authentik-ingress-tls-cert", v.Ingress.TLS[0].SecretName)
		assert.True(t, v.Postgresql.Enabled)
		assert.Equal(t, "authentik-postgresql", v.Postgresql.ExistingSecret)
		assert.Equal(t, "restic", v.Postgresql.Primary.Sidecars[0].Name)
		assert.Equal(t, "restic/restic", v.Postgresql.Primary.Sidecars[0].Image)
		assert.ElementsMatch(t, v.Postgresql.Primary.Sidecars[0].Command, []string{"/bin/sh"})
		assert.ElementsMatch(t, v.Postgresql.Primary.Sidecars[0].Args, []string{"-c", "sleep infinity"})
		assert.Equal(t, "RESTIC_REPOSITORY", v.Postgresql.Primary.Sidecars[0].Env[0].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[0].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "url", v.Postgresql.Primary.Sidecars[0].Env[0].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "RESTIC_PASSWORD", v.Postgresql.Primary.Sidecars[0].Env[1].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[1].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "password", v.Postgresql.Primary.Sidecars[0].Env[1].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_ACCESS_KEY_ID", v.Postgresql.Primary.Sidecars[0].Env[2].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[2].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "access-key-id", v.Postgresql.Primary.Sidecars[0].Env[2].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "AWS_SECRET_ACCESS_KEY", v.Postgresql.Primary.Sidecars[0].Env[3].Name)
		assert.Equal(t, "ts-storage-config", v.Postgresql.Primary.Sidecars[0].Env[3].ValueFrom.SecretKeyRef.Name)
		assert.Equal(t, "secret-access-key", v.Postgresql.Primary.Sidecars[0].Env[3].ValueFrom.SecretKeyRef.Key)
		assert.Equal(t, "/tmp/backup", v.Postgresql.Primary.Sidecars[0].VolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.Postgresql.Primary.Sidecars[0].VolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.Postgresql.Primary.Sidecars[0].VolumeMounts[1].MountPath)
		assert.Equal(t, "backup", v.Postgresql.Primary.ExtraVolumes[0].Name)
		assert.Equal(t, "1Gi", v.Postgresql.Primary.ExtraVolumes[0].EmptyDir.SizeLimit)
		assert.Equal(t, "backup", v.Postgresql.Primary.ExtraVolumeMounts[0].Name)
		assert.Equal(t, "/tmp/backup", v.Postgresql.Primary.ExtraVolumeMounts[0].MountPath)
		assert.Equal(t, "restore", v.Postgresql.Primary.ExtraVolumeMounts[1].Name)
		assert.Equal(t, "/tmp/restore", v.Postgresql.Primary.ExtraVolumeMounts[1].MountPath)
		assert.Equal(t, "1Gi", v.Postgresql.Primary.ExtraVolumes[1].EmptyDir.SizeLimit)
		assert.True(t, v.Redis.Enabled)
		assert.Equal(t, "authentik-postgresql", v.EnvValueFrom.AuthentikPostgresqlPassword.SecretKeyRef.Name)
		assert.Equal(t, "postgresql-postgres-password", v.EnvValueFrom.AuthentikPostgresqlPassword.SecretKeyRef.Key)
		assert.Equal(t, "authentik-bootstrap", v.EnvValueFrom.AuthentikBootstrapToken.SecretKeyRef.Name)
		assert.Equal(t, "api-token", v.EnvValueFrom.AuthentikBootstrapToken.SecretKeyRef.Key)
	})
}

func TestBackup(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Authentik{}).backup(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"authentik-postgresql-0",
			"authentik-postgresql",
			`PGPASSWORD=$POSTGRES_PASSWORD pg_dump -U $POSTGRES_USER -F c -b -v -f /tmp/backup/authentik-postgresql $POSTGRES_DB`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"authentik-postgresql-0",
			"restic",
			`restic check; if [ "$?" == "1" ]; then restic init; fi; restic backup /tmp/backup/authentik-postgresql`,
			"test",
		},
	)
}

func TestRestore(t *testing.T) {
	dispatcher := client.NewFakeDispatcher()
	if err := (&Authentik{}).restore(dispatcher, "test"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][0],
		[]string{
			"authentik-postgresql-0",
			"restic",
			`restic restore latest --target /tmp/restore --include /tmp/backup/authentik-postgresql`,
			"test",
		},
	)
	assert.ElementsMatch(
		t,
		dispatcher.MockCalls["ExecCommand"][1],
		[]string{
			"authentik-postgresql-0",
			"authentik-postgresql",
			`PGPASWORD=$POSTGRES_PASSWORD pg_restore -U $POSTGRES_USER /tmp/restore/tmp/backup/authentik-postgresql`,
			"test",
		},
	)
}

func TestGetOIDCDiscoveryURL(t *testing.T) {
	t.Run("insecure-endpoint", func(t *testing.T) {
		assert.Equal(
			t,
			"http://authentik.local.gd:8081/application/o/test/",
			GetOIDCDiscoveryURL("local.gd", "test", 8081, true),
		)
	})

	t.Run("secure-endpoint", func(t *testing.T) {
		assert.Equal(
			t,
			"https://authentik.test.trustacks.io:443/application/o/test/",
			GetOIDCDiscoveryURL("test.trustacks.io", "test", 443, false),
		)
	})
}

func TestCreateAPIToken(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	namespace := "test"
	if err := createAPIToken(namespace, "test-token", clientset); err != nil {
		t.Fatal(err)
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), "authentik-bootstrap", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "test-token", strings.TrimSpace(string(secret.Data["api-token"])), "got an unexpected token value")

	// check idempotence.
	if err := createAPIToken(namespace, "test-token", clientset); err != nil {
		t.Fatal(err)
	}
}

func TestGetAPIToken(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "authentik-bootstrap",
		},
		Data: map[string][]byte{
			"api-token": []byte("test-token"),
		},
	}
	namespace := "test"
	_, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	token, err := getAPIToken(namespace, clientset)
	if err != nil {
		t.Fatal(err)
	}
	if token != "test-token" {
		t.Fatal("got an unexpected token value")
	}
}

func TestGetPropertyMappings(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"results": [
			{"managed": "goauthentik.io/providers/oauth2/scope-email", "pk": "pk1"},
			{"managed": "goauthentik.io/providers/oauth2/scope-openid", "pk": "pk2"},
			{"managed": "goauthentik.io/providers/oauth2/scope-profile", "pk": "pk3"}
		]}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	pm, err := getPropertyMappings(ts.URL, "test-token")
	if err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, pm, []string{"pk1", "pk2", "pk3"}, "got unexpected property mapping identifiers")
}

func TestGetAuthoroizationFlow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"results": [{"pk": "123", "slug": "default-provider-authorization-explicit-consent"}]}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	pk, err := getAuthorizationFlow(ts.URL, "test-token")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", pk, "got an unexpected flow pk")
}

func TestGetSigningKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"results": [{"pk": "123", "name": "authentik Self-signed Certificate"}]}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	pk, err := getCertificateKeypair(ts.URL, "test-token")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", pk, "got an unexpected certificate keypair pk")
}

func TestCreateOIDCProvider(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{"pk": 123}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	mappings := []string{
		"225abbd5-1a2b-44a8-b21d-df7f3c9be735",
		"faee6fea-8b07-40f7-bda8-0c02159cd608",
		"32c4a700-5352-44d2-880a-c5dfb08328f9",
	}
	flow := "c53f70da-aa78-42c1-950a-f0c7e7e324a1"
	signingKey := "62b33e8b-033b-4dc7-9580-0de0a3f457e6"
	pk, id, secret, err := createOIDCProvider("test", ts.URL, "test-token", flow, signingKey, mappings)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, id, 40, "expected a 40 character id")
	assert.Len(t, secret, 128, "expected a 128 character secret")
	assert.Equal(t, 123, pk, "got an unexpected provider pk")
}

func TestCreateApplication(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(`{}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()
	if err := createApplication(123, "test", ts.URL, "test-token"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateGroups(t *testing.T) {
	getGroups := make([]string, 0)
	postGroups := make([]string, 0)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			name := r.URL.Query().Get("name")
			// add group to get groups.
			getGroups = append(getGroups, name)
			if name == "admins" {
				// return a result to simulate an existing 'admins' group.
				if _, err := w.Write([]byte(`{"results": [{}]}`)); err != nil {
					t.Fatal(err)
				}
			} else {
				// return no result to invoke group creation.
				if _, err := w.Write([]byte(`{"results": []}`)); err != nil {
					t.Fatal(err)
				}
			}
		case "POST":
			data, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			g := &group{}
			if err := json.Unmarshal(data, &g); err != nil {
				t.Fatal(err)
			}
			// add group to post groups.
			postGroups = append(postGroups, g.Name)
			if _, err := w.Write([]byte(`{}`)); err != nil {
				t.Fatal(err)
			}
		}
	}))
	defer ts.Close()
	if err := createGroups(ts.URL, "test-token"); err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, getGroups, []string{"admins", "editors", "viewers"})
	assert.ElementsMatch(t, postGroups, []string{"editors", "viewers"})
}

func TestHealthCheckService(t *testing.T) {
	// test the health check with a malforned URL.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := healthCheckService("http://test.trustacks.local", 1, ctx); err == nil {
		t.Fatal("expected a timeout error")
	}
	// test the health check with a valid URL.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if err := healthCheckService(ts.URL, 1, context.TODO()); err != nil {
		t.Fatal(err)
	}
}
