package clock

import (
	"sync/atomic"
)

// Loclock represents a lamports logical clock in a distributed system
type Loclock struct {
	counter uint64
}

// A little bit inspiration from https://github.com/hashicorp/serf/blob/master/serf/lamport.go

// Tick increments the logical clock
func (loclock *Loclock) Tick() {
	atomic.AddUint64(&loclock.counter, 1)
}

// Update the clock to a new value
// drawback of using accepting and returning uint64 is that callers may pass any value or tend to do arithmetic on the value
func (loclock *Loclock) Update(newTime uint64) {
	cur := atomic.LoadUint64(&loclock.counter)

	if newTime < loclock.counter {
		return
	}

	// Why should compare and swap fail
	atomic.CompareAndSwapUint64(&loclock.counter, cur, newTime+1)
}

// Get provides an instance of time;
func (loclock *Loclock) Get() uint64 {
	return atomic.LoadUint64(&loclock.counter)
}
