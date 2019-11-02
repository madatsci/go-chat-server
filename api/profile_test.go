package api

import (
	"encoding/json"
	"github.com/madatsci/go-chat-server/internal/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestProfile(t *testing.T) {
	suite := newTestSuite(t, http.MethodGet, nil, nil)
	suite.authorize()
	defer suite.close()

	err := suite.api.Profile(suite.context)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, suite.recorder.Code)
	checkUser := new(models.User)
	decodeErr := json.NewDecoder(suite.recorder.Body).Decode(checkUser)
	require.NoError(t, decodeErr)
	require.Equal(t, suite.user, checkUser)
}
