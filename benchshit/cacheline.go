package benchshit

import (
	"fmt"
	"time"
)

const size = 33554432 // 32 MB
func Cpy(arr1 []uint8, arr2 []uint8, slice int) {
	for i := 0; i < len(arr1); i += slice {
		arr2[i] = arr1[i]
	}
}

func AverageMinMax(f func() float64) (float64, float64, float64) {
	var sum float64
	var minimum float64
	var maximum float64

	for i := 0; i < 10; i++ {
		arr1 = make([]uint8, size)
		arr2 = make([]uint8, size)

		v := f()
		sum += v
		if i == 0 || v < minimum {
			minimum = v
		}
		if i == 0 || v > maximum {
			maximum = v
		}
	}
	return sum / 10, minimum, maximum
}

var arr1 []uint8
var arr2 []uint8

func cachelineRun(size int, slice int) float64 {
	start := time.Now()
	times := 10
	for i := 0; i < times*slice; i++ {
		Cpy(arr1, arr2, slice)
	}
	end := time.Now()
	dur := float64(end.Sub(start)) / float64(times*slice)
	return dur
}

func cacheline() {
	for slice := 16; slice <= 4096; slice *= 2 {
		a, m, M := AverageMinMax(func() float64 { return cachelineRun(size, slice-1) })
		fmt.Printf("%10d: %10.1f GB/s [%4.1f - %4.1f]\n", slice-1, float64(size)/a, float64(size)/M, float64(size)/m)
	}
}
