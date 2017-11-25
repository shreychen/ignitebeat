package beater

import (
	//	"encoding/json"
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
						"type":   b.Info.Name,
						"metric": node_metric,
					},
				}
				bt.client.Publish(event)
				logp.Info("Metric event sent")
			}
		}

		//              counter++
	}
}

func (bt *Ignitebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

