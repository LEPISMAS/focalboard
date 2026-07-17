package tests

import (
	"bytes"
	"net"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/app"
	"github.com/mattermost/focalboard/server/auth"
	"github.com/mattermost/focalboard/server/model"
	"github.com/mattermost/focalboard/server/services/audit"
	"github.com/mattermost/focalboard/server/services/config"
	"github.com/mattermost/focalboard/server/services/metrics"
	"github.com/mattermost/focalboard/server/services/store/mockstore"
	"github.com/mattermost/focalboard/server/services/webhook"
	"github.com/mattermost/focalboard/server/ws"

	mmModel "github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

type mockPermissionsService struct {
	hasPermissionTo        func(userID string, permission *mmModel.Permission) bool
	hasPermissionToTeam    func(userID, teamID string, permission *mmModel.Permission) bool
	hasPermissionToChannel func(userID, channelID string, permission *mmModel.Permission) bool
	hasPermissionToBoard   func(userID, boardID string, permission *mmModel.Permission) bool
}

func (m *mockPermissionsService) HasPermissionTo(userID string, permission *mmModel.Permission) bool {
	if m.hasPermissionTo != nil {
		return m.hasPermissionTo(userID, permission)
	}
	return true
}

func (m *mockPermissionsService) HasPermissionToTeam(userID, teamID string, permission *mmModel.Permission) bool {
	if m.hasPermissionToTeam != nil {
		return m.hasPermissionToTeam(userID, teamID, permission)
	}
	return true
}

func (m *mockPermissionsService) HasPermissionToChannel(userID, channelID string, permission *mmModel.Permission) bool {
	if m.hasPermissionToChannel != nil {
		return m.hasPermissionToChannel(userID, channelID, permission)
	}
	return true
}

func (m *mockPermissionsService) HasPermissionToBoard(userID, boardID string, permission *mmModel.Permission) bool {
	if m.hasPermissionToBoard != nil {
		return m.hasPermissionToBoard(userID, boardID, permission)
	}
	return true
}

type TestAPIHelper struct {
	Router             *mux.Router
	Store              *mockstore.MockStore
	PermissionsService *mockPermissionsService
	SingleUserToken    string
}

func SetupTestAPI(t *testing.T) (*TestAPIHelper, func()) {
	ctrl := gomock.NewController(t)

	cfg := &config.Configuration{
		EnablePublicSharedBoards: true,
	}

	store := mockstore.NewMockStore(ctrl)

	// Default mock expectations for background store calls made by the app layer.
	// These cover category management, board members, subscriptions, etc.
	defaultCategory := model.Category{ID: "default-category-id", Name: "Boards", UserID: "single-user", TeamID: "test-team", Type: model.CategoryTypeSystem}
	store.EXPECT().GetUserCategoryBoards(gomock.Any(), gomock.Any()).Return([]model.CategoryBoards{
		{Category: defaultCategory, BoardMetadata: []model.CategoryBoardMetadata{}},
	}, nil).AnyTimes()
	store.EXPECT().GetMembersForBoard(gomock.Any()).Return([]*model.BoardMember{}, nil).AnyTimes()
	store.EXPECT().CreateCategory(gomock.Any()).Return(nil).AnyTimes()
	store.EXPECT().GetCategory(gomock.Any()).Return(&defaultCategory, nil).AnyTimes()
	store.EXPECT().AddUpdateCategoryBoard(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	store.EXPECT().GetMemberForBoard(gomock.Any(), gomock.Any()).Return(&model.BoardMember{}, nil).AnyTimes()
	store.EXPECT().GetUsersByTeam(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.User{}, nil).AnyTimes()

	logger, _ := mlog.NewLogger()
	authService := auth.New(cfg, store, nil)
	wsserver := ws.NewServer(authService, "TESTTOKEN", false, logger, store)
	metricsService := metrics.NewMetrics(metrics.InstanceInfo{})
	webhookClient := webhook.NewClient(cfg, logger)

	perms := &mockPermissionsService{}

	appServices := app.Services{
		Auth:             authService,
		Store:            store,
		Metrics:          metricsService,
		Webhook:          webhookClient,
		Logger:           logger,
		SkipTemplateInit: true,
		Permissions:      perms,
	}

	application := app.New(cfg, wsserver, appServices)

	auditService, _ := audit.NewAudit()

	singleUserToken := "TESTTOKEN"
	testAPI := api.NewAPI(application, singleUserToken, "local", perms, logger, auditService)

	router := mux.NewRouter()
	testAPI.RegisterRoutes(router)

	tearDown := func() {
		application.Shutdown()
		if logger != nil {
			_ = logger.Shutdown()
		}
		ctrl.Finish()
	}

	return &TestAPIHelper{
		Router:             router,
		Store:              store,
		PermissionsService: perms,
		SingleUserToken:    singleUserToken,
	}, tearDown
}

func (h *TestAPIHelper) NewRequest(method, url string, body []byte) *http.Request {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", "Bearer "+h.SingleUserToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	return req
}

// Dummy type to implement net.Conn if needed for local unix connections (not used unless testing admin set password)
type DummyUnixConn struct {
	net.Conn
}

func (c *DummyUnixConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "local", Net: "unix"}
}
