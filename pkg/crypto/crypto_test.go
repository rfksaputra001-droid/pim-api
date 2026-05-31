package crypto_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"pim-api-go/pkg/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	plaintext := "081234567890"
	encrypted, err := crypto.Encrypt(plaintext)
	assert.NoError(t, err)
	assert.Contains(t, encrypted, ":")

	decrypted, err := crypto.Decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestDecryptBadInput(t *testing.T) {
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	_, err := crypto.Decrypt("notvalidformat")
	assert.Error(t, err)
}
