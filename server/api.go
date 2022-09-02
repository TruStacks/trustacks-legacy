package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitwurx/jrpc2"
	"github.com/go-git/go-git/v5"
	"github.com/trustacks/trustacks/pkg/toolchain"
	"gopkg.in/yaml.v2"
)

const (
	// defaultCatalogSource is the url of the default toolchain catalog.
	defaultCatalogSource = "https://github.com/trustacks/toolchain"
	// defaultDatabaseProvider is the default database.
	defaultDatabaseProvider = "firebase"
)

var (
	// muxHandlers are json-rpc 2.0 mux handlers mapped to url paths.
	muxHandlers map[string]*jrpc2.MuxHandler
	// defaultParameters is the path to the default toolchain parameters.
	defaultParameters = os.Getenv("DEFAULT_PARAMETERS")
	// cloud firestore client.
	databaseProviders = map[string]func() (databaseProvider, error){
		"firebase": newFirebaseProvider,
	}
	database databaseProvider
)

type databaseProvider interface {
	checkToolchainExists(string) (bool, error)
	createToolchain(string, map[string]interface{}) error
}

// generateToolchainConfig joins the default parameters with the
// provided parameters and returns the stored config filesystem path.
func (api *apiV1) generateToolchainConfig(name, source, defaultParameters string, parameters map[string]interface{}) (string, map[string]interface{}, error) {
	configFile, err := os.CreateTemp("", "install-toolchain-config")
	if err != nil {
		return "", nil, fmt.Errorf("error creating the temp file: %s", err)
	}
	defer configFile.Close()
	defaults, err := os.ReadFile(defaultParameters)
	if err != nil {
		return "", nil, fmt.Errorf("error reading the default parameters: %s", err)
	}
	if err := yaml.Unmarshal(defaults, &parameters); err != nil {
		return "", nil, fmt.Errorf("error unmarshalling the default parameters: %s", err)
	}
	// add the toolchain name as the domain prefix.
	parameters["domain"] = fmt.Sprintf("%s.%s", name, parameters["domain"])
	config := map[string]interface{}{
		"name":       name,
		"source":     source,
		"parameters": parameters,
	}
	data, err := json.Marshal(config)
	if err != nil {
		return "", nil, fmt.Errorf("error marshalling the toolchain config")
	}
	if _, err := configFile.Write(data); err != nil {
		return "", nil, fmt.Errorf("error writing the config to the temp file: %s", err)
	}
	return configFile.Name(), config, nil
}

// apiV1 is the version 1 api.
type apiV1 struct{}

// muxHandler returns the api version 1 mux handler for the json-rpc
// 2.0 server.
func (api *apiV1) muxHandler() *jrpc2.MuxHandler {
	handler := jrpc2.NewMuxHandler()
	handler.Register("install-toolchain", jrpc2.Method{Method: api.installToolchain})
	return handler
}

// apiV1InstallToolchainParams contains the toolchain name and
// installation parameters.
type apiV1InstallToolchainParams struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// FromPositional parses the parameter positional arguments.
func (p *apiV1InstallToolchainParams) FromPositional(params []interface{}) error {
	return nil
}

// installToolchain generates the toolchain configuration file and
// runs the toolchain installation.
func (api *apiV1) installToolchain(params json.RawMessage) (interface{}, *jrpc2.ErrorObject) {
	p := new(apiV1InstallToolchainParams)
	if err := jrpc2.ParseParams(params, p); err != nil {
		return nil, err
	}
	exists, err := api.checkToolchainExists(p.Name)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    -32000,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	if exists {
		return nil, &jrpc2.ErrorObject{
			Code:    -32000,
			Message: jrpc2.ServerErrorMsg,
			Data:    fmt.Sprintf("error: toolchain '%s' already exists", p.Name),
		}
	}
	configPath, config, err := api.generateToolchainConfig(p.Name, defaultCatalogSource, defaultParameters, p.Parameters)
	if err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    -32000,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	defer os.Remove(configPath)
	if err := api.createToolchain(p.Name, config); err != nil {
		return nil, &jrpc2.ErrorObject{
			Code:    -32000,
			Message: jrpc2.ServerErrorMsg,
			Data:    err.Error(),
		}
	}
	go func() {
		if err := toolchain.Install(configPath, false, git.PlainClone); err != nil {
			log.Println(err)
		}
	}()
	return nil, nil
}

func (api *apiV1) checkToolchainExists(name string) (bool, error) {
	return database.checkToolchainExists(name)
}

func (api *apiV1) createToolchain(name string, data map[string]interface{}) error {
	return database.createToolchain(name, data)
}

func initDatabase() {
	dbProvider := os.Getenv("DATABASE_PROVIDER")
	if dbProvider == "" {
		dbProvider = defaultDatabaseProvider
	}
	var err error
	database, err = databaseProviders[dbProvider]()
	if err != nil {
		log.Fatalf("error initializing the database: %v\n", err)
	}
}

func init() {
	muxHandlers = make(map[string]*jrpc2.MuxHandler)
	addMuxHandler((&apiV1{}).muxHandler(), "/rpc/v1")
}
