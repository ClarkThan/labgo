package promshit

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func kickStats() {
	go func() {
		for {
			// helloKick.Inc()
			seed := rand.Intn(100)
			doAlloc(seed)
			if seed < 5 || seed > 15 {
				time.Sleep(time.Duration(seed) * time.Second)
			} else {
				time.Sleep(time.Duration(seed*15) * time.Second)
			}
		}
	}()
}

func doAlloc(n int) {
	array := make([][]byte, n)
	for i := 0; i < n; i++ {
		array[i] = make([]byte, i+1)
	}
}

func fakeAllocate(w http.ResponseWriter, req *http.Request) {
	// get query parameters n from request
	n := req.URL.Query().Get("n")
	num, err := strconv.Atoi(n)
	if err != nil || num < 0 || num > 5000 {
		num = 500
	}

	doAlloc(num)

	fmt.Fprintf(w, "ok")
}

func registry_demo() {
	// 创建一个没有任何 label 标签的 gauge 指标
	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "home_temperature_celsius",
		Help: "The current temperature in degrees Celsius.",
	})

	// 在默认的注册表中注册该指标
	prometheus.MustRegister(temp)

	// 设置 gauge 的值为 39
	temp.Set(39)
}
