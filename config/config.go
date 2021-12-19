package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/jbowes/cling"
)

const DEFAULT_CACHE = ".protoc-cache"

// Config represents configuration for the protoc tool
type Config struct {
	Go           *Go          `hcl:"go,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
	cache        *string      `hcl:"cache,attr"`
	Cache        string
}

// Go represents go.mod configuration for the protoc tool
type Go struct {
	Path    string    `hcl:"path"`
	Ignores *[]string `hcl:"ignore,attr"`
}

// Dependency represents a dependency for the protoc tool
type Dependency struct {
	Type    string `hcl:"type,label"`
	Name    string `hcl:"name,label"`
	Version string `hcl:"version"`
}

// String returns a string representation of the dependency
func (d *Dependency) String() string {
	return fmt.Sprintf("%s %s %s", d.Type, d.Name, d.Version)
}

func Parse(path string) (*Config, error) {
	cfg := &Config{}

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

	if cfg.cache == nil || *cfg.cache == "" {
		defaultPath := filepath.Join(dir, DEFAULT_CACHE)
		cfg.Cache = defaultPath
	} else {
		if !strings.HasPrefix(*cfg.cache, "/") {
			cache := filepath.Join(dir, *cfg.cache)
			cfg.Cache = cache
		}
	}

	return cfg, nil
}
