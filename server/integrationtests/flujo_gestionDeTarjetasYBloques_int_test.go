package integrationtests

import (
	"fmt"
	"testing"
	"time"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/require"
)

// Helper helper function to find blocks in a slice by ID
func findBlockInSlice(blocks []*model.Block, id string) *model.Block {
	for _, b := range blocks {
		if b.ID == id {
			return b
		}
	}
	return nil
}

// INT-03-01: Crear una tarjeta dentro de un tablero vía API y verificar que se persista como Block de tipo card
func TestINT0301CrearTarjetaEnTableroPersisteComoBloqueTypeCard(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-01: Crear tarjeta vía API y verificar persistencia como block    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-01: Verifica API POST /boards/{boardID}/cards -> App.CreateCard -> Store, confirmando type=card y persistencia.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	fmt.Printf("  → Tablero de prueba creado con ID: %s por usuario: %s\n", board.ID, user1.Username)

	card := &model.Card{
		BoardID: board.ID,
		Title:   "Tarjeta INT-03-01",
	}

	fmt.Println("  → Enviando petición POST para crear tarjeta...")
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, cardNew)
	require.NotEmpty(t, cardNew.ID)
	require.Equal(t, board.ID, cardNew.BoardID)
	require.Equal(t, "Tarjeta INT-03-01", cardNew.Title)
	fmt.Printf("  ✓ Tarjeta creada con ID: %s vía API REST\n", cardNew.ID)

	// Verificar persistencia en base de datos como Block de tipo card
	fmt.Println("  → Verificando persistencia directa en el Store...")
	block, err := th.Server.App().GetBlockByID(cardNew.ID)
	require.NoError(t, err)
	require.NotNil(t, block)
	require.Equal(t, cardNew.ID, block.ID)
	require.Equal(t, board.ID, block.BoardID)
	require.Equal(t, model.BlockType(model.TypeCard), block.Type)
	require.Equal(t, "Tarjeta INT-03-01", block.Title)
	require.Equal(t, user1.ID, block.CreatedBy)

	fmt.Println("  ✓ El bloque se persistió correctamente con type=card y boardId correcto")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-01: PASSED")
	fmt.Println()
}

// INT-03-02: Obtener tarjetas de un tablero y verificar que incluyan propiedades personalizadas
func TestINT0302ObtenerTarjetasConPropiedades(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-02: Obtener tarjetas y verificar propiedades personalizadas     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-02: Verifica GET /boards/{boardID}/cards -> App.GetCards -> Store, comprobando mapeo de properties.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)

	props := map[string]any{
		"priority":  "high",
		"milestone": "v1.0",
		"assignee":  "user1",
	}

	card := &model.Card{
		BoardID:    board.ID,
		Title:      "Tarjeta con Propiedades",
		Properties: props,
	}

	fmt.Println("  → Creando tarjeta con propiedades personalizadas...")
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, cardNew)

	fmt.Println("  → Solicitando tarjetas del tablero vía API...")
	cardsList, resp := th.Client.GetCards(board.ID, 0, -1)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotEmpty(t, cardsList)

	// Encontrar nuestra tarjeta y comprobar sus propiedades
	var foundCard *model.Card
	for _, c := range cardsList {
		if c.ID == cardNew.ID {
			foundCard = c
			break
		}
	}
	require.NotNil(t, foundCard, "No se encontró la tarjeta creada en la lista de tarjetas")

	fmt.Println("  → Validando mapeo y contenido de las propiedades...")
	require.Equal(t, "high", foundCard.Properties["priority"])
	require.Equal(t, "v1.0", foundCard.Properties["milestone"])
	require.Equal(t, "user1", foundCard.Properties["assignee"])

	fmt.Println("  ✓ Propiedades personalizadas recuperadas correctamente desde el Store")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-02: PASSED")
	fmt.Println()
}

