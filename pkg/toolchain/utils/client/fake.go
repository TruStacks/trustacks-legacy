package client

import (
	"time"

	"k8s.io/client-go/kubernetes"
)

// FakeDispatcher is a fake client disptacher for use in tests.
type FakeDispatcher struct {
	MockCalls map[string][]interface{}
}

func (d *FakeDispatcher) Clientset() kubernetes.Interface {
	return nil
}

func (d *FakeDispatcher) CreateNamespace() error {
	return nil
}

func (d *FakeDispatcher) DeleteNamespace() error {
	return nil
}

func (d *FakeDispatcher) InstallChart(string, interface{}, time.Duration, string) error {
	return nil
}

func (d *FakeDispatcher) UpgradeChart(string, interface{}, time.Duration, string) error {
	return nil
}

func (d *FakeDispatcher) RollbackRelease(string, interface{}, time.Duration, string) error {
	return nil
}

func (d *FakeDispatcher) UninstallChart(string) error {
	return nil
}

func (d *FakeDispatcher) ExecCommand(pod, container, command, namespace string) error {
	d.MockCalls["ExecCommand"] = append(d.MockCalls["ExecCommand"], []string{pod, container, command, namespace})
	return nil
}

func NewFakeDispatcher() *FakeDispatcher {
	return &FakeDispatcher{MockCalls: make(map[string][]interface{})}
}
