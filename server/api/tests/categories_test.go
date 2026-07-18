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

func TestCategoriesEndpoints(t *testing.T) {
	t.Run("POST /teams/{teamID}/categories creates category", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().CreateCategory(gomock.Any()).Return(nil)
		th.Store.EXPECT().GetCategory("category1").Return(&model.Category{
			ID:     "category1",
			UserID: "single-user",
			TeamID: "team1",
			Name:   "Category 1",
			Type:   model.CategoryTypeCustom,
		}, nil)

		body := `{"id": "category1", "userID": "single-user", "teamID": "team1", "name": "Category 1", "type": "custom"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/categories", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Category 1")
	})

	t.Run("PUT /teams/{teamID}/categories/{categoryID} updates category", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		existingCat := &model.Category{
			ID:     "category1",
			UserID: "single-user",
			TeamID: "team1",
			Name:   "Old Name",
			Type:   model.CategoryTypeCustom,
		}
		th.Store.EXPECT().GetCategory("category1").Return(existingCat, nil).Times(2)
		th.Store.EXPECT().UpdateCategory(gomock.Any()).Return(nil)

		body := `{"id": "category1", "userID": "single-user", "teamID": "team1", "name": "New Name", "type": "custom"}`
		req, _ := http.NewRequest(http.MethodPut, "/api/v2/teams/team1/categories/category1", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("DELETE /teams/{teamID}/categories/{categoryID} deletes category", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		existingCat := &model.Category{
			ID:     "category1",
			UserID: "single-user",
			TeamID: "team1",
			Name:   "Category 1",
			Type:   model.CategoryTypeCustom,
		}
		th.Store.EXPECT().GetCategory("category1").Return(existingCat, nil).Times(2)

		categoryBoards := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:   "defaultCat",
					Name: "Boards",
					Type: model.CategoryTypeSystem,
				},
				BoardMetadata: []model.CategoryBoardMetadata{},
			},
			{
				Category:      *existingCat,
				BoardMetadata: []model.CategoryBoardMetadata{},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categoryBoards, nil)
		th.Store.EXPECT().DeleteCategory("category1", "single-user", "team1").Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/teams/team1/categories/category1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("GET /teams/{teamID}/categories returns user category boards", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		categoryBoards := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:   "cat1",
					Name: "Boards",
					Type: model.CategoryTypeSystem,
				},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categoryBoards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/categories", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "Boards")
	})

	t.Run("POST /teams/{teamID}/categories/{categoryID}/boards/{boardID} updates category board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().AddUpdateCategoryBoard("single-user", "category1", []string{"board1"}).Return(nil)

		categoryBoards := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:     "default_boards",
					Name:   "Boards",
					UserID: "single-user",
					TeamID: "team1",
					Type:   model.CategoryTypeSystem,
				},
			},
			{
				Category: model.Category{
					ID:     "category1",
					UserID: "single-user",
					TeamID: "team1",
					Type:   model.CategoryTypeCustom,
				},
				BoardMetadata: []model.CategoryBoardMetadata{
					{BoardID: "board1"},
				},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categoryBoards, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/categories/category1/boards/board1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "ok", resp.Body.String())
	})

	t.Run("PUT /teams/{teamID}/categories/reorder reorders categories", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		existingCats := []model.Category{
			{ID: "cat1", UserID: "single-user", TeamID: "team1", Type: model.CategoryTypeCustom},
			{ID: "cat2", UserID: "single-user", TeamID: "team1", Type: model.CategoryTypeCustom},
		}
		th.Store.EXPECT().GetUserCategories("single-user", "team1").Return(existingCats, nil)
		th.Store.EXPECT().ReorderCategories("single-user", "team1", []string{"cat2", "cat1"}).Return([]string{"cat2", "cat1"}, nil)

		body := `["cat2", "cat1"]`
		req, _ := http.NewRequest(http.MethodPut, "/api/v2/teams/team1/categories/reorder", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("PUT /teams/{teamID}/categories/{categoryID}/boards/reorder reorders boards in category", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetCategory("category1").Return(&model.Category{
			ID:     "category1",
			UserID: "single-user",
			TeamID: "team1",
			Type:   model.CategoryTypeCustom,
		}, nil)

		categoryBoards := []model.CategoryBoards{
			{
				Category: model.Category{
					ID:     "default_boards",
					Name:   "Boards",
					UserID: "single-user",
					TeamID: "team1",
					Type:   model.CategoryTypeSystem,
				},
			},
			{
				Category: model.Category{
					ID:     "category1",
					UserID: "single-user",
					TeamID: "team1",
					Type:   model.CategoryTypeCustom,
				},
				BoardMetadata: []model.CategoryBoardMetadata{
					{BoardID: "board1"},
					{BoardID: "board2"},
				},
			},
		}
		th.Store.EXPECT().GetUserCategoryBoards("single-user", "team1").Return(categoryBoards, nil)
		th.Store.EXPECT().ReorderCategoryBoards("category1", []string{"board2", "board1"}).Return([]string{"board2", "board1"}, nil)

		body := `["board2", "board1"]`
		req, _ := http.NewRequest(http.MethodPut, "/api/v2/teams/team1/categories/category1/boards/reorder", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("PUT /teams/{teamID}/categories/{categoryID}/boards/{boardID}/hide hides board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().SetBoardVisibility("single-user", "category1", "board1", false).Return(nil)

		req, _ := http.NewRequest(http.MethodPut, "/api/v2/teams/team1/categories/category1/boards/board1/hide", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("PUT /teams/{teamID}/categories/{categoryID}/boards/{boardID}/unhide unhides board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().SetBoardVisibility("single-user", "category1", "board1", true).Return(nil)

		req, _ := http.NewRequest(http.MethodPut, "/api/v2/teams/team1/categories/category1/boards/board1/unhide", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
	})
}
