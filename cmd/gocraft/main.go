package main

import (
	"os"

	"github.com/nduyhai/gocraft/internal/adapters/inbound/cli"
)

func main() {
	root := cli.NewRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
