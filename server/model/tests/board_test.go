package tests

import (
	"bytes"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/stretchr/testify/assert"
)

func TestBoardGetPropertyString(t *testing.T) {
	b := &model.Board{
		Properties: map[string]interface{}{
			"strProp":  "hello",
			"intProp":  123,
			"boolProp": true,
		},
	}

	// Case 1: property not found
	val, err := b.GetPropertyString("missing")
	assert.Error(t, err)
	assert.Empty(t, val)

	// Case 2: property exists but not a string
	val, err = b.GetPropertyString("intProp")
	assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
	assert.Empty(t, val)

	// Case 3: property is a string
	val, err = b.GetPropertyString("strProp")
	assert.NoError(t, err)
	assert.Equal(t, "hello", val)
}

func TestBoardJSONHelpers(t *testing.T) {
	t.Run("BoardFromJSON", func(t *testing.T) {
		data := []byte(`{"id":"b1","title":"My Board"}`)
		board := model.BoardFromJSON(bytes.NewReader(data))
		assert.NotNil(t, board)
		assert.Equal(t, "b1", board.ID)
		assert.Equal(t, "My Board", board.Title)
	})

	t.Run("BoardsFromJSON", func(t *testing.T) {
		data := []byte(`[{"id":"b1"},{"id":"b2"}]`)
		boards := model.BoardsFromJSON(bytes.NewReader(data))
		assert.Len(t, boards, 2)
		assert.Equal(t, "b1", boards[0].ID)
		assert.Equal(t, "b2", boards[1].ID)
	})

	t.Run("BoardMemberFromJSON", func(t *testing.T) {
		data := []byte(`{"boardId":"b1","userId":"u1","roles":"editor"}`)
		member := model.BoardMemberFromJSON(bytes.NewReader(data))
		assert.NotNil(t, member)
		assert.Equal(t, "b1", member.BoardID)
		assert.Equal(t, "u1", member.UserID)
		assert.Equal(t, "editor", member.Roles)
	})

	t.Run("BoardMembersFromJSON", func(t *testing.T) {
		data := []byte(`[{"userId":"u1"},{"userId":"u2"}]`)
		members := model.BoardMembersFromJSON(bytes.NewReader(data))
		assert.Len(t, members, 2)
		assert.Equal(t, "u1", members[0].UserID)
		assert.Equal(t, "u2", members[1].UserID)
	})

	t.Run("BoardMetadataFromJSON", func(t *testing.T) {
		data := []byte(`{"boardId":"b1","createdBy":"user1"}`)
		metadata := model.BoardMetadataFromJSON(bytes.NewReader(data))
		assert.NotNil(t, metadata)
		assert.Equal(t, "b1", metadata.BoardID)
		assert.Equal(t, "user1", metadata.CreatedBy)
	})
}

func TestBoardPatch(t *testing.T) {
	bTypeOpen := model.BoardTypeOpen
	bTypePrivate := model.BoardTypePrivate
	roleEditor := model.BoardRoleEditor
	roleAdmin := model.BoardRoleAdmin

	b := &model.Board{
		ID:              "b1",
		Type:            bTypeOpen,
		Title:           "Original Title",
		MinimumRole:     roleEditor,
		Description:     "Original Desc",
		Icon:            "original-icon",
		ShowDescription: false,
		ChannelID:       "chan1",
		Properties: map[string]interface{}{
			"keep":   "value1",
			"delete": "value2",
		},
		CardProperties: []map[string]interface{}{
			{"id": "prop1", "name": "Prop 1"},
			{"id": "prop2", "name": "Prop 2"},
			{"id": 1234, "name": "Bad Prop"}, // bad id type, will be skipped
		},
	}

	newTitle := "New Title"
	newDesc := "New Desc"
	newIcon := "new-icon"
	newShowDesc := true
	newChanID := "chan2"

	patch := &model.BoardPatch{
		Type:            &bTypePrivate,
		MinimumRole:     &roleAdmin,
		Title:           &newTitle,
		Description:     &newDesc,
		Icon:            &newIcon,
		ShowDescription: &newShowDesc,
		ChannelID:       &newChanID,
		UpdatedProperties: map[string]interface{}{
			"keep": "updated1",
			"new":  "newvalue",
		},
		DeletedProperties: []string{"delete"},
		UpdatedCardProperties: []map[string]interface{}{
			{"id": "prop1", "name": "Prop 1 Updated"},
			{"id": "prop3", "name": "Prop 3"},
			{"id": 5678, "name": "Bad Updated Prop"}, // bad id type, skipped
		},
		DeletedCardProperties: []string{"prop2"},
	}

	patched := patch.Patch(b)
	assert.Equal(t, bTypePrivate, patched.Type)
	assert.Equal(t, roleAdmin, patched.MinimumRole)
	assert.Equal(t, "New Title", patched.Title)
	assert.Equal(t, "New Desc", patched.Description)
	assert.Equal(t, "new-icon", patched.Icon)
	assert.True(t, patched.ShowDescription)
	assert.Equal(t, "chan2", patched.ChannelID)

	assert.Equal(t, "updated1", patched.Properties["keep"])
	assert.Equal(t, "newvalue", patched.Properties["new"])
	assert.NotContains(t, patched.Properties, "delete")

	assert.Len(t, patched.CardProperties, 2)
	assert.Equal(t, "prop1", patched.CardProperties[0]["id"])
	assert.Equal(t, "Prop 1 Updated", patched.CardProperties[0]["name"])
	assert.Equal(t, "prop3", patched.CardProperties[1]["id"])
	assert.Equal(t, "Prop 3", patched.CardProperties[1]["name"])
}

