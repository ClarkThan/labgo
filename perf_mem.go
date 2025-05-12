package main

import (
	"fmt"
	"runtime"
	"strconv"
)

var initAlloc = func() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}()

func printAlloc(prefix string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	size := (m.Alloc - initAlloc) / (1024)
	print(prefix + ": heap " + strconv.Itoa(int(size)) + "\n")
}

type Client struct {
	id   uint64
	body [40]byte
}

func main() {
	print("values are KB\n")
	printAlloc("initial state     ")

	m := make(map[int]Client)
	printAlloc("after declaration ")

	for i := range 1000 {
		m[i] = Client{id: uint64(i)}
	}
	runtime.GC()
	fmt.Println("len: ", len(m))
	printAlloc("after insertion 1000   ")

	for i := 1000; i < 10000; i++ {
		m[i] = Client{id: uint64(i)}
	}
	runtime.GC()
	fmt.Println("len: ", len(m))
	printAlloc("after insertion 10000   ")

	for i := range 9000 {
		delete(m, i)
	}
	runtime.GC()
	fmt.Println("len: ", len(m))
	printAlloc("after deletion 9/10")
}
