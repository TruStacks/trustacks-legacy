package toolchain

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
	// configMapName is the name of the desired state configuration
	// config map.
	configMapName = "desired-state-config"
	// configMapKey is the name of the key that contains the
	// desired state config json object inside the config map.
	configMapKey = "config"
)

type DesiredStateConfiguration struct {
	listeners map[string]listener
}

type listener interface {
	Reconcile(interface{}) error
}

// Save creates a config map to store the desired state configuration.
func (dsc *DesiredStateConfiguration) Save(toolchain string, config interface{}, clientset kubernetes.Interface) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: configMapName},
		Data:       map[string]string{configMapKey: string(data)},
	}
	if _, err := clientset.CoreV1().ConfigMaps(toolchain).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			patch := map[string]interface{}{"data": map[string]string{configMapKey: string(data)}}
			data, err = json.Marshal(patch)
			if err != nil {
				return err
			}
			clientset.CoreV1().ConfigMaps(toolchain).Patch(context.TODO(), configMapName, types.MergePatchType, data, metav1.PatchOptions{})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// Save loads the desired state configuration from the config map.
func (dsc *DesiredStateConfiguration) Load(toolchain string, config interface{}, clientset kubernetes.Interface) error {
	cm, err := clientset.CoreV1().ConfigMaps(toolchain).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cm.Data["config"]), config)
}

// Reconcile calls the reconcile method of the listener that is
// watching for changes to the desired state configuration.
func (dsc *DesiredStateConfiguration) Reconcile(toolchain string, config interface{}) error {
	listener, ok := dsc.listeners[toolchain]
	if !ok {
		return nil
	}
	return listener.Reconcile(config)
}

// AddListener adds the listener to the desired state configuration.
func (dsc *DesiredStateConfiguration) AddListener(toolchain string, l listener) {
	dsc.listeners[toolchain] = l
}

// NewDesiredStateConfiguration create a new instance of the desired
// state configuration.
func NewDesiredStateConfiguration() *DesiredStateConfiguration {
	return &DesiredStateConfiguration{listeners: map[string]listener{}}
}
