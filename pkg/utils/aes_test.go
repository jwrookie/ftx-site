package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	text := "hello world!hello world!hello world!hello world!hello world!"

	encryptText, err := Base64AESCBCEncrypt(text)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptText)
	t.Log(encryptText)

	decryptText, err := Base64AESCBCDecrypt(encryptText)
	assert.NoError(t, err)
	assert.NotEmpty(t, decryptText)
	assert.Equal(t, text, decryptText)
}
