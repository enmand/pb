package mod

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jbowes/cling"
)

type Dependencies struct {
	Go           *Go          `hcl:"go,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
}

type Go struct {
	Path    string   `hcl:"path"`
	Ignores []string `hcl:"ignores"`
}

type Dependency struct {
	Type    string            `hcl:"type,label"`
	Require map[string]string `hcl:"require,attr"`
	Cache   string            `hcl:"cache"`
}

func Parse(path string) (*Dependencies, error) {
	cfg := &Dependencies{}

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, cling.Wrap(err, "unable to get absolute path")
	}
	dir := filepath.Dir(abs)

	err = hclsimple.DecodeFile(path, nil, cfg)
	if err != nil {
		return nil, cling.Wrap(err, "unable to parse config")
	}

	if cfg.Go != nil {
		if !strings.HasPrefix(cfg.Go.Path, "/") {
			cfg.Go.Path = filepath.Join(dir, cfg.Go.Path)
		}
	}

	for i, dep := range cfg.Dependencies {
		if dep.Cache == "" {
			dep.Cache = filepath.Join(dir, ".cache", dep.Type)
		} else {
			if !strings.HasPrefix(dep.Cache, "/") {
				dep.Cache = filepath.Join(dir, dep.Cache)
			}
		}

		cfg.Dependencies[i] = dep
	}

	return cfg, nil
}
