package standard

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/teris-io/shortid"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/authentik"
	docker "github.com/trustacks/trustacks/pkg/toolchain/standard/docker"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/loki"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/profile"
	"github.com/trustacks/trustacks/pkg/toolchain/standard/temporal"
	"github.com/trustacks/trustacks/pkg/toolchain/state"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/backend"
	"github.com/trustacks/trustacks/pkg/toolchain/utils/client"
)

const (
	statusInstalled = "installed"
	statusUpgraded  = "upgraded"
	statusReverted  = "reverted"
)

type component interface {
	Install(client.Dispatcher, string) error
	Uninstall(client.Dispatcher, string) error
	Upgrade(client.Dispatcher, string) error
	Rollback(client.Dispatcher, string) error
}

// Install installs the toolchain components.
func Install(namespace string, config *DesiredConfig, stateManager *state.StateManager, dispatcher client.Dispatcher) {
	cid, err := shortid.Generate()
	if err != nil {
		log.Println(err)
		return
	}
	defer revert(cid, namespace, config, stateManager, dispatcher)
	names, components := getComponents(config)
	for i, name := range names {
		status, err := stateManager.Get(fmt.Sprintf("%s.status", name))
		if err != nil {
			panic(err)
		}
		if status == "" {
			if err := components[i].Install(dispatcher, namespace); err != nil {
				panic(err)
			}
			if err := stateManager.Set(fmt.Sprintf("%s.status", name), statusInstalled); err != nil {
				panic(err)
			}
			if err := stateManager.Set(fmt.Sprintf("%s.cid", name), cid); err != nil {
				panic(err)
			}
		}
	}
}

// Upgrade upgrades the toolchain components.
func Upgrade(namespace string, config *DesiredConfig, stateManager *state.StateManager, dispatcher client.Dispatcher) {
	storageConfig := backend.StorageConfig{
		URL:             os.Getenv("STORAGE_URL"),
		AccessKeyID:     os.Getenv("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		log.Println(err)
		return
	}
	cid, err := shortid.Generate()
	if err != nil {
		log.Println(err)
		return
	}
	defer revert(cid, namespace, config, stateManager, dispatcher)
	names, components := getComponents(config)
	for i, name := range names {
		status, err := stateManager.Get(fmt.Sprintf("%s.status", name))
		if err != nil {
			panic(err)
		}
		if status != "" {
			if err := components[i].Upgrade(dispatcher, namespace); err != nil {
				panic(err)
			}
			if err := stateManager.Set(fmt.Sprintf("%s.status", name), statusUpgraded); err != nil {
				panic(err)
			}
			if err := stateManager.Set(fmt.Sprintf("%s.cid", name), cid); err != nil {
				panic(err)
			}
		}
	}
}

// Uninstall uninstalls the toolchain components.
func Uninstall(namespace string, config *DesiredConfig, stateManager *state.StateManager, dispatcher client.Dispatcher) {
	storageConfig := backend.StorageConfig{
		URL:             os.Getenv("STORAGE_URL"),
		AccessKeyID:     os.Getenv("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
	}
	if err := backend.NewStorageConfig(storageConfig, namespace, dispatcher.Clientset()); err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	names, components := getComponents(config)
	for i, name := range names {
		wg.Add(1)
		go func(name string, c component, wg *sync.WaitGroup) {
			status, err := stateManager.Get(fmt.Sprintf("%s.status", name))
			if err != nil {
				panic(err)
			}
			if status != "" {
				if err := c.Uninstall(dispatcher, namespace); err != nil {
					panic(err)
				}
			}
			wg.Done()
		}(name, components[i], &wg)
	}
	wg.Wait()
	if err := dispatcher.DeleteNamespace(); err != nil {
		log.Println(err)
	}
}

// getComponents returns the list of actionable components.
func getComponents(config *DesiredConfig) ([]string, []component) {
	names := []string{
		"docker",
		"authentik",
		"loki",
		"temporal",
	}
	components := []component{
		authentik.New(config.Profile),
		docker.New(config.Profile),
		loki.New(config.Profile),
		temporal.New(config.Profile),
	}
	return names, components
}

// revert rolls back components that were modified during the current
// change operation.
func revert(cid, namespace string, config *DesiredConfig, stateManager *state.StateManager, dispatcher client.Dispatcher) {
	if r := recover(); r != nil {
		fmt.Println(r)
		if err := func() error {
			newCID, err := shortid.Generate()
			if err != nil {
				return err
			}
			names, components := getComponents(config)
			for i, name := range names {
				status, err := stateManager.Get(fmt.Sprintf("%s.status", name))
				if err != nil {
					return err
				}
				changeID, err := stateManager.Get(fmt.Sprintf("%s.cid", name))
				if err != nil {
					return err
				}
				if status == statusUpgraded && changeID == cid {
					if err := components[i].Rollback(dispatcher, namespace); err != nil {
						return err
					}
					if err := stateManager.Set(fmt.Sprintf("%s.status", name), statusReverted); err != nil {
						return err
					}
					if err := stateManager.Set(fmt.Sprintf("%s.cid", name), newCID); err != nil {
						return err
					}
				}
			}
			return nil
		}(); err != nil {
			log.Println(err)
			return
		}
	}
}

type DesiredConfig struct {
	Profile profile.Profile `json:"profile"`
}

// NewDesiredConfig creates a standard toolchain configuration
// instance.
func NewDesiredConfig(profileMap map[string]interface{}) (*DesiredConfig, error) {
	conf := &DesiredConfig{}
	data, err := json.Marshal(profileMap)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &conf.Profile); err != nil {
		return nil, err
	}
	if err := conf.Profile.Validate(); err != nil {
		return nil, err
	}
	return conf, nil
}
