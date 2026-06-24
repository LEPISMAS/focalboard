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

func TestSharingEndpoints(t *testing.T) {
	t.Run("GET /boards/{boardID}/sharing returns sharing info", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boardID := "board1"
		sharing := &model.Sharing{
			ID:      boardID,
			Enabled: true,
		}

		th.Store.EXPECT().GetSharing(boardID).Return(sharing, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/sharing", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), `"enabled":true`)
	})

	t.Run("POST /boards/{boardID}/sharing updates sharing info", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().UpsertSharing(gomock.Any()).Return(nil)

		body := `{"enabled": true}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/sharing", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})
}
