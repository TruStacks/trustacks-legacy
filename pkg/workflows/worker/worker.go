package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	hostPort         = "temporal-frontend-headless:7233"
	runRetryMax      = 10
	runRetryInterval = 3
)

type definitions struct {
	workflow interface{}
	activity interface{}
}

var workflowDefinitions map[string]definitions

func RegisterDefinitions(name string, workflow interface{}, activity interface{}) {
	workflowDefinitions[name] = definitions{workflow, activity}
}

func New(application, kind string) error {
	definitions, ok := workflowDefinitions[kind]
	if !ok {
		return fmt.Errorf("workflow '%s' has no registered definitions", kind)
	}
	nsc, err := client.NewNamespaceClient(client.Options{HostPort: hostPort})
	if err != nil {
		return err
	}
	defer nsc.Close()
	retentionPeriod := time.Duration(time.Hour * 72)
	if err := nsc.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        application,
		WorkflowExecutionRetentionPeriod: &retentionPeriod,
	}); err != nil && !strings.Contains(err.Error(), "Namespace already exists") {
		return err
	}
	c, err := client.Dial(client.Options{
		Namespace: application,
		HostPort:  hostPort,
	})
	if err != nil {
		return err
	}
	w := worker.New(c, application, worker.Options{})
	w.RegisterWorkflow(definitions.workflow)
	w.RegisterActivity(definitions.activity)
	return w.Run(worker.InterruptCh())
}

func init() {
	workflowDefinitions = make(map[string]definitions)
}
