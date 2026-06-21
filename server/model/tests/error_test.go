package tests

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	pluginapi "github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
)

func TestErrorStructures(t *testing.T) {
	err1 := model.NewErrNotFound("entity1")
	assert.Equal(t, "{entity1} not found", err1.Error())

	err2 := model.NewErrNotAllFound("entity2", []string{"id1", "id2"})
	assert.Equal(t, "not all instances of {entity2} in {id1, id2} found", err2.Error())

	err3 := model.NewErrBadRequest("bad reason")
	assert.Equal(t, "bad reason", err3.Error())

	err4 := model.NewErrUnauthorized("unauthorized reason")
	assert.Equal(t, "unauthorized reason", err4.Error())

	err5 := model.NewErrPermission("permission reason")
	assert.Equal(t, "permission reason", err5.Error())

	err6 := model.NewErrForbidden("forbidden reason")
	assert.Equal(t, "forbidden reason", err6.Error())

	err7 := model.NewErrInvalidCategory("invalid category reason")
	assert.Equal(t, "invalid category reason", err7.Error())

	err8 := model.NewErrNotImplemented("not implemented reason")
	assert.Equal(t, "not implemented reason", err8.Error())
}

func TestIsErrHandlers(t *testing.T) {
	t.Run("IsErrBadRequest", func(t *testing.T) {
		assert.False(t, model.IsErrBadRequest(nil))
		assert.True(t, model.IsErrBadRequest(model.NewErrBadRequest("bad")))
		assert.True(t, model.IsErrBadRequest(model.ErrViewsLimitReached))
		assert.True(t, model.IsErrBadRequest(model.NewErrAuthParam("auth")))
		assert.True(t, model.IsErrBadRequest(model.NewErrInvalidCategory("cat")))
		assert.True(t, model.IsErrBadRequest(model.ErrBoardMemberIsLastAdmin))
		assert.True(t, model.IsErrBadRequest(model.ErrBoardIDMismatch))
		assert.True(t, model.IsErrBadRequest(model.ErrBlockTitleSizeLimitExceeded))
		assert.True(t, model.IsErrBadRequest(model.ErrBlockFieldsSizeLimitExceeded))
		assert.False(t, model.IsErrBadRequest(errors.New("other")))
	})

	t.Run("IsErrUnauthorized", func(t *testing.T) {
		assert.False(t, model.IsErrUnauthorized(nil))
		assert.True(t, model.IsErrUnauthorized(model.NewErrUnauthorized("unauth")))
		assert.False(t, model.IsErrUnauthorized(errors.New("other")))
	})

	t.Run("IsErrForbidden", func(t *testing.T) {
		assert.False(t, model.IsErrForbidden(nil))
		assert.True(t, model.IsErrForbidden(model.NewErrForbidden("forbidden")))
		assert.True(t, model.IsErrForbidden(model.NewErrPermission("permission")))
		assert.True(t, model.IsErrForbidden(model.ErrPatchUpdatesLimitedCards))
		assert.True(t, model.IsErrForbidden(model.ErrCategoryPermissionDenied))
		assert.False(t, model.IsErrForbidden(errors.New("other")))
	})

	t.Run("IsErrNotFound", func(t *testing.T) {
		assert.False(t, model.IsErrNotFound(nil))
		assert.True(t, model.IsErrNotFound(model.NewErrNotFound("entity")))
		assert.True(t, model.IsErrNotFound(model.NewErrNotAllFound("entity", []string{"id"})))
		assert.True(t, model.IsErrNotFound(sql.ErrNoRows))
		assert.True(t, model.IsErrNotFound(pluginapi.ErrNotFound))
		assert.True(t, model.IsErrNotFound(model.ErrCategoryDeleted))

		appErrNotFound := &mmModel.AppError{StatusCode: http.StatusNotFound}
		assert.True(t, model.IsErrNotFound(appErrNotFound))

		appErrInternal := &mmModel.AppError{StatusCode: http.StatusInternalServerError}
		assert.False(t, model.IsErrNotFound(appErrInternal))

		assert.False(t, model.IsErrNotFound(errors.New("other")))
	})

	t.Run("IsErrRequestEntityTooLarge", func(t *testing.T) {
		assert.True(t, model.IsErrRequestEntityTooLarge(model.ErrRequestEntityTooLarge))
		assert.False(t, model.IsErrRequestEntityTooLarge(errors.New("other")))
	})

	t.Run("IsErrNotImplemented", func(t *testing.T) {
		assert.False(t, model.IsErrNotImplemented(nil))
		assert.True(t, model.IsErrNotImplemented(model.NewErrNotImplemented("ni")))
		assert.True(t, model.IsErrNotImplemented(model.ErrInsufficientLicense))
		assert.False(t, model.IsErrNotImplemented(errors.New("other")))
	})
}
