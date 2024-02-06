package encrypt

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/zeromicro/go-zero/core/codec"
	"strings"
)

const (
	passwordEncryptSeed = "(beyond)@#$"
	mobileAesKey        = "5A2E746B08D846502F37A6E2D85D583B"
)

func EncPassword(pwd string) string {
	return Md5Sum([]byte(strings.TrimSpace(pwd + passwordEncryptSeed)))
}

func EncMobile(mobile string) (string, error) {
	data, err := codec.EcbEncrypt([]byte(mobileAesKey), []byte(mobile))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecMobile(mobile string) (string, error) {
	originalData, err := base64.StdEncoding.DecodeString(mobileAesKey)
	if err != nil {
		return "", err
	}
	data, err := codec.EcbDecrypt([]byte(mobileAesKey), originalData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	t := make([]byte, 16) // 使用容量为16的切片
	for i := range in {
		t[i] = in[i]
	}
	return t
}