func TestBoardValidation(t *testing.T) {
	t.Run("IsBoardTypeValid", func(t *testing.T) {
		assert.True(t, model.IsBoardTypeValid(model.BoardTypeOpen))
		assert.True(t, model.IsBoardTypeValid(model.BoardTypePrivate))
		assert.False(t, model.IsBoardTypeValid(model.BoardType("invalid")))
	})

	t.Run("IsBoardMinimumRoleValid", func(t *testing.T) {
		assert.True(t, model.IsBoardMinimumRoleValid(model.BoardRoleNone))
		assert.True(t, model.IsBoardMinimumRoleValid(model.BoardRoleAdmin))
		assert.True(t, model.IsBoardMinimumRoleValid(model.BoardRoleEditor))
		assert.True(t, model.IsBoardMinimumRoleValid(model.BoardRoleCommenter))
		assert.True(t, model.IsBoardMinimumRoleValid(model.BoardRoleViewer))
		assert.False(t, model.IsBoardMinimumRoleValid(model.BoardRole("invalid")))
	})

	t.Run("BoardPatch.IsValid", func(t *testing.T) {
		badType := model.BoardType("bad")
		badRole := model.BoardRole("bad")
		goodType := model.BoardTypeOpen
		goodRole := model.BoardRoleEditor

		p1 := &model.BoardPatch{Type: &badType}
		assert.Error(t, p1.IsValid())

		p2 := &model.BoardPatch{MinimumRole: &badRole}
		assert.Error(t, p2.IsValid())

		p3 := &model.BoardPatch{Type: &goodType, MinimumRole: &goodRole}
		assert.NoError(t, p3.IsValid())
	})

	t.Run("Board.IsValid", func(t *testing.T) {
		b := &model.Board{
			TeamID:      "team1",
			Type:        model.BoardTypeOpen,
			MinimumRole: model.BoardRoleEditor,
		}
		assert.NoError(t, b.IsValid())

		// empty team ID
		b.TeamID = ""
		err := b.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "empty-team-id", err.Error())

		// invalid board type
		b.TeamID = "team1"
		b.Type = model.BoardType("invalid")
		err = b.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "invalid-board-type", err.Error())

		// invalid minimum role
		b.Type = model.BoardTypeOpen
		b.MinimumRole = model.BoardRole("invalid")
		err = b.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "invalid-board-minimum-role", err.Error())
	})
}

func TestBoardSearchFieldFromString(t *testing.T) {
	f, err := model.BoardSearchFieldFromString("title")
	assert.NoError(t, err)
	assert.Equal(t, model.BoardSearchFieldTitle, f)

	f, err = model.BoardSearchFieldFromString("property_name")
	assert.NoError(t, err)
	assert.Equal(t, model.BoardSearchFieldPropertyName, f)

	f, err = model.BoardSearchFieldFromString("invalid")
	assert.Error(t, err)
	assert.Equal(t, model.BoardSearchFieldNone, f)
}

