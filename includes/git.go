package includes

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jbowes/cling"
)

func Git(deps map[string]string, cache string) (Modules, error) {
	// iterate over the deps where the key is the repository URL and the value is the version
	// (commit or tag) to clone each repository using github.com/go-git/go-git/v5
	ms := Modules{}
	for name, version := range deps {
		path := filepath.Join(cache, name)
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

			if ref.Name().String() == version || ref.Name().Short() == version || ref.Hash().String() == version {
				// the version is the same, so we can skip this repository
				continue
			}

			// the version is different, so we need to remove the old repository
			err = os.RemoveAll(path)
			if err != nil {
				return nil, cling.Wrap(err, "unable to remove old repository")
			}
		}

		// clone the repo
		r, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      fmt.Sprintf("git://git@%s", name),
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, cling.Wrap(err, "unable to clone repo")
		}

		// get the commit hash for the version
		commit, err := r.ResolveRevision(plumbing.Revision(version))
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
	}

	for repo := range deps {
		path := filepath.Join(cache, repo)
		ms = append(ms, Module{
			Name: repo,
			Path: path,
		})
	}

	return ms, nil
}
