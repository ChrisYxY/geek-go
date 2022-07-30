package rolling_window_demo

import (
	"sync"
	"time"
)

type Window struct {
	timeGapSec int64
	bucketSize int
	buckets map[int64]*Bucket
	lock   *sync.RWMutex
}

func NewWindow(buckets int, timeGap int64) *Window  {
	return &Window{
		bucketSize: buckets,
		buckets: make(map[int64]*Bucket, buckets),
		timeGapSec: timeGap,
		lock:   &sync.RWMutex{},
	}
}

// Increment increments the number in current timeBucket.
func (w *Window) Increment(i int64) {
	if i == 0 {
		return
	}

	w.lock.Lock()
	defer w.lock.Unlock()

	w.getCurrentBucket().Add(i)
	if len(w.buckets) >= w.bucketSize {
		w.removeOldBuckets()
	}
}

// Sum sums the values over the buckets in the last time gap seconds.
func (w *Window) Sum(now time.Time) int64 {
	var sum int64

	w.lock.RLock()
	defer w.lock.RUnlock()

	for timestamp, bucket := range w.buckets {
		if timestamp >= now.Unix()-w.timeGapSec {
			sum += bucket.GetValue()
		}
	}
	return sum
}

// Max returns the maximum value seen in the last time gap seconds.
func (w *Window) Max(now time.Time) int64 {
	var max int64

	w.lock.RLock()
	defer w.lock.RUnlock()

	for timestamp, bucket := range w.buckets {
		if timestamp >= now.Unix()-w.timeGapSec {
			val := bucket.GetValue()
			if val > max {
				max = val
			}
		}
	}
	return max
}

// Avg returns the avg value seen in the last time gap seconds.
func (w *Window) Avg(now time.Time) float64 {
	return float64(w.Sum(now)) / float64(w.timeGapSec)
}

// getCurrentBucket returns the bucket in current timeBucket
func (w *Window) getCurrentBucket() *Bucket {
	now := time.Now().Unix()
	var bucket *Bucket
	var ok bool

	if bucket, ok = w.buckets[now]; !ok {
		bucket = NewBucket()
		w.buckets[now] = bucket
	}
	return bucket
}

func (w *Window) removeOldBuckets() {
	now := time.Now().Unix() - w.timeGapSec
	for timestamp := range w.buckets {
		if timestamp <= now {
			delete(w.buckets, timestamp)
		}
	}
}