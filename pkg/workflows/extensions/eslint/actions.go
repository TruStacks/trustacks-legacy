package eslint

import (
	"context"
	"strings"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
)

type RunArgs struct {
	ESLintRC string
}

// Run executes the eslint linter.
func Run(name string, project *utils.Project, client *dagger.Client, args RunArgs) (*dagger.Container, error) {
	script := func(container *dagger.Container) (*dagger.Container, error) {
		ctx := context.Background()
		entries, err := project.Source.Entries(ctx)
		if err != nil {
			return nil, err
		}
		var eslintExists bool
		for _, f := range entries {
			if strings.Contains(f, ".eslintrc") {
				eslintExists = true
			}
		}
		if !eslintExists {
			container = container.WithNewFile("/src/.eslint", dagger.ContainerWithNewFileOpts{Contents: args.ESLintRC})
		}
		container = container.WithExec([]string{"npm", "install", "eslint-html-reporter"})
		container, err = utils.WithTrapExec(
			container,
			[]string{"npx", "-y", "eslint", "./", "-f", "node_modules/eslint-html-reporter/reporter.js", "-o", utils.ArtifactPath(name, "report.html")},
		)
		if err != nil {
			return container, err
		}
		return container, nil
	}
	action := &utils.Action{
		Image:  "node:alpine",
		Script: script,
	}
	return action.Exec(name, client, project, &utils.RunOpts{
		Caches: []string{"/src/node_modules"},
	})
}
