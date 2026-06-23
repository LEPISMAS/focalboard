package app

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mmModel "github.com/mattermost/mattermost/server/public/model"

	"github.com/mattermost/focalboard/server/model"
)

type blockError struct {
	msg string
}

func (be blockError) Error() string {
	return be.msg
}

func TestInsertBlock(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	t.Run("success scenario", func(t *testing.T) {
		boardID := testBoardID
		block := &model.Block{BoardID: boardID}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(block, "user-id-1").Return(nil)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil)
		err := th.App.InsertBlock(block, "user-id-1")
		require.NoError(t, err)
	})

	t.Run("error scenario", func(t *testing.T) {
		boardID := testBoardID
		block := &model.Block{BoardID: boardID}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(block, "user-id-1").Return(blockError{"error"})
		err := th.App.InsertBlock(block, "user-id-1")
		require.Error(t, err, "error")
	})
}

func TestPatchBlocks(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	t.Run("patchBlocks success scenario", func(t *testing.T) {
		blockPatches := model.BlockPatchBatch{
			BlockIDs: []string{"block1"},
			BlockPatches: []model.BlockPatch{
				{Title: mmModel.NewString("new title")},
			},
		}

		block1 := &model.Block{ID: "block1"}
		th.Store.EXPECT().GetBlocksByIDs([]string{"block1"}).Return([]*model.Block{block1}, nil)
		th.Store.EXPECT().PatchBlocks(gomock.Eq(&blockPatches), gomock.Eq("user-id-1")).Return(nil)
		th.Store.EXPECT().GetBlock("block1").Return(block1, nil)
		// this call comes from the WS server notification
		th.Store.EXPECT().GetMembersForBoard(gomock.Any()).Times(1)
		err := th.App.PatchBlocks("team-id", &blockPatches, "user-id-1")
		require.NoError(t, err)
	})

	t.Run("patchBlocks error scenario", func(t *testing.T) {
		blockPatches := model.BlockPatchBatch{BlockIDs: []string{}}
		th.Store.EXPECT().GetBlocksByIDs([]string{}).Return(nil, sql.ErrNoRows)
		err := th.App.PatchBlocks("team-id", &blockPatches, "user-id-1")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("cloud limit error scenario", func(t *testing.T) {
		t.Skipf("The Cloud Limits feature has been disabled")

		th.App.SetCardLimit(5)

		fakeLicense := &mmModel.License{
			Features: &mmModel.Features{Cloud: mmModel.NewBool(true)},
		}

		blockPatches := model.BlockPatchBatch{
			BlockIDs: []string{"block1"},
			BlockPatches: []model.BlockPatch{
				{Title: mmModel.NewString("new title")},
			},
		}

		block1 := &model.Block{
			ID:       "block1",
			Type:     model.TypeCard,
			ParentID: "board-id",
			BoardID:  "board-id",
			UpdateAt: 100,
		}

		board1 := &model.Board{
			ID:   "board-id",
			Type: model.BoardTypeOpen,
		}

		th.Store.EXPECT().GetBlocksByIDs([]string{"block1"}).Return([]*model.Block{block1}, nil)
		th.Store.EXPECT().GetBoard("board-id").Return(board1, nil)
		th.Store.EXPECT().GetLicense().Return(fakeLicense)
		th.Store.EXPECT().GetCardLimitTimestamp().Return(int64(150), nil)
		err := th.App.PatchBlocks("team-id", &blockPatches, "user-id-1")
		require.ErrorIs(t, err, model.ErrPatchUpdatesLimitedCards)
	})
}

func TestDeleteBlock(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	t.Run("success scenario", func(t *testing.T) {
		boardID := testBoardID
		board := &model.Board{ID: boardID}
		block := &model.Block{
			ID:      "block-id",
			BoardID: board.ID,
		}
		th.Store.EXPECT().GetBlock(gomock.Eq("block-id")).Return(block, nil)
		th.Store.EXPECT().DeleteBlock(gomock.Eq("block-id"), gomock.Eq("user-id-1")).Return(nil)
		th.Store.EXPECT().GetBoard(gomock.Eq(testBoardID)).Return(board, nil)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil)
		err := th.App.DeleteBlock("block-id", "user-id-1")
		require.NoError(t, err)
	})

	t.Run("error scenario", func(t *testing.T) {
		boardID := testBoardID
		board := &model.Board{ID: boardID}
		block := &model.Block{
			ID:      "block-id",
			BoardID: board.ID,
		}
		th.Store.EXPECT().GetBlock(gomock.Eq("block-id")).Return(block, nil)
		th.Store.EXPECT().DeleteBlock(gomock.Eq("block-id"), gomock.Eq("user-id-1")).Return(blockError{"error"})
		th.Store.EXPECT().GetBoard(gomock.Eq(testBoardID)).Return(board, nil)
		err := th.App.DeleteBlock("block-id", "user-id-1")
		require.Error(t, err, "error")
	})
}

func TestUndeleteBlock(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	t.Run("success scenario", func(t *testing.T) {
		boardID := testBoardID
		board := &model.Board{ID: boardID}
		block := &model.Block{
			ID:      "block-id",
			BoardID: board.ID,
		}
		th.Store.EXPECT().GetBlockHistory(
			gomock.Eq("block-id"),
			gomock.Eq(model.QueryBlockHistoryOptions{Limit: 1, Descending: true}),
		).Return([]*model.Block{block}, nil)
		th.Store.EXPECT().UndeleteBlock(gomock.Eq("block-id"), gomock.Eq("user-id-1")).Return(nil)
		th.Store.EXPECT().GetBlock(gomock.Eq("block-id")).Return(block, nil)
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil)
		_, err := th.App.UndeleteBlock("block-id", "user-id-1")
		require.NoError(t, err)
	})

	t.Run("error scenario", func(t *testing.T) {
		block := &model.Block{
			ID: "block-id",
		}
		th.Store.EXPECT().GetBlockHistory(
			gomock.Eq("block-id"),
			gomock.Eq(model.QueryBlockHistoryOptions{Limit: 1, Descending: true}),
		).Return([]*model.Block{block}, nil)
		th.Store.EXPECT().UndeleteBlock(gomock.Eq("block-id"), gomock.Eq("user-id-1")).Return(blockError{"error"})
		_, err := th.App.UndeleteBlock("block-id", "user-id-1")
		require.Error(t, err, "error")
	})
}

