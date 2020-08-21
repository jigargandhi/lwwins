package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressInitialize(t *testing.T) {
	registrar := NewRegistrar("1", createKey())
	assert.NotNil(t, registrar, "Make should not return nil")
}

func TestAddressHandler(t *testing.T) {
	registrar := &Registrar{}
	registrar.address = make(map[string]string)
	registrar.NewAddress = make(chan string)
	registrar.token = createKey()
	go func() {
		<-registrar.NewAddress
	}()
	data := encrypt(registrar.token, "2|171.0.0.0")
	registrar.handleAddress(data)

	val, ok := registrar.address["2"]
	assert.True(t, ok, "registrar contains a new address")
	assert.Equal(t, "171.0.0.0", val, "expected 172.0.0.0 go %v ", val)
}

func TestAddressHandler_DecryptionFailed_AddressNotCreated(t *testing.T) {
	registrar := &Registrar{}
	registrar.address = make(map[string]string)
	registrar.NewAddress = make(chan string)
	registrar.token = createKey()
	go func() {
		<-registrar.NewAddress
	}()
	data := encrypt(createKey(), "2|171.0.0.0")
	registrar.handleAddress(data)

	_, ok := registrar.address["2"]
	assert.True(t, !ok, "registrar should not add address when token not same")
}
