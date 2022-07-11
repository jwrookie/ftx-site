package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var (
	defaultKey = []byte("fSDaJKG1R1AGBsdjfklf7942SDFfKJFg")
)

// PKCS7Padding 以块长度的整数倍填充明文
func _PKCS7Padding(p []byte, blockSize int) []byte {
	pad := blockSize - len(p)%blockSize
	padtext := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(p, padtext...)
}

// PKCS7UnPadding 从明文的尾部删除填充数据
func _PKCS7UnPadding(p []byte) []byte {
	length := len(p)
	paddLen := int(p[length-1])
	return p[:(length - paddLen)]
}

// AESCBCEncrypt 在CBC模式下使用AES算法加密数据
// 注意，要选择AES-128、AES-192或AES-256，密钥长度必须为16、24或32字节
// 注意，AES块大小为16字节
func _AESCBCEncrypt(p, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	p = _PKCS7Padding(p, block.BlockSize())
	ciphertext := make([]byte, len(p))
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(ciphertext, p)
	return ciphertext, nil
}

// AESCBCDecrypt ：在CBC模式下用AES算法解密密码文本
// 注意，要选择AES-128、AES-192或AES-256，密钥长度必须为16、24或32字节
// 注意，AES块大小为16字节
func _AESCBCDecrypt(c, key []byte) (p []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("err: %v", e)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(c))
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(plaintext, c)
	return _PKCS7UnPadding(plaintext), nil
}

// Base64AESCBCEncrypt 使用CBC模式的AES算法加密数据，并使用base64进行编码
func Base64AESCBCEncrypt(p string) (string, error) {
	c, err := _AESCBCEncrypt([]byte(p), defaultKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(c), nil
}

// Base64AESCBCDecrypt 用CBC模式下的AES算法对base64编码的密码文本进行解密
func Base64AESCBCDecrypt(c string) (string, error) {
	oriCipher, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return "", err
	}
	p, err := _AESCBCDecrypt(oriCipher, defaultKey)
	if err != nil {
		return "", err
	}
	return string(p), nil
}
