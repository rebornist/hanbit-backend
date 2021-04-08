package mixins

import (
	"encoding/base32"

	"github.com/rebornist/hanbit/config"
)

func Signing(s string) string {
	byteStr, _ := config.EncryptionAESKey(s)
	signedtext := base32.StdEncoding.EncodeToString(byteStr)
	return signedtext
}

func Unsigning(s string) string {
	unsignedByte, _ := base32.StdEncoding.DecodeString(s)
	byteStr, _ := config.DecryptionAESKey(string(unsignedByte))
	return string(byteStr)
}
