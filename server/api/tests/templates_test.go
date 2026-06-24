package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestTemplatesEndpoints(t *testing.T) {
	t.Run("GET /teams/{teamID}/templates returns template boards", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boards := []*model.Board{
			{ID: "board1", Title: "Template 1", Type: model.BoardTypeOpen, IsTemplate: true},
		}
		th.Store.EXPECT().GetTemplateBoards("team1", "single-user").Return(boards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/templates", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Template 1")
	})
}
