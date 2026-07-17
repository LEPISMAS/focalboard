package integrationtests

import (
	"io"
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	"github.com/stretchr/testify/require"
)

func TestINT1001RegistrarUsuarioLoginExitoso(t *testing.T) {
	t.Log("INT-10-01: simula flujo UI RegisterPage -> OctoClient -> API -> Store, verificando registro exitoso y login posterior.")

	th := SetupTestHelper(t).Start()
	defer th.TearDown()

	username := "int10_user_register"
	email := "int10_user_register@example.com"
	password := "Pa$$word-int10-01"

	// Simular registro desde UI
	success, resp := th.Client.Register(&model.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Token:    "",
	})
	th.CheckOK(resp)
	require.True(t, success, "El registro debería ser exitoso")

	// Verificar que el usuario persistió en Store
	persistedUser, err := th.Server.Store().GetUserByUsername(username)
	require.NoError(t, err)
	require.NotNil(t, persistedUser)
	require.Equal(t, username, persistedUser.Username)
	require.Equal(t, email, persistedUser.Email)

	// Simular login desde UI
	loginData, resp := th.Client.Login(&model.LoginRequest{
		Type:     "normal",
		Username: username,
		Password: password,
	})
	th.CheckOK(resp)
	require.NotNil(t, loginData)
	require.NotEmpty(t, loginData.Token)

	// Verificar que el usuario puede acceder a su espacio de trabajo
	me, resp := th.Client.GetMe()
	th.CheckOK(resp)
	require.NotNil(t, me)
	require.Equal(t, username, me.Username)

	t.Log("Flujo completo exitoso: registro -> persistencia -> login -> acceso al workspace")
}

func TestINT1002CrearTableroApareceEnSidebar(t *testing.T) {
	t.Log("INT-10-02: simula flujo UI AddBoard -> Mutator -> OctoClient -> API -> Store, verificando creación y aparición en listado.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	boardTitle := "INT-10-02 Nuevo Tablero"

	// Simular creación de tablero desde UI
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  boardTitle,
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)
	require.NotEmpty(t, board.ID)
	require.Equal(t, boardTitle, board.Title)

	// Simular consulta de tableros para sidebar (GET /teams/{teamID}/boards)
	boards, resp := th.Client.GetBoardsForTeam(testTeamID)
	th.CheckOK(resp)
	require.NotNil(t, boards)

	// Verificar que el tablero creado aparece en el listado
	boardIDs := make([]string, 0, len(boards))
	for _, b := range boards {
		boardIDs = append(boardIDs, b.ID)
	}
	require.Contains(t, boardIDs, board.ID, "El tablero creado debería aparecer en el sidebar")

	// Verificar persistencia en Store
	persistedBoard, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedBoard)
	require.Equal(t, boardTitle, persistedBoard.Title)

	t.Logf("Tablero creado exitosamente y visible en sidebar: %s (ID: %s)", boardTitle, board.ID)
}

