package lastwriterwins

import (
	"fmt"
	"time"
)

// Payload for the writer
type Payload struct {
	value int
	time  int64
}

func (payload *Payload) Update(val int) {
	payload.value = val
	payload.time = time.Now().UnixNano()
}

func (payload *Payload) Value() int {
	return payload.value
}

func (payload *Payload) Merge(remoteValue int, time int64) {
	if payload.time < time {
		payload.value = remoteValue
	} else {
		fmt.Println("Merge failed ")
	}

}
