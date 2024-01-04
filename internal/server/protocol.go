package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var KEY = []byte("AES256Key-32Characters1234567890")

func Encrypt(data []byte) ([]byte, error) {
	// 创建 AES 分组
	block, err := aes.NewCipher(KEY)
	if err != nil {
		return nil, err
	}

	// 创建一个加密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成一个随机的 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 使用 AES-GCM 进行加密
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	// 创建 AES 分组
	block, err := aes.NewCipher(KEY)
	if err != nil {
		return nil, err
	}

	// 创建一个解密器
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 解密数据
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
