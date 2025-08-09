package main

import (
	"os"

	"github.com/nduyhai/go-clean-arch-starter/internal/adapters/inbound/cli"
)

func main() {
	root := cli.NewRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
