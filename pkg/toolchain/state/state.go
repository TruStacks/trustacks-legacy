package state

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	// desiredStateConfigMapName is the name of the desired state
	// config map.
	desiredStateConfigMapName = "ts-desired-state"
	// activeStateConfigMapName is the name of the active state
	// config map.
	activeStateConfigMapName = "ts-active-state"
	// configMapKey is the name of the key that contains the
	// state config json object inside the config map.
	configMapKey = "config"
)

type DesiredState struct{}

// save creates a config map to store the desired state configuration.
func (ds *DesiredState) save(toolchain string, config interface{}, clientset kubernetes.Interface) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: desiredStateConfigMapName},
		Data:       map[string]string{configMapKey: string(data)},
	}
	if _, err := clientset.CoreV1().ConfigMaps(toolchain).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			patch := map[string]interface{}{"data": map[string]string{configMapKey: string(data)}}
			data, err = json.Marshal(patch)
			if err != nil {
				return err
			}
			clientset.CoreV1().ConfigMaps(toolchain).Patch(context.TODO(), desiredStateConfigMapName, types.MergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// load loads the desired state configuration from the config map.
func (ds *DesiredState) load(toolchain string, config interface{}, clientset kubernetes.Interface) error {
	cm, err := clientset.CoreV1().ConfigMaps(toolchain).Get(context.TODO(), desiredStateConfigMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cm.Data["config"]), config)
}

type ActiveState struct{}

// set updates the active state key with the provided value.
func (as *ActiveState) set(toolchain, key, value string, clientset kubernetes.Interface) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: activeStateConfigMapName},
		Data:       map[string]string{key: value},
	}
	if _, err := clientset.CoreV1().ConfigMaps(toolchain).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			patch := map[string]interface{}{"data": map[string]string{key: value}}
			data, err := json.Marshal(patch)
			if err != nil {
				return err
			}
			clientset.CoreV1().ConfigMaps(toolchain).Patch(context.TODO(), activeStateConfigMapName, types.MergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// get fetches the value of the key from the active state.
func (as *ActiveState) get(toolchain, key string, clientset kubernetes.Interface) (string, error) {
	cm, err := clientset.CoreV1().ConfigMaps(toolchain).Get(context.TODO(), activeStateConfigMapName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return cm.Data[key], nil
}

type StateManager struct {
	toolchain string
	clientset kubernetes.Interface
	ds        *DesiredState
	as        *ActiveState
}

// Save stores the desired state to a config map.
func (sm *StateManager) Save(config interface{}) error {
	return sm.ds.save(sm.toolchain, config, sm.clientset)
}

// Load fetches the desired state from the stored config map.
func (sm *StateManager) Load(config interface{}) error {
	return sm.ds.load(sm.toolchain, config, sm.clientset)
}

// Set updates the active state key with the provided value.
func (sm *StateManager) Set(key, value string) error {
	return sm.as.set(sm.toolchain, key, value, sm.clientset)
}

// Get fetches the value from the active state.
func (sm *StateManager) Get(key string) (string, error) {
	return sm.as.get(sm.toolchain, key, sm.clientset)
}

// NewStateManager creates an instance of the state manager.
func NewStateManager(toolchain string, clientset kubernetes.Interface) *StateManager {
	return &StateManager{toolchain, clientset, &DesiredState{}, &ActiveState{}}
}
