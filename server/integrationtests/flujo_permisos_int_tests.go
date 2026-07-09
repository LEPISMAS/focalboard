package integrationtests

import (
	"fmt"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/require"
)

// Helper function to find a board member by user ID
func findBoardMember(members []*model.BoardMember, userID string) *model.BoardMember {
	for _, m := range members {
		if m.UserID == userID {
			return m
		}
	}
	return nil
}

// INT-04-01: Un usuario viewer intenta editar un tablero y verificar que el sistema rechace la operación
func TestINT0401ViewerNoPuedeEditarTablero(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-01: Usuario 'viewer' intenta editar y es rechazado (403)         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-01: Verifica que rol viewer impida modificar el tablero vía API PATCH /boards/{id}.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero como user1 (creador será SchemeAdmin por defecto)
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero creado por %s con ID: %s\n", user1.Username, board.ID)

	// Agregar a user2 como SchemeViewer (viewer)
	_, err := th.Server.App().AddMemberToBoard(&model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeViewer: true,
	})
	require.NoError(t, err)
	fmt.Printf("  → Usuario %s agregado al tablero con rol 'viewer'\n", user2.Username)

	// Intentar editar el tablero usando Client2 (user2)
	fmt.Println("  → Intentando editar título de tablero como user2...")
	newTitle := "Título no Autorizado"
	updatedBoard, resp := th.Client2.PatchBoard(board.ID, &model.BoardPatch{Title: &newTitle})

	// Verificar rechazo 403 Forbidden
	th.CheckForbidden(resp)
	require.Nil(t, updatedBoard)
	fmt.Println("  ✓ Petición rechazada correctamente con HTTP 403 Forbidden")

	// Verificar que el título original siga intacto en el Store
	boardDB, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.Equal(t, board.Title, boardDB.Title)
	fmt.Printf("  ✓ Confirmado: el título sigue siendo '%s'\n", boardDB.Title)

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-01: PASSED")
	fmt.Println()
}

// INT-04-02: Un usuario editor edita un tablero y verificar que el sistema permita la operación
func TestINT0402EditorSiPuedeEditarTablero(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-02: Usuario 'editor' edita tablero exitosamente (200)            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-02: Verifica que rol editor permita modificar el tablero vía API PATCH /boards/{id}.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero como user1
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero creado por %s con ID: %s\n", user1.Username, board.ID)

	// Agregar a user2 como SchemeEditor (editor)
	_, err := th.Server.App().AddMemberToBoard(&model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeEditor: true,
	})
	require.NoError(t, err)
	fmt.Printf("  → Usuario %s agregado al tablero con rol 'editor'\n", user2.Username)

	// Intentar editar el tablero usando Client2 (user2)
	fmt.Println("  → Intentando editar título de tablero como user2...")
	newTitle := "Título Modificado por Editor"
	updatedBoard, resp := th.Client2.PatchBoard(board.ID, &model.BoardPatch{Title: &newTitle})

	// Verificar éxito 200 OK
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, updatedBoard)
	require.Equal(t, newTitle, updatedBoard.Title)
	fmt.Println("  ✓ Petición aceptada correctamente con HTTP 200 OK")

	// Verificar persistencia directa en Store
	boardDB, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.Equal(t, newTitle, boardDB.Title)
	fmt.Printf("  ✓ Confirmado en base de datos: el nuevo título es '%s'\n", boardDB.Title)

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-02: PASSED")
	fmt.Println()
}

