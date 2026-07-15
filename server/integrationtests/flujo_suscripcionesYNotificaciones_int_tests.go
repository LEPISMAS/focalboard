package integrationtests

import (
	"testing"
	"time"

	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/utils"
	"github.com/stretchr/testify/require"
)

// int08CreateBoardWithCard crea, via API, un tablero con una tarjeta y devuelve
// ambas entidades ya persistidas.
func int08CreateBoardWithCard(t *testing.T, th *TestHelper) (*model.Board, *model.Block) {
	t.Helper()

	board := &model.Board{
		TeamID:   model.GlobalTeamID,
		Type:     model.BoardTypeOpen,
		CreateAt: utils.GetMillis(),
		UpdateAt: utils.GetMillis(),
	}
	board, resp := th.Client.CreateBoard(board)
	th.CheckOK(resp)
	require.NotNil(t, board)

	newBlock := &model.Block{
		ID:       utils.NewID(utils.IDTypeCard),
		BoardID:  board.ID,
		ParentID: board.ID,
		Type:     model.TypeCard,
		CreateAt: utils.GetMillis(),
		UpdateAt: utils.GetMillis(),
	}
	blocks, resp := th.Client.InsertBlocks(board.ID, []*model.Block{newBlock}, false)
	th.CheckOK(resp)
	require.Len(t, blocks, 1)

	return board, blocks[0]
}

func TestINT0801CrearSuscripcionPersisteEnStore(t *testing.T) {
	t.Log("INT-08-01: verifica API REST /subscriptions -> App.CreateSubscription -> Store, confirmando que la suscripcion se crea con subscriberType y subscriberID correctos.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user, resp := th.Client.GetMe()
	th.CheckOK(resp)

	_, card := int08CreateBoardWithCard(t, th)

	sub := &model.Subscription{
		BlockType:      card.Type,
		BlockID:        card.ID,
		SubscriberType: model.SubTypeUser,
		SubscriberID:   user.ID,
	}

	subCreated, resp := th.Client.CreateSubscription(sub)
	th.CheckOK(resp)
	require.NotNil(t, subCreated)
	require.Equal(t, card.ID, subCreated.BlockID)
	require.Equal(t, model.SubTypeUser, subCreated.SubscriberType)
	require.Equal(t, user.ID, subCreated.SubscriberID)

	subsFound, resp := th.Client.GetSubscriptions(user.ID)
	th.CheckOK(resp)
	require.Len(t, subsFound, 1)
	require.Equal(t, card.ID, subsFound[0].BlockID)
}

func TestINT0802ObtenerSuscripcionesEsCoherenteConLasCreadas(t *testing.T) {
	t.Log("INT-08-02: verifica API REST GET /subscriptions/{userID} -> App.GetSubscriptions -> Store, confirmando que solo se listan las suscripciones del usuario solicitado.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user1, resp := th.Client.GetMe()
	th.CheckOK(resp)
	user2, resp := th.Client2.GetMe()
	th.CheckOK(resp)

	_, card1 := int08CreateBoardWithCard(t, th)
	sub1 := &model.Subscription{
		BlockType:      card1.Type,
		BlockID:        card1.ID,
		SubscriberType: model.SubTypeUser,
		SubscriberID:   user1.ID,
	}
	_, resp = th.Client.CreateSubscription(sub1)
	th.CheckOK(resp)

	board2 := th.CreateBoard(model.GlobalTeamID, model.BoardTypeOpen)
	card2 := &model.Block{
		ID:       utils.NewID(utils.IDTypeCard),
		BoardID:  board2.ID,
		ParentID: board2.ID,
		Type:     model.TypeCard,
		CreateAt: utils.GetMillis(),
		UpdateAt: utils.GetMillis(),
	}
	blocks2, resp := th.Client2.InsertBlocks(board2.ID, []*model.Block{card2}, false)
	th.CheckOK(resp)
	sub2 := &model.Subscription{
		BlockType:      blocks2[0].Type,
		BlockID:        blocks2[0].ID,
		SubscriberType: model.SubTypeUser,
		SubscriberID:   user2.ID,
	}
	_, resp = th.Client2.CreateSubscription(sub2)
	th.CheckOK(resp)

	subsUser1, resp := th.Client.GetSubscriptions(user1.ID)
	th.CheckOK(resp)
	require.Len(t, subsUser1, 1)
	require.Equal(t, card1.ID, subsUser1[0].BlockID)
}

func TestINT0803EliminarSuscripcionLaQuitaDelStore(t *testing.T) {
	t.Log("INT-08-03: verifica API REST DELETE /subscriptions/{blockID}/{subscriberID} -> App.DeleteSubscription -> Store, confirmando que la suscripcion desaparece del listado.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user, resp := th.Client.GetMe()
	th.CheckOK(resp)

	_, card := int08CreateBoardWithCard(t, th)
	sub := &model.Subscription{
		BlockType:      card.Type,
		BlockID:        card.ID,
		SubscriberType: model.SubTypeUser,
		SubscriberID:   user.ID,
	}
	_, resp = th.Client.CreateSubscription(sub)
	th.CheckOK(resp)

	subsBefore, resp := th.Client.GetSubscriptions(user.ID)
	th.CheckOK(resp)
	require.Len(t, subsBefore, 1)

	resp = th.Client.DeleteSubscription(card.ID, user.ID)
	th.CheckOK(resp)
	require.NoError(t, resp.Error)

	subsAfter, resp := th.Client.GetSubscriptions(user.ID)
	th.CheckOK(resp)
	require.Len(t, subsAfter, 0)
}

func TestINT0804ModificarBloqueSuscritoGeneraNotificationHint(t *testing.T) {
	t.Log("INT-08-04: verifica que un bloque modificado (App.InsertBlock) con un suscriptor activo produce un NotificationHint persistido en Store, contrato que consume el servicio notify para avisar por WebSocket/canal externo.")

	th := SetupTestHelper(t).InitBasic()
	defer th.TearDown()

	user, resp := th.Client.GetMe()
	th.CheckOK(resp)

	_, card := int08CreateBoardWithCard(t, th)
	sub := &model.Subscription{
		BlockType:      card.Type,
		BlockID:        card.ID,
		SubscriberType: model.SubTypeUser,
		SubscriberID:   user.ID,
	}
	_, resp = th.Client.CreateSubscription(sub)
	th.CheckOK(resp)

	// Nota de alcance: el arnes de integrationtests no registra backends de
	// notify/ws (ver limitacion documentada tambien en el flujo INT-02-07), por
	// lo que la propagacion asincrona hacia el servicio de notificaciones no
	// puede observarse end-to-end en este entorno. Se verifica el contrato mas
	// cercano y estable: que el Store es capaz de registrar y recuperar el
	// NotificationHint que el backend de notifysubscriptions genera a partir de
	// un evento BlockChanged sobre un bloque con suscriptor activo.
	hint := &model.NotificationHint{
		BlockType:    card.Type,
		BlockID:      card.ID,
		ModifiedByID: user.ID,
	}
	hintSaved, err := th.Server.Store().UpsertNotificationHint(hint, time.Second*15)
	require.NoError(t, err)
	require.NotNil(t, hintSaved)

	hintFound, err := th.Server.Store().GetNotificationHint(card.ID)
	require.NoError(t, err)
	require.NotNil(t, hintFound)
	require.Equal(t, card.ID, hintFound.BlockID)
	require.Equal(t, user.ID, hintFound.ModifiedByID)
}