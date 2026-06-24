package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestComplianceEndpoints(t *testing.T) {
	trueVal := true
	complianceLicense := &mmModel.License{
		Features: &mmModel.Features{
			Compliance: &trueVal,
		},
	}

	t.Run("GET /admin/boards returns unauthorized if permission missing", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Permissions.HasPermissionToFunc = func(userID string, permission *mmModel.Permission) bool {
			if permission.Id == mmModel.PermissionManageSystem.Id {
				return false
			}
			return true
		}

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/boards", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), "access denied Compliance Export getAllBoards")
	})

	t.Run("GET /admin/boards returns not implemented if license missing", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetLicense().Return(nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/boards", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
		require.Contains(t, resp.Body.String(), "insufficient license Compliance Export getAllBoards")
	})

	t.Run("GET /admin/boards returns bad request if team not found", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetLicense().Return(complianceLicense)
		th.Store.EXPECT().GetTeam("invalid-team").Return(nil, errors.New("team not found"))

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/boards?team_id=invalid-team", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), "invalid team id: invalid-team")
	})

	t.Run("GET /admin/boards returns boards list on success", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetLicense().Return(complianceLicense)
		th.Store.EXPECT().GetTeam("team1").Return(&model.Team{ID: "team1"}, nil)

		boards := []*model.Board{
			{ID: "board1", TeamID: "team1"},
		}
		opts := model.QueryBoardsForComplianceOptions{
			TeamID:  "team1",
			Page:    0,
			PerPage: 60,
		}
		th.Store.EXPECT().GetBoardsForCompliance(opts).Return(boards, false, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/boards?team_id=team1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("GET /admin/boards_history returns boards compliance history", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetLicense().Return(complianceLicense)

		history := []*model.BoardHistory{
			{ID: "board1", TeamID: "team1"},
		}
		opts := model.QueryBoardsComplianceHistoryOptions{
			ModifiedSince:  1000,
			IncludeDeleted: true,
			Page:           0,
			PerPage:        60,
		}
		th.Store.EXPECT().GetBoardsComplianceHistory(opts).Return(history, false, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/boards_history?modified_since=1000&include_deleted=true", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")
	})

	t.Run("GET /admin/blocks_history returns blocks compliance history", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.Store.EXPECT().GetLicense().Return(complianceLicense)

		history := []*model.BlockHistory{
			{ID: "block1", BoardID: "board1"},
		}
		opts := model.QueryBlocksComplianceHistoryOptions{
			ModifiedSince:  1000,
			IncludeDeleted: false,
			BoardID:        "board1",
			Page:           0,
			PerPage:        60,
		}
		th.Store.EXPECT().GetBlocksComplianceHistory(opts).Return(history, false, nil)
		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1"}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/admin/blocks_history?modified_since=1000&board_id=board1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "block1")
	})
}
