package workflows

import (
	"context"
	"os"
	"time"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/extensions/reactscripts"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
	"go.temporal.io/sdk/workflow"
)

type ReactWorkflowParams struct {
	Name string
	Repo string
}

func ReactWorkflowDefinition(ctx workflow.Context, params ReactWorkflowParams) error {
	var activities *ReactWorkflowActivities
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: 3600 * time.Second})
	future := workflow.ExecuteActivity(ctx, activities.TestDefinition, ReactWorkflowActivitiesParams{params.Name, params.Repo})
	return future.Get(ctx, nil)
}

type ReactWorkflowActivitiesParams struct {
	Name string
	Repo string
}

type ReactWorkflowActivities struct{}

func (a *ReactWorkflowActivities) TestDefinition(ctx context.Context, params ReactWorkflowActivitiesParams) error {
	path, err := os.MkdirTemp("", params.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)
	commitHash, err := utils.CloneSource(params.Repo, path)
	if err != nil {
		return err
	}
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(utils.NewLokiLogger(params.Name, commitHash, "loki:3100")))
	if err != nil {
		return err
	}
	defer client.Close()
	project := utils.NewProject(params.Name, path, client)
	if _, err := reactscripts.Test(project, client); err != nil {
		return err
	}
	return nil
}
