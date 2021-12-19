package config

import (
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
	Plugins      []Plugin     `hcl:"plugin,block"`
	cache        *string      `hcl:"cache,attr"`
	Cache        string
}

type Plugin struct {
	Name        string  `hcl:"name,label"`
	PathResolve *string `hcl:"path,attr"`
	Path        string
	RelPath     *string            `hcl:"rel_path,attr"`
	Options     *map[string]string `hcl:"options,attr"`
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

	resolved := []Plugin{}
	for _, p := range cfg.Plugins {
		if p.PathResolve != nil && !strings.HasPrefix(*p.PathResolve, "/") {
			p.Path = filepath.Join(dir, *p.PathResolve)
		}
		resolved = append(resolved, p)
	}
	cfg.Plugins = resolved

	return cfg, nil
}
