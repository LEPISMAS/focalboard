package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/app"
	"github.com/mattermost/focalboard/server/auth"
	"github.com/mattermost/focalboard/server/services/audit"
	"github.com/mattermost/focalboard/server/services/config"
	"github.com/mattermost/focalboard/server/services/metrics"
	"github.com/mattermost/focalboard/server/services/store/mockstore"
	"github.com/mattermost/focalboard/server/services/webhook"
	"github.com/mattermost/focalboard/server/ws"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	filestoreMocks "github.com/mattermost/mattermost/server/v8/platform/shared/filestore/mocks"
)

const (
	TestSingleUserToken = "test-single-user-token"
)

type TestAPIHelper struct {
	API          *api.API
	App          *app.App
	Store        *mockstore.MockStore
	Router       *mux.Router
	Permissions  *MockPermissionsService
	FilesBackend *filestoreMocks.FileBackend
	Audit        *audit.Audit
}

func SetupTestAPIHelper(t *testing.T, singleUserToken string) (*TestAPIHelper, func()) {
	ctrl := gomock.NewController(t)
	cfg := &config.Configuration{
		SessionExpireTime:        2592000,
		SessionRefreshTime:       1800,
		EnablePublicSharedBoards: true,
	}
	
	store := mockstore.NewMockStore(ctrl)
	// Add default mocks for websocket side-effects and message posting
	store.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	
	filesBackend := &filestoreMocks.FileBackend{}
	
	// Create real auth but since it calls store, store is mocked
	authService := auth.New(cfg, store, nil)
	
	logger := mlog.CreateConsoleTestLogger(t)
	
	sessionToken := "test-session-token"
	wsserver := ws.NewServer(authService, sessionToken, false, logger, store)
	webhookClient := webhook.NewClient(cfg, logger)
	metricsService := metrics.NewMetrics(metrics.InstanceInfo{})
	
	permissions := &MockPermissionsService{}
	
	appServices := app.Services{
		Auth:             authService,
		Store:            store,
		FilesBackend:     filesBackend,
		Webhook:          webhookClient,
		Metrics:          metricsService,
		Logger:           logger,
		SkipTemplateInit: true,
		Permissions:      permissions,
	}
	
	appInstance := app.New(cfg, wsserver, appServices)
	
	auditLogger, _ := audit.NewAudit()
	
	// Instantiate API with given singleUserToken
	testAPI := api.NewAPI(appInstance, singleUserToken, "native", permissions, logger, auditLogger)
	
	router := mux.NewRouter()
	testAPI.RegisterRoutes(router)
	
	// Register Admin routes too
	testAPI.RegisterAdminRoutes(router)
	
	tearDown := func() {
		appInstance.Shutdown()
		_ = auditLogger.Shutdown()
		ctrl.Finish()
	}
	
	return &TestAPIHelper{
		API:          testAPI,
		App:          appInstance,
		Store:        store,
		Router:       router,
		Permissions:  permissions,
		FilesBackend: filesBackend,
		Audit:        auditLogger,
	}, tearDown
}

func doRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