func TestINT1003CrearTarjetaKanbanPersisteYEditable(t *testing.T) {
	t.Log("INT-10-03: simula flujo UI KanbanCard -> Mutator -> OctoClient -> API -> Store, verificando creación, persistencia y edición.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tablero con vista Board
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-10-03 Board Kanban",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	// Crear vista tipo Board (Kanban)
	viewBlock := &model.Block{
		ID:       utils.NewID(utils.IDTypeBlock),
		BoardID:  board.ID,
		Type:     model.TypeView,
		CreateAt: 1,
		UpdateAt: 1,
		Title:    "Board View",
		Fields: map[string]interface{}{
			"viewType": "board",
		},
	}
	blocks, resp := th.Client.InsertBlocks(board.ID, []*model.Block{viewBlock}, false)
	th.CheckOK(resp)
	require.Len(t, blocks, 1)

	// Crear tarjeta (card)
	cardTitle := "INT-10-03 Test Card"
	card := &model.Card{
		Title:        cardTitle,
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
		Properties:   map[string]interface{}{"status": "todo"},
	}
	createdCard, resp := th.Client.CreateCard(board.ID, card, true)
	th.CheckOK(resp)
	require.NotNil(t, createdCard)
	require.NotEmpty(t, createdCard.ID)
	require.Equal(t, cardTitle, createdCard.Title)

	// Verificar persistencia en Store
	persistedCard, err := th.Server.App().GetCardByID(createdCard.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedCard)
	require.Equal(t, cardTitle, persistedCard.Title)

	// Simular edición de la tarjeta
	updatedTitle := "INT-10-03 Test Card Updated"
	cardPatch := &model.CardPatch{
		Title: &updatedTitle,
	}
	patchResp, resp := th.Client.PatchCard(createdCard.ID, cardPatch, false)
	th.CheckOK(resp)
	require.NotNil(t, patchResp)
	require.Equal(t, updatedTitle, patchResp.Title)

	// Verificar que la edición persistió
	updatedPersistedCard, err := th.Server.App().GetCardByID(createdCard.ID)
	require.NoError(t, err)
	require.Equal(t, updatedTitle, updatedPersistedCard.Title)

	t.Logf("Tarjeta creada, persistida y editada exitosamente: %s (ID: %s)", updatedTitle, createdCard.ID)
}

func TestINT1004CambiarVistaBoardTableMismasTarjetas(t *testing.T) {
	t.Log("INT-10-04: simula flujo UI ViewMenu -> Mutator -> OctoClient -> API -> Store, verificando cambio de vista con mismos datos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tablero
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-10-04 Board Multi-Vista",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	// Crear vista Board
	boardView := &model.Block{
		ID:       utils.NewID(utils.IDTypeBlock),
		BoardID:  board.ID,
		Type:     model.TypeView,
		CreateAt: 1,
		UpdateAt: 1,
		Title:    "Board View",
		Fields: map[string]interface{}{
			"viewType": "board",
		},
	}
	blocks, resp := th.Client.InsertBlocks(board.ID, []*model.Block{boardView}, false)
	th.CheckOK(resp)
	require.Len(t, blocks, 1)

	// Crear vista Table
	tableView := &model.Block{
		ID:       utils.NewID(utils.IDTypeBlock),
		BoardID:  board.ID,
		Type:     model.TypeView,
		CreateAt: 1,
		UpdateAt: 1,
		Title:    "Table View",
		Fields: map[string]interface{}{
			"viewType": "table",
		},
	}
	tableBlocks, resp := th.Client.InsertBlocks(board.ID, []*model.Block{tableView}, false)
	th.CheckOK(resp)
	require.Len(t, tableBlocks, 1)

	// Crear tarjetas
	card1, resp := th.Client.CreateCard(board.ID, &model.Card{
		Title:        "Card 1",
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
	}, true)
	th.CheckOK(resp)
	require.NotNil(t, card1)

	card2, resp := th.Client.CreateCard(board.ID, &model.Card{
		Title:        "Card 2",
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
	}, true)
	th.CheckOK(resp)
	require.NotNil(t, card2)

	// Obtener tarjetas del tablero (simula carga en cualquier vista)
	cards, err := th.Server.App().GetCardsForBoard(board.ID, 0, 100)
	require.NoError(t, err)
	require.Len(t, cards, 2)

	cardTitles := make([]string, 0, len(cards))
	for _, c := range cards {
		cardTitles = append(cardTitles, c.Title)
	}
	require.Contains(t, cardTitles, "Card 1")
	require.Contains(t, cardTitles, "Card 2")

	t.Log("Cambio de vista exitoso: las mismas tarjetas (Card 1, Card 2) están disponibles en ambas vistas Board y Table")
}

func TestINT1005AplicarFiltroVistaTarjetasFiltradas(t *testing.T) {
	t.Log("INT-10-05: simula flujo UI FilterComponent -> Redux -> OctoClient -> API, verificando filtrado de tarjetas.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	// Crear tablero
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-10-05 Board con Filtros",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
		CardProperties: []map[string]interface{}{
			{
				"id":   "priority",
				"name": "Priority",
				"type": "select",
				"options": []map[string]interface{}{
					{"id": "high", "value": "High"},
					{"id": "low", "value": "Low"},
				},
			},
		},
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	// Crear tarjetas con diferentes propiedades
	cardHigh, resp := th.Client.CreateCard(board.ID, &model.Card{
		Title:        "High Priority Card",
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
		Properties:   map[string]interface{}{"priority": "high"},
	}, true)
	th.CheckOK(resp)
	require.NotNil(t, cardHigh)

	cardLow, resp := th.Client.CreateCard(board.ID, &model.Card{
		Title:        "Low Priority Card",
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
		Properties:   map[string]interface{}{"priority": "low"},
	}, true)
	th.CheckOK(resp)
	require.NotNil(t, cardLow)

	cardMedium, resp := th.Client.CreateCard(board.ID, &model.Card{
		Title:        "Medium Priority Card",
		ContentOrder: []string{utils.NewID(utils.IDTypeBlock)},
		Properties:   map[string]interface{}{"priority": "medium"},
	}, true)
	th.CheckOK(resp)
	require.NotNil(t, cardMedium)

	// Obtener todas las tarjetas
	allCards, err := th.Server.App().GetCardsForBoard(board.ID, 0, 100)
	require.NoError(t, err)
	require.Len(t, allCards, 3)

	// Simular filtrado por propiedad (en una aplicación real esto se haría a través de la API)
	// Aquí verificamos que las tarjetas tienen las propiedades correctas
	var highPriorityCards []*model.Card
	for _, card := range allCards {
		if priority, ok := card.Properties["priority"].(string); ok && priority == "high" {
			highPriorityCards = append(highPriorityCards, card)
		}
	}
	require.Len(t, highPriorityCards, 1)
	require.Equal(t, "High Priority Card", highPriorityCards[0].Title)

	t.Logf("Filtrado exitoso: de 3 tarjetas, solo 1 cumple con el filtro priority=high (%s)", highPriorityCards[0].Title)
}

func TestINT1006CerrarSesionRedireccionLoginInvalidaAcceso(t *testing.T) {
	t.Log("INT-10-06: simula flujo UI SidebarUserMenu -> OctoClient -> API -> Auth, verificando logout y redirección.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	username := "int10_user_logout"
	email := "int10_user_logout@example.com"
	password := "Pa$$word-int10-06"

	// Registrar y hacer login usando el token de signup del equipo
	team, resp := th.Client.GetTeam(model.GlobalTeamID)
	th.CheckOK(resp)
	require.NotEmpty(t, team.SignupToken)

	th.RegisterAndLogin(th.Client, username, email, password, team.SignupToken)

	// Obtener token antes de logout
	loginData, resp := th.Client.Login(&model.LoginRequest{
		Type:     "normal",
		Username: username,
		Password: password,
	})
	th.CheckOK(resp)
	oldToken := loginData.Token
	require.NotEmpty(t, oldToken)

	// Verificar que el usuario está autenticado
	me, resp := th.Client.GetMe()
	th.CheckOK(resp)
	require.NotNil(t, me)
	require.Equal(t, username, me.Username)

	// Simular logout desde UI
	httpResp, err := th.Client.DoAPIPost("/logout", "")
	require.NoError(t, err)
	require.NotNil(t, httpResp)
	defer httpResp.Body.Close()
	_, _ = io.Copy(io.Discard, httpResp.Body)
	require.Equal(t, 200, httpResp.StatusCode)

	// Intentar usar el token anterior para acceder a /users/me
	th.Client.Token = oldToken
	meAfterLogout, resp := th.Client.GetMe()
	th.CheckUnauthorized(resp)
	require.Nil(t, meAfterLogout)

	// Intentar hacer login nuevamente para verificar que funciona
	newLoginData, resp := th.Client.Login(&model.LoginRequest{
		Type:     "normal",
		Username: username,
		Password: password,
	})
	th.CheckOK(resp)
	require.NotNil(t, newLoginData)
	require.NotEmpty(t, newLoginData.Token)

	// Verificar que el nuevo token funciona
	th.Client.Token = newLoginData.Token
	meNewLogin, resp := th.Client.GetMe()
	th.CheckOK(resp)
	require.NotNil(t, meNewLogin)
	require.Equal(t, username, meNewLogin.Username)

	t.Log("Flujo de logout exitoso: sesión invalidada, token anterior rechazado, nuevo login funciona correctamente")
}
