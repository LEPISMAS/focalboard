#!/bin/bash
cd "$(dirname "$0")/../.." || exit
echo "Running tests for ws/tests..."
go test -coverpkg=github.com/mattermost/focalboard/server/ws -coverprofile=coverage.out ./ws/tests
go tool cover -func=coverage.out
