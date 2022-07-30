package rolling_window_demo

import "sync/atomic"

type Bucket struct {
	value int64
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (b *Bucket) Add(delta int64) {
	atomic.AddInt64(&b.value, delta)
}

func (b *Bucket) GetValue() int64 {
	return atomic.LoadInt64(&b.value)
}