// INT-04-03: Un usuario admin del tablero agrega un nuevo miembro y verificar la persistencia de la membresía
func TestINT0403AdminAgregaMiembro(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-03: Admin agrega miembro a tablero y se persiste en Store        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-03: Verifica API POST /boards/{id}/members -> App.AddMemberToBoard -> Store, comprobando persistencia de membresía.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero como user1 (admin por defecto)
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero creado por %s con ID: %s\n", user1.Username, board.ID)

	// Agregar user2 como miembro vía API (usando el cliente de user1 que es admin)
	fmt.Println("  → Agregando a user2 al tablero como 'editor' vía API...")
	memberParam := &model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeEditor: true,
	}
	memberRes, resp := th.Client.AddMemberToBoard(memberParam)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, memberRes)
	require.Equal(t, user2.ID, memberRes.UserID)
	require.True(t, memberRes.SchemeEditor)
	fmt.Println("  ✓ Miembro agregado con éxito vía API")

	// Verificar persistencia en base de datos
	fmt.Println("  → Consultando membresías en el Store...")
	members, err := th.Server.App().GetMembersForBoard(board.ID)
	require.NoError(t, err)

	memberDB := findBoardMember(members, user2.ID)
	require.NotNil(t, memberDB, "No se encontró el miembro en base de datos")
	require.True(t, memberDB.SchemeEditor)
	fmt.Println("  ✓ Membresía y roles persistidos correctamente en la base de datos")

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-03: PASSED")
	fmt.Println()
}

// INT-04-04: Un usuario no miembro intenta acceder a un tablero privado y verificar rechazo
func TestINT0404NoMiembroNoAccedeATableroPrivado(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-04: No miembro intenta acceder a tablero privado y es rechazado  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-04: Verifica que GET /boards/{boardID} por un no miembro retorne 403 Forbidden.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero privado como user1
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero privado creado por %s con ID: %s\n", user1.Username, board.ID)
	fmt.Printf("  → Usuario %s no tiene membresía en este tablero\n", user2.Username)

	// Intentar obtener el tablero como user2
	fmt.Println("  → Intentando acceder al tablero privado como user2...")
	retrievedBoard, resp := th.Client2.GetBoard(board.ID, "")

	// Verificar rechazo 403 Forbidden
	th.CheckForbidden(resp)
	require.Nil(t, retrievedBoard)
	fmt.Println("  ✓ Petición rechazada con HTTP 403 Forbidden correctamente")

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-04: PASSED")
	fmt.Println()
}

// INT-04-05: Cambiar el rol de un miembro de viewer a editor y verificar que los permisos se actualicen
func TestINT0405CambiarRolMiembroViewerAEditor(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-05: Cambiar rol de miembro de 'viewer' a 'editor' y verificar    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-05: Verifica API PUT /boards/{id}/members/{userID} -> App.UpdateBoardMember -> Store, comprobando escalación de permisos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero creado por %s con ID: %s\n", user1.Username, board.ID)

	// Agregar user2 como SchemeViewer (viewer)
	_, err := th.Server.App().AddMemberToBoard(&model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeViewer: true,
	})
	require.NoError(t, err)

	// Verificar que user2 NO puede editar inicialmente (debería retornar 403)
	newTitle1 := "Fallo Esperado"
	_, resp1 := th.Client2.PatchBoard(board.ID, &model.BoardPatch{Title: &newTitle1})
	th.CheckForbidden(resp1)
	fmt.Println("  ✓ Precondición cumplida: user2 tiene rol 'viewer' y no puede editar")

	// user1 (admin) actualiza el rol de user2 a SchemeEditor (editor) vía API
	fmt.Println("  → Actualizando rol de user2 a 'editor' vía API PUT...")
	updateMemberParam := &model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeEditor: true,
	}
	updatedMember, resp2 := th.Client.UpdateBoardMember(updateMemberParam)
	th.CheckOK(resp2)
	require.NoError(t, resp2.Error)
	require.NotNil(t, updatedMember)
	require.True(t, updatedMember.SchemeEditor)
	require.False(t, updatedMember.SchemeViewer)
	fmt.Println("  ✓ Rol actualizado a 'editor' vía API")

	// Verificar que user2 AHORA SÍ puede editar el tablero
	fmt.Println("  → Intentando editar título de tablero como user2 nuevamente...")
	newTitle2 := "Éxito Esperado por Editor"
	updatedBoard, resp3 := th.Client2.PatchBoard(board.ID, &model.BoardPatch{Title: &newTitle2})
	th.CheckOK(resp3)
	require.NoError(t, resp3.Error)
	require.NotNil(t, updatedBoard)
	require.Equal(t, newTitle2, updatedBoard.Title)
	fmt.Println("  ✓ user2 ahora puede editar (HTTP 200 OK)")

	// Verificar persistencia directa en Store
	boardDB, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.Equal(t, newTitle2, boardDB.Title)
	fmt.Println("  ✓ Cambios correctamente persistidos en el Store")

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-05: PASSED")
	fmt.Println()
}

