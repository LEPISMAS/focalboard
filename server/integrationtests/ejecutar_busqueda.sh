#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-09 Búsqueda"
cd "$(dirname "$0")"
go test -v . -run TestINT09 -tags sqlite_json
