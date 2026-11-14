package hasher

import (
	"crypto/sha256"
	"testing"

	"golang.org/x/crypto/argon2"
)

var result []byte

func BenchmarkSHA256(b *testing.B) {
	password := []byte("12345678")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash := sha256.Sum256(password)
		result = hash[:]
	}
}

func BenchmarkArgon2id(b *testing.B) {
	password := []byte("12345678")
	salt := []byte("0123456789abcdef")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
		result = hash
	}
}