func TestBoardsAndBlocks(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		bab := &model.BoardsAndBlocks{}

		// No boards
		err := bab.IsValid()
		assert.ErrorIs(t, err, model.ErrNoBoardsInBoardsAndBlocks)

		// No blocks
		bab.Boards = []*model.Board{{ID: "b1"}}
		err = bab.IsValid()
		assert.ErrorIs(t, err, model.ErrNoBlocksInBoardsAndBlocks)

		// Block doesn't belong to any board
		bab.Blocks = []*model.Block{{ID: "bl1", BoardID: "b2"}}
		err = bab.IsValid()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "block bl1 doesn't belong to any board")

		// Valid case
		bab.Blocks = []*model.Block{{ID: "bl1", BoardID: "b1"}}
		err = bab.IsValid()
		assert.NoError(t, err)
	})

	t.Run("NewDeleteBoardsAndBlocksFromBabs", func(t *testing.T) {
		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{{ID: "b1"}},
			Blocks: []*model.Block{{ID: "bl1", BoardID: "b1"}},
		}
		dbab := model.NewDeleteBoardsAndBlocksFromBabs(bab)
		assert.Len(t, dbab.Boards, 1)
		assert.Equal(t, "b1", dbab.Boards[0])
		assert.Len(t, dbab.Blocks, 1)
		assert.Equal(t, "bl1", dbab.Blocks[0])
	})

	t.Run("DeleteBoardsAndBlocks.IsValid", func(t *testing.T) {
		dbab := &model.DeleteBoardsAndBlocks{}
		err := dbab.IsValid()
		assert.ErrorIs(t, err, model.ErrNoBoardsInBoardsAndBlocks)

		dbab.Boards = []string{"b1"}
		err = dbab.IsValid()
		assert.NoError(t, err)
	})

	t.Run("PatchBoardsAndBlocks.IsValid", func(t *testing.T) {
		pbab := &model.PatchBoardsAndBlocks{}

		// No board IDs
		err := pbab.IsValid()
		assert.ErrorIs(t, err, model.ErrNoBoardsInBoardsAndBlocks)

		// Board IDs and Patches mismatch
		pbab.BoardIDs = []string{"b1"}
		pbab.BoardPatches = []*model.BoardPatch{}
		err = pbab.IsValid()
		assert.ErrorIs(t, err, model.ErrBoardIDsAndPatchesMissmatchInBoardsAndBlocks)

		// Block IDs and Patches mismatch
		pbab.BoardPatches = []*model.BoardPatch{{}}
		pbab.BlockIDs = []string{"bl1"}
		pbab.BlockPatches = []*model.BlockPatch{}
		err = pbab.IsValid()
		assert.ErrorIs(t, err, model.ErrBlockIDsAndPatchesMissmatchInBoardsAndBlocks)

		// Valid
		pbab.BlockPatches = []*model.BlockPatch{{}}
		err = pbab.IsValid()
		assert.NoError(t, err)
	})

	t.Run("GenerateBoardsAndBlocksIDs", func(t *testing.T) {
		bab := &model.BoardsAndBlocks{
			Boards: []*model.Board{{ID: "b1"}},
			Blocks: []*model.Block{{ID: "bl1", BoardID: "b1", Type: model.TypeCard}},
		}

		newBab, err := model.GenerateBoardsAndBlocksIDs(bab, &mlog.Logger{})
		assert.NoError(t, err)
		assert.NotNil(t, newBab)
		assert.Len(t, newBab.Boards, 1)
		assert.Len(t, newBab.Blocks, 1)

		assert.NotEqual(t, "b1", newBab.Boards[0].ID)
		assert.Equal(t, newBab.Boards[0].ID, newBab.Blocks[0].BoardID)
		assert.NotEqual(t, "bl1", newBab.Blocks[0].ID)

		// Error case (invalid babs)
		invalidBab := &model.BoardsAndBlocks{}
		newBab2, err := model.GenerateBoardsAndBlocksIDs(invalidBab, &mlog.Logger{})
		assert.Error(t, err)
		assert.Nil(t, newBab2)
	})

	t.Run("BoardsAndBlocksFromJSON", func(t *testing.T) {
		data := []byte(`{"boards":[{"id":"b1"}],"blocks":[{"id":"bl1"}]}`)
		bab := model.BoardsAndBlocksFromJSON(bytes.NewReader(data))
		assert.NotNil(t, bab)
		assert.Len(t, bab.Boards, 1)
		assert.Equal(t, "b1", bab.Boards[0].ID)
		assert.Len(t, bab.Blocks, 1)
		assert.Equal(t, "bl1", bab.Blocks[0].ID)
	})
}
