package commitizen

import (
	"context"
	"regexp"
	"strings"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
)

// Version runs a "soft" bump to generate the next semantic version
// using conventional commits.
func Version(name string, project *utils.Project, client *dagger.Client) (*dagger.Container, error) {
	script := func(container *dagger.Container) (*dagger.Container, error) {
		ctx := context.Background()
		f := container.File("/src/.cz.json")
		_, err := f.Contents(ctx)
		if err != nil {
			if !strings.Contains(err.Error(), "no such file or directory") {
				return container, err
			}
			container = container.WithNewFile("/src/.cz.json", dagger.ContainerWithNewFileOpts{
				Contents: `{"commitizen": {"version": "0.0.0"}}`,
			})
		}
		container = container.WithExec([]string{"apk", "add", "git"})
		container = container.WithExec([]string{"pip", "install", "commitizen"})
		container = container.WithExec([]string{"cz", "bump", "--dry-run", "--yes"})
		versionOutput, err := container.Stdout(ctx)
		if err != nil {
			return container, err
		}
		regex, err := regexp.Compile(`tag to create: (\d.\d.\d)`)
		if err != nil {
			return container, err
		}
		container = container.WithNewFile(
			utils.ArtifactPath(name, "version"),
			dagger.ContainerWithNewFileOpts{Contents: regex.FindStringSubmatch(versionOutput)[1]},
		)
		return container, nil
	}
	action := &utils.Action{
		Image:  "python:alpine",
		Script: script,
	}
	return action.Exec(name, client, project, nil)
}
