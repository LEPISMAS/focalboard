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

func TestBoardsAndBlocksEndpoints(t *testing.T) {
	t.Run("POST /boards-and-blocks creates boards and blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{
				{ID: "board1", TeamID: "team1", Title: "Board 1"},
			},
			Blocks: []*model.Block{
				{ID: "block1", BoardID: "board1", Type: "text", CreateAt: 10, UpdateAt: 10},
			},
		}

		th.Store.EXPECT().CreateBoardsAndBlocksWithAdmin(gomock.Any(), "single-user").Return(bab, []*model.BoardMember{}, nil)

		// Mock category board setup
		categories := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:     "cat1",
					Name:   "Boards",
					Type:   "system",
					UserID: "single-user",
					TeamID: "team1",
				},
				BoardMetadata: []model.CategoryBoardMetadata{},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categories, nil).AnyTimes()
		th.Store.EXPECT().AddUpdateCategoryBoard("single-user", "cat1", []string{"board1"}).Return(nil).AnyTimes()
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		body := `{"boards": [{"id": "board1", "teamId": "team1", "title": "Board 1"}], "blocks": [{"id": "block1", "boardId": "board1", "type": "text", "createAt": 10, "updateAt": 10}]}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards-and-blocks", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("PATCH /boards-and-blocks patches boards and blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		block := &model.Block{ID: "block1", BoardID: "board1"}

		th.Store.EXPECT().GetBoard("board1").Return(board, nil)
		th.Store.EXPECT().GetBlock("block1").Return(block, nil)
		th.Store.EXPECT().GetBlocksByIDs([]string{"block1"}).Return([]*model.Block{block}, nil)
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{board},
			Blocks: []*model.Block{block},
		}
		th.Store.EXPECT().PatchBoardsAndBlocks(gomock.Any(), "single-user").Return(bab, nil)

		body := `{"boardIDs": ["board1"], "boardPatches": [{}], "blockIDs": ["block1"], "blockPatches": [{}]}`
		req, _ := http.NewRequest(http.MethodPatch, "/api/v2/boards-and-blocks", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("DELETE /boards-and-blocks deletes boards and blocks", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		block := &model.Block{ID: "block1", BoardID: "board1"}

		th.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()
		th.Store.EXPECT().GetBlock("block1").Return(block, nil).AnyTimes()
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().DeleteBoardsAndBlocks(gomock.Any(), "single-user").Return(nil)

		body := `{"boards": ["board1"], "blocks": ["block1"]}`
		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/boards-and-blocks", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})
}
