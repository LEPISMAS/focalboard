package tests

import (
	"bytes"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/audit"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlocksFromJSON(t *testing.T) {
	data := []byte(`[{"id":"bl1","title":"Block 1"},{"id":"bl2","title":"Block 2"}]`)
	blocks := model.BlocksFromJSON(bytes.NewReader(data))
	assert.Len(t, blocks, 2)
	assert.Equal(t, "bl1", blocks[0].ID)
	assert.Equal(t, "Block 1", blocks[0].Title)
	assert.Equal(t, "bl2", blocks[1].ID)
	assert.Equal(t, "Block 2", blocks[1].Title)
}

func TestBlockIsValid(t *testing.T) {
	t.Run("empty board id", func(t *testing.T) {
		b := &model.Block{
			ID:      "bl1",
			BoardID: "",
		}
		err := b.IsValid()
		assert.ErrorIs(t, err, model.ErrBlockEmptyBoardID)
	})

	t.Run("title too long", func(t *testing.T) {
		titleTooLong := make([]byte, model.BlockTitleMaxBytes+10)
		for i := range titleTooLong {
			titleTooLong[i] = 'a'
		}
		b := &model.Block{
			ID:      "bl1",
			BoardID: "board1",
			Title:   string(titleTooLong),
		}
		err := b.IsValid()
		assert.ErrorIs(t, err, model.ErrBlockTitleSizeLimitExceeded)
	})

	t.Run("fields marshal error", func(t *testing.T) {
		b := &model.Block{
			ID:      "bl1",
			BoardID: "board1",
			Fields: map[string]interface{}{
				"invalid": make(chan int), // channels cannot be marshalled to JSON
			},
		}
		err := b.IsValid()
		assert.Error(t, err)
	})

	t.Run("fields too long", func(t *testing.T) {
		// we construct a very large string value in fields
		largeVal := make([]byte, model.BlockFieldsMaxRunes+10)
		for i := range largeVal {
			largeVal[i] = 'x'
		}
		b := &model.Block{
			ID:      "bl1",
			BoardID: "board1",
			Fields: map[string]interface{}{
				"large": string(largeVal),
			},
		}
		err := b.IsValid()
		assert.ErrorIs(t, err, model.ErrBlockFieldsSizeLimitExceeded)
	})

	t.Run("valid block", func(t *testing.T) {
		b := &model.Block{
			ID:      "bl1",
			BoardID: "board1",
			Title:   "Valid Title",
			Fields: map[string]interface{}{
				"key": "value",
			},
		}
		err := b.IsValid()
		assert.NoError(t, err)
	})
}

func TestBlockLogClone(t *testing.T) {
	b := &model.Block{
		ID:       "bl1",
		ParentID: "parent1",
		BoardID:  "board1",
		Type:     model.TypeCard,
		Title:    "Do not clone title",
	}

	clone := b.LogClone()
	assert.NotNil(t, clone)

	// it returns an anonymous struct with ID, ParentID, BoardID, Type
	// let's verify using reflection or marshalling
	val, ok := clone.(struct {
		ID       string
		ParentID string
		BoardID  string
		Type     model.BlockType
	})
	require.True(t, ok)
	assert.Equal(t, "bl1", val.ID)
	assert.Equal(t, "parent1", val.ParentID)
	assert.Equal(t, "board1", val.BoardID)
	assert.Equal(t, model.BlockType(model.TypeCard), val.Type)
}

func TestBlockPatch(t *testing.T) {
	parentID := "new-parent"
	schema := int64(2)
	bType := model.BlockType(model.TypeView)
	title := "New Title"

	b := &model.Block{
		ID:       "bl1",
		ParentID: "old-parent",
		Schema:   1,
		Type:     model.TypeCard,
		Title:    "Old Title",
		Fields: map[string]interface{}{
			"keep":   "val1",
			"delete": "val2",
		},
	}

	patch := &model.BlockPatch{
		ParentID: &parentID,
		Schema:   &schema,
		Type:     &bType,
		Title:    &title,
		UpdatedFields: map[string]interface{}{
			"keep": "updated1",
			"new":  "newvalue",
		},
		DeletedFields: []string{"delete"},
	}

	patched := patch.Patch(b)
	assert.Equal(t, "new-parent", patched.ParentID)
	assert.Equal(t, int64(2), patched.Schema)
	assert.Equal(t, model.BlockType(model.TypeView), patched.Type)
	assert.Equal(t, "New Title", patched.Title)
	assert.Equal(t, "updated1", patched.Fields["keep"])
	assert.Equal(t, "newvalue", patched.Fields["new"])
	assert.NotContains(t, patched.Fields, "delete")
}

func TestBlockStampModificationMetadata(t *testing.T) {
	t.Run("single user", func(t *testing.T) {
		b := &model.Block{ID: "bl1"}
		blocks := []*model.Block{b}
		auditRec := &audit.Record{}

		model.StampModificationMetadata(model.SingleUser, blocks, auditRec)
		assert.Equal(t, "", b.ModifiedBy)
		assert.NotEmpty(t, b.UpdateAt)
		assert.Len(t, auditRec.Meta, 1)
		assert.Equal(t, "block_0", auditRec.Meta[0].K)
	})

	t.Run("regular user", func(t *testing.T) {
		b := &model.Block{ID: "bl1"}
		blocks := []*model.Block{b}

		model.StampModificationMetadata("user123", blocks, nil)
		assert.Equal(t, "user123", b.ModifiedBy)
		assert.NotEmpty(t, b.UpdateAt)
	})
}

