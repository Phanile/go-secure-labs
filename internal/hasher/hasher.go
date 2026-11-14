package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 64 * 1024 // 64 МБ памяти
	argonIterations  uint32 = 3         // 3 прохода
	argonParallelism uint8  = 2         // 2 потока
	argonSaltLength  uint32 = 16        // 16 байт соли
	argonKeyLength   uint32 = 32        // 32 байта длина итогового ключа
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
	return subtle.ConstantTimeCompare([]byte(passwordHash), []byte(hash)) == 1, nil
}

type Argon2idHasher struct {
}

func (Argon2idHasher) Hash(password string) (string, error) {
	// 1. Генерируем криптографически стойкую случайную соль
	salt := make([]byte, argonSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Вычисляем хеш Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonKeyLength,
	)

	// 3. Кодируем соль и хеш в Base64 (без символов '=' в конце)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 4. Формируем стандартную строку, содержащую все параметры для будущей проверки
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func (Argon2idHasher) Compare(password, hash string) (bool, error) {
	memory, iterations, parallelism, salt, decodedHash, err := parseArgon2idHash(hash)
	if err != nil {
		return false, err
	}

	keyLength := uint32(len(decodedHash))

	// 2. Вычисляем хеш от введенного пароля с теми же параметрами и солью
	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	// 3. Сравниваем за константное время (защита от Timing Attack)
	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}

func parseArgon2idHash(encoded string) (uint32, uint32, uint8, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid argon2id hash format")
	}

	version, err := parseArgonParameter(parts[2], "v")
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}
	if version != uint64(argon2.Version) {
		return 0, 0, 0, nil, nil, fmt.Errorf("incompatible argon2 version: %d", version)
	}

	parameters := strings.Split(parts[3], ",")
	if len(parameters) != 3 {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid argon2 parameters")
	}
	memory, err := parseArgonParameter(parameters[0], "m")
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}
	iterations, err := parseArgonParameter(parameters[1], "t")
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	parallelism, err := parseArgonParameter(parameters[2], "p")
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	// Keep malformed or attacker-controlled parameter values from causing an
	// unbounded CPU/memory allocation during login.
	if memory < 8*uint64(parallelism) || memory > 1024*1024 || iterations == 0 || iterations > 10 || parallelism == 0 || parallelism > 16 {
		return 0, 0, 0, nil, nil, fmt.Errorf("argon2 parameters are outside safe bounds")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 8 || len(salt) > 64 {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid argon2 salt")
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(decodedHash) < 16 || len(decodedHash) > 64 {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid argon2 hash key")
	}

	return uint32(memory), uint32(iterations), uint8(parallelism), salt, decodedHash, nil
}

func parseArgonParameter(value, name string) (uint64, error) {
	prefix := name + "="
	if !strings.HasPrefix(value, prefix) {
		return 0, fmt.Errorf("invalid argon2 parameter %q", value)
	}
	parsed, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid argon2 parameter %q: %w", value, err)
	}
	return parsed, nil
}
