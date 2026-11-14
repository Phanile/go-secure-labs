package hasher

import "testing"

func TestArgon2idHasherHashAndCompare(t *testing.T) {
	hasher := Argon2idHasher{}
	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}

	if ok, err := hasher.Compare("correct horse battery staple", hash); err != nil || !ok {
		t.Fatalf("Compare(correct password) = %v, %v", ok, err)
	}
	if ok, err := hasher.Compare("wrong password", hash); err != nil || ok {
		t.Fatalf("Compare(wrong password) = %v, %v", ok, err)
	}
}

func TestArgon2idHasherRejectsMalformedHash(t *testing.T) {
	if ok, err := (Argon2idHasher{}).Compare("password", "not-a-hash"); err == nil || ok {
		t.Fatalf("Compare(malformed hash) = %v, %v; want an error", ok, err)
	}
}
