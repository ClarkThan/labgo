package poolshit

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/ClarkThan/labgo/utils"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

func demo_lru1() {
	// make cache with 10ms TTL and 5 max keys
	cache := expirable.NewLRU[string, string](5, nil, time.Millisecond*100)

	// set value under key1.
	cache.Add("key1", "val1")

	// get value under key1
	r, ok := cache.Get("key1")

	time.Sleep(time.Millisecond * 50)

	cache.Add("key2", "val2")

	// check for OK value
	if ok {
		fmt.Printf("value before expiration is found: %v, value: %q\n", ok, r)
	}

	// wait for cache to expire
	time.Sleep(time.Millisecond * 60)

	// get value under key1 after key expiration
	r, ok = cache.Get("key1")
	fmt.Printf("value after expiration is found: %v, value: %q\n", ok, r)

	r, ok = cache.Get("key2")
	fmt.Printf("key2 is found: %v, value: %q\n", ok, r)

	fmt.Printf("Cache len: %d\n", cache.Len())

	// set value under key2, would evict old entry because it is already expired.
	cache.Add("key3", "val3")

	fmt.Printf("Cache len: %d\n", cache.Len())
	// Output:
	// value before expiration is found: true, value: "val1"
	// value after expiration is found: false, value: ""
	// Cache len: 1
}

func demo_lru2() {
	l, _ := lru.New[int, int](16)
	for i := 1; i <= 32; i++ {
		l.Add(i, i+100)
	}
	if l.Len() != 16 {
		panic(fmt.Sprintf("bad len: %v", l.Len()))
	}
	l.Add(5, 55)

	// fmt.Println(l.Keys())

	val, ok := l.Get(5)
	fmt.Println(ok, val)
	oldKey, oldVal, ok := l.GetOldest()
	fmt.Println(oldKey, oldVal, ok)

	val, ok = l.Get(99)
	fmt.Println(ok, val)
}

func demo_lru3() {
	utils.TraceMemStats()
	l, _ := lru.New[string, bool](1 << 18)
	utils.TraceMemStats()
	for i := 0; i < 1<<18; i++ {
		l.Add(genMD5(strconv.Itoa(i+100)), true)
	}
	runtime.GC()
	time.Sleep(time.Second)
	utils.TraceMemStats()
	fmt.Println("keys:", l.Len())
}

func demo_lru4() {
	o := md5.New().Sum(utils.String2Bytes("1111"))
	x := hex.EncodeToString(o[:16])
	fmt.Println(x)

	y := hex.EncodeToString(md5.New().Sum([]byte("1111sddwewvfxxxxxxxxxxx"))[:16])
	fmt.Println(y)

	h := md5.New()
	h.Write([]byte("1111sddwewvfxxxxxxxxxxx"))
	ret := h.Sum(nil)
	fmt.Println(hex.EncodeToString(ret))
}

func genMD5(data string) string {
	h := md5.New()
	h.Write(utils.String2Bytes(data))
	return hex.EncodeToString(h.Sum(nil))
}
