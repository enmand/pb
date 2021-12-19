package main

import (
	"flag"
	"os"

	"github.com/enmand/pb/internal/config"
	"github.com/enmand/pb/internal/protoc"
)

func main() {
	path := flag.String("proto-path", "", "proto file(s) to compile")
	configFile := flag.String("config", "", "protoc tool config file")
	flag.Parse()

	if *path == "" || path == nil {
		p := "*.proto"
		path = &p
	}

	cfg, err := config.Parse(*configFile)
	if err != nil {
		panic(err)
	}

	cmpl := protoc.New(cfg)

	cmd, err := cmpl.Exec(*path)
	if err != nil {
		panic(err)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
