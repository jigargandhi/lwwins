package lastwriterwins

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNotNilPayloadCreated(t *testing.T) {
	payload := New(10)
	assert.NotNil(t, payload)
}

func TestUpdate(t *testing.T) {
	payload := New(1)
	payload.Update(2)
	assert.Equal(t, 2, payload.value)
}

func TestValue(t *testing.T) {
	payload := New(1)
	val := payload.Value()
	assert.Equal(t, 1, val)
}

func TestMerge_WhenOldValue(t *testing.T) {
	payload := New(1)
	d, _ := time.ParseDuration("-1s")
	payload.time = time.Now().Add(d).UnixNano()
	payload.Merge(3, time.Now().UnixNano())
	assert.Equal(t, 3, payload.value)
}
