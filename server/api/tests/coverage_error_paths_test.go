package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func executeAPIRequest(t *testing.T, method, path, body string, configure func(*ExtendedTestAPIHelper)) *http.Response {
	t.Helper()

	th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
	t.Cleanup(tearDown)
	if configure != nil {
		configure(th)
	}

	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
	req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

	recorder := doRequest(th.Router, req)
	return recorder.Result()
}

func denyTeamPermission(th *ExtendedTestAPIHelper) {
	th.Permissions.HasPermissionToTeamFunc = func(_, _ string, _ *mmModel.Permission) bool { return false }
}

func denyBoardPermission(th *ExtendedTestAPIHelper) {
	th.Permissions.HasPermissionToBoardFunc = func(_, _ string, _ *mmModel.Permission) bool { return false }
}

func TestCoverageMalformedJSON(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		path      string
		configure func(*ExtendedTestAPIHelper)
	}{
		{"create board", http.MethodPost, "/api/v2/boards", nil},
		{"patch board", http.MethodPatch, "/api/v2/boards/board1", func(th *ExtendedTestAPIHelper) {
			th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)
		}},
		{"create boards and blocks", http.MethodPost, "/api/v2/boards-and-blocks", nil},
		{"patch boards and blocks", http.MethodPatch, "/api/v2/boards-and-blocks", nil},
		{"delete boards and blocks", http.MethodDelete, "/api/v2/boards-and-blocks", nil},
		{"post blocks", http.MethodPost, "/api/v2/boards/board1/blocks", nil},
		{"patch block", http.MethodPatch, "/api/v2/boards/board1/blocks/block1", func(th *ExtendedTestAPIHelper) {
			th.Store.EXPECT().GetBlock("block1").Return(&model.Block{ID: "block1", BoardID: "board1"}, nil)
		}},
		{"patch blocks", http.MethodPatch, "/api/v2/boards/board1/blocks", nil},
		{"add member", http.MethodPost, "/api/v2/boards/board1/members", func(th *ExtendedTestAPIHelper) {
			th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil)
		}},
		{"update member", http.MethodPut, "/api/v2/boards/board1/members/user2", nil},
		{"create category", http.MethodPost, "/api/v2/teams/team1/categories", nil},
		{"update category", http.MethodPut, "/api/v2/teams/team1/categories/category1", nil},
		{"reorder categories", http.MethodPut, "/api/v2/teams/team1/categories/reorder", nil},
		{"reorder category boards", http.MethodPut, "/api/v2/teams/team1/categories/category1/boards/reorder", func(th *ExtendedTestAPIHelper) {
			th.Store.EXPECT().GetCategory("category1").Return(&model.Category{ID: "category1", UserID: "single-user", TeamID: "team1"}, nil)
		}},
		{"list users", http.MethodPost, "/api/v2/users", nil},
		{"update user config", http.MethodPut, "/api/v2/users/single-user/config", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := executeAPIRequest(t, tc.method, tc.path, `{`, tc.configure)
			defer resp.Body.Close()
			require.GreaterOrEqual(t, resp.StatusCode, http.StatusBadRequest)
		})
	}
}

func TestCoveragePermissionDenied(t *testing.T) {
	teamRoutes := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"get boards", http.MethodGet, "/api/v2/teams/team1/boards", ""},
		{"create category", http.MethodPost, "/api/v2/teams/team1/categories", `{"id":"cat1","userID":"single-user","teamID":"team1","name":"Category","type":"custom"}`},
		{"update category", http.MethodPut, "/api/v2/teams/team1/categories/cat1", `{"id":"cat1","userID":"single-user","teamID":"team1","name":"Category","type":"custom"}`},
		{"delete category", http.MethodDelete, "/api/v2/teams/team1/categories/cat1", ""},
		{"get categories", http.MethodGet, "/api/v2/teams/team1/categories", ""},
		{"update category board", http.MethodPost, "/api/v2/teams/team1/categories/cat1/boards/board1", ""},
		{"reorder categories", http.MethodPut, "/api/v2/teams/team1/categories/reorder", `[]`},
		{"reorder category boards", http.MethodPut, "/api/v2/teams/team1/categories/cat1/boards/reorder", `[]`},
		{"hide board", http.MethodPut, "/api/v2/teams/team1/categories/cat1/boards/board1/hide", ""},
		{"unhide board", http.MethodPut, "/api/v2/teams/team1/categories/cat1/boards/board1/unhide", ""},
	}

	for _, tc := range teamRoutes {
		t.Run(tc.name, func(t *testing.T) {
			resp := executeAPIRequest(t, tc.method, tc.path, tc.body, denyTeamPermission)
			defer resp.Body.Close()
			require.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	}

	boardRoutes := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"get members", http.MethodGet, "/api/v2/boards/board1/members", ""},
		{"leave board", http.MethodPost, "/api/v2/boards/board1/leave", ""},
		{"delete block", http.MethodDelete, "/api/v2/boards/board1/blocks/block1", ""},
		{"patch block", http.MethodPatch, "/api/v2/boards/board1/blocks/block1", `{}`},
	}

	for _, tc := range boardRoutes {
		t.Run(tc.name, func(t *testing.T) {
			resp := executeAPIRequest(t, tc.method, tc.path, tc.body, denyBoardPermission)
			defer resp.Body.Close()
			require.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	}
}

func TestCoverageRequestValidation(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"empty users list", http.MethodPost, "/api/v2/users", `[]`},
		{"boards and blocks requires board", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[],"blocks":[]}`},
		{"boards must use one team", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"},{"id":"b2","teamId":"t2"}]}`},
		{"boards require IDs", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"},{"teamId":"t1"}]}`},
		{"new block requires type", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"}],"blocks":[{"id":"x","boardId":"b1","createAt":1,"updateAt":1}]}`},
		{"new block requires create time", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"}],"blocks":[{"id":"x","boardId":"b1","type":"card","updateAt":1}]}`},
		{"new block requires update time", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"}],"blocks":[{"id":"x","boardId":"b1","type":"card","createAt":1}]}`},
		{"new block references created board", http.MethodPost, "/api/v2/boards-and-blocks", `{"boards":[{"id":"b1","teamId":"t1"}],"blocks":[{"id":"x","boardId":"other","type":"card","createAt":1,"updateAt":1}]}`},
		{"block requires type", http.MethodPost, "/api/v2/boards/board1/blocks", `[{"id":"x","boardId":"board1","createAt":1,"updateAt":1}]`},
		{"block requires create time", http.MethodPost, "/api/v2/boards/board1/blocks", `[{"id":"x","boardId":"board1","type":"card","updateAt":1}]`},
		{"block requires update time", http.MethodPost, "/api/v2/boards/board1/blocks", `[{"id":"x","boardId":"board1","type":"card","createAt":1}]`},
		{"block board must match URL", http.MethodPost, "/api/v2/boards/board1/blocks", `[{"id":"x","boardId":"other","type":"card","createAt":1,"updateAt":1}]`},
		{"category team must match URL", http.MethodPost, "/api/v2/teams/team1/categories", `{"id":"cat1","userID":"single-user","teamID":"other","name":"Category","type":"custom"}`},
		{"category ID must match URL", http.MethodPut, "/api/v2/teams/team1/categories/cat1", `{"id":"other","userID":"single-user","teamID":"team1","name":"Category","type":"custom"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := executeAPIRequest(t, tc.method, tc.path, tc.body, nil)
			defer resp.Body.Close()
			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
