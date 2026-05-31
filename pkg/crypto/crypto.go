package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func getKey() ([]byte, error) {
	hexKey := os.Getenv("ENCRYPTION_KEY")
	if len(hexKey) != 64 {
		return nil, fmt.Errorf("ENCRYPTION_KEY harus 64 karakter hex")
	}
	return hex.DecodeString(hexKey)
}

func Encrypt(plaintext string) (string, error) {
	key, err := getKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	encrypted := make([]byte, len(padded))
	mode.CryptBlocks(encrypted, padded)

	return hex.EncodeToString(iv) + ":" + hex.EncodeToString(encrypted), nil
}

func Decrypt(ciphertext string) (string, error) {
	parts := strings.SplitN(ciphertext, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("format ciphertext tidak valid")
	}

	key, err := getKey()
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	return string(pkcs7Unpad(decrypted)), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	padding := int(data[length-1])
	return data[:length-padding]
}
