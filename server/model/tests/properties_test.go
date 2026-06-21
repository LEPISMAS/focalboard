package tests

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropDefGetValue(t *testing.T) {
	t.Run("select type", func(t *testing.T) {
		pd := model.PropDef{
			ID:   "prop1",
			Type: "select",
			Options: map[string]model.PropDefOption{
				"opt1": {ID: "opt1", Value: "Option One"},
			},
		}

		// non-string value
		val, err := pd.GetValue(123, nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// option not found
		val, err = pd.GetValue("opt2", nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValue)
		assert.Empty(t, val)

		// option found
		val, err = pd.GetValue("opt1", nil)
		assert.NoError(t, err)
		assert.Equal(t, "OPTION ONE", val)
	})

	t.Run("date type", func(t *testing.T) {
		pd := model.PropDef{
			ID:   "prop2",
			Type: "date",
		}

		// non-string value
		val, err := pd.GetValue(123, nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// valid date
		val, err = pd.GetValue(`{"from":1642161600000}`, nil)
		assert.NoError(t, err)
		assert.Contains(t, val, "2022")
	})

	t.Run("person type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockResolver := mocks.NewMockPropValueResolver(ctrl)

		pd := model.PropDef{
			ID:   "prop3",
			Type: "person",
		}

		// non-string value
		val, err := pd.GetValue(123, mockResolver)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// resolver nil
		val, err = pd.GetValue("user1", nil)
		assert.NoError(t, err)
		assert.Equal(t, "user1", val)

		// resolver returns error
		mockResolver.EXPECT().GetUserByID("user1").Return(nil, errors.New("db error"))
		val, err = pd.GetValue("user1", mockResolver)
		assert.Error(t, err)
		assert.Empty(t, val)

		// resolver returns nil user
		mockResolver.EXPECT().GetUserByID("user1").Return(nil, nil)
		val, err = pd.GetValue("user1", mockResolver)
		assert.NoError(t, err)
		assert.Equal(t, "user1", val)

		// resolver returns valid user
		mockResolver.EXPECT().GetUserByID("user1").Return(&model.User{Username: "john_doe"}, nil)
		val, err = pd.GetValue("user1", mockResolver)
		assert.NoError(t, err)
		assert.Equal(t, "john_doe", val)
	})

	t.Run("multiPerson type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockResolver := mocks.NewMockPropValueResolver(ctrl)

		pd := model.PropDef{
			ID:   "prop4",
			Type: "multiPerson",
		}

		// non-slice value
		val, err := pd.GetValue("not-a-slice", mockResolver)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// resolver nil
		val, err = pd.GetValue([]interface{}{"user1", "user2"}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "[user1 user2]", val)

		// resolver returns error for first user
		mockResolver.EXPECT().GetUserByID("user1").Return(nil, errors.New("db error"))
		val, err = pd.GetValue([]interface{}{"user1", "user2"}, mockResolver)
		assert.Error(t, err)
		assert.Empty(t, val)

		// resolver returns nil for first user, and valid user for second
		mockResolver.EXPECT().GetUserByID("user1").Return(nil, nil)
		mockResolver.EXPECT().GetUserByID("user2").Return(&model.User{Username: "jane_doe"}, nil)
		val, err = pd.GetValue([]interface{}{"user1", "user2"}, mockResolver)
		assert.NoError(t, err)
		assert.Equal(t, "user1, jane_doe", val)
	})

	t.Run("multiSelect type", func(t *testing.T) {
		pd := model.PropDef{
			ID:   "prop5",
			Type: "multiSelect",
			Options: map[string]model.PropDefOption{
				"opt1": {ID: "opt1", Value: "Red"},
				"opt2": {ID: "opt2", Value: "Green"},
			},
		}

		// non-slice value
		val, err := pd.GetValue("not-a-slice", nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// slice with non-string elements
		val, err = pd.GetValue([]interface{}{"opt1", 123}, nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValueType)
		assert.Empty(t, val)

		// option not found
		val, err = pd.GetValue([]interface{}{"opt3"}, nil)
		assert.ErrorIs(t, err, model.ErrInvalidPropertyValue)
		assert.Empty(t, val)

		// valid options
		val, err = pd.GetValue([]interface{}{"opt1", "opt2"}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "RED, GREEN", val)
	})

	t.Run("other type", func(t *testing.T) {
		pd := model.PropDef{
			ID:   "prop6",
			Type: "text",
		}

		val, err := pd.GetValue(100.5, nil)
		assert.NoError(t, err)
		assert.Equal(t, "100.5", val)
	})
}

func TestParseDate(t *testing.T) {
	pd := model.PropDef{ID: "date-prop"}

	// Invalid JSON string
	val, err := pd.ParseDate("{invalid json}")
	assert.Error(t, err)
	assert.Equal(t, "{invalid json}", val)

	// Missing from key
	val, err = pd.ParseDate(`{"to":1642248000000}`)
	assert.ErrorIs(t, err, model.ErrInvalidDate)
	assert.Equal(t, `{"to":1642248000000}`, val)

	// Valid single date
	val, err = pd.ParseDate(`{"from":1642161600000}`)
	assert.NoError(t, err)
	assert.Contains(t, val, "January 14, 2022")

	// Valid date range
	val, err = pd.ParseDate(`{"from":1642161600000, "to":1642248000000}`)
	assert.NoError(t, err)
	assert.Contains(t, val, "January 14, 2022 -> January 15, 2022")
}

func TestParsePropertySchema(t *testing.T) {
	t.Run("empty board properties", func(t *testing.T) {
		b := &model.Board{}
		schema, err := model.ParsePropertySchema(b)
		assert.NoError(t, err)
		assert.Empty(t, schema)
	})

	t.Run("valid card properties", func(t *testing.T) {
		b := &model.Board{
			CardProperties: []map[string]interface{}{
				{
					"id":   "prop1",
					"name": "Priority",
					"type": "select",
					"options": []interface{}{
						map[string]interface{}{
							"id":    "opt1",
							"value": "High",
							"color": "red",
						},
					},
				},
			},
		}

		schema, err := model.ParsePropertySchema(b)
		assert.NoError(t, err)
		assert.Len(t, schema, 1)

		pd, ok := schema["prop1"]
		require.True(t, ok)
		assert.Equal(t, "prop1", pd.ID)
		assert.Equal(t, 0, pd.Index)
		assert.Equal(t, "Priority", pd.Name)
		assert.Equal(t, "select", pd.Type)
		assert.Len(t, pd.Options, 1)

		opt, ok := pd.Options["opt1"]
		require.True(t, ok)
		assert.Equal(t, "opt1", opt.ID)
		assert.Equal(t, 0, opt.Index)
		assert.Equal(t, "High", opt.Value)
		assert.Equal(t, "red", opt.Color)
	})

	t.Run("invalid options interface", func(t *testing.T) {
		b := &model.Board{
			CardProperties: []map[string]interface{}{
				{
					"id":      "prop1",
					"options": "not-a-slice",
				},
			},
		}

		schema, err := model.ParsePropertySchema(b)
		assert.ErrorIs(t, err, model.ErrInvalidPropSchema)
		assert.Nil(t, schema)
	})

	t.Run(
		"invalid option type",
		func(t *testing.T) {
			b := &model.Board{
				CardProperties: []map[string]interface{}{
					{
						"id": "prop1",
						"options": []interface{}{
							"not-a-map",
						},
					},
				},
			}

			schema, err := model.ParsePropertySchema(b)
			assert.ErrorIs(t, err, model.ErrInvalidPropSchema)
			assert.Nil(t, schema)
		},
	)

	t.Run(
		"non-string value in property map",
		func(t *testing.T) {
			b := &model.Board{
				CardProperties: []map[string]interface{}{
					{
						"id":   123, // not a string
						"name": "Prop",
					},
				},
			}

			schema, err := model.ParsePropertySchema(b)
			assert.NoError(t, err)
			assert.Len(t, schema, 1)
			assert.Contains(t, schema, "") // getMapString returned "" for ID
		},
	)
}

func TestParseProperties(t *testing.T) {
	schema := model.PropSchema{
		"prop1": {
			ID:    "prop1",
			Name:  "Priority",
			Type:  "select",
			Index: 0,
			Options: map[string]model.PropDefOption{
				"opt1": {ID: "opt1", Value: "High"},
			},
		},
	}

	t.Run("nil block", func(t *testing.T) {
		props, err := model.ParseProperties(nil, schema, nil)
		assert.NoError(t, err)
		assert.Empty(t, props)
	})

	t.Run("missing properties field", func(t *testing.T) {
		block := &model.Block{
			Fields: map[string]interface{}{},
		}
		props, err := model.ParseProperties(block, schema, nil)
		assert.NoError(t, err)
		assert.Empty(t, props)
	})

	t.Run("invalid properties field type", func(t *testing.T) {
		block := &model.Block{
			Fields: map[string]interface{}{
				"properties": "wrong-type",
			},
		}
		props, err := model.ParseProperties(block, schema, nil)
		assert.ErrorIs(t, err, model.ErrInvalidProperty)
		assert.Empty(t, props)
	})

	t.Run("empty properties map", func(t *testing.T) {
		block := &model.Block{
			Fields: map[string]interface{}{
				"properties": map[string]interface{}{},
			},
		}
		props, err := model.ParseProperties(block, schema, nil)
		assert.NoError(t, err)
		assert.Empty(t, props)
	})

	t.Run("valid properties parsing", func(t *testing.T) {
		block := &model.Block{
			Fields: map[string]interface{}{
				"properties": map[string]interface{}{
					"prop1": "opt1",
					"prop2": "custom-value",
				},
			},
		}

		props, err := model.ParseProperties(block, schema, nil)
		assert.NoError(t, err)
		assert.Len(t, props, 2)

		p1, ok := props["prop1"]
		require.True(t, ok)
		assert.Equal(t, "prop1", p1.ID)
		assert.Equal(t, "Priority", p1.Name)
		assert.Equal(t, "HIGH", p1.Value)
		assert.Equal(t, 0, p1.Index)

		p2, ok := props["prop2"]
		require.True(t, ok)
		assert.Equal(t, "prop2", p2.ID)
		assert.Equal(t, "prop2", p2.Name)
		assert.Equal(t, "custom-value", p2.Value)
		assert.Equal(t, 0, p2.Index)
	})

	t.Run("error parsing property value", func(t *testing.T) {
		block := &model.Block{
			Fields: map[string]interface{}{
				"properties": map[string]interface{}{
					"prop1": 123, // select expects string, will fail GetValue
				},
			},
		}

		props, err := model.ParseProperties(block, schema, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse property value")
		assert.Empty(t, props)
	})
}
