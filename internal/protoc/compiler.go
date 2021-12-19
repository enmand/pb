package protoc

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jbowes/cling"
	"github.com/unerror/id-hub/tools/protoc/internal/config"
	"github.com/unerror/id-hub/tools/protoc/internal/includes"
)

type Compiler struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Compiler {
	return &Compiler{
		cfg: cfg,
	}
}

func (c *Compiler) includes() (includes.Modules, error) {
	mods := includes.Modules{
		{Path: "."},
	}

	if c.cfg.Go != nil {
		gms, err := includes.GoModDependencies(*c.cfg.Go)
		if err != nil {
			return nil, cling.Wrap(err, "unable to get proto includes from go modules")
		}

		mods = append(mods, gms...)
	}

	for _, d := range c.cfg.Dependencies {
		switch d.Type {
		case "git":
			gm, err := includes.GitDependencies(d, c.cfg.Cache)
			if err != nil {
				return nil, cling.Wrap(err, "unable to get proto includes from git")
			}
			mods = append(mods, *gm)
		case "local":
			lm, err := includes.Local(d)
			if err != nil {
				return nil, cling.Wrap(err, "unable to get proto includes from local")
			}
			mods = append(mods, lm)
		default:
			return nil, cling.Errorf("unable to get proto includes from %s", d.Type)
		}
	}

	return mods, nil
}

func (c *Compiler) Exec(path string) (*exec.Cmd, error) {
	args := []string{}

	protoPaths, err := c.includes()
	if err != nil {
		return nil, cling.Wrap(err, "unable to get proto includes")
	}

	args = append(args, "-I", protoPaths.String())

	for _, p := range c.cfg.Plugins {
		var opts string
		optsKvp := []string{}

		if p.Options != nil {
			for k, v := range *p.Options {
				optsKvp = append(optsKvp, fmt.Sprintf("%s=%s", k, v))
			}
			if len(optsKvp) > 0 {
				opts = fmt.Sprintf("%s:", strings.Join(optsKvp, ","))
			}
		}

		switch {
		case p.Path != "" && p.RelPath != nil:
			return nil, cling.Errorf("plugin %s has both path and rel_path", p.Name)
		case p.Path == "" && p.RelPath == nil:
			return nil, cling.Errorf("plugin %s has neither path nor rel_path", p.Name)
		case p.Path != "":
			args = append(args, fmt.Sprintf("--%s_out=%s%s", p.Name, opts, p.Path))
		case p.RelPath != nil:
			args = append(args, fmt.Sprintf("--%s_out=%s%s", p.Name, opts, *p.RelPath))
		}
	}

	args = append(args, path)

	return exec.Command("protoc", args...), nil
}