// INT-03-03: Actualizar propiedades de una tarjeta vía PATCH y verificar que el Store refleje los cambios
func TestINT0303ActualizarPropiedadesViaPATCH(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-03: Actualizar propiedades vía PATCH y verificar persistencia    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-03: Verifica PATCH /cards/{cardID} -> App.PatchCard -> Store, confirmando actualización de properties.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)

	props := map[string]any{
		"status": "todo",
		"owner":  "daniel",
	}

	card := &model.Card{
		BoardID:    board.ID,
		Title:      "Tarjeta original PATCH",
		Properties: props,
	}
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	fmt.Println("  → Tarjeta creada. Enviando PATCH para actualizar propiedades...")
	newProps := map[string]any{
		"status": "in_progress",
		"owner":  "daniel_updated",
		"newKey": "value",
	}
	patch := &model.CardPatch{
		UpdatedProperties: newProps,
	}

	cardPatched, resp := th.Client.PatchCard(cardNew.ID, patch, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, cardPatched)

	fmt.Println("  → Verificando respuesta del PATCH...")
	require.Equal(t, "in_progress", cardPatched.Properties["status"])
	require.Equal(t, "daniel_updated", cardPatched.Properties["owner"])
	require.Equal(t, "value", cardPatched.Properties["newKey"])

	fmt.Println("  → Verificando persistencia directa en base de datos...")
	block, err := th.Server.App().GetBlockByID(cardNew.ID)
	require.NoError(t, err)
	require.NotNil(t, block)

	// properties en model.Block están mapeadas en Fields["properties"]
	blockProps, ok := block.Fields["properties"].(map[string]any)
	require.True(t, ok, "El campo fields['properties'] debería ser un mapa")
	require.Equal(t, "in_progress", blockProps["status"])
	require.Equal(t, "daniel_updated", blockProps["owner"])
	require.Equal(t, "value", blockProps["newKey"])

	fmt.Println("  ✓ Las propiedades fueron actualizadas correctamente en base de datos y la API retornó los cambios")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-03: PASSED")
	fmt.Println()
}

