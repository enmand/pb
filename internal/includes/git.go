package includes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jbowes/cling"
	"github.com/unerror/id-hub/tools/protoc/internal/config"
)

func GitDependencies(dep config.Dependency, cache string) (*Module, error) {
	// iterate over the deps where the key is the repository URL and the value is the version
	// (commit or tag) to clone each repository using github.com/go-git/go-git/v5

	path := filepath.Join(cache, dep.Name)
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, cling.Wrap(err, "unable to stat repository")
		}

		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, cling.Wrap(err, "unable to create repository directory")
		}
	} else {
		// path exists, so we need to check if the version is the same
		r, err := git.PlainOpen(path)
		if err != nil {
			return nil, cling.Wrap(err, "unable to open git repository")
		}

		ref, err := r.Head()
		if err != nil {
			return nil, cling.Wrap(err, "unable to get head")
		}

		if ref.Name().String() == dep.Version || ref.Name().Short() == dep.Version || ref.Hash().String() == dep.Version {
			// the version is the same, so we can skip this repository
			return &Module{Name: dep.Name, Path: path}, nil
		}

		// the version is different, so we need to remove the old repository
		err = os.RemoveAll(path)
		if err != nil {
			return nil, cling.Wrap(err, "unable to remove old repository")
		}
	}

	// clone the repo
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:          fmt.Sprintf("git://git@%s", dep.Name),
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
	})
	if err != nil {
		return nil, cling.Wrap(err, "unable to clone repo")
	}

	r.RepackObjects(&git.RepackConfig{
		UseRefDeltas:             false,
		OnlyDeletePacksOlderThan: time.Now(),
	})

	// get the commit hash for the version
	commit, err := r.ResolveRevision(plumbing.Revision(dep.Version))
	if err != nil {
		return nil, cling.Wrap(err, "unable to resolve revision")
	}

	// checkout the commit
	w, err := r.Worktree()
	if err != nil {
		return nil, cling.Wrap(err, "unable to get worktree")
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: *commit,
	})
	if err != nil {
		return nil, cling.Wrap(err, "unable to checkout commit")
	}

	return &Module{
		Name: dep.Name,
		Path: path,
	}, nil
}
