//go:build ignore

package integrationtests

import (
	"net/http"
	"testing"

	"github.com/mattermost/focalboard/server/client"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/require"
)

func TestINT0501HabilitarComparticion(t *testing.T) {
	t.Log("INT-05-01: Habilitar compartición pública de un tablero y verificar que se genere un token de compartición en el Store")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	token := utils.NewID(utils.IDTypeToken)

	th.Server.Config().EnablePublicSharedBoards = true
	sharing := model.Sharing{
		ID:      board.ID,
		Token:   token,
		Enabled: true,
	}

	success, resp := th.Client.PostSharing(&sharing)
	th.CheckOK(resp)
	require.True(t, success)

	// Verificar directamente en la base de datos (Store)
	persistedSharing, err := th.Server.Store().GetSharing(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedSharing)
	require.True(t, persistedSharing.Enabled)
	require.Equal(t, token, persistedSharing.Token)
}

func TestINT0502AccederTableroCompartidoValido(t *testing.T) {
	t.Log("INT-05-02: Acceder a un tablero compartido con token válido sin autenticación y verificar respuesta")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	token := utils.NewID(utils.IDTypeToken)

	th.Server.Config().EnablePublicSharedBoards = true
	sharing := model.Sharing{
		ID:      board.ID,
		Token:   token,
		Enabled: true,
	}

	success, resp := th.Client.PostSharing(&sharing)
	th.CheckOK(resp)
	require.True(t, success)

	// Cliente sin autenticación
	anonClient := client.NewClient(th.Server.Config().ServerRoot, "")
	sharedBoard, anonResp := anonClient.GetBoard(board.ID, token)
	require.NoError(t, anonResp.Error)
	require.Equal(t, http.StatusOK, anonResp.StatusCode)
	require.NotNil(t, sharedBoard)
	require.Equal(t, board.ID, sharedBoard.ID)
}

func TestINT0503AccederTableroCompartidoInvalido(t *testing.T) {
	t.Log("INT-05-03: Acceder a un tablero compartido con token inválido y verificar rechazo")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	token := utils.NewID(utils.IDTypeToken)

	th.Server.Config().EnablePublicSharedBoards = true
	sharing := model.Sharing{
		ID:      board.ID,
		Token:   token,
		Enabled: true,
	}

	success, resp := th.Client.PostSharing(&sharing)
	th.CheckOK(resp)
	require.True(t, success)

	// Intentar acceder con un token inválido sin autenticación
	anonClient := client.NewClient(th.Server.Config().ServerRoot, "")
	invalidToken := "invalid-token-value-12345"
	sharedBoard, anonResp := anonClient.GetBoard(board.ID, invalidToken)
	require.Error(t, anonResp.Error)
	require.Equal(t, http.StatusUnauthorized, anonResp.StatusCode)
	require.Nil(t, sharedBoard)
}

func TestINT0504DeshabilitarComparticion(t *testing.T) {
	t.Log("INT-05-04: Deshabilitar compartición y verificar que el acceso público se revoque")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	token := utils.NewID(utils.IDTypeToken)

	th.Server.Config().EnablePublicSharedBoards = true
	sharing := model.Sharing{
		ID:      board.ID,
		Token:   token,
		Enabled: true,
	}

	success, resp := th.Client.PostSharing(&sharing)
	th.CheckOK(resp)
	require.True(t, success)

	// Deshabilitar compartición pública
	sharing.Enabled = false
	success, resp = th.Client.PostSharing(&sharing)
	th.CheckOK(resp)
	require.True(t, success)

	// Verificar estado en la base de datos (Store)
	persistedSharing, err := th.Server.Store().GetSharing(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedSharing)
	require.False(t, persistedSharing.Enabled)

	// Verificar rechazo de acceso público mediante el token
	anonClient := client.NewClient(th.Server.Config().ServerRoot, "")
	sharedBoard, anonResp := anonClient.GetBoard(board.ID, token)
	require.Error(t, anonResp.Error)
	require.Equal(t, http.StatusUnauthorized, anonResp.StatusCode)
	require.Nil(t, sharedBoard)
}
