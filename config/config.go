// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period      time.Duration `config:"period"`
	Server      string        `config:"server"`
	NodeMetric  bool          `config:"node_metric"`
	CacheMetric bool          `config:"cache_metric"`
	AllCache    bool          `config:"all_cache"`
	CacheList   []string      `config:"cache_list"`
}

var DefaultConfig = Config{
	Period:      10 * time.Second,
	Server:      "http://localhost:8080",
	NodeMetric:  true,
	CacheMetric: true,
	AllCache:    true,
}
