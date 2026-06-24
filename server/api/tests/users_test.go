package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestUsersEndpoints(t *testing.T) {
	th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
	defer tearDown()

	t.Run("POST /users list of users (single user)", func(t *testing.T) {
		th.Store.EXPECT().GetTeam("0").Return(&model.Team{ID: "0"}, nil)
		// Mock GetUserByID to check if current user is guest
		th.Store.EXPECT().GetUserByID("single-user").Return(&model.User{ID: "single-user", IsGuest: false}, nil)

		body := `["single-user"]`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/users", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "single-user")
	})

	t.Run("GET /users/me/memberships", func(t *testing.T) {
		members := []*model.BoardMember{
			{BoardID: "board1", UserID: "single-user", SchemeAdmin: true},
		}
		th.Store.EXPECT().GetMembersForUser("single-user").Return(members, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/me/memberships", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("GET /users/{userID}", func(t *testing.T) {
		th.Store.EXPECT().GetUserByID("single-user").Return(&model.User{ID: "single-user", IsGuest: false}, nil).Times(2) // Once for CanSeeUser, once for GetUser

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/single-user", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "single-user")
	})

	t.Run("PUT /users/{userID}/config", func(t *testing.T) {
		prefs := []mmModel.Preference{
			{UserId: "single-user", Category: "theme", Name: "dark", Value: "true"},
		}
		th.Store.EXPECT().PatchUserPreferences("single-user", gomock.Any()).Return(prefs, nil)

		body := `{"updatedFields":{"theme":"dark"}}`
		req, _ := http.NewRequest(http.MethodPut, "/api/v2/users/single-user/config", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "theme")
	})

	t.Run("GET /users/me/config", func(t *testing.T) {
		prefs := []mmModel.Preference{
			{UserId: "single-user", Category: "theme", Name: "dark", Value: "true"},
		}
		th.Store.EXPECT().GetUserPreferences("single-user").Return(prefs, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/users/me/config", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "theme")
	})
}
