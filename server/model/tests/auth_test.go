package tests

import (
	"bytes"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthParamError(t *testing.T) {
	err := model.NewErrAuthParam("some message")
	assert.Equal(t, "some message", err.Error())
}

func TestLoginResponseFromJSON(t *testing.T) {
	t.Run("valid login response", func(t *testing.T) {
		data := []byte(`{"token":"my-session-token"}`)
		resp, err := model.LoginResponseFromJSON(bytes.NewReader(data))
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "my-session-token", resp.Token)
	})

	t.Run("invalid json", func(t *testing.T) {
		data := []byte(`{invalid`)
		resp, err := model.LoginResponseFromJSON(bytes.NewReader(data))
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestRegisterRequestIsValid(t *testing.T) {
	t.Run("valid registration request", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "user123",
			Email:    "user@example.com",
			Password: "securepassword",
			Token:    "token123",
		}
		assert.NoError(t, req.IsValid())
	})

	t.Run("empty username", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "   ",
			Email:    "user@example.com",
			Password: "securepassword",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "username is required", err.Error())
	})

	t.Run("empty email", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "user123",
			Email:    "",
			Password: "securepassword",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("invalid email format", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "user123",
			Email:    "invalid-email",
			Password: "securepassword",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "invalid email format", err.Error())
	})

	t.Run("empty password", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "user123",
			Email:    "user@example.com",
			Password: "",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "password is required", err.Error())
	})

	t.Run("password too short", func(t *testing.T) {
		req := &model.RegisterRequest{
			Username: "user123",
			Email:    "user@example.com",
			Password: "short",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "password must be at least 8 characters", err.Error())
	})
}

func TestChangePasswordRequestIsValid(t *testing.T) {
	t.Run("valid change password request", func(t *testing.T) {
		req := &model.ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "newpassword",
		}
		assert.NoError(t, req.IsValid())
	})

	t.Run("empty old password", func(t *testing.T) {
		req := &model.ChangePasswordRequest{
			OldPassword: "",
			NewPassword: "newpassword",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "old password is required", err.Error())
	})

	t.Run("empty new password", func(t *testing.T) {
		req := &model.ChangePasswordRequest{
			OldPassword: "oldpassword",
			NewPassword: "",
		}
		err := req.IsValid()
		assert.Error(t, err)
		assert.Equal(t, "new password is required", err.Error())
	})
}
