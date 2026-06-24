package tests

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestCSRFAndAuthMiddleware(t *testing.T) {
	th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
	defer tearDown()

	t.Run("Missing CSRF header returns 400 Bad Request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/me", nil)
		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), "checkCSRFToken FAILED")
	})

	t.Run("Invalid single user token returns 401 Unauthorized", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/me", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), "invalid single user token")
	})

	t.Run("Valid single user token passes authentication", func(t *testing.T) {
		// Mock GetTeam("0") since GetRootTeam calls GetTeam("0")
		th.Store.EXPECT().GetTeam("0").Return(&model.Team{ID: "0", SignupToken: "abc", UpdateAt: 123456}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/me", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)
		resp := doRequest(th.Router, req)
		
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "single-user")
	})
}

func TestPanicHandler(t *testing.T) {
	th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
	defer tearDown()

	t.Run("Recover from handler panic and return 500 Internal Server Error", func(t *testing.T) {
		// Mock GetBoardsForUserAndTeam to panic
		th.Store.EXPECT().
			GetBoardsForUserAndTeam(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(userID, teamID string, includeTemplates bool) (interface{}, error) {
				panic("simulated database panic")
			})

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/boards", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusInternalServerError, resp.Code)
		require.Contains(t, resp.Body.String(), "internal server error")
	})
}
