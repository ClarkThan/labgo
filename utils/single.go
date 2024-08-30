package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 简化版本的SingleFlight
type SingleFlight struct {
	ok  atomic.Bool
	mu  sync.Mutex
	ret any
	err error
}

func (s *SingleFlight) Do(key string, fn func() (any, error)) (any, error) {
	if s.ok.Load() {
		return s.ret, s.err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ok.Load() {
		return s.ret, s.err
	}
	s.ret, s.err = fn()
	s.ok.Store(true)
	return s.ret, s.err
}

func test_singleflight() {
	var single SingleFlight
	var wg sync.WaitGroup
	key := "foo"
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			ret, err := single.Do(key, func() (any, error) {
				time.Sleep(500 * time.Millisecond)
				return rand.Int31n(100), nil
			})
			fmt.Println(x, ret, err)
		}(i)
	}

	wg.Wait()
}
