package inputstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"filippo.io/age"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	identityConfigMap = "inputStoreIdentity"
)

// GenerateIdentity creates an age identity and recipient key for
// input store encryption and decryption.
func GenerateIdentity(namespace string, clientset kubernetes.Interface) error {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return err
	}
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      identityConfigMap,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"identity":  []byte(identity.String()),
			"recipient": []byte(identity.Recipient().String()),
		},
	}
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// Write stores the data into the input store config map.
func Write(namespace, application string, data []byte, clientset kubernetes.Interface) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-input-store", application),
			Namespace: namespace,
		},
		Data: map[string]string{
			"inputs": string(data),
		},
	}
	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Read(namespace, application string, inputs interface{}, clienset kubernetes.Interface) error {
	return nil
}

// Encrypt sops encrypts the input store with the provided age
// public recipient key.
func Encrypt(recipient string, inputStore interface{}) ([]byte, error) {
	data, err := json.Marshal(inputStore)
	if err != nil {
		return nil, err
	}
	tf, err := os.CreateTemp("", "input-store")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tf.Name())
	if _, err := tf.Write(data); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	cmd := exec.Command("sops", "--encrypt", "--age", recipient, tf.Name())
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func Decrypt(identity string) error {
	return nil
}
