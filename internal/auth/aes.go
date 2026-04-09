package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// AESEncrypt encrypts plaintext using AES/CBC/PKCS5Padding per the open platform spec.
// Key = MD5(openKey) → 32 hex chars used as bytes → AES-256.
// IV  = MD5(openId)[8:24] → 16 hex chars used as bytes.
func AESEncrypt(plaintext, openId, openKey string) (string, error) {
	key := deriveAESKey(openKey)
	iv := deriveAESIV(openId)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}

	padded := pkcs5Pad([]byte(plaintext), block.BlockSize())
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt decrypts base64-encoded ciphertext using AES/CBC/PKCS5Padding.
func AESDecrypt(ciphertextB64, openId, openKey string) (string, error) {
	key := deriveAESKey(openKey)
	iv := deriveAESIV(openId)

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	unpadded, err := pkcs5Unpad(plaintext)
	if err != nil {
		return "", err
	}
	return string(unpadded), nil
}

// deriveAESKey returns 32 bytes from MD5(openKey) hex string.
// MD5 produces 32 hex chars; all 32 chars used as bytes → AES-256.
func deriveAESKey(openKey string) []byte {
	hash := md5.Sum([]byte(openKey))
	hexStr := hex.EncodeToString(hash[:])
	return []byte(hexStr)
}

// deriveAESIV returns 16 bytes from MD5(openId)[8:24] hex chars.
// MD5 produces 32 hex chars; we take chars [8:24] = 16 hex chars as bytes.
func deriveAESIV(openId string) []byte {
	hash := md5.Sum([]byte(openId))
	hexStr := hex.EncodeToString(hash[:])
	return []byte(hexStr[8:24])
}

func pkcs5Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func pkcs5Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}
