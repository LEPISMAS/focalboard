package tests

import (
<<<<<<< HEAD
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
=======
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestBoardsEndpoints(t *testing.T) {
	t.Run("GET /teams/{teamID}/boards returns list of boards", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boards := []*model.Board{
			{ID: "board1", TeamID: "team1", Title: "Board 1"},
		}
		th.Store.EXPECT().GetBoardsForUserAndTeam("single-user", "team1", true).Return(boards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/boards", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("POST /boards creates a board successfully", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Permissions.HasPermissionToTeamFunc = func(userID, teamID string, permission *mmModel.Permission) bool {
			return true
		}
		
		createdBoard := &model.Board{ID: "newboard123", TeamID: "team1", Title: "New Board"}
		createdMember := &model.BoardMember{BoardID: "newboard123", UserID: "single-user", SchemeAdmin: true}
		th.Store.EXPECT().InsertBoardWithAdmin(gomock.Any(), "single-user").Return(createdBoard, createdMember, nil)

		// addBoardsToDefaultCategory and AddUpdateUserCategoryBoard mock calls
		categories := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:   "category1",
					Name: "Boards", // default category name
				},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categories, nil).Times(2)
		th.Store.EXPECT().AddUpdateCategoryBoard("single-user", "category1", []string{"newboard123"}).Return(nil)
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		body := `{"teamId": "team1", "title": "New Board", "type": "O"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "newboard123")
	})

	t.Run("GET /boards/{boardID} returns single board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("PATCH /boards/{boardID} updates board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen}, nil)
		th.Store.EXPECT().PatchBoard("board1", gomock.Any(), "single-user").Return(&model.Board{ID: "board1", TeamID: "team1", Title: "Patched Title"}, nil)
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		body := `{"title": "Patched Title"}`
		req, _ := http.NewRequest(http.MethodPatch, "/api/v2/boards/board1", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Patched Title")
	})

	t.Run("DELETE /boards/{boardID} deletes board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil).Times(2)
		th.Store.EXPECT().DeleteBoard("board1", "single-user").Return(nil)
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/boards/board1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("GET /boards/{boardID}/members returns board members", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		members := []*model.BoardMember{
			{BoardID: "board1", UserID: "user2", SchemeAdmin: false},
		}
		th.Store.EXPECT().GetMembersForBoard("board1").Return(members, nil)
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/members", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user2")
	})

	t.Run("POST /boards/{boardID}/members adds board member", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1", IsTemplate: true}, nil).Times(2)
		th.Store.EXPECT().GetMemberForBoard("board1", "user2").Return(nil, model.NewErrNotFound("member"))
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
		
		savedMember := &model.BoardMember{BoardID: "board1", UserID: "user2", SchemeAdmin: false}
		th.Store.EXPECT().SaveMember(gomock.Any()).Return(savedMember, nil)

		body := `{"userId": "user2"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/members", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user2")
	})

	t.Run("POST /boards/{boardID}/duplicate duplicates board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(2)

		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{
				{ID: "board2", TeamID: "team1", Title: "Duplicated Board"},
			},
		}
		th.Store.EXPECT().DuplicateBoard("board1", "single-user", "team1", false).Return(bab, []*model.BoardMember{}, nil)

		// mock DefaultCategory board addition
		categories := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:     "cat1",
					Name:   "Boards",
					Type:   "system",
					UserID: "single-user",
					TeamID: "team1",
				},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categories, nil).AnyTimes()
		th.Store.EXPECT().AddUpdateCategoryBoard("single-user", "cat1", []string{"board2"}).Return(nil).AnyTimes()
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/duplicate", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Duplicated Board")
	})

	t.Run("POST /boards/{boardID}/undelete undeletes board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoardHistory("board1", gomock.Any()).Return([]*model.Board{board}, nil)
		th.Store.EXPECT().UndeleteBoard("board1", "single-user").Return(nil)
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/undelete", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("GET /boards/{boardID}/metadata returns board metadata", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		boolTrue := true
		license := &mmModel.License{
			Features: &mmModel.Features{
				Compliance: &boolTrue,
			},
		}
		th.Store.EXPECT().GetLicense().Return(license)

		board := &model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen, CreatedBy: "user1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)

		// getBoardDescendantModifiedInfo mock calls (earliest & latest)
		th.Store.EXPECT().GetBoardHistory("board1", gomock.Any()).Return([]*model.Board{board}, nil).Times(2)
		th.Store.EXPECT().GetBlockHistoryDescendants("board1", gomock.Any()).Return([]*model.Block{}, nil).Times(2)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/metadata", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
>>>>>>> bb9efff28a5b47309a8fe28b149a8c9533a9a2e6
	})
}
