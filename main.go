package main

import (
	"os"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/shreychen/ignitebeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		logp.Err(err.Error())
		os.Exit(1)
	}
}