func TestInsertBlocks(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	t.Run("success scenario", func(t *testing.T) {
		boardID := testBoardID
		block := &model.Block{BoardID: boardID}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(block, "user-id-1").Return(nil)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil)
		_, err := th.App.InsertBlocks([]*model.Block{block}, "user-id-1")
		require.NoError(t, err)
	})

	t.Run("error scenario", func(t *testing.T) {
		boardID := testBoardID
		block := &model.Block{BoardID: boardID}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(block, "user-id-1").Return(blockError{"error"})
		_, err := th.App.InsertBlocks([]*model.Block{block}, "user-id-1")
		require.Error(t, err, "error")
	})

	t.Run("create view within limits", func(t *testing.T) {
		t.Skipf("The Cloud Limits feature has been disabled")

		boardID := testBoardID
		block := &model.Block{
			Type:     model.TypeView,
			ParentID: "parent_id",
			BoardID:  boardID,
		}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(block, "user-id-1").Return(nil)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil)

		// setting up mocks for limits
		fakeLicense := &mmModel.License{
			Features: &mmModel.Features{Cloud: mmModel.NewBool(true)},
		}
		th.Store.EXPECT().GetLicense().Return(fakeLicense)

		th.Store.EXPECT().GetUsedCardsCount().Return(1, nil)
		th.Store.EXPECT().GetCardLimitTimestamp().Return(int64(1), nil)
		th.Store.EXPECT().GetBlocksWithParentAndType("test-board-id", "parent_id", "view").Return([]*model.Block{{}}, nil)

		_, err := th.App.InsertBlocks([]*model.Block{block}, "user-id-1")
		require.NoError(t, err)
	})

	t.Run("create view exceeding limits", func(t *testing.T) {
		t.Skipf("The Cloud Limits feature has been disabled")

		boardID := testBoardID
		block := &model.Block{
			Type:     model.TypeView,
			ParentID: "parent_id",
			BoardID:  boardID,
		}
		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)

		// setting up mocks for limits
		fakeLicense := &mmModel.License{
			Features: &mmModel.Features{Cloud: mmModel.NewBool(true)},
		}
		th.Store.EXPECT().GetLicense().Return(fakeLicense)

		th.Store.EXPECT().GetUsedCardsCount().Return(1, nil)
		th.Store.EXPECT().GetCardLimitTimestamp().Return(int64(1), nil)
		th.Store.EXPECT().GetBlocksWithParentAndType("test-board-id", "parent_id", "view").Return([]*model.Block{{}, {}}, nil)

		_, err := th.App.InsertBlocks([]*model.Block{block}, "user-id-1")
		require.Error(t, err)
	})

	t.Run("creating multiple views, reaching limit in the process", func(t *testing.T) {
		t.Skipf("Will be fixed soon")

		boardID := testBoardID
		view1 := &model.Block{
			Type:     model.TypeView,
			ParentID: "parent_id",
			BoardID:  boardID,
		}

		view2 := &model.Block{
			Type:     model.TypeView,
			ParentID: "parent_id",
			BoardID:  boardID,
		}

		board := &model.Board{ID: boardID}
		th.Store.EXPECT().GetBoard(boardID).Return(board, nil)
		th.Store.EXPECT().InsertBlock(view1, "user-id-1").Return(nil).Times(2)
		th.Store.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{}, nil).Times(2)

		// setting up mocks for limits
		fakeLicense := &mmModel.License{
			Features: &mmModel.Features{Cloud: mmModel.NewBool(true)},
		}
		th.Store.EXPECT().GetLicense().Return(fakeLicense).Times(2)

		th.Store.EXPECT().GetUsedCardsCount().Return(1, nil).Times(2)
		th.Store.EXPECT().GetCardLimitTimestamp().Return(int64(1), nil).Times(2)
		th.Store.EXPECT().GetBlocksWithParentAndType("test-board-id", "parent_id", "view").Return([]*model.Block{{}}, nil).Times(2)

		_, err := th.App.InsertBlocks([]*model.Block{view1, view2}, "user-id-1")
		require.Error(t, err)
	})
}

