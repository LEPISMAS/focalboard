#!/bin/bash
# run_tests_api.sh - Ejecuta las pruebas unitarias y muestra la cobertura del modulo api

# Obtener la ruta del directorio del script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SERVER_DIR="$( cd "$SCRIPT_DIR/../.." && pwd )"

cd "$SERVER_DIR"

echo "=== Ejecutando pruebas unitarias para server/api/tests ==="
go test -v -coverprofile=coverage.out -coverpkg=github.com/mattermost/focalboard/server/api ./api/tests

echo ""
echo "=== Resumen de cobertura de sentencias ==="
go tool cover -func=coverage.out | tail -n 1
