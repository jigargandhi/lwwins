package lastwriterwins

import (
	"testing"

	"github.com/jigargandhi/lwwins/clock"
	"github.com/stretchr/testify/assert"
)

func TestNotNilPayloadCreated(t *testing.T) {
	var loClock clock.Loclock
	payload := New(&loClock, 10)
	assert.NotNil(t, payload)
}

func TestUpdate(t *testing.T) {
	var loClock clock.Loclock
	payload := New(&loClock, 1)
	payload.Update(2)
	assert.Equal(t, 2, payload.value)
}

func TestValue(t *testing.T) {
	var loClock clock.Loclock
	payload := New(&loClock, 1)
	val := payload.Value()
	assert.Equal(t, 1, val)
}

func TestMerge_WhenOldValue(t *testing.T) {
	var loClock clock.Loclock
	payload := New(&loClock, 1)
	loClock.Tick()
	loClock.Tick()
	payload.Merge(3, 4)
	assert.Equal(t, 3, payload.value)
}
