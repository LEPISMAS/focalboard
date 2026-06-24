package tests

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestArchiveEndpoints(t *testing.T) {
	t.Run("GET /boards/{boardID}/archive/export exports board archive", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(3)
		th.Store.EXPECT().GetBlocksForBoard("board1").Return([]*model.Block{}, nil)
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/archive/export", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "application/octet-stream", resp.Header().Get("Content-Type"))
	})

	t.Run("POST /teams/{teamID}/archive/import imports board archive", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		// create a minimal valid zip archive in memory
		zipBuf := &bytes.Buffer{}
		zw := zip.NewWriter(zipBuf)
		w, err := zw.Create("version.json")
		require.NoError(t, err)
		_, err = w.Write([]byte(`{"version": 2, "date": 123456789}`))
		require.NoError(t, err)
		err = zw.Close()
		require.NoError(t, err)

		// Create a multipart form request
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "archive.boardarchive")
		require.NoError(t, err)
		_, err = part.Write(zipBuf.Bytes())
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/archive/import", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("GET /teams/{teamID}/archive/export exports team boards archive in standalone mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		boards := []*model.Board{
			{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen},
		}
		th.Store.EXPECT().GetTeamsForUser("single-user").Return([]*model.Team{{ID: "team1"}}, nil).AnyTimes()
		th.Store.EXPECT().GetBoardsForUserAndTeam("single-user", "team1", true).Return(boards, nil)
		th.Store.EXPECT().GetBoard("board1").Return(boards[0], nil).Times(2)
		th.Store.EXPECT().GetBlocksForBoard("board1").Return([]*model.Block{}, nil)
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/archive/export", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "application/octet-stream", resp.Header().Get("Content-Type"))
	})

	t.Run("GET /teams/{teamID}/archive/export returns Not Implemented in plugin mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/archive/export", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
	})
}
