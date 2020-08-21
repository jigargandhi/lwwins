package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressInitialize(t *testing.T) {
	registrar := NewRegistrar(1, "token")
	assert.NotNil(t, registrar, "Make should not return nil")
}

func TestAddressHandler(t *testing.T) {
	registrar := &Registrar{}
	registrar.address = make(map[int]string)
	registrar.NewAddress = make(chan string)
	registrar.token = "token"
	go func() {
		<-registrar.NewAddress
	}()
	registrar.handleAddress("2|token|171.0.0.0")

	val, ok := registrar.address[2]
	assert.True(t, ok, "registrar contains a new address")
	assert.Equal(t, "171.0.0.0", val, "expected 172.0.0.0 go %v ", val)
}

func TestAddressHandler_TokenNotSame_AddressNotCreated(t *testing.T) {
	registrar := &Registrar{}
	registrar.address = make(map[int]string)
	registrar.NewAddress = make(chan string)
	registrar.token = "token1"
	go func() {
		<-registrar.NewAddress
	}()
	registrar.handleAddress("2|token|171.0.0.0")

	_, ok := registrar.address[2]
	assert.True(t, !ok, "registrar should not add address when token not same")
}
