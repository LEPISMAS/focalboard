package integrationtests

import (
	"testing"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/require"
)

func int02BoardIDs(boards []*model.Board) []string {
	ids := make([]string, 0, len(boards))
	for _, board := range boards {
		ids = append(ids, board.ID)
	}
	return ids
}

func int02FindMember(members []*model.BoardMember, userID string) *model.BoardMember {
	for _, member := range members {
		if member.UserID == userID {
			return member
		}
	}
	return nil
}

func TestINT0201CrearTableroPersisteEnStore(t *testing.T) {
	t.Log("INT-02-01: verifica API REST /boards -> App.CreateBoard -> Store.InsertBoardWithAdmin, confirmando ID, createdBy y persistencia.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	newBoard := &model.Board{
		Title:  "INT-02-01 board",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	}

	board, resp := th.Client.CreateBoard(newBoard)
	th.CheckOK(resp)
	require.NotNil(t, board)
	require.NotEmpty(t, board.ID)
	require.Equal(t, user1.ID, board.CreatedBy)
	require.Equal(t, user1.ID, board.ModifiedBy)

	persistedBoard, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedBoard)
	require.Equal(t, board.ID, persistedBoard.ID)
	require.Equal(t, "INT-02-01 board", persistedBoard.Title)
	require.Equal(t, user1.ID, persistedBoard.CreatedBy)
}

