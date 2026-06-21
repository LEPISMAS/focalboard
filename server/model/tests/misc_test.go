package tests

import (
	"bytes"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	logger := getTestLogger()
	// Should not panic
	model.LogServerInfo(logger)
}

func TestTeam(t *testing.T) {
	t.Run("TeamFromJSON", func(t *testing.T) {
		data := []byte(`{"id":"team1","title":"Team 1"}`)
		team := model.TeamFromJSON(bytes.NewReader(data))
		assert.NotNil(t, team)
		assert.Equal(t, "team1", team.ID)
		assert.Equal(t, "Team 1", team.Title)
	})

	t.Run("TeamsFromJSON", func(t *testing.T) {
		data := []byte(`[{"id":"team1"},{"id":"team2"}]`)
		teams := model.TeamsFromJSON(bytes.NewReader(data))
		assert.Len(t, teams, 2)
		assert.Equal(t, "team1", teams[0].ID)
		assert.Equal(t, "team2", teams[1].ID)
	})
}

func TestFile(t *testing.T) {
	fi := model.NewFileInfo("test.png")
	assert.NotNil(t, fi)
	assert.Equal(t, "test.png", fi.Name)
	assert.Equal(t, ".png", fi.Extension)
}

func TestSharing(t *testing.T) {
	t.Run("SharingFromJSON", func(t *testing.T) {
		data := []byte(`{"id":"board1","enabled":true}`)
		sh := model.SharingFromJSON(bytes.NewReader(data))
		assert.Equal(t, "board1", sh.ID)
		assert.True(t, sh.Enabled)
	})
}

func TestCategory(t *testing.T) {
	t.Run("Hydrate and IsValid", func(t *testing.T) {
		c := &model.Category{
			SortOrder: -1,
		}
		c.Hydrate()
		assert.NotEmpty(t, c.ID)
		assert.NotZero(t, c.CreateAt)
		assert.NotZero(t, c.UpdateAt)
		assert.Equal(t, 0, c.SortOrder)
		assert.Equal(t, model.CategoryTypeCustom, c.Type)

		assert.Error(t, c.IsValid())

		c.Name = "My Category"
		c.UserID = "user1"
		c.TeamID = "team1"
		assert.NoError(t, c.IsValid())

		// Invalid Category type
		c.Type = "unknown"
		assert.Error(t, c.IsValid())

		c.Type = model.CategoryTypeSystem
		assert.NoError(t, c.IsValid())
	})

	t.Run("CategoryFromJSON", func(t *testing.T) {
		data := []byte(`{"id":"cat1","name":"Category 1"}`)
		c := model.CategoryFromJSON(bytes.NewReader(data))
		assert.NotNil(t, c)
		assert.Equal(t, "cat1", c.ID)
		assert.Equal(t, "Category 1", c.Name)
	})
}
