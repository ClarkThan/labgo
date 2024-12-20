package benchshit

import (
	"fmt"
	"os"
	"testing"
)

var fact int

func BenchmarkFactorial(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fact = 1
		for i := 1; i <= 10; i++ {
			fact *= i
		}
	}
}
func BenchmarkFactorialBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buffer := make([]int, 11)
		buffer[0] = 1
		for i := 1; i <= 10; i++ {
			buffer[i] = i * buffer[i-1]
		}
	}
	b.ReportAllocs()
}

func BenchmarkFactorialBufferLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buffer := make([]int, 100001)
		buffer[0] = 1
		for i := 1; i <= 100000; i++ {
			buffer[i] = i * buffer[i-1]
		}
	}
	b.ReportAllocs()
}

func basic() {
	pageSize := os.Getpagesize()
	fmt.Printf("Page size: %d bytes (%d KB)\n", pageSize, pageSize/1024)
	res := testing.Benchmark(BenchmarkFactorial)
	fmt.Println("BenchmarkFactorial", res)
	resmem := testing.Benchmark(BenchmarkFactorialBuffer)
	fmt.Println("BenchmarkFactorialBuffer", resmem, resmem.MemString())
	resmem = testing.Benchmark(BenchmarkFactorialBufferLarge)
	fmt.Println("BenchmarkFactorialBufferLarge", resmem, resmem.MemString())
}
