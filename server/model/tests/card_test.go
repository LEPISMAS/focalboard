package tests

import (
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCardErrors(t *testing.T) {
	err1 := model.NewErrInvalidCard("some error")
	assert.Equal(t, "invalid card, some error", err1.Error())

	// ErrInvalidFieldType is not exported directly, but we can verify its error string indirectly
	// or check that Block2Card returns an error containing the field name.
}

func TestCardPopulateAndValid(t *testing.T) {
	t.Run("Populate defaults", func(t *testing.T) {
		c := &model.Card{}
		c.Populate()
		assert.NotEmpty(t, c.ID)
		assert.NotNil(t, c.ContentOrder)
		assert.NotNil(t, c.Properties)
		assert.NotZero(t, c.CreateAt)
		assert.NotZero(t, c.UpdateAt)
	})

	t.Run("PopulateWithBoardID", func(t *testing.T) {
		c := &model.Card{}
		c.PopulateWithBoardID("board123")
		assert.Equal(t, "board123", c.BoardID)
		assert.NotEmpty(t, c.ID)
	})

	t.Run("CheckValid - missing ID", func(t *testing.T) {
		c := &model.Card{BoardID: "b1", ContentOrder: []string{}, Properties: map[string]any{}, CreateAt: 1, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - missing BoardID", func(t *testing.T) {
		c := &model.Card{ID: "c1", ContentOrder: []string{}, Properties: map[string]any{}, CreateAt: 1, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - missing ContentOrder", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", Properties: map[string]any{}, CreateAt: 1, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - invalid Icon", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", ContentOrder: []string{}, Icon: "ab", Properties: map[string]any{}, CreateAt: 1, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - missing Properties", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", ContentOrder: []string{}, CreateAt: 1, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - missing CreateAt", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", ContentOrder: []string{}, Properties: map[string]any{}, UpdateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - missing UpdateAt", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", ContentOrder: []string{}, Properties: map[string]any{}, CreateAt: 1}
		assert.Error(t, c.CheckValid())
	})

	t.Run("CheckValid - valid card", func(t *testing.T) {
		c := &model.Card{ID: "c1", BoardID: "b1", ContentOrder: []string{}, Icon: "🐙", Properties: map[string]any{}, CreateAt: 1, UpdateAt: 1}
		assert.NoError(t, c.CheckValid())
	})
}

func TestCardPatch(t *testing.T) {
	t.Run("Patch values", func(t *testing.T) {
		c := &model.Card{
			ID:           "c1",
			Title:        "Old Title",
			ContentOrder: []string{"co1"},
			Icon:         "🐶",
			Properties: map[string]any{
				"keep":   "val1",
				"change": "val2",
			},
		}

		newTitle := "New Title"
		newContentOrder := []string{"co2"}
		newIcon := "🐱"

		patch := &model.CardPatch{
			Title:        &newTitle,
			ContentOrder: &newContentOrder,
			Icon:         &newIcon,
			UpdatedProperties: map[string]any{
				"change": "updated",
				"new":    "newval",
			},
		}

		patched := patch.Patch(c)
		assert.Equal(t, "New Title", patched.Title)
		assert.Equal(t, []string{"co2"}, patched.ContentOrder)
		assert.Equal(t, "🐱", patched.Icon)
		assert.Equal(t, "val1", patched.Properties["keep"])
		assert.Equal(t, "updated", patched.Properties["change"])
		assert.Equal(t, "newval", patched.Properties["new"])
	})

	t.Run("Patch values - nil properties", func(t *testing.T) {
		c := &model.Card{
			ID:         "c1",
			Properties: nil,
		}
		patch := &model.CardPatch{
			UpdatedProperties: map[string]any{
				"new": "val",
			},
		}
		patched := patch.Patch(c)
		assert.NotNil(t, patched.Properties)
		assert.Equal(t, "val", patched.Properties["new"])
	})

	t.Run("CheckValid", func(t *testing.T) {
		p := &model.CardPatch{}
		assert.NoError(t, p.CheckValid())

		badIcon := "abc"
		p.Icon = &badIcon
		assert.Error(t, p.CheckValid())
	})
}

func TestCard2Block(t *testing.T) {
	c := &model.Card{
		ID:           "c1",
		BoardID:      "b1",
		CreatedBy:    "u1",
		ModifiedBy:   "u2",
		Title:        "Card Title",
		ContentOrder: []string{"co1"},
		Icon:         "🦊",
		IsTemplate:   true,
		Properties: map[string]any{
			"prop": "val",
		},
		CreateAt: 10,
		UpdateAt: 20,
		DeleteAt: 30,
	}

	b := model.Card2Block(c)
	assert.NotNil(t, b)
	assert.Equal(t, "c1", b.ID)
	assert.Equal(t, "b1", b.BoardID)
	assert.Equal(t, "b1", b.ParentID)
	assert.Equal(t, "u1", b.CreatedBy)
	assert.Equal(t, "u2", b.ModifiedBy)
	assert.Equal(t, model.BlockType(model.TypeCard), b.Type)
	assert.Equal(t, "Card Title", b.Title)
	assert.Equal(t, int64(10), b.CreateAt)
	assert.Equal(t, int64(20), b.UpdateAt)
	assert.Equal(t, int64(30), b.DeleteAt)

	assert.Equal(t, []string{"co1"}, b.Fields["contentOrder"])
	assert.Equal(t, "🦊", b.Fields["icon"])
	assert.Equal(t, true, b.Fields["isTemplate"])
	assert.Equal(t, "val", b.Fields["properties"].(map[string]any)["prop"])
}

func TestBlock2Card(t *testing.T) {
	t.Run("not a card block", func(t *testing.T) {
		b := &model.Block{Type: model.BlockType(model.TypeBoard)}
		c, err := model.Block2Card(b)
		assert.ErrorIs(t, err, model.ErrNotCardBlock)
		assert.Nil(t, c)
	})

	t.Run("valid conversion - string slice", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"contentOrder": []string{"co1"},
				"icon":         "🦁",
				"isTemplate":   true,
				"properties": map[string]any{
					"prop": "val",
				},
			},
		}

		c, err := model.Block2Card(b)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, "c1", c.ID)
		assert.Equal(t, []string{"co1"}, c.ContentOrder)
		assert.Equal(t, "🦁", c.Icon)
		assert.True(t, c.IsTemplate)
		assert.Equal(t, "val", c.Properties["prop"])
	})

	t.Run("valid conversion - any slice", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"contentOrder": []any{"co1", "co2"},
			},
		}

		c, err := model.Block2Card(b)
		assert.NoError(t, err)
		assert.NotNil(t, c)
		assert.Equal(t, []string{"co1", "co2"}, c.ContentOrder)
	})

	t.Run("invalid contentOrder item type", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"contentOrder": []any{123},
			},
		}

		c, err := model.Block2Card(b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contentOrder item")
		assert.Nil(t, c)
	})

	t.Run("invalid contentOrder type", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"contentOrder": "not-a-slice",
			},
		}

		c, err := model.Block2Card(b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contentOrder")
		assert.Nil(t, c)
	})

	t.Run("invalid icon type", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"icon": 123,
			},
		}

		c, err := model.Block2Card(b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "icon")
		assert.Nil(t, c)
	})

	t.Run("invalid isTemplate type", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"isTemplate": "yes",
			},
		}

		c, err := model.Block2Card(b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "isTemplate")
		assert.Nil(t, c)
	})

	t.Run("invalid properties type", func(t *testing.T) {
		b := &model.Block{
			ID:      "c1",
			Type:    model.BlockType(model.TypeCard),
			BoardID: "b1",
			Fields: map[string]interface{}{
				"properties": "not-a-map",
			},
		}

		c, err := model.Block2Card(b)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "properties")
		assert.Nil(t, c)
	})
}

func TestCardPatch2BlockPatch(t *testing.T) {
	t.Run("invalid patch", func(t *testing.T) {
		badIcon := "abc"
		cp := &model.CardPatch{Icon: &badIcon}
		bp, err := model.CardPatch2BlockPatch(cp)
		assert.Error(t, err)
		assert.Nil(t, bp)
	})

	t.Run("valid patch conversion", func(t *testing.T) {
		title := "New Title"
		co := []string{"co1"}
		icon := "🦄"
		cp := &model.CardPatch{
			Title:        &title,
			ContentOrder: &co,
			Icon:         &icon,
			UpdatedProperties: map[string]any{
				"p": "v",
			},
		}

		bp, err := model.CardPatch2BlockPatch(cp)
		require.NoError(t, err)
		require.NotNil(t, bp)

		assert.Equal(t, "New Title", *bp.Title)
		assert.Equal(t, &co, bp.UpdatedFields["contentOrder"])
		assert.Equal(t, &icon, bp.UpdatedFields["icon"])
		assert.Equal(t, cp.UpdatedProperties, bp.UpdatedFields["properties"])
	})
}
