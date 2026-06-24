package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/focalboard/server/api"
	"github.com/mattermost/focalboard/server/model"
	mmModel "github.com/mattermost/mattermost/server/public/model"
	filestoreMocks "github.com/mattermost/mattermost/server/v8/platform/shared/filestore/mocks"
)

func TestFilesEndpoints(t *testing.T) {
	t.Run("GET /files/teams/{teamID}/{boardID}/{filename}/info returns file info", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		fileInfo := &mmModel.FileInfo{
			Id:   "fileinfo-id",
			Name: "test-file.png",
		}
		th.Store.EXPECT().GetFileInfo(gomock.Any()).Return(fileInfo, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/files/teams/team1/board1/7fileinfo-id.png/info", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "test-file.png")
	})

	t.Run("GET /files/teams/{teamID}/{boardID}/{filename} serves file", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1"}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)

		fileInfo := &mmModel.FileInfo{
			Id:       "fileinfo-id",
			Name:     "test-file.png",
			Path:     "team1/board1/7fileinfo-id.png",
			MimeType: "image/png",
			Size:     100,
		}
		th.Store.EXPECT().GetFileInfo("fileinfo-id").Return(fileInfo, nil)

		mockedReadCloseSeek := &filestoreMocks.ReadCloseSeeker{}
		// Mock Reader behavior
		mockedReadCloseSeek.On("Read", mock.Anything).Return(0, nil)
		mockedReadCloseSeek.On("Seek", mock.Anything, mock.Anything).Return(int64(0), nil)
		mockedReadCloseSeek.On("Close").Return(nil)

		th.FilesBackend.On("FileExists", "team1/board1/7fileinfo-id.png").Return(true, nil)
		th.FilesBackend.On("Reader", "team1/board1/7fileinfo-id.png").Return(mockedReadCloseSeek, nil)

		req, _ := http.NewRequest(http.MethodGet, "/api/v2/files/teams/team1/board1/7fileinfo-id.png", nil)
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Equal(t, "image/png", resp.Header().Get("Content-Type"))
	})

	t.Run("POST /teams/{teamID}/boards/{boardID}/files uploads file", func(t *testing.T) {
		th, tearDown := SetupTestAPIHelper(t, TestSingleUserToken)
		defer tearDown()

		board := &model.Board{ID: "board1", TeamID: "team1", IsTemplate: false}
		th.Store.EXPECT().GetBoard("board1").Return(board, nil)
		th.Store.EXPECT().SaveFileInfo(gomock.Any()).Return(nil)

		th.FilesBackend.On("WriteFile", mock.Anything, mock.Anything).Return(int64(100), nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test-file.png")
		require.NoError(t, err)
		_, err = part.Write([]byte("fake-image-bytes"))
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/api/v2/teams/team1/board1/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set(api.HeaderRequestedWith, api.HeaderRequestedWithXML)
		req.Header.Set("Authorization", "Bearer "+TestSingleUserToken)

		resp := doRequest(th.Router, req)
		require.Equal(t, http.StatusOK, resp.Code)
		require.Contains(t, resp.Body.String(), "fileId")
	})

	t.Run("JSON helper functions in files.go", func(t *testing.T) {
		r1 := strings.NewReader(`{"fileId": "helper-id"}`)
		res1, err := api.FileUploadResponseFromJSON(r1)
		require.NoError(t, err)
		require.Equal(t, "helper-id", res1.FileID)

		r2 := strings.NewReader(`{"id": "file-info-id", "name": "helper-file.png"}`)
		res2, err := api.FileInfoResponseFromJSON(r2)
		require.NoError(t, err)
		require.Equal(t, "file-info-id", res2.Id)
		require.Equal(t, "helper-file.png", res2.Name)
	})
}
