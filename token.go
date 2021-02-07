package main

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate         int64 //固定的token放入速率, r/s
	capacity     int64 //桶的容量
	tokens       int64 //桶中当前token数量
	lastTokenSec int64 //桶上次放token的时间戳 s
	lock         sync.Mutex
}

func (l *TokenBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().Unix()
	l.tokens = l.tokens + (now-l.lastTokenSec)*l.rate

	if l.tokens > l.capacity {
		l.tokens = l.capacity
	}
	l.lastTokenSec = now
	if l.tokens > 0 {
		// 还有令牌，领取令牌
		l.tokens--
		return true
	} else {
		return false
	}
}

func (l *TokenBucket) Set(r, c int64) {
	l.rate = r
	l.capacity = c
	l.tokens = 1
	l.lastTokenSec = time.Now().Unix()
}

//令牌桶
func main() {

}
