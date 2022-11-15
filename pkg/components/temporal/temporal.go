package components

import toolchain "github.com/trustacks/trustacks/pkg/provision"

const temporalComponentChart = "https://storage.cloud.google.com/charts.trustacks.io/temporal-0.18.4.tgz"

type TemporalComponent struct {
	values map[string]interface{}
}

func (t *TemporalComponent) enableIngress(config toolchain.Config) {

}

func (t *TemporalComponent) Values() string {
	return ""
}

func NewTemporalComponent() *TemporalComponent {
	component := &TemporalComponent{}
	return component
}
