package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/app"
)

func TestSystemEndpoints(t *testing.T) {
	th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
	defer tearDown()

	t.Run("GET /hello returns 'Hello' and 200 OK", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
		resp := doRequest(th.Router, req)
		
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "Hello", resp.Body.String())
	})

	t.Run("GET /ping returns metadata and 200 OK", func(t *testing.T) {
		th.Store.EXPECT().DBType().Return("sqlite3")
		th.Store.EXPECT().DBVersion().Return("3.35.5")

		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		resp := doRequest(th.Router, req)

		require.Equal(t, http.StatusOK, resp.Code)

		var metadata app.ServerMetadata
		err := json.Unmarshal(resp.Body.Bytes(), &metadata)
		require.NoError(t, err)

		require.Equal(t, "sqlite3", metadata.DBType)
		require.Equal(t, "3.35.5", metadata.DBVersion)
		require.Equal(t, "personal_desktop", metadata.SKU) // Since TestSingleUserToken is configured on the API
	})
}
