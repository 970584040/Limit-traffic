package main

import (
	"log"
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	rate       float64 //固定每秒固定出水速率
	capacity   float64 //桶的容量
	water      float64 //桶中当前水量
	lastLeakMs int64   //桶上次漏水时间

	lock sync.Mutex
}

func (l *LeakyBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano()
	eclipse := float64((now - l.lastLeakMs)) * l.rate //先执行漏水
	l.water = l.water - eclipse                       //剩余水量
	l.water = math.Max(0, l.water)                    //桶干了

	l.lastLeakMs = now

	if (l.water + 1) < l.capacity {
		l.water++
		return true
	} else {
		return false
	}
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}

//漏桶算法
func main() {
	var wg sync.WaitGroup
	var lr LeakyBucket
	lr.Set(1, 3)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		log.Println("创建请求:", i)
		go func(i int) {
			if lr.Allow() {
				log.Println("桶中当前水量", lr.water)
				log.Println("响应请求:", i)
			}
			wg.Done()
		}(i)

		//time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
}