// INT-03-04: Insertar un bloque de contenido (texto, imagen, checkbox) dentro de una tarjeta y verificar relación padre
func TestINT0304InsertarBloqueContenidoYRelacionPadre(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-04: Insertar bloques de contenido y verificar ParentID (CardID) ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-04: Verifica POST /boards/{boardID}/blocks -> App.InsertBlocks -> Store, validando relación de parentesco.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	card := &model.Card{
		BoardID: board.ID,
		Title:   "Tarjeta Contenedora",
	}
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	fmt.Printf("  → Tarjeta contenedora creada con ID: %s\n", cardNew.ID)
	fmt.Println("  → Insertando bloque de texto, imagen y checkbox bajo la tarjeta...")

	textBlockID := utils.NewID(utils.IDTypeBlock)
	imageBlockID := utils.NewID(utils.IDTypeBlock)
	checkboxBlockID := utils.NewID(utils.IDTypeBlock)

	blocks := []*model.Block{
		{
			ID:       textBlockID,
			BoardID:  board.ID,
			ParentID: cardNew.ID,
			Type:     model.TypeText,
			Title:    "Bloque de Texto",
			CreateAt: 1,
			UpdateAt: 1,
		},
		{
			ID:       imageBlockID,
			BoardID:  board.ID,
			ParentID: cardNew.ID,
			Type:     model.TypeImage,
			Title:    "Bloque de Imagen",
			CreateAt: 1,
			UpdateAt: 1,
		},
		{
			ID:       checkboxBlockID,
			BoardID:  board.ID,
			ParentID: cardNew.ID,
			Type:     model.TypeCheckbox,
			Title:    "Bloque de Checkbox",
			CreateAt: 1,
			UpdateAt: 1,
		},
	}

	insertedBlocks, resp := th.Client.InsertBlocks(board.ID, blocks, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.Len(t, insertedBlocks, 3)

	// InsertBlocks replaces client-provided IDs with server-generated IDs while
	// preserving the order and references between blocks. Use the returned IDs
	// when checking persistence instead of the temporary request IDs.
	textBlockID = insertedBlocks[0].ID
	imageBlockID = insertedBlocks[1].ID
	checkboxBlockID = insertedBlocks[2].ID

	fmt.Println("  → Verificando relaciones de ParentID de los bloques creados...")
	for _, b := range insertedBlocks {
		require.Equal(t, cardNew.ID, b.ParentID)
		require.Equal(t, board.ID, b.BoardID)
	}

	// Comprobar persistencia directa en Store
	fmt.Println("  → Consultando bloques de contenido en el Store...")
	allBlocks, resp := th.Client.GetAllBlocksForBoard(board.ID)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	textB := findBlockInSlice(allBlocks, textBlockID)
	require.NotNil(t, textB)
	require.Equal(t, cardNew.ID, textB.ParentID)
	require.Equal(t, model.BlockType(model.TypeText), textB.Type)

	imageB := findBlockInSlice(allBlocks, imageBlockID)
	require.NotNil(t, imageB)
	require.Equal(t, cardNew.ID, imageB.ParentID)
	require.Equal(t, model.BlockType(model.TypeImage), imageB.Type)

	checkboxB := findBlockInSlice(allBlocks, checkboxBlockID)
	require.NotNil(t, checkboxB)
	require.Equal(t, cardNew.ID, checkboxB.ParentID)
	require.Equal(t, model.BlockType(model.TypeCheckbox), checkboxB.Type)

	fmt.Println("  ✓ Todos los bloques de contenido tienen como parentId el ID de la tarjeta contenedora")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-04: PASSED")
	fmt.Println()
}

// INT-03-05: Eliminar una tarjeta y verificar que sus bloques hijos también se marquen como eliminados
func TestINT0305EliminarTarjetaYBloquesHijos(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-05: Eliminar tarjeta y verificar eliminación de hijos cascada   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-05: Verifica DELETE /boards/{boardID}/blocks/{blockID} -> App.DeleteBlock -> Store, confirmando cascada deleteAt > 0.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	card := &model.Card{
		BoardID: board.ID,
		Title:   "Tarjeta a Eliminar",
	}
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	childBlockID := utils.NewID(utils.IDTypeBlock)
	blocks := []*model.Block{
		{
			ID:       childBlockID,
			BoardID:  board.ID,
			ParentID: cardNew.ID,
			Type:     model.TypeText,
			Title:    "Hijo de Tarjeta",
			CreateAt: 1,
			UpdateAt: 1,
		},
	}
	insertedBlocks, resp := th.Client.InsertBlocks(board.ID, blocks, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.Len(t, insertedBlocks, 1)
	childBlockID = insertedBlocks[0].ID

	// Para evitar colisiones en insert_at en el historial
	time.Sleep(15 * time.Millisecond)

	fmt.Printf("  → Tarjeta contenedora ID: %s y bloque hijo ID: %s creados\n", cardNew.ID, childBlockID)
	fmt.Println("  → Eliminando tarjeta contenedora...")
	success, resp := th.Client.DeleteBlock(board.ID, cardNew.ID, false)
	th.CheckOK(resp)
	require.True(t, success)

	fmt.Println("  → Verificando que tarjeta e hijos ya no estén en bloques activos...")
	activeBlocks, resp := th.Client.GetBlocksForBoard(board.ID)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.Nil(t, findBlockInSlice(activeBlocks, cardNew.ID))
	require.Nil(t, findBlockInSlice(activeBlocks, childBlockID))

	fmt.Println("  → Verificando persistencia de eliminación en el historial (deleteAt > 0)...")
	cardHistory, err := th.Server.Store().GetBlockHistory(cardNew.ID, model.QueryBlockHistoryOptions{Limit: 1, Descending: true})
	require.NoError(t, err)
	require.NotEmpty(t, cardHistory)
	require.Greater(t, cardHistory[0].DeleteAt, int64(0))
	fmt.Printf("  ✓ Tarjeta en historial marcada con deleteAt: %d\n", cardHistory[0].DeleteAt)

	childHistory, err := th.Server.Store().GetBlockHistory(childBlockID, model.QueryBlockHistoryOptions{Limit: 1, Descending: true})
	require.NoError(t, err)
	require.NotEmpty(t, childHistory)
	require.Greater(t, childHistory[0].DeleteAt, int64(0))
	fmt.Printf("  ✓ Bloque hijo en historial marcado con deleteAt: %d\n", childHistory[0].DeleteAt)

	fmt.Println("  ✓ La tarjeta y su bloque de contenido hijo fueron eliminados de forma lógica en cascada")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-05: PASSED")
	fmt.Println()
}

// INT-03-06: Restaurar una tarjeta eliminada y verificar que se recupere junto con sus bloques hijos
func TestINT0306RestaurarTarjetaYBloquesHijos(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-06: Restaurar tarjeta y verificar restauración de hijos cascada   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-06: Verifica POST /boards/{boardID}/blocks/{blockID}/undelete -> App.UndeleteBlock -> Store, confirmando cascada deleteAt = 0.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board := th.CreateBoard(testTeamID, model.BoardTypeOpen)
	card := &model.Card{
		BoardID: board.ID,
		Title:   "Tarjeta a Restaurar",
	}
	cardNew, resp := th.Client.CreateCard(board.ID, card, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	childBlockID := utils.NewID(utils.IDTypeBlock)
	blocks := []*model.Block{
		{
			ID:       childBlockID,
			BoardID:  board.ID,
			ParentID: cardNew.ID,
			Type:     model.TypeText,
			Title:    "Hijo a Restaurar",
			CreateAt: 1,
			UpdateAt: 1,
		},
	}
	insertedBlocks, resp := th.Client.InsertBlocks(board.ID, blocks, false)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.Len(t, insertedBlocks, 1)
	childBlockID = insertedBlocks[0].ID

	time.Sleep(15 * time.Millisecond)

	fmt.Println("  → Eliminando tarjeta para preparar la precondición...")
	success, resp := th.Client.DeleteBlock(board.ID, cardNew.ID, false)
	th.CheckOK(resp)
	require.True(t, success)

	time.Sleep(15 * time.Millisecond)

	fmt.Println("  → Restaurando tarjeta eliminada...")
	success, resp = th.Client.UndeleteBlock(board.ID, cardNew.ID)
	th.CheckOK(resp)
	require.True(t, success)

	fmt.Println("  → Verificando que tarjeta e hijo vuelvan a estar activos...")
	activeBlocks, resp := th.Client.GetBlocksForBoard(board.ID)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	activeCard := findBlockInSlice(activeBlocks, cardNew.ID)
	require.NotNil(t, activeCard, "La tarjeta no volvió a estar activa")
	require.Equal(t, int64(0), activeCard.DeleteAt)

	activeChild := findBlockInSlice(activeBlocks, childBlockID)
	require.NotNil(t, activeChild, "El bloque hijo no volvió a estar activo")
	require.Equal(t, int64(0), activeChild.DeleteAt)

	fmt.Println("  → Verificando base de datos a través del historial...")
	cardHistory, err := th.Server.Store().GetBlockHistory(cardNew.ID, model.QueryBlockHistoryOptions{Limit: 1, Descending: true})
	require.NoError(t, err)
	require.Equal(t, int64(0), cardHistory[0].DeleteAt)

	childHistory, err := th.Server.Store().GetBlockHistory(childBlockID, model.QueryBlockHistoryOptions{Limit: 1, Descending: true})
	require.NoError(t, err)
	require.Equal(t, int64(0), childHistory[0].DeleteAt)

	fmt.Println("  ✓ La tarjeta y su bloque hijo fueron recuperados exitosamente (deleteAt = 0)")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-06: PASSED")
	fmt.Println()
}

// INT-03-07: Crear un tablero con bloques en lote (CreateBoardsAndBlocks) y verificar atomicidad de la operación
func TestINT0307CrearTableroYBloquesLote(t *testing.T) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  INT-03-07: Crear en lote vía API y verificar atomicidad / transacción   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════╝")
	t.Log("INT-03-07: Verifica POST /boards-and-blocks -> App.CreateBoardsAndBlocks -> Store (transacción), validando atomicidad.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	boardID := utils.NewID(utils.IDTypeBoard)
	blockID1 := utils.NewID(utils.IDTypeBlock)
	timestamp := model.GetMillis()

	// Caso 1: Envío inválido para forzar error de validación del lote
	fmt.Println("  → Enviando lote inválido (bloque pertenece a tablero inexistente)...")
	invalidBab := &model.BoardsAndBlocks{
		Boards: []*model.Board{
			{
				ID:     boardID,
				TeamID: testTeamID,
				Title:  "Tablero Lote Inválido",
				Type:   model.BoardTypeOpen,
			},
		},
		Blocks: []*model.Block{
			{
				ID:       blockID1,
				BoardID:  "nonexistent-board-id", // Provocará fallo en validación
				Type:     model.TypeCard,
				Title:    "Tarjeta Huérfana",
				CreateAt: timestamp,
				UpdateAt: timestamp,
			},
		},
	}

	babRes, resp := th.Client.CreateBoardsAndBlocks(invalidBab)
	th.CheckBadRequest(resp)
	require.Nil(t, babRes)
	fmt.Println("  ✓ Petición rechazada con Bad Request (400) exitosamente")

	// Verificar atomicidad: el tablero no debió haber sido guardado en Store
	fmt.Println("  → Verificando atomicidad (rollback): consultando tablero en Store...")
	boardDB, err := th.Server.App().GetBoard(boardID)
	require.Error(t, err)
	require.True(t, model.IsErrNotFound(err))
	require.Nil(t, boardDB)
	fmt.Println("  ✓ Confirmado: el tablero NO fue persistido (Rollback exitoso)")

	// Caso 2: Envío de lote válido para verificar que sí se crea todo atómicamente si no hay errores
	fmt.Println("  → Enviando lote válido...")
	validBoardID := utils.NewID(utils.IDTypeBoard)
	validBlockID := utils.NewID(utils.IDTypeBlock)

	validBab := &model.BoardsAndBlocks{
		Boards: []*model.Board{
			{
				ID:     validBoardID,
				TeamID: testTeamID,
				Title:  "Tablero Lote Válido",
				Type:   model.BoardTypeOpen,
			},
		},
		Blocks: []*model.Block{
			{
				ID:       validBlockID,
				BoardID:  validBoardID,
				Type:     model.TypeCard,
				Title:    "Tarjeta en Lote",
				CreateAt: timestamp,
				UpdateAt: timestamp,
			},
		},
	}

	babResValid, resp := th.Client.CreateBoardsAndBlocks(validBab)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)
	require.NotNil(t, babResValid)
	require.Len(t, babResValid.Boards, 1)
	require.Len(t, babResValid.Blocks, 1)

	// El lote se reasigna ID internamente, obtenemos el ID retornado
	createdBoardID := babResValid.Boards[0].ID
	createdBlockID := babResValid.Blocks[0].ID

	// Verificar persistencia de todo el lote
	boardDBValid, err := th.Server.App().GetBoard(createdBoardID)
	require.NoError(t, err)
	require.NotNil(t, boardDBValid)

	blockDBValid, err := th.Server.App().GetBlockByID(createdBlockID)
	require.NoError(t, err)
	require.NotNil(t, blockDBValid)

	fmt.Println("  ✓ Confirmado: todos los elementos del lote válido se crearon correctamente")
	fmt.Println("  ════════════════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ INT-03-07: PASSED")
	fmt.Println()
}
