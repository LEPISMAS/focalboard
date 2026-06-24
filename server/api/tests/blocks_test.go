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

func TestBlocksEndpoints(t *testing.T) {
	t.Run("GET /boards/{boardID}/blocks returns blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen}, nil)
		
		blocks := []*model.Block{
			{ID: "block1", BoardID: "board1", Type: model.TypeCard, Title: "Block 1"},
		}
		th.Store.EXPECT().GetBlocksWithParent("board1", "").Return(blocks, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/blocks", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "block1")
	})

	t.Run("POST /boards/{boardID}/blocks inserts blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)
		th.Store.EXPECT().InsertBlock(gomock.Any(), "single-user").Return(nil)

		body := `[{"id": "block1", "boardId": "board1", "type": "card", "createAt": 100, "updateAt": 100}]`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/blocks", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("DELETE /boards/{boardID}/blocks/{blockID} deletes block", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil).Times(2)
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)
		th.Store.EXPECT().DeleteBlock("block1", "single-user").Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/boards/board1/blocks/block1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /boards/{boardID}/blocks/{blockID}/undelete undeletes block", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil).Times(2)
		
		history := []*model.Block{
			{ID: "block1", BoardID: "board1"},
		}
		th.Store.EXPECT().GetBlockHistory("block1", model.QueryBlockHistoryOptions{Limit: 1, Descending: true}).Return(history, nil).Times(2)
		th.Store.EXPECT().UndeleteBlock("block1", "single-user").Return(nil)
		th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/blocks/block1/undelete", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("PATCH /boards/{boardID}/blocks/{blockID} patches block", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil).Times(3)
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)
		th.Store.EXPECT().PatchBlock("block1", gomock.Any(), "single-user").Return(nil)

		body := `{"title": "Patched Title"}`
		req, _ := http.NewRequest(http.MethodPatch, "/api/v2/boards/board1/blocks/block1", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("PATCH /boards/{boardID}/blocks patches batch of blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil).Times(2)
		th.Store.EXPECT().GetBlocksByIDs([]string{"block1"}).Return([]*model.Block{{ID: "block1", BoardID: "board1"}}, nil)
		th.Store.EXPECT().PatchBlocks(gomock.Any(), "single-user").Return(nil)

		body := `{"block_ids": ["block1"], "block_patches": [{"title": "Batch Patched"}]}`
		req, _ := http.NewRequest(http.MethodPatch, "/api/v2/boards/board1/blocks", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("POST /boards/{boardID}/blocks/{blockID}/duplicate duplicates block", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil).Times(3)
		th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil)
		
		duplicatedBlocks := []*model.Block{
			{ID: "dup1", BoardID: "board1", Type: model.TypeCard},
		}
		th.Store.EXPECT().DuplicateBlock("board1", "block1", "single-user", false).Return(duplicatedBlocks, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/blocks/block1/duplicate", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "dup1")
	})
}
