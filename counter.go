package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Counter struct {
	rate  int           //计数周期最大允许请求次数
	start time.Time     //计数开始时间
	cycle time.Duration //计数周期
	count int           //计数周期内累计收到的请求数
	lock  sync.Mutex
}

func (l *Counter) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.count == l.rate-1 {
		now := time.Now()
		fmt.Println("now:", now)
		fmt.Println("l.start:", l.start, ",sub:", now.Sub(l.start))
		if now.Sub(l.start) >= l.cycle {
			l.Reset(now)
			return true
		} else {
			return false
		}
	} else {
		l.count++
		return true
	}
}

func (l *Counter) Set(r int, cycle time.Duration) {
	l.rate = r
	l.start = time.Now()
	l.cycle = cycle
	l.count = 0
}

func (l *Counter) Reset(t time.Time) {
	l.start = t
	l.count = 0
}

//计数器算法
func main() {
	var wg sync.WaitGroup
	var lr Counter
	lr.Set(3, time.Second) // 1s内最多请求3次
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求:", i)
		go func(i int) {
			if lr.Allow() {
				log.Println("响应请求:", i)
			}
			wg.Done()
		}(i)

		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}
