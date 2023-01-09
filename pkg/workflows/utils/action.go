package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

const artifactsRoot = "/tmp/artifacts"

type Action struct {
	Image     string
	Script    RunCmd
	Artifacts map[string]string
}

type RunCmd func(container *dagger.Container) (*dagger.Container, error)

type RunOpts struct {
	Caches []string
}

// Exec configure caches, artifacts, and sources and runs the task.
func (a *Action) Exec(name string, client *dagger.Client, project *Project, opts *RunOpts) (*dagger.Container, error) {
	artifactsPath := fmt.Sprintf("%s/%s", artifactsRoot, name)
	container := client.Container().From(a.Image)
	if opts != nil {
		for _, path := range opts.Caches {
			container = container.WithMountedCache(project.CreateCache(path, client))
		}
	}
	container = container.WithMountedDirectory("/src", project.Source).WithWorkdir("/src")
	container = container.WithExec([]string{"mkdir", "-p", artifactsPath})
	container, err := a.Script(container)
	if err != nil {
		fmt.Println(err)
	}
	_, exportErr := container.Directory(artifactsRoot).Export(context.Background(), "artifacts/")
	fmt.Println(exportErr)
	if exportErr != nil {
		return nil, exportErr
	}
	return container, err
}

// ArtifactPath returns the filesystem path for the provided
// artifact.
func ArtifactPath(name, artifact string) string {
	return fmt.Sprintf("/tmp/artifacts/%s/%s", name, artifact)
}

// WithTrapExec runs the command with error trapping to prevent
// failure from bricking the dagger container.
func WithTrapExec(container *dagger.Container, cmd []string) (*dagger.Container, error) {
	container = container.WithExec(
		[]string{"/bin/sh", "-c", fmt.Sprintf("%s > /tmp/stderr 2>&1 || echo $? > /tmp/code", strings.Join(cmd, " "))},
	)
	ctx := context.Background()
	code, err := container.File("/tmp/code").Contents(ctx)
	if err != nil {
		// if the command succeeds /tmp/code will not be created
		if strings.Contains(err.Error(), "no such file or directory") {
			return container, nil
		}
		return container, err
	}
	if code != "0" {
		stderr, err := container.File("/tmp/stderr").Contents(ctx)
		if err != nil {
			return container, err
		}
		return container, errors.New(stderr)
	}
	return container, nil
}
