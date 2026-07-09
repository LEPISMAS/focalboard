package integrationtests

import (
	"testing"

	"github.com/mattermost/focalboard/server/model"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"github.com/stretchr/testify/require"
)

func TestINT0901BuscarTablerosPorTituloRetornaSoloAccesibles(t *testing.T) {
	t.Log("INT-09-01: verifica API REST /teams/{teamID}/boards/search -> App.SearchBoardsForUser -> Store.SearchBoardsForUser, confirmando que solo se retornan tableros accesibles al usuario.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tableros accesibles para user1
	user1Board1, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-01 Proyecto Alpha",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user1Board1)

	user1Board2, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-01 Proyecto Beta",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user1Board2)

	// Crear tablero privado para user2 (no accesible para user1)
	user2PrivateBoard, resp := th.Client2.CreateBoard(&model.Board{
		Title:  "INT-09-01 Proyecto Gamma",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user2PrivateBoard)

	// Crear tablero público (accesible para todos)
	publicBoard, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-01 Proyecto Delta",
		Type:   model.BoardTypeOpen,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, publicBoard)

	// Buscar tableros con término "Proyecto"
	boards, resp := th.Client.SearchBoardsForUser(testTeamID, "Proyecto", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.NotNil(t, boards)

	// Verificar que user1 ve sus tableros privados y el público
	boardIDs := make([]string, 0, len(boards))
	for _, board := range boards {
		boardIDs = append(boardIDs, board.ID)
	}

	require.Contains(t, boardIDs, user1Board1.ID, "user1 debería ver su propio tablero privado")
	require.Contains(t, boardIDs, user1Board2.ID, "user1 debería ver su segundo tablero privado")
	require.Contains(t, boardIDs, publicBoard.ID, "user1 debería ver tableros públicos")
	require.NotContains(t, boardIDs, user2PrivateBoard.ID, "user1 NO debería ver tableros privados de user2")

	t.Logf("Búsqueda exitosa: user1 encontró %d tableros accesibles (2 privados + 1 público), excluyendo 1 privado de user2", len(boards))
}

func TestINT0902BuscarTerminoSinCoincidenciasRetornaVacio(t *testing.T) {
	t.Log("INT-09-02: verifica que una búsqueda con término que no coincide devuelve array vacío con HTTP 200.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear algunos tableros
	_, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-02 Tablero existente",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)

	// Buscar con término que no existe
	boards, resp := th.Client.SearchBoardsForUser(testTeamID, "TerminoInexistenteXYZ123", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.NotNil(t, boards)
	require.Empty(t, boards, "La búsqueda debería retornar array vacío cuando no hay coincidencias")

	t.Log("Búsqueda exitosa: se retornó array vacío para término sin coincidencias")
}

func TestINT0903BuscarTablerosRespetanPermisos(t *testing.T) {
	t.Log("INT-09-03: verifica que la búsqueda respeta permisos y no muestra tableros privados ajenos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tablero privado para user1
	user1Board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-03 Privado User1",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user1Board)

	// Crear tablero privado para user2
	user2Board, resp := th.Client2.CreateBoard(&model.Board{
		Title:  "INT-09-03 Privado User2",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user2Board)

	// Crear tablero público
	publicBoard, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-03 Publico",
		Type:   model.BoardTypeOpen,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, publicBoard)

	// Buscar como user1
	boardsUser1, resp := th.Client.SearchBoardsForUser(testTeamID, "Privado", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.NotNil(t, boardsUser1)

	boardIDsUser1 := make([]string, 0, len(boardsUser1))
	for _, board := range boardsUser1 {
		boardIDsUser1 = append(boardIDsUser1, board.ID)
	}

	require.Contains(t, boardIDsUser1, user1Board.ID, "user1 debería ver su tablero privado")
	require.NotContains(t, boardIDsUser1, user2Board.ID, "user1 NO debería ver tablero privado de user2")

	// Buscar como user2
	boardsUser2, resp := th.Client2.SearchBoardsForUser(testTeamID, "Privado", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.NotNil(t, boardsUser2)

	boardIDsUser2 := make([]string, 0, len(boardsUser2))
	for _, board := range boardsUser2 {
		boardIDsUser2 = append(boardIDsUser2, board.ID)
	}

	require.Contains(t, boardIDsUser2, user2Board.ID, "user2 debería ver su tablero privado")
	require.NotContains(t, boardIDsUser2, user1Board.ID, "user2 NO debería ver tablero privado de user1")

	// Verificar que ambos pueden ver el público
	boardsPublicUser1, resp := th.Client.SearchBoardsForUser(testTeamID, "Publico", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.Len(t, boardsPublicUser1, 1)
	require.Equal(t, publicBoard.ID, boardsPublicUser1[0].ID)

	boardsPublicUser2, resp := th.Client2.SearchBoardsForUser(testTeamID, "Publico", model.BoardSearchFieldTitle)
	th.CheckOK(resp)
	require.Len(t, boardsPublicUser2, 1)
	require.Equal(t, publicBoard.ID, boardsPublicUser2[0].ID)

	t.Log("Permisos verificados: cada usuario solo ve sus tableros privados y tableros públicos")
}

func TestINT0904BuscarCaracteresEspecialesSinErroresSQL(t *testing.T) {
	t.Log("INT-09-04: verifica que la búsqueda con caracteres especiales no causa errores SQL en la capa de Store.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tablero para tener contexto
	_, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-09-04 Tablero normal",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)

	// Probar búsqueda con caracteres especiales que podrían causar inyección SQL
	specialTerms := []string{
		"' OR '1'='1",
		"'; DROP TABLE boards; --",
		"<script>alert('xss')</script>",
		"../../etc/passwd",
		"%",
		"_",
		"*",
		"[]",
		"{}",
	}

	for _, term := range specialTerms {
		boards, resp := th.Client.SearchBoardsForUser(testTeamID, term, model.BoardSearchFieldTitle)
		// La respuesta debe ser válida (200 OK) sin errores de servidor
		th.CheckOK(resp)
		require.NotNil(t, boards, "La búsqueda con caracteres especiales debería retornar un array (vacío o con resultados)")
		// El array puede estar vacío o tener resultados, pero no debe causar error
		t.Logf("Búsqueda con término especial '%s' completada sin errores SQL, resultados: %d", term, len(boards))
	}

	t.Log("Búsqueda con caracteres especiales completada exitosamente sin errores SQL")
}
