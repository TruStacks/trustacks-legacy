package utils

import (
	"context"
	"strings"

	"dagger.io/dagger"
)

type Action struct {
	Image string
	Run   []string
}

type RunOpts struct {
	Caches []string
}

func (action *Action) Exec(client *dagger.Client, project *Project, opts *RunOpts) (*dagger.Container, error) {
	container := client.Container().From(action.Image)
	for _, path := range opts.Caches {
		container = container.WithMountedCache(project.CreateCache(path, client))
	}
	container = container.WithMountedDirectory("/src", project.Source).WithWorkdir("/src")
	container = container.Exec(dagger.ContainerExecOpts{Args: []string{"mkdir", "-p", "/artifacts"}})
	container = container.Exec(dagger.ContainerExecOpts{Args: []string{"/bin/sh", "-c", strings.Join(action.Run, "; ")}})

	_, err := container.Directory("/artifacts").Export(context.Background(), "artifacts/")
	if err != nil {
		return nil, err
	}

	return container, nil
}