func TestGetBlocksEmptyBoardID(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	blocks, err := th.App.GetBlocks("", "", "")
	require.NoError(t, err)
	require.Len(t, blocks, 0)
}

func TestGetBlocksByParentAndType(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	expected := []*model.Block{{ID: "b1"}}

	th.Store.EXPECT().
		GetBlocksWithParentAndType(
			"board1",
			"parent1",
			"card",
		).
		Return(expected, nil)

	blocks, err := th.App.GetBlocks(
		"board1",
		"parent1",
		"card",
	)

	require.NoError(t, err)
	require.Equal(t, expected, blocks)
}

func TestGetBlocksByType(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	expected := []*model.Block{{ID: "b1"}}

	th.Store.EXPECT().
		GetBlocksWithType(
			"board1",
			"card",
		).
		Return(expected, nil)

	blocks, err := th.App.GetBlocks(
		"board1",
		"",
		"card",
	)

	require.NoError(t, err)
	require.Equal(t, expected, blocks)
}

func TestGetBlocksByParent(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	expected := []*model.Block{{ID: "b1"}}

	th.Store.EXPECT().
		GetBlocksWithParent(
			"board1",
			"parent1",
		).
		Return(expected, nil)

	blocks, err := th.App.GetBlocks(
		"board1",
		"parent1",
		"",
	)

	require.NoError(t, err)
	require.Equal(t, expected, blocks)
}

func TestGetBlockByID(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	expected := &model.Block{ID: "block1"}

	th.Store.EXPECT().
		GetBlock("block1").
		Return(expected, nil)

	block, err := th.App.GetBlockByID("block1")

	require.NoError(t, err)
	require.Equal(t, expected, block)
}

