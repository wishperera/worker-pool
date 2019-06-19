package worker_pool

import "github.com/prometheus/client_golang/prometheus"

var (
	activeJobs prometheus.Gauge
	workerBuffer prometheus.GaugeFunc
	processLatency prometheus.Summary
)

func initMetrics(pool *Pool){
	activeJobs = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: pool.conf.Metrics.NameSpace,
		Subsystem: pool.conf.Metrics.SubSystem,
		Name: "active_job_count",
		Help: "number of active jobs in the pool",
	})

	workerBuffer = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace:pool.conf.Metrics.NameSpace,
		Subsystem: pool.conf.Metrics.SubSystem,
		Name: "worker_buffer_size",
		Help: "worker buffer size in terms of jobs",
	}, func() float64 {
		var s = 0
		for _,w := range pool.manager.workers{
			s = s+(len(w.buffer)-cap(w.buffer))
		}
		return float64(s)
	})

	processLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: pool.conf.Metrics.NameSpace,
		Subsystem: pool.conf.Metrics.SubSystem,
		Name:      "process_latency",
		Help:      "process latency in milliseconds",
	})
}
