package service

import "sync/atomic"

type Counter struct {
	count atomic.Uint64
}

func (c *Counter) Count() int64 {
	return int64(c.count.Load())
}

func (c *Counter) Increment() {
	c.count.Add(1)
}

func (c *Counter) Add(delta uint64) {
	c.count.Add(delta)
}
