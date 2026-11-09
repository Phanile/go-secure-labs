package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256Hasher struct {
}

func (SHA256Hasher) Hash(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:]), nil
}

func (SHA256Hasher) Compare(password, hash string) (bool, error) {
	passwordHash, err := SHA256Hasher{}.Hash(password)
	if err != nil {
		return false, err
	}
	return passwordHash == hash, nil
}

type Argon2idHasher struct {
}

func (Argon2idHasher) Hash(password string) (string, error) {
	return "", nil
}

func (Argon2idHasher) Compare(password, hash string) (bool, error) {
	return false, nil
}
