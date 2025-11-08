package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func padKey(key []byte) []byte {
	if len(key) >= 32 {
		return key[:32]
	}
	newKey := make([]byte, 32)
	copy(newKey, key)
	for i := len(key); i < 32; i++ {
		newKey[i] = byte(32 - i)
	}
	return newKey
}

// AESEncrypt encrypts plaintext using AES with the given key.
func AESEncrypt(plaintext string, key string) (string, error) {
	// Convert key and plaintext to byte slices
	keyBytes := padKey([]byte(key))
	plainBytes := []byte(plaintext)

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Generate a random IV (Initialization Vector)
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Use CFB mode for encryption
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainBytes)

	// Return the encrypted text as a hex-encoded string
	return hex.EncodeToString(ciphertext), nil
}

// AESDecrypt decrypts the ciphertext using AES with the given key.
func AESDecrypt(ciphertext string, key string) (string, error) {
	// Convert key and ciphertext to byte slices
	keyBytes := padKey([]byte(key))
	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Extract the IV from the ciphertext
	iv := cipherBytes[:aes.BlockSize]
	cipherBytes = cipherBytes[aes.BlockSize:]

	// Use CFB mode for decryption
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherBytes, cipherBytes)

	// Return the decrypted plaintext as a string
	return string(cipherBytes), nil
}

// ComparePassword .
func ComparePassword(uid int64, encodePW string, password string) bool {
	// 全部转小写
	elp := EncryptPassword(uid, password)
	if elp == encodePW {
		return true
	}

	// TODO 移除兼容区分大小写密码
	ep := EncryptPasswordRaw(uid, password)
	return ep == encodePW
}

// EncryptPassword .
// 密码要全部转小写
func EncryptPassword(uid int64, password string) string {
	lp := strings.ToLower(password)
	return EncryptPasswordRaw(uid, lp)
}

// EncryptPasswordRaw .
func EncryptPasswordRaw(uid int64, password string) string {
	p := strconv.FormatInt(uid, 10) + password + strconv.FormatInt(uid, 10)
	hash := md5.Sum([]byte(p))
	return hex.EncodeToString(hash[:])
}

// PKCS7 Padding Function
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := append(data, make([]byte, padding)...)
	for i := len(data); i < len(data)+padding; i++ {
		padText[i] = byte(padding)
	}
	return padText
}

// PKCS7 Unpadding Function
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("decryption error: input data is empty")
	}
	padding := int(data[length-1]) // Last byte represents the padding size
	if padding > length || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-padding], nil
}

func RijndaelEncrypt(plaintext string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Pad the plaintext using PKCS7
	plainBytes := pkcs7Padding([]byte(plaintext), aes.BlockSize)

	// Encrypt using ECB mode
	ciphertext := make([]byte, len(plainBytes))
	for i := 0; i < len(plainBytes); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], plainBytes[i:i+aes.BlockSize])
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RijndaelDecrypt 适配CS的RijndaelEncrypt，用于解密CS的RijndaelEncrypt加密的字符串
func RijndaelDecrypt(ciphertext string, key string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(encrypted)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	// Decrypt using ECB mode
	decrypted := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], encrypted[i:i+aes.BlockSize])
	}

	// Remove PKCS7 padding
	unpadded, err := pkcs7Unpadding(decrypted)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}
