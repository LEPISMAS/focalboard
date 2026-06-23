@echo off
cd %~dp0\..\..
echo Running tests for ws and ws/tests...
go test -coverpkg=github.com/mattermost/focalboard/server/ws -coverprofile=coverage.out ./ws ./ws/tests
go tool cover -func=coverage.out
