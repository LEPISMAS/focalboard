package tests

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestOnboardingEndpoints(t *testing.T) {
	t.Run("POST /teams/{teamID}/onboard prepares onboarding tour", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		// 1. GetTemplateBoards returns welcome board template
		templates := []*model.Board{
			{ID: "welcome-board-template", Title: "Welcome to Boards!", TeamID: "0"},
		}
		th.Store.EXPECT().GetTemplateBoards("0", "").Return(templates, nil)

		// GetBoard("welcome-board-template") is called during duplication
		th.Store.EXPECT().GetBoard("welcome-board-template").Return(&model.Board{ID: "welcome-board-template", TeamID: "0"}, nil).AnyTimes()

		// 2. DuplicateBoard duplicates welcome board
		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{
				{ID: "new-welcome-board", TeamID: "team1"},
			},
		}
		th.Store.EXPECT().DuplicateBoard("welcome-board-template", "single-user", "team1", false).Return(bab, []*model.BoardMember{}, nil)

		// GetUserCategoryBoards is called in AddBoardToCategory, and needs a default system category
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
		th.Store.EXPECT().CreateCategory(gomock.Any()).Return(nil).AnyTimes()
		th.Store.EXPECT().GetCategory(gomock.Any()).Return(&model.Category{ID: "cat1"}, nil).AnyTimes()
		th.Store.EXPECT().AddUpdateCategoryBoard("single-user", "cat1", []string{"new-welcome-board"}).Return(nil).AnyTimes()

		// GetMembersForUser is called during PatchBoard permission verification/checking
		th.Store.EXPECT().GetMembersForUser("single-user").Return([]*model.BoardMember{}, nil).AnyTimes()
		th.Store.EXPECT().GetBoardsForUserAndTeam("single-user", "team1", false).Return([]*model.Board{}, nil).AnyTimes()

		// 3. PatchBoard patches the duplicated board
		th.Store.EXPECT().GetBoard("new-welcome-board").Return(&model.Board{ID: "new-welcome-board", TeamID: "team1"}, nil).AnyTimes()
		th.Store.EXPECT().PatchBoard("new-welcome-board", gomock.Any(), "single-user").Return(&model.Board{ID: "new-welcome-board", TeamID: "team1"}, nil)
		th.Store.EXPECT().GetMembersForBoard("new-welcome-board").Return([]*model.BoardMember{}, nil).AnyTimes()

		// 4. PatchUserPreferences sets onboarding preferences
		th.Store.EXPECT().PatchUserPreferences("single-user", gomock.Any()).Return(mmModel.Preferences{}, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/onboard", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "new-welcome-board")
	})
}
