package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/require"

	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestGetBoards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		boards := []*model.Board{
			{
				ID:     "board1",
				TeamID: "team1",
			},
		}

		helper.Store.EXPECT().GetBoardsForUserAndTeam("single-user", "team1", true).Return(boards, nil)

		req := helper.NewRequest("GET", "/api/v2/teams/team1/boards", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp []*model.Board
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp, 1)
		require.Equal(t, "board1", resp[0].ID)
	})

	t.Run("permission denied", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		helper.PermissionsService.hasPermissionToTeam = func(userID, teamID string, permission *mmModel.Permission) bool {
			return false
		}

		req := helper.NewRequest("GET", "/api/v2/teams/team1/boards", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestCreateBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		newBoard := &model.Board{
			TeamID: "team1",
			Title:  "New Board",
			Type:   model.BoardTypeOpen,
		}
		body, _ := json.Marshal(newBoard)

		createdBoard := &model.Board{
			ID:     "board1",
			TeamID: "team1",
			Title:  "New Board",
			Type:   model.BoardTypeOpen,
		}
		member := &model.BoardMember{
			BoardID: "board1",
			UserID:  "single-user",
		}

		helper.Store.EXPECT().InsertBoardWithAdmin(gomock.Any(), gomock.Eq("single-user")).Return(createdBoard, member, nil)

		req := helper.NewRequest("POST", "/api/v2/boards", body)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp model.Board
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "board1", resp.ID)
	})

	t.Run("invalid json", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		req := helper.NewRequest("POST", "/api/v2/boards", []byte("invalid json"))
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		board := &model.Board{
			ID:     "board1",
			TeamID: "team1",
		}

		helper.Store.EXPECT().GetBoard("board1").Return(board, nil)

		req := helper.NewRequest("GET", "/api/v2/boards/board1", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp model.Board
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "board1", resp.ID)
	})

	t.Run("not found", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		helper.Store.EXPECT().GetBoard("board1").Return(nil, model.NewErrNotFound("board"))

		req := helper.NewRequest("GET", "/api/v2/boards/board1", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPatchBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		patch := &model.BoardPatch{
			Title: &[]string{"Updated Title"}[0],
		}
		body, _ := json.Marshal(patch)

		board := &model.Board{
			ID:     "board1",
			TeamID: "team1",
			Title:  "Old Title",
			Type:   model.BoardTypeOpen,
		}
		updatedBoard := &model.Board{
			ID:     "board1",
			TeamID: "team1",
			Title:  "Updated Title",
			Type:   model.BoardTypeOpen,
		}

		helper.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()
		helper.Store.EXPECT().PatchBoard("board1", patch, "single-user").Return(updatedBoard, nil)

		req := helper.NewRequest("PATCH", "/api/v2/boards/board1", body)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp model.Board
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", resp.Title)
	})
}

func TestDeleteBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		board := &model.Board{
			ID:     "board1",
			TeamID: "team1",
			Type:   model.BoardTypeOpen,
		}

		helper.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()
		helper.Store.EXPECT().DeleteBoard("board1", "single-user").Return(nil)

		req := helper.NewRequest("DELETE", "/api/v2/boards/board1", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDuplicateBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		board := &model.Board{
			ID:     "board1",
			TeamID: "team1",
			Type:   model.BoardTypeOpen,
		}

		respBoards := &model.BoardsAndBlocks{
			Boards: []*model.Board{{ID: "new_board", TeamID: "team1"}},
		}

		helper.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()
		helper.Store.EXPECT().DuplicateBoard("board1", "single-user", "team1", false).Return(respBoards, []*model.BoardMember{}, nil)

		req := helper.NewRequest("POST", "/api/v2/boards/board1/duplicate?asTemplate=false", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp model.BoardsAndBlocks
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp.Boards, 1)
		require.Equal(t, "new_board", resp.Boards[0].ID)
	})
}

func TestUndeleteBoard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		board := &model.Board{
			ID:       "board1",
			TeamID:   "team1",
			Type:     model.BoardTypeOpen,
			DeleteAt: 12345,
		}

		helper.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()
		helper.Store.EXPECT().GetBoardHistory("board1", gomock.Any()).Return([]*model.Board{board}, nil)
		helper.Store.EXPECT().UndeleteBoard("board1", "single-user").Return(nil)

		req := helper.NewRequest("POST", "/api/v2/boards/board1/undelete", nil)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateBoardsAndBlocks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		helper, tearDown := SetupTestAPI(t)
		defer tearDown()

		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{{ID: "board1", TeamID: "team1", Title: "Board 1", Type: model.BoardTypeOpen}},
			Blocks: []*model.Block{{ID: "block1", BoardID: "board1", Type: model.TypeCard, CreateAt: 1234, UpdateAt: 1234}},
		}
		body, _ := json.Marshal(bab)

		helper.Store.EXPECT().CreateBoardsAndBlocksWithAdmin(gomock.Any(), gomock.Eq("single-user")).Return(bab, []*model.BoardMember{}, nil)

		req := helper.NewRequest("POST", "/api/v2/boards-and-blocks", body)
		w := httptest.NewRecorder()
		helper.Router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp model.BoardsAndBlocks
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp.Boards, 1)
	})
}
