#!/bin/bash
# Script to run Focalboard model layer unit tests and show coverage

# Exit on error
set -e

# Navigate to the server root directory relative to the script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Running unit tests for Focalboard model package..."
go test -v -coverpkg=github.com/mattermost/focalboard/server/model -coverprofile=model_coverage.out ./...

echo ""
echo "=========================================="
echo "         Coverage Summary Report"
echo "=========================================="
go tool cover -func=model_coverage.out
