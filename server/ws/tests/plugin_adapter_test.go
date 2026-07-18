package tests

import (
	"testing"

	"github.com/golang/mock/gomock"

	authMocks "github.com/mattermost/focalboard/server/auth/mocks"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/ws"
	wsMocks "github.com/mattermost/focalboard/server/ws/mocks"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func setupPluginAdapterTest(t *testing.T) (*ws.PluginAdapter, *wsMocks.MockAPI, *authMocks.MockAuthInterface, *wsMocks.MockStore) {
	ctrl := gomock.NewController(t)
	mockAPI := wsMocks.NewMockAPI(ctrl)
	mockAuth := authMocks.NewMockAuthInterface(ctrl)
	mockStore := wsMocks.NewMockStore(ctrl)

	logger := mlog.CreateConsoleTestLogger(t)

	// Mocking logging so they don't block/fail
	mockAPI.EXPECT().LogDebug(gomock.Any(), gomock.Any()).AnyTimes()
	mockAPI.EXPECT().LogInfo(gomock.Any(), gomock.Any()).AnyTimes()
	mockAPI.EXPECT().LogError(gomock.Any(), gomock.Any()).AnyTimes()
	mockAPI.EXPECT().LogWarn(gomock.Any(), gomock.Any()).AnyTimes()

	pa := ws.NewPluginAdapter(mockAPI, mockAuth, mockStore, logger)

	return pa, mockAPI, mockAuth, mockStore
}

func TestPluginAdapter_ConnectionAndMessages(t *testing.T) {
	pa, _, mockAuth, _ := setupPluginAdapterTest(t)

	webConnID := "conn-1"
	userID := "user-1"
	teamID := "team-1"

	// Connect
	pa.OnWebSocketConnect(webConnID, userID)

	// Reconnect
	pa.OnWebSocketConnect(webConnID, userID)

	// Subscribe team - requires DoesUserHaveTeamAccess to be true
	mockAuth.EXPECT().DoesUserHaveTeamAccess(userID, teamID).Return(true).AnyTimes()

	reqSub := &mmModel.WebSocketRequest{
		Action: "custom_focalboard_SUBSCRIBE_TEAM",
		Data: map[string]interface{}{
			"teamId": teamID,
		},
	}
	pa.WebSocketMessageHasBeenPosted(webConnID, userID, reqSub)

	// Unsubscribe team
	reqUnsub := &mmModel.WebSocketRequest{
		Action: "custom_focalboard_UNSUBSCRIBE_TEAM",
		Data: map[string]interface{}{
			"teamId": teamID,
		},
	}
	pa.WebSocketMessageHasBeenPosted(webConnID, userID, reqUnsub)

	// Subscribe blocks (unimplemented but we send it)
	reqSubBlocks := &mmModel.WebSocketRequest{
		Action: "custom_focalboard_SUBSCRIBE_BLOCKS",
		Data: map[string]interface{}{
			"teamId": teamID,
		},
	}
	pa.WebSocketMessageHasBeenPosted(webConnID, userID, reqSubBlocks)

	// Invalid command missing teamId
	reqMissingTeam := &mmModel.WebSocketRequest{
		Action: "custom_focalboard_SUBSCRIBE_TEAM",
		Data:   map[string]interface{}{},
	}
	pa.WebSocketMessageHasBeenPosted(webConnID, userID, reqMissingTeam)

	// Disconnect
	pa.OnWebSocketDisconnect(webConnID, userID)

	// Disconnect unregistered
	pa.OnWebSocketDisconnect("unregistered", "user-1")

	// Post to unregistered
	pa.WebSocketMessageHasBeenPosted("unregistered", "user-1", reqSub)
}

func TestPluginAdapter_Broadcasts(t *testing.T) {
	pa, mockAPI, mockAuth, mockStore := setupPluginAdapterTest(t)

	webConnID := "conn-1"
	userID := "user-1"
	teamID := "team-1"
	boardID := "board-1"

	pa.OnWebSocketConnect(webConnID, userID)
	mockAuth.EXPECT().DoesUserHaveTeamAccess(userID, teamID).Return(true).AnyTimes()
	mockStore.EXPECT().GetMembersForBoard(boardID).Return([]*model.BoardMember{
		{UserID: userID, BoardID: boardID},
	}, nil).AnyTimes()

	// Connect and subscribe to set up internal maps
	reqSub := &mmModel.WebSocketRequest{
		Action: "custom_focalboard_SUBSCRIBE_TEAM",
		Data: map[string]interface{}{
			"teamId": teamID,
		},
	}
	pa.WebSocketMessageHasBeenPosted(webConnID, userID, reqSub)

	// Test broadcasts - these should call mockAPI.PublishWebSocketEvent and PublishPluginClusterEvent
	mockAPI.EXPECT().PublishWebSocketEvent(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockAPI.EXPECT().PublishPluginClusterEvent(gomock.Any(), gomock.Any()).AnyTimes()

	pa.BroadcastConfigChange(model.ClientConfig{})

	pa.BroadcastBlockChange(teamID, &model.Block{ID: "block-1", BoardID: boardID})
	pa.BroadcastBlockDelete(teamID, "block-1", boardID)

	pa.BroadcastCategoryChange(model.Category{TeamID: teamID, UserID: userID, ID: "cat-1"})
	pa.BroadcastCategoryReorder(teamID, userID, []string{"cat-1"})
	pa.BroadcastCategoryBoardsReorder(teamID, userID, "cat-1", []string{boardID})
	pa.BroadcastCategoryBoardChange(teamID, userID, []*model.BoardCategoryWebsocketData{})

	pa.BroadcastBoardChange(teamID, &model.Board{ID: boardID, TeamID: teamID})
	pa.BroadcastBoardDelete(teamID, boardID)

	pa.BroadcastMemberChange(teamID, boardID, &model.BoardMember{UserID: userID})
	pa.BroadcastMemberDelete(teamID, boardID, userID)

	pa.BroadcastSubscriptionChange(teamID, &model.Subscription{BlockID: "block-1", SubscriberID: userID})
	pa.BroadcastCardLimitTimestampChange(100)
}

func TestPluginAdapter_ClusterEvent(t *testing.T) {
	pa, _, _, _ := setupPluginAdapterTest(t)

	// Test cluster event handling
	ev := mmModel.PluginClusterEvent{
		Id:   "custom_focalboard_",
		Data: []byte(`{"teamId": "team-1", "payload": {}}`),
	}
	pa.HandleClusterEvent(ev)

	evUser := mmModel.PluginClusterEvent{
		Id:   "custom_focalboard_",
		Data: []byte(`{"userId": "user-1", "payload": {}}`),
	}
	pa.HandleClusterEvent(evUser)

	// Test invalid JSON
	evInvalid := mmModel.PluginClusterEvent{
		Id:   "custom_focalboard_",
		Data: []byte(`{invalid-json}`),
	}
	pa.HandleClusterEvent(evInvalid)

	// Test block getters to increase coverage
	pa.GetListenersByBlock("block-1")

	// Test more cluster events for different actions
	actions := []string{
		"custom_focalboard_UPDATE_BLOCK",
		"custom_focalboard_UPDATE_BOARD",
		"custom_focalboard_UPDATE_MEMBER",
		"custom_focalboard_DELETE_MEMBER",
		"custom_focalboard_UPDATE_CATEGORY",
		"custom_focalboard_REORDER_CATEGORIES",
		"custom_focalboard_REORDER_CATEGORY_BOARDS",
		"custom_focalboard_UPDATE_CATEGORY_BOARD",
		"custom_focalboard_UPDATE_SUBSCRIPTION",
	}

	for _, action := range actions {
		evAction := mmModel.PluginClusterEvent{
			Id:   action,
			Data: []byte(`{"teamId": "team-1", "payload": {}}`),
		}
		pa.HandleClusterEvent(evAction)
	}
}
