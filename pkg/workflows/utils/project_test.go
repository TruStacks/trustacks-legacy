package utils

import (
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"gotest.tools/v3/assert"
)

func TestProjectCreateCache(t *testing.T) {
	if testing.Short() {
		t.Skip("[skipped]")
	}
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	project := &Project{caches: make(map[string]*dagger.CacheVolume)}

	path, cache := project.CreateCache("/tmp/test", client)
	assert.Equal(t, path, "/tmp/test")

	_, storedCached := project.CreateCache("/tmp/test", client)
	assert.Equal(t, cache, storedCached)
}

func TestNewProject(t *testing.T) {
	if testing.Short() {
		t.Skip("[skipped]")
	}
	client, err := dagger.Connect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	mockSource, err := os.MkdirTemp("", "mock-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(mockSource)
	project := NewProject("test", mockSource, client)
	_, err = project.Source.Entries(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
