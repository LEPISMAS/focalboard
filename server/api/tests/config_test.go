package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
)

func TestConfigEndpoints(t *testing.T) {
	t.Run("GET /clientConfig returns config", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/clientConfig", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "enablePublicSharedBoards")
	})
}
