package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"strings"
	"time"
)

type GoogleAuth struct{}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

// 生成密钥
func (auth *GoogleAuth) GenerateSecret() string {
	// 生成随机字节数组作为密钥
	secret, _ := RandByte(20)

	// 将密钥进行base32编码
	encoded := base32.StdEncoding.EncodeToString(secret)
	// 移除编码后的末尾填充字符=
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}

// 验证动态码
func (auth *GoogleAuth) VerifyOTP(secretKey, otp string, windowSize int) bool {
	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 计算时间戳的时间步数
	counter := uint64(math.Floor(float64(timestamp) / 30))

	// 解码密钥
	secret, _ := base32.StdEncoding.DecodeString(secretKey)

	// 遍历时间步数窗口
	for i := -windowSize; i <= windowSize; i++ {
		// 计算当前时间步数
		currentCounter := counter + uint64(i)

		// 使用HMAC-SHA1算法计算哈希值
		hash := hmac.New(sha1.New, secret)
		counterBytes := make([]byte, 8)
		for j := 7; j >= 0; j-- {
			counterBytes[j] = byte(currentCounter & 0xff)
			currentCounter = currentCounter >> 8
		}
		hash.Write(counterBytes)
		hmacHash := hash.Sum(nil)

		// 获取动态码的起始位置
		offset := int(hmacHash[len(hmacHash)-1] & 0x0f)

		// 从哈希值中获取动态码
		binary := int(hmacHash[offset]&0x7f)<<24 |
			int(hmacHash[offset+1])<<16 |
			int(hmacHash[offset+2])<<8 |
			int(hmacHash[offset+3])

		// 格式化动态码
		currentOTP := fmt.Sprintf("%06d", binary%1000000)

		// 如果验证通过，返回true
		if currentOTP == otp {
			return true
		}
	}
	return false
}

func (auth *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}
