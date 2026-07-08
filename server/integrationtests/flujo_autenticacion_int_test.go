package integrationtests

import (
	"io"
	"net/http"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/require"
)

func int01RegisterUser(t *testing.T, th *TestHelper, username, email, password, token string) {
	t.Helper()

	success, resp := th.Client.Register(&model.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Token:    token,
	})
	th.CheckOK(resp)
	require.True(t, success)
}

func int01LoginUser(t *testing.T, th *TestHelper, username, password string) *model.LoginResponse {
	t.Helper()

	data, resp := th.Client.Login(&model.LoginRequest{
		Type:     "normal",
		Username: username,
		Password: password,
	})
	th.CheckOK(resp)
	require.NotNil(t, data)
	require.NotEmpty(t, data.Token)
	return data
}

func int01LogoutViaAPI(t *testing.T, th *TestHelper) {
	t.Helper()

	resp, err := th.Client.DoAPIPost("/logout", "")
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestINT0101RegistrarUsuarioNuevoPersisteEnStore(t *testing.T) {
	t.Log("INT-01-01: verifica API REST /register -> App.RegisterUser -> Store.CreateUser, confirmando persistencia del usuario.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_register"
	email := "int01_user_register@example.com"
	userPassword := "Pa$$word-int01-01"

	int01RegisterUser(t, th, username, email, userPassword, "")

	persistedUser, err := th.Server.Store().GetUserByUsername(username)
	require.NoError(t, err)
	require.NotNil(t, persistedUser)
	require.Equal(t, username, persistedUser.Username)
	require.Equal(t, email, persistedUser.Email)
	require.NotEmpty(t, persistedUser.ID)
}

func TestINT0102LoginGeneraSesionValida(t *testing.T) {
	t.Log("INT-01-02: verifica API REST /login -> App.Login -> Auth/Store.CreateSession, confirmando token y sesion persistida.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_login"
	email := "int01_user_login@example.com"
	userPassword := "Pa$$word-int01-02"

	int01RegisterUser(t, th, username, email, userPassword, "")
	loginData := int01LoginUser(t, th, username, userPassword)

	storedUser, err := th.Server.Store().GetUserByUsername(username)
	require.NoError(t, err)
	require.NotNil(t, storedUser)

	session, err := th.Server.App().GetSession(loginData.Token)
	require.NoError(t, err)
	require.NotNil(t, session)
	require.Equal(t, storedUser.ID, session.UserID)
	require.Equal(t, loginData.Token, session.Token)
}

func TestINT0103EndpointProtegidoConTokenValido(t *testing.T) {
	t.Log("INT-01-03: verifica que un token valido atraviesa API REST protegida /users/me -> Auth.GetSession -> Store.GetSession.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_me"
	email := "int01_user_me@example.com"
	userPassword := "Pa$$word-int01-03"

	int01RegisterUser(t, th, username, email, userPassword, "")
	int01LoginUser(t, th, username, userPassword)

	me, resp := th.Client.GetMe()
	th.CheckOK(resp)
	require.NotNil(t, me)
	require.Equal(t, username, me.Username)
}

func TestINT0104EndpointProtegidoSinTokenRechaza(t *testing.T) {
	t.Log("INT-01-04: verifica que /users/me sin sesion es rechazado por la capa API/Auth con 401 Unauthorized.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	me, resp := th.Client.GetMe()
	th.CheckUnauthorized(resp)
	require.Nil(t, me)
}

func TestINT0105CambiarPasswordPermiteLoginNuevo(t *testing.T) {
	t.Log("INT-01-05: verifica /users/{id}/changepassword -> App.ChangePassword -> Store.UpdateUserPasswordByID y login posterior con la nueva clave.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_change_password"
	email := "int01_user_change_password@example.com"
	oldPassword := "Pa$$word-int01-05"
	newPassword := "N3wPa$$word-int01-05"

	int01RegisterUser(t, th, username, email, oldPassword, "")
	int01LoginUser(t, th, username, oldPassword)

	me, resp := th.Client.GetMe()
	th.CheckOK(resp)
	require.NotNil(t, me)

	success, resp := th.Client.UserChangePassword(me.ID, &model.ChangePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
	th.CheckOK(resp)
	require.True(t, success)

	int01LogoutViaAPI(t, th)
	loginData := int01LoginUser(t, th, username, newPassword)
	require.NotEmpty(t, loginData.Token)
}

func TestINT0106LogoutInvalidaToken(t *testing.T) {
	t.Log("INT-01-06: verifica /logout -> App.Logout -> Store.DeleteSession, confirmando que el token anterior ya no autoriza /users/me.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_logout"
	email := "int01_user_logout@example.com"
	userPassword := "Pa$$word-int01-06"

	int01RegisterUser(t, th, username, email, userPassword, "")
	loginData := int01LoginUser(t, th, username, userPassword)
	oldToken := loginData.Token

	int01LogoutViaAPI(t, th)
	th.Client.Token = oldToken

	me, resp := th.Client.GetMe()
	th.CheckUnauthorized(resp)
	require.Nil(t, me)
}

func TestINT0107RegistroDuplicadoRechazado(t *testing.T) {
	t.Log("INT-01-07: verifica que un registro duplicado con token de invitacion valido llega a App/Store y es rechazado sin duplicar usuarios.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int01_user_duplicate"
	email := "int01_user_duplicate@example.com"
	userPassword := "Pa$$word-int01-07"

	int01RegisterUser(t, th, username, email, userPassword, "")
	int01LoginUser(t, th, username, userPassword)

	team, resp := th.Client.GetTeam(model.GlobalTeamID)
	th.CheckOK(resp)
	require.NotNil(t, team)
	require.NotEmpty(t, team.SignupToken)

	beforeCount, err := th.Server.App().GetRegisteredUserCount()
	require.NoError(t, err)

	success, resp := th.Client2.Register(&model.RegisterRequest{
		Username: username,
		Email:    email,
		Password: userPassword,
		Token:    team.SignupToken,
	})
	th.CheckBadRequest(resp)
	require.False(t, success)

	afterCount, err := th.Server.App().GetRegisteredUserCount()
	require.NoError(t, err)
	require.Equal(t, beforeCount, afterCount)

	persistedUser, err := th.Server.Store().GetUserByUsername(username)
	require.NoError(t, err)
	require.NotNil(t, persistedUser)
	require.Equal(t, email, persistedUser.Email)
}
