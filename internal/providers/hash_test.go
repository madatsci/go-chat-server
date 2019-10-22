package providers

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcryptHashHash(t *testing.T) {
	v := viper.New()
	v.Set("hash.complexity", 2)

	hasher := NewBcryptHasher(v)

	tests := []string{"12345", "test", "new_password"}

	for _, test := range tests {
		hash, err := hasher.Hash(test)
		require.NoError(t, err)
		require.True(t, hasher.Compare(test, hash))
	}
}
