package sseshit

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func sseHandler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	memT := time.NewTicker(time.Second * 2)
	defer memT.Stop()

	cpuT := time.NewTicker(time.Second * 2)
	defer memT.Stop()

	done := r.Context().Done()
	rc := http.NewResponseController(w)

	for {
		select {
		case <-done:
			fmt.Println("client disconnected")
			return
		case <-memT.C:
			m, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("get memory usage error: %v", err)
				return
			}
			// w.Write([]byte("event: ping\n\n"))
			_, err = fmt.Fprintf(w, "event:mem\ndata:Total:%d, Used:%d, Perc:%.2f%%\n\n", m.Total, m.Used, m.UsedPercent)
			if err != nil {
				log.Printf("write memory usage to SSE error: %v", err)
				return
			}
			rc.Flush()
		case <-cpuT.C:
			c, err := cpu.Times(false)
			if err != nil {
				log.Printf("get CPU usage error: %v", err)
				return
			}
			// w.Write([]byte("event: ping\n\n"))
			_, err = fmt.Fprintf(w, "event:cpu\ndata:User:%.2f%%, Sys:%.2f%%, Idle:%.2f%%\n\n", c[0].User, c[0].System, c[0].Idle)
			if err != nil {
				log.Printf("write CPU usage to SSE error: %v", err)
				return
			}
			rc.Flush()
		}
	}
}
