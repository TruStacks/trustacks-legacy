package toolchain

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSave(t *testing.T) {
	dsc := NewDesiredStateConfiguration()
	clientset := fake.NewSimpleClientset()
	if err := dsc.Save("test", map[string]interface{}{"key": "value1"}, clientset); err != nil {
		t.Fatal(err)
	}
	if err := dsc.Save("test", map[string]interface{}{"key": "value2"}, clientset); err != nil {
		t.Fatal(err)
	}
	if err := dsc.Save("test", map[string]interface{}{"key": "value3"}, clientset); err != nil {
		t.Fatal(err)
	}
	cm, err := clientset.CoreV1().ConfigMaps("test").Get(context.TODO(), "desired-state-config", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(cm.Data["config"]), &config); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "value3", config["key"])
}

func TestLoad(t *testing.T) {
	dsc := NewDesiredStateConfiguration()
	clientset := fake.NewSimpleClientset()
	if err := dsc.Save("test", map[string]interface{}{"key": "value1"}, clientset); err != nil {
		t.Fatal(err)
	}
	config := map[string]interface{}{}
	if err := dsc.Load("test", &config, clientset); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "value1", config["key"])
}

type fakeConfig struct {
	Fake   string `json:"fake"`
	Faker  string `json:"faker"`
	Fakest string `json:"fakest"`
}

type fakeListener struct {
	config fakeConfig
	called bool
}

func (l *fakeListener) Reconcile(config interface{}) error {
	l.config = config.(fakeConfig)
	l.called = true
	return nil
}

func TestAddListener(t *testing.T) {
	l := &fakeListener{}
	dsc := NewDesiredStateConfiguration()
	dsc.AddListener("test", l)
	assert.Equal(t, l, dsc.listeners["test"])
}

func TestReconcile(t *testing.T) {
	l := &fakeListener{}
	config := fakeConfig{
		Fake:   "value1",
		Faker:  "value2",
		Fakest: "value3",
	}
	dsc := NewDesiredStateConfiguration()
	// test that no error is raised if the listener is not available.
	if err := dsc.Reconcile("test", config); err != nil {
		t.Fatal(err)
	}
	dsc.AddListener("test", l)
	if err := dsc.Reconcile("test", config); err != nil {
		t.Fatal(err)
	}
	assert.True(t, l.called)
	assert.Equal(t, "value1", l.config.Fake)
	assert.Equal(t, "value2", l.config.Faker)
	assert.Equal(t, "value3", l.config.Fakest)
}
