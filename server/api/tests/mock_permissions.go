package tests

import (
	mmModel "github.com/mattermost/mattermost/server/public/model"
)

type MockPermissionsService struct {
	HasPermissionToFunc        func(userID string, permission *mmModel.Permission) bool
	HasPermissionToTeamFunc    func(userID, teamID string, permission *mmModel.Permission) bool
	HasPermissionToChannelFunc func(userID, channelID string, permission *mmModel.Permission) bool
	HasPermissionToBoardFunc   func(userID, boardID string, permission *mmModel.Permission) bool
}

func (m *MockPermissionsService) HasPermissionTo(userID string, permission *mmModel.Permission) bool {
	if m.HasPermissionToFunc != nil {
		return m.HasPermissionToFunc(userID, permission)
	}
	return true
}

func (m *MockPermissionsService) HasPermissionToTeam(userID, teamID string, permission *mmModel.Permission) bool {
	if m.HasPermissionToTeamFunc != nil {
		return m.HasPermissionToTeamFunc(userID, teamID, permission)
	}
	return true
}

func (m *MockPermissionsService) HasPermissionToChannel(userID, channelID string, permission *mmModel.Permission) bool {
	if m.HasPermissionToChannelFunc != nil {
		return m.HasPermissionToChannelFunc(userID, channelID, permission)
	}
	return true
}

func (m *MockPermissionsService) HasPermissionToBoard(userID, boardID string, permission *mmModel.Permission) bool {
	if m.HasPermissionToBoardFunc != nil {
		return m.HasPermissionToBoardFunc(userID, boardID, permission)
	}
	return true
}
