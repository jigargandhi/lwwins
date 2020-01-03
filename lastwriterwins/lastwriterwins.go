package lastwriterwins

import (
	"fmt"

	"github.com/jigargandhi/lwwins/clock"
)

// Payload for the writer
type Payload struct {
	value int
	clock *clock.Loclock
}

// New Initializes a new payload
func New(clock *clock.Loclock, value int) *Payload {
	return &Payload{value: value, clock: clock}
}

// Update assigns new value to the payload
func (payload *Payload) Update(val int) {
	payload.value = val
}

// Value returns current register value
func (payload *Payload) Value() int {
	return payload.value
}

// Merge merges existing value with new value based on the logical clock
func (payload *Payload) Merge(remoteValue int, time uint64) {
	if payload.clock.Get() < time {
		payload.value = remoteValue
	} else {
		fmt.Println("Merge failed ")
	}
}
