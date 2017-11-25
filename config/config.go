// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period     time.Duration `config:"period"`
	Server     string        `config:"server"`
	NodeMetric bool          `config:"node_metric"`
}

var DefaultConfig = Config{
	Period:     10 * time.Second,
	Server:     "http://localhost:8080",
	NodeMetric: true,
}
