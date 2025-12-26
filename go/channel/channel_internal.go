package main

import (
	"fmt"
	"sync"
)

type waiter struct {
	ch chan struct{} // 只是用来阻塞 没有传数据
	v  interface{}
}

type Chan struct {
	buf  []interface{} // 底层用的是ring buffer
	size int

	sendq []waiter // 等待发送
	recvq []waiter // 等待接收

	mu sync.Mutex
}

func NewChan(size int) *Chan {
	return &Chan{
		buf:  make([]interface{}, 0, size),
		size: size,
	}
}

func (c *Chan) Send(v interface{}) {
	c.mu.Lock()

	// 有 recv 在等，直接交接数据给recv
	if len(c.recvq) > 0 {
		w := c.recvq[0]
		c.recvq = c.recvq[1:]
		w.v = v     // 在这里传输数据
		close(w.ch) // 唤醒接收者
		c.mu.Unlock()
		return
	}

	// 有缓冲且没满 则放入缓冲
	if len(c.buf) < c.size {
		c.buf = append(c.buf, v)
		c.mu.Unlock()
		return
	}

	// 否则阻塞等待
	w := waiter{ch: make(chan struct{}), v: v}
	c.sendq = append(c.sendq, w)

	c.mu.Unlock()

	<-w.ch // 阻塞
}

func (c *Chan) Recv() interface{} {
	c.mu.Lock()

	// 先拿缓冲里的数据 拿出来之后
	if len(c.buf) > 0 {
		v := c.buf[0]
		c.buf = c.buf[1:]

		// 有send 在等 然后把等待的数据推进缓冲区
		if len(c.sendq) > 0 {
			w := c.sendq[0]
			c.sendq = c.sendq[1:]

			c.buf = append(c.buf, w.v)
			close(w.ch)
		}

		c.mu.Unlock()
		return v
	}

	// buffer 空，又有 sender 阻塞 直接拿数据
	if len(c.sendq) > 0 {
		w := c.sendq[0]
		c.sendq = c.sendq[1:]
		v := w.v
		close(w.ch)
		c.mu.Unlock()
		return v
	}

	// 否则阻塞
	w := waiter{ch: make(chan struct{})}
	c.recvq = append(c.recvq, w)

	c.mu.Unlock()

	<-w.ch // 阻塞
	return w.v
}

func main() {
	ch := NewChan(1)

	go func() {
		for i := 0; i < 5; i++ {
			ch.Send(i)
			fmt.Println("send:", i)
		}
	}()

	for i := 0; i < 5; i++ {
		v := ch.Recv()
		fmt.Println("recv:", v)
	}
}
