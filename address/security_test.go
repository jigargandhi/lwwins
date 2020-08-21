package address

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createKey() string {
	key := make([]byte, 32)
	rand.Reader.Read(key)

	return base64.StdEncoding.EncodeToString(key)
}

func TestEncryptionDecryption_Valid(t *testing.T) {
	data := "nodeName|120.0.0.1"
	key := createKey()
	cipherText := encrypt(key, data)

	plaintext, err := decrypt(key, cipherText)

	assert.Equal(t, data, plaintext, "Plaintext and ciphertext should be equal")
	assert.Nil(t, err, "Error should be nil on successful decryption")
}

func TestUnsuccessfulDecryption(t *testing.T) {
	data := "nodeName|120.0.0.1"
	key := createKey()
	cipherText := encrypt(key, data)

	_, err := decrypt(createKey(), cipherText)

	assert.NotNil(t, err, "Error should not be nil on un successful decryption")
}
