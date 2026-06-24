package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestStatisticsEndpoints(t *testing.T) {
	t.Run("GET /statistics returns Not Implemented in standalone mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/statistics", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
		require.Contains(t, resp.Body.String(), "not permitted in standalone mode")
	})

	t.Run("GET /statistics returns permission denied if permission missing", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true
		th.Permissions.HasPermissionToFunc = func(userID string, permission *mmModel.Permission) bool {
			return false
		}

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/statistics", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusForbidden, resp.Code)
		require.Contains(t, resp.Body.String(), "access denied System Statistics")
	})

	t.Run("GET /statistics returns server statistics", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		th.Store.EXPECT().GetBoardCount().Return(int64(10), nil)
		th.Store.EXPECT().GetUsedCardsCount().Return(int(50), nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/statistics", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), `"board_count":10`)
		require.Contains(t, resp.Body.String(), `"card_count":50`)
	})
}
