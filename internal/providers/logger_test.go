package providers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := NewLogger()
	require.NoError(t, err)
	require.NotNil(t, logger)
}
