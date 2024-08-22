package ctxshit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	concurrentNum int = 5
)

var (
	errStop      = errors.New("stop poll")
	errDataQuery = errors.New("query data fail")
)

func poll(ctx context.Context) error {
	userChan := make(chan string, concurrentNum)

	go func() {
		defer func() {
			fmt.Println("produce over!")
			close(userChan)
		}()

		var i int
		for {
			select {
			case <-ctx.Done():
				return
			default:
				userChan <- strconv.Itoa(i)
				time.Sleep(time.Duration(100+i*10) * time.Millisecond)
				i++
				if i >= 100 {
					return
				}
			}
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < concurrentNum; i++ {
		wg.Add(1)
		go func(ctx context.Context, userChan chan string) {
			defer wg.Done()
			for userID := range userChan {
				select {
				case <-ctx.Done():
					fmt.Println("stop at ", userID)
					return
				default:
					fmt.Println("process ", userID)
				}
			}
		}(ctx, userChan)
	}

	fmt.Println("shit...")
	wg.Wait()
	fmt.Println("fuck...")

	select {
	case <-ctx.Done():
		fmt.Println("close ctx:", ctx.Err())
		return ctx.Err()
	default:
	}

	return nil
}

func loop(ctx context.Context) {
	for {
		if err := poll(ctx); err != nil {
			fmt.Printf("got err: %v\n", err)
			break
		}
	}
	fmt.Println("over loop ...")
}

type Poller struct {
	stop chan struct{}
	wg   *sync.WaitGroup
}

func NewPoller() *Poller {
	return &Poller{
		stop: make(chan struct{}),
		wg:   new(sync.WaitGroup),
	}
}

func (p *Poller) poll(ctx context.Context) error {
	todoChan := make(chan string, 10)

	go func() {
		defer func() {
			fmt.Println("produce over!")
			close(todoChan)
		}()

		var i int
		for { // 不停地产生任务
			select {
			case <-p.stop:
				log.Println("fetcher recv stop signal")
				return
			default:
				todoChan <- strconv.Itoa(i)
				time.Sleep(time.Duration(100) * time.Millisecond)
				i++
				if i >= 50 {
					return
				}
			}
		}
	}()

	startTS := time.Now()

	limiter := make(chan struct{}, concurrentNum)
	for userID := range todoChan {
		limiter <- struct{}{}
		p.wg.Add(1)
		go func(ctx context.Context, userID string, limiter chan struct{}) {
			log.Println("start a new goroutine", userID)
			defer func() {
				p.wg.Done()
				<-limiter
				// log.Println("over goroutine", userID)
			}()

			select {
			case <-p.stop:
				return
			default:
				select {
				case <-p.stop:
					log.Println("processor recv stop signal", userID)
					return
				default:
					fmt.Println("processing ...", userID)
					time.Sleep(300 * time.Millisecond)
					fmt.Println("processed ", userID)
				}
			}
		}(ctx, userID, limiter)
	}

	p.wg.Wait() // 等待已经开启的 goroutine 处理完才关闭
	log.Println("over processing...")

	select {
	case <-p.stop:
		log.Println("main recv stop signal")
		return errStop
	default:
	}

	if time.Since(startTS) < 2*time.Second {
		log.Println("处理得太快了, 先 sleep 5s")
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (p *Poller) Stop() {
	close(p.stop)
}

func (p *Poller) HasStop() bool {
	select {
	case <-p.stop:
		return true
	default:
	}

	return false
}

func (p *Poller) Start() bool {
	if !p.HasStop() {
		return false
	}

	p.stop = make(chan struct{})
	go p.PollLoop()
	return true
}

func (p *Poller) PollLoop() {
	ctx := context.Background()

	go func() {
		for {
			if err := p.poll(ctx); err != nil {
				fmt.Printf("got err: %v\n", err)
				break
			}
		}

		fmt.Println("over loop ...")
	}()
}

func Main() {
	poller := NewPoller()
	poller.PollLoop()

	http.HandleFunc("/stop", func(w http.ResponseWriter, req *http.Request) {
		poller.Stop()
		w.Write([]byte(`emit the stop signal`))
	})

	http.ListenAndServe(":8090", nil)
}
