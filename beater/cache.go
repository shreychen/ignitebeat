package beater

import (
	//  "net"
	//  "net/http"
	//  "net/url"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/elastic/beats/libbeat/logp"
)

const modePartitioned = "PARTITIONED"
const CACHE_METRIC = "/ignite?cmd=cache&cacheName="

type Cache struct {
	Name      string `json:"name"`
	Mode      string `json:"mode"`
	SqlSchema string `json:"sqlSchema"`
}

type CacheMetric struct {
	Read  int64
	Write int64
	Hit   int64
	Miss  int64
}

type CacheMetricResponse struct {
	successStatus  int
	affinityNodeId string
	Err            string
	sessionToken   string
	response       CacheMetric
}

func (ib *Ignitebeat) GetCacheList() ([]Cache, error) {
	var cache_list []Cache

	myip, _ := GetMyIP()
	node_metric_url := fmt.Sprintf("%s%s&ip=%s", ib.config.Server, NODE_METRIC, myip)
	logp.Info("read statistic info from %s", node_metric_url)

	body, err := OpenURL(node_metric_url)
	if err != nil {
		logp.Info(err.Error())
		return cache_list, err
	}

	logp.Debug(selectorDetail, "body[%s]", string(body))

	r, _ := regexp.Compile(`(?ms)"caches":[^\]]*]`)
	sub_body := r.Find(body)
	if sub_body == nil {
		return cache_list, fmt.Errorf("can't get cache list")
	}

	sub_body = sub_body[9:]

	err = json.Unmarshal(sub_body, &cache_list)
	return cache_list, err
}

func (ib *Ignitebeat) GetCacheMetric(c *Cache) (CacheMetric, error) {
	metric := CacheMetric{}

	cache_metric_url := fmt.Sprintf("%s%s%s", ib.config.Server, CACHE_METRIC, c.Name)
	logp.Info("Get metrics of Cache %s from %s", c.Name, cache_metric_url)

	body, err := OpenURL(cache_metric_url)
	if err != nil {
		logp.Info(err.Error())
		return metric, err
	}

	logp.Debug(selectorDetail, "body[%s]", string(body))

	cache_rsp := CacheMetricResponse{}
	err = json.Unmarshal(body, &cache_rsp)

	if err != nil {
		logp.Err(err.Error())
		return metric, err
	}

	return cache_rsp.response, err
}
