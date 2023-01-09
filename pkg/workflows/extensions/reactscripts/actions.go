package reactscripts

import (
	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
)

// Test runs jest unit tests.
func Test(name string, project *utils.Project, client *dagger.Client) (*dagger.Container, error) {
	script := func(container *dagger.Container) (*dagger.Container, error) {
		container = container.WithExec([]string{"npm", "install"})
		container = container.WithExec([]string{"npm", "install", "jest-html-reporter"})
		container = container.WithEnvVariable("CI", "true")
		container, err := utils.WithTrapExec(container, []string{"npx", "react-scripts", "test", "--testResultsProcessor=./node_modules/jest-html-reporter", "--coverage"})
		container = container.WithExec([]string{"/bin/sh", "-c", "ls -al"})
		container = container.WithExec([]string{"mv", "test-report.html", utils.ArtifactPath(name, "test-report.html")})
		container = container.WithExec([]string{"mv", "coverage", utils.ArtifactPath(name, "coverage")})
		return container, err
	}
	action := &utils.Action{
		Image:  "node:alpine",
		Script: script,
	}
	return action.Exec(name, client, project, &utils.RunOpts{
		Caches: []string{"/src/node_modules"},
	})
}
