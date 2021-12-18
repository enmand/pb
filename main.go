package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jbowes/cling"
	"github.com/unerror/id-hub/tools/protoc/includes"
	"github.com/unerror/id-hub/tools/protoc/mod"
)

func main() {
	path := flag.String("proto-path", "", "proto file(s) to compile")
	swaggerOut := flag.String("swagger-out", ".", "output path")
	includeRemote := flag.String("config", "", "protoc tool config file")
	flag.Parse()

	if swaggerOut != nil {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		out := fmt.Sprintf("%s/%s", cwd, filepath.Clean(*swaggerOut))
		swaggerOut = &out
	}

	if *path == "" || path == nil {
		p := "*.proto"
		path = &p
	}

	mods, err := getIncludes(*includeRemote)
	if err != nil {
		panic(err)
	}
	defer mods.Cleanup()

	cmd := exec.Command(
		"protoc",
		"-I",
		".",
		"-I",
		mods.String(),
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		fmt.Sprintf("--swagger_out=%s/openapi", *swaggerOut),
		"--grpc-gateway_out=allow_patch_feature=false,paths=source_relative:.",
		*path,
	)
	fmt.Println(cmd)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func getIncludes(remotes string) (includes.Modules, error) {
	f, err := mod.Parse(remotes)
	if err != nil {
		return nil, cling.Wrap(err, "unable to parse go.mod")
	}

	mods := includes.Modules{}

	if f.Go != nil {
		gms, err := includes.GoMod()
		if err != nil {
			return nil, cling.Wrap(err, "unable to get proto includes from go modules")
		}

		mods = append(mods, gms...)
	}

	for _, d := range f.Dependencies {
		switch d.Type {
		case "git":
			gms, err := includes.Git(d.Require, d.Cache)
			if err != nil {
				return nil, cling.Wrap(err, "unable to get proto includes from git")
			}
			mods = append(mods, gms...)
		case "local":
			gms, err := includes.Local(d.Require)
			if err != nil {
				return nil, cling.Wrap(err, "unable to get proto includes from local")
			}
			mods = append(mods, gms...)
		default:
			return nil, cling.Errorf("unable to get proto includes from %s", d.Type)
		}
	}

	return mods, nil
}