// INT-04-06: Eliminar membresía de un usuario y verificar que pierde acceso inmediatamente
func TestINT0406EliminarMembresiaPierdeAcceso(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-06: Eliminar membresía de usuario y verificar pérdida de acceso  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-06: Verifica API DELETE /boards/{id}/members/{userID} -> App.DeleteBoardMember -> Store, comprobando denegación posterior.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero privado
	board := th.CreateBoard(testTeamID, model.BoardTypePrivate)
	fmt.Printf("  → Tablero privado creado por %s con ID: %s\n", user1.Username, board.ID)

	// Agregar user2 como SchemeViewer (viewer)
	_, err := th.Server.App().AddMemberToBoard(&model.BoardMember{
		BoardID:      board.ID,
		UserID:       user2.ID,
		SchemeViewer: true,
	})
	require.NoError(t, err)

	// Verificar acceso inicial (espera 200 OK)
	retrievedBoard, resp1 := th.Client2.GetBoard(board.ID, "")
	th.CheckOK(resp1)
	require.NotNil(t, retrievedBoard)
	fmt.Println("  ✓ user2 tiene acceso inicial al tablero (con membresía activa)")

	// user1 (admin) elimina la membresía de user2
	fmt.Println("  → Eliminando membresía de user2 vía API DELETE...")
	deleteMemberParam := &model.BoardMember{
		BoardID: board.ID,
		UserID:  user2.ID,
	}
	success, resp2 := th.Client.DeleteBoardMember(deleteMemberParam)
	th.CheckOK(resp2)
	require.True(t, success)
	fmt.Println("  ✓ Membresía eliminada vía API")

	// Verificar que user2 ya no puede acceder al tablero (espera 403)
	fmt.Println("  → Intentando acceder al tablero como user2 tras la eliminación...")
	retrievedBoard2, resp3 := th.Client2.GetBoard(board.ID, "")
	th.CheckForbidden(resp3)
	require.Nil(t, retrievedBoard2)
	fmt.Println("  ✓ Acceso denegado con HTTP 403 Forbidden de forma inmediata")

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-06: PASSED")
	fmt.Println()
}

// INT-04-07: Verificar que un tablero público (tipo Open) sea accesible para cualquier miembro del equipo
func TestINT0407TableroPublicoAccesibleMiembrosEquipo(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-04-07: Tablero público 'Open' es accesible por miembros del equipo ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-04-07: Verifica que un tablero Open en un equipo sea visible por otro miembro del equipo sin membresía explícita.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	user2 := th.GetUser2()

	// Crear tablero público (Open) como user1 en testTeamID
	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	fmt.Printf("  → Tablero Open creado por %s con ID: %s en el equipo: %s\n", user1.Username, board.ID, testTeamID)
	fmt.Printf("  → user2 es miembro del equipo pero no tiene membresía explícita en este tablero\n")

	// user2 intenta acceder al tablero público Open vía API
	fmt.Println("  → Intentando acceder al tablero público como user2...")
	retrievedBoard, resp := th.Client2.GetBoard(board.ID, "")

	// Debería permitir el acceso (200 OK) porque el tablero es Open
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, retrievedBoard)
	require.Equal(t, board.ID, retrievedBoard.ID)
	require.Equal(t, model.BoardTypeOpen, retrievedBoard.Type)
	fmt.Printf("  ✓ Acceso concedido exitosamente (HTTP 200 OK). Título del tablero: '%s'\n", retrievedBoard.Title)

	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-04-07: PASSED")
	fmt.Println()
}
