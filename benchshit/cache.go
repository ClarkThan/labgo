package benchshit

import (
	"fmt"
	"math/bits"
	"time"
)

func Shuffle(arr []uint32) {
	seed := uint64(1234)
	for i := len(arr) - 1; i > 0; i-- {
		seed += 0x9E3779B97F4A7C15
		hi, _ := bits.Mul64(seed, uint64(i+1))
		j := int(hi)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func averageMinMax(f func() float64) (float64, float64, float64) {
	var sum float64
	var minimum float64
	var maximum float64

	for i := 0; i < 10; i++ {
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

func cacheRun(size int) float64 {
	arr := make([]uint32, size)

	for i := range arr {
		arr[i] = uint32(i + 1)
	}
	start := time.Now()
	end := time.Now()
	times := 0
	for ; end.Sub(start) < 100_000_000; times++ {
		Shuffle(arr)
		end = time.Now()
	}
	dur := float64(end.Sub(start)) / float64(times)
	return dur / float64(size)
}

func cache() {
	for size := 4096; size <= 33554432; size *= 2 {
		fmt.Printf("%20d KB ", size/1024*4)
		a, m, M := averageMinMax(func() float64 { return cacheRun(size) })
		fmt.Printf(" %.2f [%.2f, %.2f]\n", a, m, M)
	}
}
