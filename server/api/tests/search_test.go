package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestSearchEndpoints(t *testing.T) {
	t.Run("GET /teams/{teamID}/channels returns Not Implemented in standalone mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/channels?search=chan", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
		require.Contains(t, resp.Body.String(), "not permitted in standalone mode")
	})

	t.Run("GET /teams/{teamID}/channels returns channels in plugin mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		channels := []*mmModel.Channel{
			{Id: "chan1", DisplayName: "Channel 1"},
		}
		th.Store.EXPECT().SearchUserChannels("team1", "single-user", "chan").Return(channels, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/channels?search=chan", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Channel 1")
	})

	t.Run("GET /teams/{teamID}/boards/search returns matching boards", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boards := []*model.Board{
			{ID: "board1", Title: "Searchable Board"},
		}
		// SearchBoardsForUser parameters: term, searchField, userID, includePublicBoards
		th.Store.EXPECT().SearchBoardsForUser("search_term", model.BoardSearchFieldTitle, "single-user", true).Return(boards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/boards/search?q=search_term", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Searchable Board")
	})

	t.Run("GET /teams/{teamID}/boards/search/linkable returns linkable boards in plugin mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		boards := []*model.Board{
			{ID: "board1", Title: "Linkable Board"},
		}
		th.Store.EXPECT().SearchBoardsForUserInTeam("team1", "link_term", "single-user").Return(boards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/boards/search/linkable?q=link_term", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Linkable Board")
	})

	t.Run("GET /boards/search returns matching boards across teams", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boards := []*model.Board{
			{ID: "board2", Title: "Global Board"},
		}
		th.Store.EXPECT().SearchBoardsForUser("global_term", model.BoardSearchFieldTitle, "single-user", true).Return(boards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/search?q=global_term", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Global Board")
	})
}
