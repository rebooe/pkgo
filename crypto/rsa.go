package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

func ParsePublicKey(key string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	// block, _ := pem.Decode(key)
	// if block == nil {
	// 	return nil, errors.New("failed to decode PEM block containing public key")
	// }
	pubKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("provided key is not an RSA public key")
	}
	return rsaPubKey, nil
}

func ParsePrivateKey(key string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	// block, _ := pem.Decode(keyBytes)
	// if block == nil {
	// 	return nil, errors.New("failed to decode PEM block containing private key")
	// }
	privKey, err := x509.ParsePKCS8PrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("provided key is not an RSA private key")
	}
	return rsaPrivKey, nil
}

// RSAEncrypt RSA 加密
func RSAEncrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// RSADecrypt RSA 解密
func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}
