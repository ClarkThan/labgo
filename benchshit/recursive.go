package benchshit

import (
	"fmt"
	"testing"
)

var array = make([]int, 1000)

func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	return n * Factorial(n-1)
}

func FactorialLoop(n int) int {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return result
}

func BenchmarkFillNoinline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 1; i < 1000; i++ {
			array[i] = Factorial(i)
		}
	}
}

func BenchmarkFillInline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 1; i < 1000; i++ {
			array[i] = FactorialLoop(i)
		}
	}
}

func recursive() {
	res1 := testing.Benchmark(BenchmarkFillNoinline)
	fmt.Println("BenchmarkFillNoinline", res1)
	res2 := testing.Benchmark(BenchmarkFillInline)
	fmt.Println("BenchmarkFillInline", res2)
	fmt.Println(float64(res1.NsPerOp()) / float64(res2.NsPerOp()))
}
