package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/auth"
	"github.com/mattermost/focalboard/server/utils"
)

func TestAuthEndpoints(t *testing.T) {
	th, tearDown := SetupTestAPIHelper(t, "") // No singleUserToken
	defer tearDown()

	t.Run("POST /login successful", func(t *testing.T) {
		hashed := auth.HashPassword("password123")
		user := &model.User{
			ID:          "user1",
			Username:    "testuser",
			Password:    hashed,
			AuthService: "native",
		}
		th.Store.EXPECT().GetUserByUsername("testuser").Return(user, nil)
		th.Store.EXPECT().CreateSession(gomock.Any()).Return(nil)

		body := `{"username": "testuser", "password": "password123", "type": "normal"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/login", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "token")
	})

	t.Run("POST /login incorrect password", func(t *testing.T) {
		hashed := auth.HashPassword("password123")
		user := &model.User{
			ID:          "user1",
			Username:    "testuser",
			Password:    hashed,
			AuthService: "native",
		}
		th.Store.EXPECT().GetUserByUsername("testuser").Return(user, nil)

		body := `{"username": "testuser", "password": "wrongpassword", "type": "normal"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/login", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
		require.Contains(t, resp.Body.String(), "incorrect login")
	})

	t.Run("POST /register successful", func(t *testing.T) {
		// Mock GetRootTeam to get signup token
		th.Store.EXPECT().GetTeam("0").Return(&model.Team{ID: "0", SignupToken: "signuptoken"}, nil)
		th.Store.EXPECT().GetUserByUsername("newuser").Return(nil, model.NewErrNotFound("user"))
		th.Store.EXPECT().GetUserByEmail("newuser@example.com").Return(nil, model.NewErrNotFound("user"))
		th.Store.EXPECT().CreateUser(gomock.Any()).Return(&model.User{ID: "newuser1"}, nil)

		body := `{"username": "newuser", "email": "newuser@example.com", "password": "securepassword", "token": "signuptoken"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/register", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /logout successful", func(t *testing.T) {
		// Setup mock session with fresh UpdateAt
		th.Store.EXPECT().GetSession("test-token", gomock.Any()).Return(&model.Session{
			ID:          "session1",
			Token:       "test-token",
			UserID:      "user1",
			AuthService: "native",
			UpdateAt:    utils.GetMillis(),
		}, nil)
		th.Store.EXPECT().DeleteSession("session1").Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer test-token")

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /users/{userID}/changepassword successful", func(t *testing.T) {
		th.Store.EXPECT().GetSession("test-token", gomock.Any()).Return(&model.Session{
			ID:          "session1",
			Token:       "test-token",
			UserID:      "user1",
			AuthService: "native",
			UpdateAt:    utils.GetMillis(),
		}, nil)
		hashed := auth.HashPassword("oldpassword")
		th.Store.EXPECT().GetUserByID("user1").Return(&model.User{ID: "user1", Password: hashed}, nil)
		th.Store.EXPECT().UpdateUserPasswordByID("user1", gomock.Any()).Return(nil)

		body := `{"oldPassword": "oldpassword", "newPassword": "newsecurepassword"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/users/user1/changepassword", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer test-token")

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("POST /login not permitted in plugin mode", func(t *testing.T) {
		th.API.MattermostAuth = true
		defer func() { th.API.MattermostAuth = false }()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/login", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
	})

	t.Run("POST /login not permitted in single-user mode", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "token123")
		defer tearDownLocal()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/login", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("POST /login invalid login type", func(t *testing.T) {
		body := `{"username": "testuser", "password": "password123", "type": "invalid_type"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/login", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("POST /register not permitted in plugin mode", func(t *testing.T) {
		th.API.MattermostAuth = true
		defer func() { th.API.MattermostAuth = false }()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/register", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
	})

	t.Run("POST /register not permitted in single-user mode", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "token123")
		defer tearDownLocal()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/register", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})



	t.Run("POST /register invalid token", func(t *testing.T) {
		th.Store.EXPECT().GetTeam("0").Return(&model.Team{ID: "0", SignupToken: "realtoken"}, nil)

		body := `{"username": "newuser", "email": "newuser@example.com", "password": "securepassword", "token": "wrongtoken"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/register", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("POST /register token empty but user(s) already exist", func(t *testing.T) {
		th.Store.EXPECT().GetRegisteredUserCount().Return(int(5), nil)

		body := `{"username": "newuser", "email": "newuser@example.com", "password": "securepassword", "token": ""}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/register", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
	t.Run("POST /users/{userID}/changepassword not permitted in plugin mode", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "")
		defer tearDownLocal()
		thLocal.API.MattermostAuth = true

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/users/user1/changepassword", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Mattermost-User-Id", "someuser")

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
	})

	t.Run("POST /users/{userID}/changepassword not permitted in single-user mode", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "token123")
		defer tearDownLocal()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/users/user1/changepassword", bytes.NewBufferString(`{}`))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("attachSession required invalid single user token", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "token123")
		defer tearDownLocal()

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer wrongtoken")

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("attachSession MattermostAuth and header set", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "")
		defer tearDownLocal()
		thLocal.API.MattermostAuth = true

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Mattermost-User-Id", "mmuser1")

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
	})

	t.Run("attachSession GetSession fails and required true", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "")
		defer tearDownLocal()
		thLocal.Store.EXPECT().GetSession("badtoken", gomock.Any()).Return(nil, model.NewErrUnauthorized("invalid session"))

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer badtoken")

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("attachSession AuthService mismatch", func(t *testing.T) {
		thLocal, tearDownLocal := SetupTestAPIHelper(t, "")
		defer tearDownLocal()
		thLocal.Store.EXPECT().GetSession("mismatchtoken", gomock.Any()).Return(&model.Session{
			ID:          "session1",
			Token:       "mismatchtoken",
			UserID:      "user1",
			AuthService: "oauth",
			UpdateAt:    utils.GetMillis(),
		}, nil)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/logout", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer mismatchtoken")

		resp := doRequest(thLocal.Router, req)
		require.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}
