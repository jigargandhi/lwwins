package address

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddressInitialize(t *testing.T) {
	registrar := Make(1, "token")
	assert.NotNil(t, registrar, "Make should not return nil")
}

func TestAddressHandler_WhenGivenCorrectAddress_AddsAndNotifiesChannel(t *testing.T) {
	t.Skip("test skipped until I learn")
	registrar := Make(1, "token")
	handleAddress("2|token|171.0.0.0", registrar)
	val, ok := registrar.address[2]
	go func() {
		<-registrar.NewAddress
	}()
	assert.True(t, ok, "registrar contains a new address")
	assert.Equal(t, "171.0.0.0", val, "expected 172.0.0.0 go %v ", val)

}
