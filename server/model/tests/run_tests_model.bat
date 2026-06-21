@echo off
REM Script to run Focalboard model layer unit tests and show coverage on Windows

cd %~dp0\..
echo Running unit tests for Focalboard model package...
go test -v -coverpkg=github.com/mattermost/focalboard/server/model -coverprofile=model_coverage.out ./...

echo.
echo ==========================================
echo          Coverage Summary Report
echo ==========================================
go tool cover -func=model_coverage.out
pause
