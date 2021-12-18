package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type includesFlag []string

func (i *includesFlag) String() string {
	return strings.Join(*i, ",")
}

func (i *includesFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var includeRemote includesFlag

func main() {
	path := flag.String("proto-path", "", "proto file(s) to compile")
	swaggerOut := flag.String("swagger-out", ".", "output path")
	flag.Var(&includeRemote, "include-remote-git", "include remote modules")
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

	mods, err := protoIncludes()
	if err != nil {
		panic(err)
	}
	//defer mods.Cleanup()

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
