#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-07 Importacion y Exportacion"
cd "$(dirname "$0")"
go test -v . -run TestINT07 -tags sqlite_json
