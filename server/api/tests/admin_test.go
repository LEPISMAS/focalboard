package tests

import (
	"bytes"
	"net"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
)

func TestAdminEndpoints(t *testing.T) {
	t.Run("POST /admin/users/{username}/password resets password with unix connection", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().UpdateUserPassword("user1", gomock.Any()).Return(nil)

		body := `{"password": "newpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/admin/users/user1/password", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		// Simulate local Unix socket connection
		ctx := api.SetContextConn(req.Context(), &net.UnixConn{})
		req = req.WithContext(ctx)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /admin/users/{username}/password returns unauthorized for non-unix connection", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		body := `{"password": "newpassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/admin/users/user1/password", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		// No unix connection set in context

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), "not a local unix connection")
	})
}
