package crypto

import (
	"regexp"
	"testing"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"

func Test_AESEncryptDecryptMessage(t *testing.T) {
	key := []byte("0123456789abcdef")
	ciphertext, err := AESEncryptMessage(key, "test message")

	if err != nil {
		t.Errorf("AESEncryptMessage(%s) failed: %s", key, err)
	}

	if match, err := regexp.Match(isBase64, []byte(ciphertext)); err != nil || !match {
		t.Errorf("AESEncryptMessage(%s) return a non base64 cipher", ciphertext)
	}

	plaintext, _ := AESDecryptMessage(key, ciphertext)
	if plaintext != "test message" {
		t.Errorf("AESEncryptMessage() plaintext = %v, want %v", plaintext, "test message")
	}
}
