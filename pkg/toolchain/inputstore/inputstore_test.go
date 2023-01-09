package inputstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGenerateKeys(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if err := GenerateIdentity("test", clientset); err != nil {
		t.Fatal(err)
	}
	secret, err := clientset.CoreV1().Secrets("test").Get(context.Background(), identityConfigMap, v1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, secret.Data["identity"], 74)
	assert.Len(t, secret.Data["recipient"], 62)
	assert.Contains(t, string(secret.Data["identity"]), "AGE-SECRET-KEY")
	assert.Contains(t, string(secret.Data["recipient"]), "age")
}

type TestInputs struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

func TestWrite(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if err := Write("test", "test", []byte(`{"value1": "x", "value2": "y", "value3": "z"}`), clientset); err != nil {
		t.Fatal(err)
	}
	cm, err := clientset.CoreV1().ConfigMaps("test").Get(context.Background(), "test-input-store", v1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	cmData := map[string]string{}
	if err := json.Unmarshal([]byte(cm.Data["inputs"]), &cmData); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "x", cmData["value1"])
	assert.Equal(t, "y", cmData["value2"])
	assert.Equal(t, "z", cmData["value3"])
}

func TestRead(t *testing.T) {
}

func TestEncryption(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	data, err := Encrypt(identity.Recipient().String(), &TestInputs{"x", "y", "z"})
	if err != nil {
		t.Fatal(err)
	}
	tf, err := os.CreateTemp("", "input-store-enc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tf.Name())
	if _, err := tf.Write(data); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	cmd := exec.Command("sops", "-d", tf.Name())
	cmd.Env = append(cmd.Env, fmt.Sprintf("SOPS_AGE_KEY=%s", identity))
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	inputs := map[string]string{}
	if err := json.Unmarshal(buf.Bytes(), &inputs); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, inputs["value1"], "x")
	assert.Equal(t, inputs["value2"], "y")
	assert.Equal(t, inputs["value3"], "z")
}

func TestDecrypt(t *testing.T) {
}