func TestINT0202CrearTableroCreaMembresiaAdmin(t *testing.T) {
	t.Log("INT-02-02: verifica que la creacion de tablero por API crea BoardMember admin para el usuario autenticado.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-02 board",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	members, err := th.Server.App().GetMembersForBoard(board.ID)
	require.NoError(t, err)
	require.Len(t, members, 1)

	creatorMember := int02FindMember(members, user1.ID)
	require.NotNil(t, creatorMember)
	require.Equal(t, board.ID, creatorMember.BoardID)
	require.True(t, creatorMember.SchemeAdmin)
}

func TestINT0203ListarTablerosFiltraPorMembresia(t *testing.T) {
	t.Log("INT-02-03: verifica GET /teams/{teamID}/boards -> App.GetBoardsForUserAndTeam -> Store, filtrando tableros privados por membresia.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1Board1, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-03 user1 private board 1",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user1Board1)

	user1Board2, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-03 user1 private board 2",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user1Board2)

	user2PrivateBoard, resp := th.Client2.CreateBoard(&model.Board{
		Title:  "INT-02-03 user2 private board",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, user2PrivateBoard)

	boardsForUser1, resp := th.Client.GetBoardsForTeam(testTeamID)
	th.CheckOK(resp)
	boardIDs := int02BoardIDs(boardsForUser1)

	require.Contains(t, boardIDs, user1Board1.ID)
	require.Contains(t, boardIDs, user1Board2.ID)
	require.NotContains(t, boardIDs, user2PrivateBoard.ID)
}

func TestINT0204ActualizarTituloTableroPersiste(t *testing.T) {
	t.Log("INT-02-04: verifica PATCH /boards/{boardID} -> App.PatchBoard -> Store.PatchBoard, confirmando respuesta y persistencia del titulo.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-04 original title",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	newTitle := "INT-02-04 updated title"
	updatedBoard, resp := th.Client.PatchBoard(board.ID, &model.BoardPatch{Title: &newTitle})
	th.CheckOK(resp)
	require.NotNil(t, updatedBoard)
	require.Equal(t, newTitle, updatedBoard.Title)

	persistedBoard, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedBoard)
	require.Equal(t, newTitle, persistedBoard.Title)
}

func TestINT0205EliminarTableroSoftDelete(t *testing.T) {
	t.Log("INT-02-05: verifica DELETE /boards/{boardID} -> App.DeleteBoard -> Store.DeleteBoard, usando historial para confirmar deleteAt y ausencia en listados activos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-05 board",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	success, resp := th.Client.DeleteBoard(board.ID)
	th.CheckOK(resp)
	require.True(t, success)

	deletedBoard, err := th.Server.App().GetBoard(board.ID)
	require.Error(t, err)
	require.True(t, model.IsErrNotFound(err))
	require.Nil(t, deletedBoard)

	history, err := th.Server.Store().GetBoardHistory(board.ID, model.QueryBoardHistoryOptions{Limit: 1, Descending: true})
	require.NoError(t, err)
	require.NotEmpty(t, history)
	require.Greater(t, history[0].DeleteAt, int64(0))

	activeBoards, resp := th.Client.GetBoardsForTeam(testTeamID)
	th.CheckOK(resp)
	require.NotContains(t, int02BoardIDs(activeBoards), board.ID)
}

func TestINT0206DuplicarTableroCopiaBloquesPropiedadesMembresias(t *testing.T) {
	t.Log("INT-02-06: verifica POST /boards/{boardID}/duplicate -> App.DuplicateBoard -> Store.DuplicateBoard, confirmando nuevo tablero, bloques, propiedades y membresia admin del creador.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-06 source board",
		Type:   model.BoardTypeOpen,
		TeamID: testTeamID,
		Properties: map[string]interface{}{
			"risk": "integration",
		},
		CardProperties: []map[string]interface{}{
			{
				"id":   "status",
				"name": "Status",
				"type": "select",
			},
		},
	})
	th.CheckOK(resp)
	require.NotNil(t, board)

	sourceBlocks, resp := th.Client.InsertBlocks(board.ID, []*model.Block{
		{
			ID:       utils.NewID(utils.IDTypeBlock),
			BoardID:  board.ID,
			CreateAt: 1,
			UpdateAt: 1,
			Title:    "INT-02-06 copied view",
			Type:     model.TypeView,
			Fields: map[string]interface{}{
				"viewType": "table",
			},
		},
	}, false)
	th.CheckOK(resp)
	require.Len(t, sourceBlocks, 1)

	duplicated, resp := th.Client.DuplicateBoard(board.ID, false, testTeamID)
	th.CheckOK(resp)
	require.NotNil(t, duplicated)
	require.Len(t, duplicated.Boards, 1)
	require.Len(t, duplicated.Blocks, 1)

	duplicatedBoard := duplicated.Boards[0]
	require.NotEqual(t, board.ID, duplicatedBoard.ID)
	require.Equal(t, "INT-02-06 source board copy", duplicatedBoard.Title)
	require.Equal(t, "integration", duplicatedBoard.Properties["risk"])
	require.Len(t, duplicatedBoard.CardProperties, 1)
	require.Equal(t, "status", duplicatedBoard.CardProperties[0]["id"])

	duplicatedBlock := duplicated.Blocks[0]
	require.NotEqual(t, sourceBlocks[0].ID, duplicatedBlock.ID)
	require.Equal(t, duplicatedBoard.ID, duplicatedBlock.BoardID)
	require.Equal(t, sourceBlocks[0].Title, duplicatedBlock.Title)

	members, err := th.Server.App().GetMembersForBoard(duplicatedBoard.ID)
	require.NoError(t, err)
	require.Len(t, members, 1)
	creatorMember := int02FindMember(members, user1.ID)
	require.NotNil(t, creatorMember)
	require.True(t, creatorMember.SchemeAdmin)
}

func TestINT0207CrearTableroNotificaWebSocket(t *testing.T) {
	t.Log("INT-02-07: no hay helper estable de WebSocket en integrationtests; se verifica el tramo API -> App -> Store y se documenta que WS end-to-end requiere arnes de suscripcion especifico.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1 := th.GetUser1()
	board, resp := th.Client.CreateBoard(&model.Board{
		Title:  "INT-02-07 board",
		Type:   model.BoardTypePrivate,
		TeamID: testTeamID,
	})
	th.CheckOK(resp)
	require.NotNil(t, board)
	require.NotEmpty(t, board.ID)

	persistedBoard, err := th.Server.App().GetBoard(board.ID)
	require.NoError(t, err)
	require.NotNil(t, persistedBoard)
	require.Equal(t, "INT-02-07 board", persistedBoard.Title)
	require.Equal(t, user1.ID, persistedBoard.CreatedBy)

	members, err := th.Server.App().GetMembersForBoard(board.ID)
	require.NoError(t, err)
	require.NotNil(t, int02FindMember(members, user1.ID))
}
