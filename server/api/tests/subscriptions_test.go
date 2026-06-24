package tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
)

func TestSubscriptionsEndpoints(t *testing.T) {
	t.Run("POST /subscriptions creates subscription", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		// UserID must match subscriberID (which is "single-user" from TestSingleUserToken)
		sub := &model.Subscription{
			BlockID:        "block1",
			BlockType:      model.TypeBoard,
			SubscriberID:   "single-user",
			SubscriberType: model.SubTypeUser,
		}

		// check for valid block in handler: a.app.GetBlockByID
		block := &model.Block{ID: "block1", BoardID: "board1"}
		th.Store.EXPECT().GetBlock("block1").Return(block, nil).Times(1)

		th.Store.EXPECT().CreateSubscription(gomock.Any()).Return(sub, nil)

		body := `{"blockId": "block1", "blockType": "board", "subscriberId": "single-user", "subscriberType": "user"}`
		req, _ := http.NewRequest(http.MethodPost, "/api/v2/subscriptions", bytes.NewBufferString(body))
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "single-user")
	})

	t.Run("DELETE /subscriptions/{blockID}/{subscriberID} deletes subscription", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		sub := &model.Subscription{
			BlockID:        "block1",
			BlockType:      model.TypeBoard,
			SubscriberID:   "single-user",
			SubscriberType: model.SubTypeUser,
		}
		th.Store.EXPECT().GetSubscription("block1", "single-user").Return(sub, nil)
		th.Store.EXPECT().DeleteSubscription("block1", "single-user").Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/api/v2/subscriptions/block1/single-user", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "{}", resp.Body.String())
	})

	t.Run("GET /subscriptions/{subscriberID} returns subscriptions", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		subs := []*model.Subscription{
			{BlockID: "block1", SubscriberID: "single-user"},
		}
		th.Store.EXPECT().GetSubscriptions("single-user").Return(subs, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/subscriptions/single-user", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "block1")
	})
}
