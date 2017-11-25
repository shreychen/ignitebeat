package main

import (
	"os"

	"github.com/shreychen/ignitebeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
