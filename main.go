package main

import (
	"os"

	"github.com/shreychen/ignitebeat/cmd"

	_ "github.com/shreychen/ignitebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
