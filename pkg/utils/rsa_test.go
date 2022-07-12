package utils

import (
	"testing"

	"github.com/foxdex/ftx-site/config"

	"github.com/stretchr/testify/assert"
)

func TestRsa(t *testing.T) {
	text := "123456"

	encryptTest, err := RsaEncrypt([]byte(text), config.GetPublicKey())
	assert.NoError(t, err)

	decryptText, err := RsaDecrypt(encryptTest, config.GetPrivateKey())
	assert.NoError(t, err)
	assert.Equal(t, text, string(decryptText))
}
