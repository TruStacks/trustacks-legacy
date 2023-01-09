package utils

import (
	"time"

	"dagger.io/dagger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const referenceBranch = "main"

type Project struct {
	Name     string
	Source   *dagger.Directory
	Revision string
	caches   map[string]*dagger.CacheVolume
}

// CreateCache .
func (p *Project) CreateCache(path string, client *dagger.Client) (string, *dagger.CacheVolume) {
	if _, exists := p.caches[path]; !exists {
		p.caches[path] = client.CacheVolume(p.Name + time.Now().Format(time.RFC3339))
	}
	return path, p.caches[path]
}

// NewProject .
func NewProject(name, path string, client *dagger.Client) *Project {
	return &Project{
		Name:   name,
		Source: client.Host().Directory(path),
		caches: make(map[string]*dagger.CacheVolume),
	}
}

// CloneSource .
func CloneSource(repo, path string) (string, error) {
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           repo,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(referenceBranch),
	})
	if err != nil {
		return "", err
	}
	h, err := r.ResolveRevision(plumbing.Revision(referenceBranch))
	if err != nil {
		return "", err
	}
	return h.String()[0:6], nil
}
