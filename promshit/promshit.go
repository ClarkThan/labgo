package promshit

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://zhuanlan.zhihu.com/p/592560633

var (
	// avg_over_time(x_request_total[3m])
	// rate(x_request_total[3m])  最近3分钟QPS
	reqCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "x_request_total",
		Help: "The total number of request of api",
	})

	// avg_over_time(x_queue_length[3m])  最近3分钟平均值
	queueGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "x_queue_length",
		Help: "The length of queue",
	})

	// rate(x_request_time_sum[3m])/rate(x_request_time_count[3m])  最近5分钟内，平均每次请求耗时是多少
	// histogram_quantile(0.95, sum(rate(x_request_time_bucket[3m])) by (le))
	reqTimeHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:        "x_request_time",
		Help:        "The delay of api request",
		ConstLabels: prometheus.Labels{"code": "200", "method": "GET"},
	})
	// x_request_duration_summary{quantile="0.99"}  百分之99的请求耗时范围
	reqTimeSummary = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "x_request_duration_summary",
		Help:       "The summary delay of api request",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func hello(w http.ResponseWriter, req *http.Request) {
	reqCounter.Inc()
	elapsed := float64(rand.Intn(500)) / 100.0
	log.Println("/hello", elapsed)
	reqTimeHistogram.Observe(elapsed)
	reqTimeSummary.Observe(elapsed)
	fmt.Fprintf(w, "ok")
}

func Main() {
	go func() {
		for {
			queueGauge.Set(float64(rand.Intn(100)))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()

	// Serve the default Prometheus metrics registry over HTTP on /metrics.
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/-/metrics", promhttp.Handler())
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8888", nil)
}
