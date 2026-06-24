package tests

import (
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
	})
}
