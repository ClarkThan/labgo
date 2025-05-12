package main

import (
	"fmt"
	"testing"
)

func IsOdd(i int) bool {
	return i%2 == 1
}

//go:noinline
func IsOddNoInline(i int) bool {
	return i%2 == 1
}

func BenchmarkCountOddInline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sum := 0
		for i := 1; i < 1000; i++ {
			if IsOdd(i) {
				sum += i
			}
		}
	}
}

func BenchmarkCountOddNoinline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		sum := 0
		for i := 1; i < 1000; i++ {
			if IsOddNoInline(i) {
				sum += i
			}
		}
	}
}

func main() {
	res1 := testing.Benchmark(BenchmarkCountOddInline)
	fmt.Println("BenchmarkCountOddInline", res1)
	res2 := testing.Benchmark(BenchmarkCountOddNoinline)
	fmt.Println("BenchmarkCountOddNoinline", res2)
}
