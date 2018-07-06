package beater

import (
	"fmt"
	//	"io/ioutil"
	//	"net"

	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/shreychen/ignitebeat/config"
)

type Ignitebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Ignitebeat{
		done:   make(chan struct{}),
		config: config,
	}

	logp.Debug("json", "running configuraion as below:\n %s", config.ToString())

	return bt, nil
}

func (bt *Ignitebeat) Run(b *beat.Beat) error {
	logp.Info("ignitebeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	//      counter := 1

	is_ready := true
	var cache_list []Cache
	if bt.config.CacheMetric {
		if bt.config.AllCache {
			cache_list, err = bt.GetCacheList()
			if err != nil {
				is_ready = false
			}
		} else {
			for _, name := range bt.config.CacheList {
				_cache := Cache{Name: name}
				cache_list = append(cache_list, _cache)
			}
		}
		logp.Debug("json", "cache list: %s", cache_list)
	}

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		if bt.config.NodeMetric {
			node_metric, err := bt.GetNodeMetrics()
			if err != nil {
				logp.Err("can't get metrics, caused by: %s", err.Error())
			} else {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type": b.Info.Name,
						"node": node_metric,
					},
				}
				bt.client.Publish(event)
				logp.Info("Node metric event sent")
			}
		}

		if bt.config.CacheMetric && is_ready {
			cache_metrics := make(map[string]CacheMetric)
			for _, c := range cache_list {
				if cm, err := bt.GetCacheMetric(&c); err == nil {
					cache_metrics[c.Name] = cm
				} else {
					logp.Debug("json", "can't get metric of Cache %s, caused by: %s", c.Name, err.Error())
				}
			}

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":  b.Info.Name,
					"cache": cache_metrics,
				},
			}
			bt.client.Publish(event)
			logp.Info("Cache metric event sent")
		}

		//              counter++
	}
}

func (bt *Ignitebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
