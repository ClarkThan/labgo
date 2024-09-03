package poolshit

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ClarkThan/labgo/utils"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type GoPool struct {
	jobChan chan func()
	sema    chan struct{}
	closed  chan struct{}
}

func NewGoPool(size int) *GoPool {
	p := &GoPool{
		jobChan: make(chan func()),
		sema:    make(chan struct{}, size),
		closed:  make(chan struct{}),
	}

	return p
}

func (p *GoPool) Submit(fn func()) {
	select {
	case <-p.closed:
		return
	case p.sema <- struct{}{}:
		select {
		case <-p.closed:
			// <-p.sema
			return
		default:
			go p.Execute(fn)
		}
	default:
		p.jobChan <- fn
	}
}

func (p *GoPool) Execute(fn func()) {
	defer func() {
		<-p.sema
	}()
	for {
		fn()
		select {
		case <-p.closed:
			return
		case fn = <-p.jobChan:
		}
	}
}

func (p *GoPool) Stop() {
	close(p.closed)
}

func Main() {
	p := NewGoPool(5)
	go func() {
		for i := 0; i < 25; i++ {
			n := i
			p.Submit(func() {
				// time.Sleep(time.Duration(random.Int31n(2000) * int32(time.Millisecond)))
				time.Sleep(200 * time.Millisecond)
				fmt.Println(utils.GoID(), "got", n)
			})
		}
	}()
	fmt.Println("-----launch-----")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("-----over-----")
	p.Stop()
}
