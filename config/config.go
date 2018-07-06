// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"time"

	"gopkg.in/yaml.v2"
)

type SQL struct {
	sql  string `yaml:"sql"`
	size int    `yaml:"size"`
}

type Query struct {
	CacheName string `yaml:"cache_name"`
	Sqls      []SQL  `yaml:"sqls,flow"`
}

type Config struct {
	Period      time.Duration `config:"period"`
	Server      string        `config:"server"`
	NodeMetric  bool          `config:"node_metric"`
	CacheMetric bool          `config:"cache_metric"`
	AllCache    bool          `config:"all_cache"`
	CacheList   []string      `config:"cache_list"`
	SQL         bool          `config:"sql"`
	Queries     []Query       `yaml:"queries,flow"`
}

var DefaultConfig = Config{
	Period:      10 * time.Second,
	Server:      "http://127.0.0.1:8080",
	NodeMetric:  true,
	CacheMetric: false,
	AllCache:    true,
	SQL:         false,
}

func (c *Config) ToString() string {
	s, _ := yaml.Marshal(&c)
	return string(s)
}