func TestBlockShouldBeLimited(t *testing.T) {
	b1 := &model.Block{Type: model.TypeCard, UpdateAt: 100}
	assert.True(t, b1.ShouldBeLimited(200))
	assert.False(t, b1.ShouldBeLimited(50))

	b2 := &model.Block{Type: model.TypeView, UpdateAt: 100}
	assert.False(t, b2.ShouldBeLimited(200))
}

func TestBlockGetLimited(t *testing.T) {
	b := &model.Block{
		ID:          "bl1",
		ParentID:    "parent1",
		BoardID:     "board1",
		Schema:      1,
		Type:        model.TypeCard,
		Title:       "My Card",
		CreateAt:    10,
		UpdateAt:    20,
		DeleteAt:    30,
		WorkspaceID: "ws1",
		Fields: map[string]interface{}{
			"icon": "🎨",
			"desc": "secret description",
		},
	}

	lim := b.GetLimited()
	assert.Equal(t, "bl1", lim.ID)
	assert.Equal(t, "parent1", lim.ParentID)
	assert.Equal(t, "board1", lim.BoardID)
	assert.Equal(t, int64(1), lim.Schema)
	assert.Equal(t, model.BlockType(model.TypeCard), lim.Type)
	assert.Equal(t, "My Card", lim.Title)
	assert.Equal(t, int64(10), lim.CreateAt)
	assert.Equal(t, int64(20), lim.UpdateAt)
	assert.Equal(t, int64(30), lim.DeleteAt)
	assert.Equal(t, "ws1", lim.WorkspaceID)
	assert.True(t, lim.Limited)

	assert.Contains(t, lim.Fields, "icon")
	assert.Equal(t, "🎨", lim.Fields["icon"])
	assert.NotContains(t, lim.Fields, "desc")
}

func getTestLogger() mlog.LoggerIFace {
	logger, _ := mlog.NewLogger()
	cfgJSON := `{"def":{"type":"console","options":{"out":"stdout"},"format":"plain","levels":[{"id":3,"name":"warn"}]}}`
	_ = logger.Configure("", cfgJSON, nil)
	return logger
}

func TestGenerateBlockIDsEdgeCases(t *testing.T) {
	logger := getTestLogger()

	t.Run("contentOrder invalid type", func(t *testing.T) {
		blocks := []*model.Block{
			{
				ID:   "b1",
				Type: model.TypeCard,
				Fields: map[string]interface{}{
					"contentOrder": "not-a-slice",
				},
			},
		}
		res := model.GenerateBlockIDs(blocks, logger)
		assert.Len(t, res, 1)
	})

	t.Run("contentOrder with slice of slices", func(t *testing.T) {
		blocks := []*model.Block{
			{
				ID:   "b1",
				Type: model.TypeCard,
				Fields: map[string]interface{}{
					"contentOrder": []interface{}{
						[]interface{}{"b2", "b3"},
						"b4",
					},
				},
			},
			{ID: "b2", Type: model.TypeCard},
			{ID: "b3", Type: model.TypeCard},
			{ID: "b4", Type: model.TypeCard},
		}
		res := model.GenerateBlockIDs(blocks, logger)
		assert.Len(t, res, 4)

		b1Patched := res[0]
		contentOrder := b1Patched.Fields["contentOrder"].([]interface{})
		subOrder := contentOrder[0].([]interface{})
		assert.NotEqual(t, "b2", subOrder[0].(string))
		assert.NotEqual(t, "b3", subOrder[1].(string))
		assert.NotEqual(t, "b4", contentOrder[1].(string))
	})

	t.Run("defaultTemplateId invalid type", func(t *testing.T) {
		blocks := []*model.Block{
			{
				ID:   "b1",
				Type: model.TypeCard,
				Fields: map[string]interface{}{
					"defaultTemplateId": 12345, // not a string
				},
			},
		}
		res := model.GenerateBlockIDs(blocks, logger)
		assert.Len(t, res, 1)
	})

	t.Run("cardOrder valid and invalid type", func(t *testing.T) {
		blocks := []*model.Block{
			{
				ID:       "b1",
				Type:     model.TypeCard,
				ParentID: "b2", // ensures b2 is added to referenceIDs
				Fields: map[string]interface{}{
					"cardOrder": []interface{}{"b2"},
				},
			},
			{ID: "b2", Type: model.TypeCard},
			{
				ID:   "b3",
				Type: model.TypeCard,
				Fields: map[string]interface{}{
					"cardOrder": "invalid-type",
				},
			},
		}
		res := model.GenerateBlockIDs(blocks, logger)
		assert.Len(t, res, 3)
		assert.NotEqual(t, "b2", res[0].Fields["cardOrder"].([]interface{})[0].(string))
	})
}