func TestGetLastBlockHistoryEntry(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	block := &model.Block{ID: "block1"}

	th.Store.EXPECT().
		GetBlockHistory(
			"block1",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return([]*model.Block{block}, nil)

	result, err := th.App.GetLastBlockHistoryEntry("block1")

	require.NoError(t, err)
	require.Equal(t, block, result)
}

func TestGetLastBlockHistoryEntryEmpty(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	th.Store.EXPECT().
		GetBlockHistory(
			"block1",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return([]*model.Block{}, nil)

	result, err := th.App.GetLastBlockHistoryEntry("block1")

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestGetBlockCountsByType(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	counts := map[string]int64{
		"card": 10,
	}

	th.Store.EXPECT().
		GetBlockCountsByType().
		Return(counts, nil)

	result, err := th.App.GetBlockCountsByType()

	require.NoError(t, err)
	require.Equal(t, counts, result)
}

func TestGetBlocksForBoard(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	expected := []*model.Block{
		{ID: "b1"},
	}

	th.Store.EXPECT().
		GetBlocksForBoard("board1").
		Return(expected, nil)

	result, err := th.App.GetBlocksForBoard("board1")

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestInsertBlocksMultipleBoards(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	blocks := []*model.Block{
		{BoardID: "board1"},
		{BoardID: "board2"},
	}

	_, err := th.App.InsertBlocks(
		blocks,
		"user-id",
	)

	require.ErrorIs(
		t,
		err,
		ErrBlocksFromMultipleBoards,
	)
}

func TestInsertBlocksEmpty(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	blocks, err := th.App.InsertBlocks(
		[]*model.Block{},
		"user-id",
	)

	require.NoError(t, err)
	require.Len(t, blocks, 0)
}

func TestDeleteBlockGetBlockError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	th.Store.EXPECT().
		GetBlock("block-id").
		Return(nil, errors.New("db error"))

	err := th.App.DeleteBlock(
		"block-id",
		"user-id",
	)

	require.Error(t, err)
}

func TestUndeleteBlockHistoryEmpty(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	th.Store.EXPECT().
		GetBlockHistory(
			"block-id",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return([]*model.Block{}, nil)

	block, err := th.App.UndeleteBlock(
		"block-id",
		"user-id",
	)

	require.NoError(t, err)
	require.Nil(t, block)
}

func TestGetLastBlockHistoryEntryError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	th.Store.EXPECT().
		GetBlockHistory(
			"block1",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return(nil, errors.New("history error"))

	block, err := th.App.GetLastBlockHistoryEntry("block1")

	require.Error(t, err)
	require.Nil(t, block)
}



func TestInsertBlockBoardError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	block := &model.Block{
		BoardID: "board1",
	}

	th.Store.EXPECT().
		GetBoard("board1").
		Return(nil, errors.New("board error"))

	err := th.App.InsertBlock(
		block,
		"user-id",
	)

	require.Error(t, err)
}

func TestInsertBlocksGetBoardError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	blocks := []*model.Block{
		{BoardID: "board1"},
	}

	th.Store.EXPECT().
		GetBoard("board1").
		Return(nil, errors.New("board error"))

	_, err := th.App.InsertBlocks(
		blocks,
		"user-id",
	)

	require.Error(t, err)
}

func TestPatchBlockGetBlockError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	patch := &model.BlockPatch{}

	th.Store.EXPECT().
		GetBlock("block1").
		Return(nil, errors.New("error"))

	block, err := th.App.PatchBlock(
		"block1",
		patch,
		"user-id",
	)

	require.Error(t, err)
	require.Nil(t, block)
}

func TestPatchBlockBoardError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	oldBlock := &model.Block{
		ID:      "block1",
		BoardID: "board1",
	}

	patch := &model.BlockPatch{}

	th.Store.EXPECT().
		GetBlock("block1").
		Return(oldBlock, nil)

	th.Store.EXPECT().
		GetBoard("board1").
		Return(nil, errors.New("board error"))

	block, err := th.App.PatchBlock(
		"block1",
		patch,
		"user-id",
	)

	require.Error(t, err)
	require.Nil(t, block)
}

func TestPatchBlocksPatchError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	batch := model.BlockPatchBatch{
		BlockIDs: []string{"block1"},
	}

	th.Store.EXPECT().
		GetBlocksByIDs([]string{"block1"}).
		Return([]*model.Block{
			{ID: "block1"},
		}, nil)

	th.Store.EXPECT().
		PatchBlocks(
			gomock.Any(),
			"user-id",
		).
		Return(errors.New("patch error"))

	err := th.App.PatchBlocks(
		"team-id",
		&batch,
		"user-id",
	)

	require.Error(t, err)
}

func TestUndeleteBlockHistoryError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	th.Store.EXPECT().
		GetBlockHistory(
			"block-id",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return(nil, errors.New("history error"))

	block, err := th.App.UndeleteBlock(
		"block-id",
		"user-id",
	)

	require.Error(t, err)
	require.Nil(t, block)
}

func TestUndeleteBlockBoardError(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	block := &model.Block{
		ID:      "block-id",
		BoardID: "board1",
	}

	th.Store.EXPECT().
		GetBlockHistory(
			"block-id",
			model.QueryBlockHistoryOptions{
				Limit:      1,
				Descending: true,
			},
		).
		Return([]*model.Block{block}, nil)

	th.Store.EXPECT().
		UndeleteBlock(
			"block-id",
			"user-id",
		).
		Return(nil)

	th.Store.EXPECT().
		GetBlock("block-id").
		Return(block, nil)

	th.Store.EXPECT().
		GetBoard("board1").
		Return(nil, errors.New("board error"))

	result, err := th.App.UndeleteBlock(
		"block-id",
		"user-id",
	)

	require.Error(t, err)
	require.Nil(t, result)
}