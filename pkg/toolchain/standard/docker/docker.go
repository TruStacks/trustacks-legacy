package docker

import (
	_ "embed"
	"os"
	"time"

	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/chartutils"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

const (
	// componentName is the name of the component.
	componentName = "docker"
)

var (
	// chartVersion is the version of the helm chart.
	chartVersion = "20.10.22"
	//go:embed docker-*.tgz
	chartArchive []byte
)

type Docker struct {
	profile profile.Profile
}

// GetChart returns the docker helm chart archive.
func (c *Docker) GetChart() (string, error) {
	chart, err := chartutils.NewChart(componentName, chartArchive)
	if err != nil {
		return "", err
	}
	path, err := chart.Save(chartVersion)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Install runs pre installation tasks and install the component.
func (c *Docker) Install(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	if err := dispatcher.InstallChart("docker", nil, time.Minute*5, chartPath); err != nil {
		return err
	}
	return nil
}

// Upgrade upgrades the component.
func (c *Docker) Upgrade(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	return dispatcher.UpgradeChart("docker", nil, time.Minute*5, chartPath)
}

// Rollback deploys the previous state of the component.
func (c *Docker) Rollback(dispatcher client.Dispatcher, namespace string) error {
	chartPath, err := c.GetChart()
	if err != nil {
		return err
	}
	defer os.RemoveAll(chartPath)
	if err := dispatcher.RollbackRelease("docker", nil, time.Minute*5, chartPath); err != nil {
		return err
	}
	return nil
}

// Uninstall removes the component.
func (c *Docker) Uninstall(dispatcher client.Dispatcher, namespace string) error {
	if err := dispatcher.UninstallChart("docker"); err != nil {
		return err
	}
	return nil
}

// New creates a new docker instance.
func New(prof profile.Profile) *Docker {
	return &Docker{prof}
}
