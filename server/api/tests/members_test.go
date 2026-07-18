package tests

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestMembersEndpoints(t *testing.T) {
	t.Run("GET /boards/{boardID}/members returns board members", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).AnyTimes()

		members := []*model.BoardMember{
			{UserID: "user1", BoardID: "board1"},
		}
		th.Store.EXPECT().GetMembersForBoard("board1").Return(members, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/members", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user1")
	})

	t.Run("POST /boards/{boardID}/members adds member to board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		// Setting IsTemplate: true avoids calling category updates
		board := &model.Board{ID: "board1", TeamID: "team1", IsTemplate: true}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(2)
		th.Store.EXPECT().GetMemberForBoard("board1", "user2").Return(nil, model.NewErrNotFound("membership not found"))

		member := &model.BoardMember{UserID: "user2", BoardID: "board1"}
		th.Store.EXPECT().SaveMember(gomock.Any()).Return(member, nil)

		// Asynchronous WS BroadcastMemberChange will call GetMembersForBoard
		members := []*model.BoardMember{member}
		th.Store.EXPECT().GetMembersForBoard("board1").Return(members, nil).AnyTimes()

		body := `{"userID": "user2", "boardID": "board1"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/members", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "user2")

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("PUT /boards/{boardID}/members/{userID} updates board member", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)

		oldMember := &model.BoardMember{UserID: "user2", BoardID: "board1", SchemeAdmin: false}
		th.Store.EXPECT().GetMemberForBoard("board1", "user2").Return(oldMember, nil)

		updatedMember := &model.BoardMember{UserID: "user2", BoardID: "board1", SchemeAdmin: true}
		th.Store.EXPECT().SaveMember(gomock.Any()).Return(updatedMember, nil)

		// Asynchronous WS BroadcastMemberChange will call GetMembersForBoard
		members := []*model.BoardMember{updatedMember}
		th.Store.EXPECT().GetMembersForBoard("board1").Return(members, nil).AnyTimes()

		body := `{"schemeAdmin": true}`
		req, _ := http.NewRequest(http.MethodPut, "/api/v2/boards/board1/members/user2", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("DELETE /boards/{boardID}/members/{userID} deletes member", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(2)

		oldMember := &model.BoardMember{UserID: "user2", BoardID: "board1", SchemeAdmin: false}
		th.Store.EXPECT().GetMemberForBoard("board1", "user2").Return(oldMember, nil).AnyTimes()
		th.Store.EXPECT().DeleteMember("board1", "user2").Return(nil)

		// Asynchronous WS BroadcastMemberDelete will call GetMembersForBoard
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/boards/board1/members/user2", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("POST /boards/{boardID}/join joins board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1", Type: model.BoardTypeOpen, IsTemplate: true}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(2)
		th.Store.EXPECT().GetMemberForBoard("board1", "single-user").Return(nil, model.NewErrNotFound("membership not found"))

		member := &model.BoardMember{UserID: "single-user", BoardID: "board1"}
		th.Store.EXPECT().SaveMember(gomock.Any()).Return(member, nil)

		// Asynchronous WS BroadcastMemberChange will call GetMembersForBoard
		members := []*model.BoardMember{member}
		th.Store.EXPECT().GetMembersForBoard("board1").Return(members, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/join", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("POST /boards/{boardID}/leave leaves board", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil).Times(2)

		oldMember := &model.BoardMember{UserID: "single-user", BoardID: "board1", SchemeAdmin: false}
		th.Store.EXPECT().GetMemberForBoard("board1", "single-user").Return(oldMember, nil).AnyTimes()
		th.Store.EXPECT().DeleteMember("board1", "single-user").Return(nil)

		// Asynchronous WS BroadcastMemberDelete will call GetMembersForBoard
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/leave", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})
}
