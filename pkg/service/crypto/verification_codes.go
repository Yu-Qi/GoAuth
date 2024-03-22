package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	separator = ":"
)

var (
	service *VerificationCodeService
)

// Init initializes the verification code service
func Init(secretKey, salt string, iters, expireSec int) {
	service = &VerificationCodeService{
		secretKey: secretKey,
		salt:      salt,
		iters:     iters,
		ExpireSec: expireSec,
	}
}

// GetService returns the verification code service
func GetService() *VerificationCodeService {
	return service
}

// VerificationCodeService provides the service to generate and verify verification code
type VerificationCodeService struct {
	secretKey string
	salt      string
	iters     int
	ExpireSec int
}

// GenerateCode generates a verification code with uid and timestamp
func (v *VerificationCodeService) GenerateCode(uid string) (string, error) {
	timestamp := time.Now().Add(time.Duration(v.ExpireSec) * time.Second).Unix()
	data := fmt.Sprintf("%s%s%d", uid, separator, timestamp)
	return v.Encrypt(data)
}

// VerifyCode verifies the code, checks the timestamp and returns the uid
func (v *VerificationCodeService) VerifyCode(code string) (string, error) {
	data, err := v.Decrypt(code)
	if err != nil {
		return "", err
	}

	parts := strings.Split(data, separator)
	if len(parts) != 2 {
		return "", fmt.Errorf("Decrypted data format error")
	}

	timeStamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", err
	}
	if time.Now().Unix() > timeStamp {
		return "", fmt.Errorf("Verification code expired")
	}

	return parts[0], nil
}

// Encrypt encrypts data using AES
func (v *VerificationCodeService) Encrypt(plainText string) (string, error) {
	key := pbkdf2.Key([]byte(v.secretKey), []byte(v.salt), v.iters, 32, sha3.New512)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts data using AES
func (v *VerificationCodeService) Decrypt(cipherText string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(v.secretKey), []byte(v.salt), v.iters, 32, sha3.New512)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("cipherText too short")
	}

	nonce, cipherTextData := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherTextData, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
