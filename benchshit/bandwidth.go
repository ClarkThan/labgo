package benchshit

import (
	"fmt"
	"time"
)

func run() float64 {
	bestbandwidth := 0.0
	arr := make([]uint8, 2*1024*1024*1024) // 4 GB
	for i := 0; i < len(arr); i++ {
		arr[i] = 1
	}
	for t := 0; t < 20; t++ {
		start := time.Now()
		acc := 0
		for i := 0; i < len(arr); i += 64 {
			acc += int(arr[i])
		}
		end := time.Now()
		if acc != len(arr)/64 {
			panic("!!!")
		}
		bandwidth := float64(len(arr)) / end.Sub(start).Seconds() / 1024 / 1024 / 1024
		if bandwidth > bestbandwidth {
			bestbandwidth = bandwidth
		}
	}
	return bestbandwidth
}

func bandwidth() {
	for i := 0; i < 10; i++ {
		fmt.Printf(" %.2f GB/s\n", run())
	}
}
