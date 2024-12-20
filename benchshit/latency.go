package benchshit

import (
	"fmt"
	"testing"
)

type Node struct {
	data int
	next *Node
}

func build(volume int) *Node {
	var head *Node
	for i := 0; i < volume; i++ {
		head = &Node{i, head}
	}
	return head
}

var list *Node
var N int

func BenchmarkLen(b *testing.B) {
	for n := 0; n < b.N; n++ {
		len := 0
		for p := list; p != nil; p = p.next {
			len++
		}
		if len != N {
			b.Fatalf("invalid length: %d", len)
		}
	}
}

func latency() {
	N = 1000000
	list = build(N)
	res := testing.Benchmark(BenchmarkLen)
	fmt.Println("milliseconds: ", float64(res.NsPerOp())/1e6)

	fmt.Println("nanoseconds per el.", float64(res.NsPerOp())/float64(N))
}
