package benchshit

import (
	"fmt"
	"os"
	"runtime"
)

func mem() {
	pageSize := os.Getpagesize()
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf(
		"HeapSys = %.3f MiB, HeapAlloc =  %.3f MiB,  %.3f pages\n",
		float64(m.HeapSys)/1024.0/1024.0,
		float64(m.HeapAlloc)/1024.0/1024.0,
		float64(m.HeapSys)/float64(pageSize),
	)
	i := 100
	for ; i < 1000000000; i *= 10 {
		runtime.GC()
		s := make([]byte, i)
		runtime.ReadMemStats(&m)
		fmt.Printf(
			"%.3f MiB, HeapSys = %.3f MiB, HeapAlloc =  %.3f MiB,  %.3f pages\n",
			float64(len(s))/1024.0/1024.0,
			float64(m.HeapSys)/1024.0/1024.0,
			float64(m.HeapAlloc)/1024.0/1024.0,
			float64(m.HeapSys)/float64(pageSize),
		)
	}
	for ; i >= 100; i /= 10 {
		runtime.GC()
		s := make([]byte, i)
		runtime.ReadMemStats(&m)
		fmt.Printf(
			"%.3f MiB, HeapSys = %.3f MiB, HeapAlloc =  %.3f MiB,  %.3f pages\n",
			float64(len(s))/1024.0/1024.0,
			float64(m.HeapSys)/1024.0/1024.0,
			float64(m.HeapAlloc)/1024.0/1024.0,
			float64(m.HeapSys)/float64(pageSize),
		)
	}
	runtime.GC()
	runtime.ReadMemStats(&m)
	fmt.Printf(
		"\nHeapSys = %.3f MiB, HeapAlloc =  %.3f MiB,  %.3f pages\n",
		float64(m.HeapSys)/1024.0/1024.0,
		float64(m.HeapAlloc)/1024.0/1024.0,
		float64(m.HeapSys)/float64(pageSize),
	)
}
