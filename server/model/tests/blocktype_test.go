package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/assert"
)

func TestBlockTypeString(t *testing.T) {
	bt := model.BlockType("custom-type")
	assert.Equal(t, "custom-type", bt.String())
}

func TestBlockTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected model.BlockType
		wantErr  bool
	}{
		{"board", model.TypeBoard, false},
		{"card", model.TypeCard, false},
		{"view", model.TypeView, false},
		{"text", model.TypeText, false},
		{"checkbox", model.TypeCheckbox, false},
		{"comment", model.TypeComment, false},
		{"image", model.TypeImage, false},
		{"attachment", model.TypeAttachment, false},
		{"divider", model.TypeDivider, false},
		{"BOARD", model.TypeBoard, false},
		{"invalid", model.TypeUnknown, true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := model.BlockTypeFromString(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, model.BlockType(model.TypeUnknown), got)
				var eibt model.ErrInvalidBlockType
				assert.True(t, errors.As(err, &eibt))
				assert.Equal(t, tc.input+" is an invalid block type.", err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestBlockType2IDType(t *testing.T) {
	tests := []struct {
		input    model.BlockType
		expected utils.IDType
	}{
		{model.TypeBoard, utils.IDTypeBoard},
		{model.TypeCard, utils.IDTypeCard},
		{model.TypeView, utils.IDTypeView},
		{model.TypeText, utils.IDTypeBlock},
		{model.TypeCheckbox, utils.IDTypeBlock},
		{model.TypeComment, utils.IDTypeBlock},
		{model.TypeDivider, utils.IDTypeBlock},
		{model.TypeImage, utils.IDTypeAttachment},
		{model.TypeAttachment, utils.IDTypeAttachment},
		{model.BlockType("unknown"), utils.IDTypeNone},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got := model.BlockType2IDType(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestIsErrInvalidBlockType(t *testing.T) {
	err1 := &model.ErrInvalidBlockType{Type: "badtype"}
	assert.True(t, model.IsErrInvalidBlockType(err1))

	errWrapped := fmt.Errorf("wrapped: %w", err1)
	assert.True(t, model.IsErrInvalidBlockType(errWrapped))

	errOther := errors.New("other error")
	assert.False(t, model.IsErrInvalidBlockType(errOther))
}
