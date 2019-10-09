package providers

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	// Hasher describes methods for hashing and validating password
	Hasher interface {
		Hash(password string) (string, error)
		Compare(password, hash string) bool
	}

	// BcryptHasher is a default bcrypt
	BcryptHasher struct {
		complexity int
	}
)

// NewBcryptHasher creates a new hasher that uses bcrypt under the hood
func NewBcryptHasher(config *Config) Hasher {
	return &BcryptHasher{
		complexity: config.HashComplexity,
	}
}

// Hash hashes the given password using bcrypt
func (b BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.complexity)
	return string(bytes), err
}

// Compare compares given password to given hash
func (b BcryptHasher) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
