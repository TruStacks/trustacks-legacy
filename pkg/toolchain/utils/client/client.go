package client

import (
	"context"
	"os"
	"strings"
	"time"

	helmclient "github.com/mittwald/go-helm-client"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

// inClusterNamespace is the path to the in-cluster namespace.
const inClusterNamespace = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

type Dispatcher interface {
	Clientset() kubernetes.Interface
	CreateNamespace() error
	DeleteNamespace() error
	InstallChart(string, interface{}, time.Duration, string) error
	UpgradeChart(string, interface{}, time.Duration, string) error
	RollbackRelease(string, interface{}, time.Duration, string) error
	UninstallChart(string) error
	ExecCommand(string, string, string, string) error
}

type ClientDispatcher struct {
	namespace  string
	clientset  kubernetes.Interface
	helmClient helmclient.Client
	restconfig *rest.Config
}

// CreateNamespace creates a kubernetes namespace.
func (d *ClientDispatcher) CreateNamespace() error {
	if _, err := d.clientset.CoreV1().Namespaces().Create(
		context.TODO(),
		&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: d.namespace}},
		metav1.CreateOptions{},
	); err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}
	return nil
}

// DeleteNamespace delets a kubernetes namespace.
func (d *ClientDispatcher) DeleteNamespace() error {
	return d.clientset.CoreV1().Namespaces().Delete(context.TODO(), d.namespace, metav1.DeleteOptions{})
}

type InstallChartArgs struct {
	Name      string
	Values    interface{}
	Timeout   time.Duration
	ChartPath string
}

// chartSpec creates the helm spec to use for helm chart operations.
func (d *ClientDispatcher) chartSpec(name string, values interface{}, timeout time.Duration, chartPath string) (*helmclient.ChartSpec, error) {
	spec := &helmclient.ChartSpec{
		ReleaseName: name,
		Namespace:   d.namespace,
		ChartName:   chartPath,
		Timeout:     timeout,
		Atomic:      true,
		Wait:        true,
	}
	if values != nil {
		valuesYaml, err := yaml.Marshal(values)
		if err != nil {
			return nil, err
		}
		spec.ValuesYaml = string(valuesYaml)
	}
	return spec, nil
}

// InstallChart installs a helm chart.
func (d *ClientDispatcher) InstallChart(name string, values interface{}, timeout time.Duration, chartPath string) error {
	chartSpec, err := d.chartSpec(name, values, timeout, chartPath)
	if err != nil {
		return err
	}
	if _, err := d.helmClient.InstallChart(context.Background(), chartSpec, &helmclient.GenericHelmOptions{}); err != nil {
		return err
	}
	return nil
}

// UpgradeChart updates a helm chart.
func (d *ClientDispatcher) UpgradeChart(name string, values interface{}, timeout time.Duration, chartPath string) error {
	chartSpec, err := d.chartSpec(name, values, timeout, chartPath)
	if err != nil {
		return err
	}
	if _, err := d.helmClient.UpgradeChart(context.Background(), chartSpec, &helmclient.GenericHelmOptions{}); err != nil {
		return err
	}
	return nil
}

// RollbackRelease rolls back a helm release.
func (d *ClientDispatcher) RollbackRelease(name string, values interface{}, timeout time.Duration, chartPath string) error {
	chartSpec, err := d.chartSpec(name, values, timeout, chartPath)
	if err != nil {
		return err
	}
	if err := d.helmClient.RollbackRelease(chartSpec); err != nil {
		return err
	}
	return nil
}

// UninstallChart uninstall a helm chart.
func (d *ClientDispatcher) UninstallChart(name string) error {
	chartSpec := helmclient.ChartSpec{
		ReleaseName: name,
		Namespace:   d.namespace,
	}
	if err := d.helmClient.UninstallRelease(&chartSpec); err != nil {
		return err
	}
	return nil
}

// Clientset returns the client kubernetes clientset.
func (d *ClientDispatcher) Clientset() kubernetes.Interface {
	return d.clientset
}

// ExecCommand executes a command in kubernetes pod.
func (d *ClientDispatcher) ExecCommand(pod, container, command, namespace string) error {
	req := d.clientset.
		CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command:   []string{"/bin/sh", "-c", command},
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(d.restconfig, "POST", req.URL())
	if err != nil {
		return err
	}
	if err := exec.Stream(remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	}); err != nil {
		return err
	}
	return nil
}

// NewDispatcher creates a client dispatcher for kubernetes and helm
// operations.
func NewDispatcher(namespace string) (Dispatcher, error) {
	dispatcher := &ClientDispatcher{namespace: namespace}
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	dispatcher.restconfig = restConfig
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	dispatcher.clientset = clientset
	opt := &helmclient.RestConfClientOptions{
		Options:    &helmclient.Options{Namespace: namespace},
		RestConfig: restConfig,
	}
	helmClient, err := helmclient.NewClientFromRestConf(opt)
	if err != nil {
		panic(err)
	}
	dispatcher.helmClient = helmClient
	return dispatcher, nil
}

// GetNamespace returns the namespace of an in cluster kubernetes
// client.
func GetNamespace() (string, error) {
	data, err := os.ReadFile(inClusterNamespace)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), err
}
