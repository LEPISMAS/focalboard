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

func TestCardsEndpoints(t *testing.T) {
	t.Run("POST /boards/{boardID}/cards creates card", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)
		th.Store.EXPECT().InsertBlock(gomock.Any(), "single-user").Return(nil)
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil).AnyTimes()

		body := `{"id": "card1", "boardId": "board1", "title": "New Card"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/boards/board1/cards", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "board1")

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("GET /boards/{boardID}/cards returns board cards", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		blocks := []*model.Block{
			{ID: "card1", BoardID: "board1", Type: model.TypeCard, Title: "New Card"},
		}
		opts := model.QueryBlocksOptions{
			BoardID:   "board1",
			BlockType: model.TypeCard,
			Page:      0,
			PerPage:   100,
		}
		th.Store.EXPECT().GetBlocks(opts).Return(blocks, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/boards/board1/cards", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "card1")
	})

	t.Run("PATCH /cards/{cardID} patches card", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		cardBlock := &model.Block{ID: "card1", BoardID: "board1", Type: model.TypeCard, Title: "Old Title"}
		// GetBlock is called:
		// 1. GetCardByID in handler
		// 2. PatchBlockAndNotify (oldBlock)
		th.Store.EXPECT().GetBlock("card1").Return(cardBlock, nil).Times(2)

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)

		th.Store.EXPECT().PatchBlock("card1", gomock.Any(), "single-user").Return(nil)

		// 3. PatchBlockAndNotify (newBlock)
		patchedBlock := &model.Block{ID: "card1", BoardID: "board1", Type: model.TypeCard, Title: "New Title"}
		th.Store.EXPECT().GetBlock("card1").Return(patchedBlock, nil).Times(1)

		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil).AnyTimes()

		body := `{"title": "New Title"}`
		req, _ := http.NewRequest(http.MethodPatch, "/api/v2/cards/card1", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "New Title")

		// Let background goroutines finish
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("GET /cards/{cardID} returns card by ID", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		cardBlock := &model.Block{ID: "card1", BoardID: "board1", Type: model.TypeCard, Title: "Card 1"}
		th.Store.EXPECT().GetBlock("card1").Return(cardBlock, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/cards/card1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "card1")
	})
}
