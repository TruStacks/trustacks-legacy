package react

import (
	"context"
	"os"
	"time"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/extensions/commitizen"
	"github.com/trustacks/trustacks/pkg/workflows/extensions/eslint"
	"github.com/trustacks/trustacks/pkg/workflows/extensions/reactscripts"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
	"github.com/trustacks/trustacks/pkg/workflows/worker"
	"go.temporal.io/sdk/workflow"
)

type ReactWorkflowParams struct {
	Name  string `json:"name"`
	Repo  string `json:"repo"`
	RunID string `json:"id"`
}

// ReactWorkflow .
func ReactWorkflow(ctx workflow.Context, params ReactWorkflowParams) error {
	var activities *ReactWorkflowActivities
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{StartToCloseTimeout: 3600 * time.Second})
	future := workflow.ExecuteActivity(ctx, activities.Definition, params)
	return future.Get(ctx, nil)
}

type ReactWorkflowActivities struct{}

// Definition .
func (a *ReactWorkflowActivities) Definition(ctx context.Context, params ReactWorkflowParams) error {
	path, err := os.MkdirTemp("", params.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)
	commitHash, err := utils.CloneSource(params.Repo, path)
	if err != nil {
		return err
	}
	logger := utils.NewLokiLogger(params.Name, commitHash, params.RunID, "loki:3100")
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(logger))
	if err != nil {
		return err
	}
	defer client.Close()
	project := utils.NewProject(params.Name, path, client)
	if _, err := commitizen.Version("version", project, client); err != nil {
		return err
	}
	if _, err := reactscripts.Test("test", project, client); err != nil {
		return err
	}
	if _, err := eslint.Run("lint", project, client, eslint.RunArgs{ESLintRC: `{"extends": ["react-app"]}`}); err != nil {
		return err
	}
	return nil
}

func init() {
	worker.RegisterDefinitions("react", ReactWorkflow, (&ReactWorkflowActivities{}).Definition)
}
