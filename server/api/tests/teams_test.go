package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestTeamsEndpoints(t *testing.T) {
	t.Run("GET /teams returns teams for user", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		teams := []*model.Team{
			{ID: "team1", Title: "Team 1"},
		}
		th.Store.EXPECT().GetTeamsForUser("single-user").Return(teams, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Team 1")
	})

	t.Run("GET /teams/{teamID} returns root team when MattermostAuth is false", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		team := &model.Team{ID: "0", Title: "Root Team"}
		th.Store.EXPECT().GetTeam("0").Return(team, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Root Team")
	})

	t.Run("GET /teams/{teamID} returns specified team when MattermostAuth is true", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		team := &model.Team{ID: "team1", Title: "Team 1"}
		th.Store.EXPECT().GetTeam("team1").Return(team, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Team 1")
	})

	t.Run("POST /teams/{teamID}/users returns users by ID", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		users := []*model.User{
			{ID: "user1", Username: "user1"},
		}
		// GetUsersList parameters: userIDs, showEmail, showFullName
		th.Store.EXPECT().GetUsersList([]string{"user1"}, false, false).Return(users, nil)

		body := `["user1"]`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/users", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user1")
	})

	t.Run("GET /teams/{teamID}/users searches team users", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		users := []*model.User{
			{ID: "user1", Username: "user1"},
		}
		// SearchUsersByTeam parameters: teamID, searchQuery, asGuestID, excludeBots, showEmail, showFullName
		th.Store.EXPECT().SearchUsersByTeam("team1", "search_term", "", false, false, false).Return(users, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/users?search=search_term", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user1")
	})

	t.Run("POST /teams/{teamID}/regenerate_signup_token returns Not Implemented in plugin mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/regenerate_signup_token", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
		require.Contains(t, resp.Body.String(), "not permitted in plugin mode")
	})

	t.Run("POST /teams/{teamID}/regenerate_signup_token regenerates token in standalone mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		team := &model.Team{ID: "0", SignupToken: "old-token"}
		th.Store.EXPECT().GetTeam("0").Return(team, nil)
		th.Store.EXPECT().UpsertTeamSignupToken(gomock.Any()).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/regenerate_signup_token", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})
}
