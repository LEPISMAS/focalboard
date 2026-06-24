package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

func TestChannelsEndpoints(t *testing.T) {
	t.Run("GET /teams/{teamID}/channels/{channelID} returns Not Implemented in standalone mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = false

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/channels/channel1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusNotImplemented, resp.Code)
		require.Contains(t, resp.Body.String(), "not permitted in standalone mode")
	})

	t.Run("GET /teams/{teamID}/channels/{channelID} returns channel in plugin mode", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		th.API.MattermostAuth = true

		channel := &mmModel.Channel{
			Id:     "channel1",
			TeamId: "team1",
			Type:   mmModel.ChannelTypeOpen,
		}

		th.Store.EXPECT().GetChannel("team1", "channel1").Return(channel, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/teams/team1/channels/channel1", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "channel1")
	})
}
