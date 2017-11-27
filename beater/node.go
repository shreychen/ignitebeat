package beater

import (
	//	"net"
	//	"net/http"
	//	"net/url"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/elastic/beats/libbeat/logp"
)

const selectorDetail = "json"
const NODE_METRIC = "/ignite?cmd=node&mtr=true"

type NodeMetric struct {
	Avg_active_jobs         int64   `json:"averageActiveJobs"`
	Avg_cancelled_jobs      int64   `json:"averageCancelledJobs"`
	Avg_CPU_load            float64 `json:"averageCpuLoad"`
	Avg_job_exetime         float64 `json:"averageJobExecuteTime"`
	Avg_job_waittime        float64 `json:"averageJobWaitTime"`
	Avg_rejected_jobs       int64   `json:"averageRejectedJobs"`
	Avg_waiting_jobs        int64   `json:"averageWaitingJobs"`
	Busy_time_pct           float64 `json:"busyTimePercentage"`
	Crt_active_jobs         float64 `json:"currentActiveJobs"`
	Crt_cancelled_jobs      float64 `json:"currentCancelledJobs"`
	Crt_CPU_load            float64 `json:"currentCpuLoad"`
	Crt_daemon_thread_count int64   `json:"currentDaemonThreadCount"`
	Crt_GC_CPU_load         int64   `json:"currentGcCpuLoad"`
	Crt_idle_time           int64   `json:"currentIdleTime"`
	Crt_job_exetime         int64   `json:"currentJobExecuteTime"`
	Crt_job_waittime        int64   `json:"currentJobWaitTime"`
	Crt_rejected_jobs       int64   `json:"currentRejectedJobs"`
	Crt_thread_count        int64   `json:"currentThreadCount"`
	Crt_waiting_jobs        int64   `json:"currentWaitingJobs"`
	Heap_committed          int64   `json:"heapMemoryCommitted"`
	Heap_initialized        int64   `json:"heapMemoryInitialized"`
	Heap_max                int64   `json:"heapMemoryMaximum"`
	Heap_used               int64   `json:"heapMemoryUsed"`
	Idle_time_pct           float64 `json:"idleTimePercentage"`
	Max_active_jobs         int64   `json:"maximumActiveJobs"`
	Max_cancelled_jobs      int64   `json:"maximumCancelledJobs"`
	Max_job_exetime         int64   `json:"maximumJobExecuteTime"`
	Max_job_waittime        int64   `json:"maximumJobWaitTime"`
	Max_rejected_jobs       int64   `json:"maximumRejectedJobs"`
	Max_thread_count        int64   `json:"maximumThreadCount"`
	Max_waiting_jobs        int64   `json:"maximumWaitingJobs"`
	Nonheap_committed       int64   `json:"nonHeapMemoryCommitted"`
	Nonheap_initialized     int64   `json:"nonHeapMemoryInitialized"`
	Nonheap_max             int64   `json:"nonHeapMemoryMaximum"`
	Nonheap_used            int64   `json:"nonHeapMemoryUsed"`
	Received_bytes          int64   `json:"receivedBytesCount"`
	Received_msg            int64   `json:"receivedMessagesCount"`
	Sent_bytes              int64   `json:"sentBytesCount"`
	Sent_msg                int64   `json:"sentMessagesCount"`
	Total_busy_time         int64   `json:"totalBusyTime"`
	Total_cancelled_jobs    int64   `json:"totalCancelledJobs"`
	Total_executed_jobs     int64   `json:"totalExecutedJobs"`
	Total_executed_tasks    int64   `json:"totalExecutedTasks"`
	Total_idle_time         int64   `json:"totalIdleTime"`
	Total_rejected_jobs     int64   `json:"totalRejectedJobs"`
	Total_started_thread    int64   `json:"totalStartedThreadCount"`
	Uptime                  float64 `json:"upTime"`
}

func (ib *Ignitebeat) GetNodeMetrics() (NodeMetric, error) {
	metric := NodeMetric{}
	myip, _ := GetMyIP()
	node_metric_url := fmt.Sprintf("%s%s&ip=%s", ib.config.Server, NODE_METRIC, myip)
	logp.Info("read statistic info from %s", node_metric_url)

	body, err := OpenURL(node_metric_url)
	if err != nil {
		logp.Info(err.Error())
		return metric, err
	}

	logp.Debug(selectorDetail, "body[%s]", string(body))

	r, _ := regexp.Compile(`(?ms)"metrics":[^}]*}`)
	metric_string := r.Find(body)
	if metric_string == nil {
		return metric, fmt.Errorf("no metric found")
	}

	metric_string = metric_string[10:]

	err = json.Unmarshal(metric_string, &metric)
	return metric, err

}
