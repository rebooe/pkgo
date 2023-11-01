package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSAEncrypt RSA 加密
func RSAEncrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// RSADecrypt RSA 解密
func RSADecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// 解析密钥结构
type ParseRsaKey struct {
	key []byte
	err error
}

func (r *ParseRsaKey) DecodePEM(key []byte) *ParseRsaKey {
	block, _ := pem.Decode(key)
	if block == nil {
		r.err = errors.New("failed to decode PEM block containing public key")
		return r
	}
	r.key = block.Bytes
	return r
}

func (r *ParseRsaKey) DecodeBase64(key string) *ParseRsaKey {
	decodeBytes, err := base64.StdEncoding.DecodeString(key)
	r.err = err
	r.key = decodeBytes
	return r
}

func (r *ParseRsaKey) ToPublicKey() (*rsa.PublicKey, error) {
	if r.err != nil {
		return nil, r.err
	}

	pubKey, err := x509.ParsePKIXPublicKey(r.key)
	if err != nil {
		return nil, err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("provided key is not an RSA public key")
	}
	return rsaPubKey, nil
}

func (r *ParseRsaKey) ToPrivateKey() (*rsa.PrivateKey, error) {
	if r.err != nil {
		return nil, r.err
	}

	privKey, err := x509.ParsePKCS8PrivateKey(r.key)
	if err != nil {
		return nil, err
	}
	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("provided key is not an RSA private key")
	}
	return rsaPrivKey, nil
}
