package state

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDesiredStateSave(t *testing.T) {
	ds := &DesiredState{}
	clientset := fake.NewSimpleClientset()
	if err := ds.save("test", map[string]interface{}{"key": "value1"}, clientset); err != nil {
		t.Fatal(err)
	}
	if err := ds.save("test", map[string]interface{}{"key": "value2"}, clientset); err != nil {
		t.Fatal(err)
	}
	if err := ds.save("test", map[string]interface{}{"key": "value3"}, clientset); err != nil {
		t.Fatal(err)
	}
	cm, err := clientset.CoreV1().ConfigMaps("test").Get(context.TODO(), desiredStateConfigMapName, metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(cm.Data["config"]), &config); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "value3", config["key"])
}

func TestDesiredStateLoad(t *testing.T) {
	ds := &DesiredState{}
	clientset := fake.NewSimpleClientset()
	if err := ds.save("test", map[string]interface{}{"key": "value1"}, clientset); err != nil {
		t.Fatal(err)
	}
	config := map[string]interface{}{}
	if err := ds.load("test", &config, clientset); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "value1", config["key"])
}

func TestActiveStateSetAndGet(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	as := &ActiveState{}
	if err := as.set("test", "test-key", "value1", clientset); err != nil {
		t.Fatal(err)
	}
	v, err := as.get("test", "test-key", clientset)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "value1", v)
}

type fakeConfig struct {
	Fake   string `json:"fake"`
	Faker  string `json:"faker"`
	Fakest string `json:"fakest"`
}

type fakeListener struct {
	config fakeConfig
}

func (l *fakeListener) Reconcile(name, event string, config interface{}, sm *StateManager) error {
	l.config = config.(fakeConfig)
	switch event {
	case "test-event":
		return sm.Set("called", "true")
	}
	return errors.New("unable to handle event")
}
