package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AesEncrypt 加密数据
func AesEncrypt(plainData []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 加密数据块的长度必须为 block 的块大小，因此需要填充数据
	paddedData := pkcs7Padding(plainData, block.BlockSize())

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	// 将 iv 和 ciphertext 合并成一个字节数组，方便后续处理
	encryptedData := append(iv, ciphertext...)

	// 对加密后的数据进行 base64 编码，方便传输
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// AesDecrypt 解密数据
func AesDecrypt(encryptedData string, key []byte) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedBytes) < aes.BlockSize {
		return nil, errors.New("encrypted data too short")
	}

	iv := encryptedBytes[:aes.BlockSize]
	ciphertext := encryptedBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	paddedData := make([]byte, len(ciphertext))
	mode.CryptBlocks(paddedData, ciphertext)

	// 去除填充数据
	plainData, err := pkcs7Unpadding(paddedData)
	if err != nil {
		return nil, err
	}

	return plainData, nil
}

// pkcs7Padding 填充数据
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7Unpadding 去除填充数据
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:(length - unpadding)], nil
}
