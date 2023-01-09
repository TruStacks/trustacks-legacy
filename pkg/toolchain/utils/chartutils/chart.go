package chartutils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type HookRBACRules struct {
	ApiGroups []string
	Resources []string
	Verbs     []string
}

type chartYAML struct {
	ApiVersion   string                   `yaml:"apiVersion"`
	Name         string                   `yaml:"name"`
	Version      string                   `yaml:"version"`
	Dependencies []map[string]interface{} `yaml:"dependencies"`
}

type chart struct {
	name     string
	subchart []byte
}

// Save writes the helm chart to the temporary dir.
func (c *chart) Save(version string) (string, error) {
	chartDir, err := os.MkdirTemp("", fmt.Sprintf("%s-chart-", c.name))
	if err != nil {
		return "", err
	}
	yamlData, err := yaml.Marshal(chartYAML{
		ApiVersion:   "v1",
		Name:         c.name,
		Version:      version,
		Dependencies: []map[string]interface{}{{"name": c.name}},
	})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/Chart.yaml", chartDir), yamlData, 0644); err != nil {
		return "", err
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/charts", chartDir), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(fmt.Sprintf("%s/charts/%s.tgz", chartDir, c.name), c.subchart, 0644); err != nil {
		return "", err
	}
	return chartDir, nil
}

// UniqueID creates a short id from a namespace string.
func UniqueID(namespace string) (string, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(namespace)); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:7], nil
}

// NewChart creates a new chart instance.
func NewChart(name string, subchart []byte) (*chart, error) {
	c := &chart{
		name:     name,
		subchart: subchart,
	}
	return c, nil
}
