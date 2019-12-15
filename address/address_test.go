package address

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddressInitialize(t *testing.T) {
	registrar := Make(1, "token")
	assert.NotNil(t, registrar, "Make should not return nil")
}
