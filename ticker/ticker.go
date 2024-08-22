package ticker

import (
	"fmt"
	"time"
)

type MyTicker struct {
	*time.Ticker
	stopChan chan struct{}
}

func (t *MyTicker) Stop() {
	defer t.Ticker.Stop()
	t.stopChan <- struct{}{}
	fmt.Println("overing")
}

func NewMyTicker(t time.Duration) *MyTicker {
	return &MyTicker{
		Ticker:   time.NewTicker(t),
		stopChan: make(chan struct{}),
	}
}

func test1(ticker *MyTicker) {
	time.AfterFunc(10*time.Second, func() {
		fmt.Println("ready stop")
		ticker.Stop()
	})

	for {
		select {
		case <-ticker.C:
			fmt.Println("now: ", time.Now())
		case <-ticker.stopChan:
			fmt.Println("receive stop signal")
			return
		}
	}
}

func test2(t *time.Ticker) {
	time.AfterFunc(10*time.Second, func() {
		fmt.Println("ready stop")
		t.Stop()
	})

	for range t.C {
		fmt.Println("now: ", time.Now())
	}

	fmt.Println("over!")
}

func Main() {
	// ticker := NewMyTicker(2 * time.Second)
	// test1(ticker)
	test2(time.NewTicker(2 * time.Second))
}
