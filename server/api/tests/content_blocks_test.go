package tests

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestContentBlocksEndpoints(t *testing.T) {
	t.Run("POST /content-blocks/{blockID}/moveto/{where}/{dstBlockID} moves block", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		block := &model.Block{ID: "block1", BoardID: "board1", ParentID: "card1"}
		dstBlock := &model.Block{ID: "block2", BoardID: "board1", ParentID: "card1"}
		card := &model.Block{
			ID:      "card1",
			BoardID: "board1",
			Fields: map[string]interface{}{
				"contentOrder": []interface{}{"block1", "block2"},
			},
		}

		th.Store.EXPECT().GetBlock("block1").Return(block, nil)
		th.Store.EXPECT().GetBlock("block2").Return(dstBlock, nil)

		// GetBlock("card1") is called 3 times:
		// 1. GetBlockByID(block.ParentID) in MoveContentBlock
		// 2. PatchBlockAndNotify (oldBlock)
		// 3. PatchBlockAndNotify (newBlock)
		th.Store.EXPECT().GetBlock("card1").Return(card, nil).Times(3)

		th.Store.EXPECT().GetBoard("board1").Return(&model.Board{ID: "board1", TeamID: "team1"}, nil).AnyTimes()
		th.Store.EXPECT().PatchBlock("card1", gomock.Any(), "single-user").Return(nil)
		th.Store.EXPECT().GetMembersForBoard("board1").Return([]*model.BoardMember{}, nil).AnyTimes()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/content-blocks/block1/moveto/after/block2", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /content-blocks/{blockID}/moveto/{where}/{dstBlockID} returns 400 for invalid where parameter", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		block := &model.Block{ID: "block1", BoardID: "board1", ParentID: "card1"}
		dstBlock := &model.Block{ID: "block2", BoardID: "board1", ParentID: "card1"}

		th.Store.EXPECT().GetBlock("block1").Return(block, nil)
		th.Store.EXPECT().GetBlock("block2").Return(dstBlock, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/content-blocks/block1/moveto/invalid/block2", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, resp.Body.String(), "invalid where parameter")
	})
}
