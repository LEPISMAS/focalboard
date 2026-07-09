#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-06 Categorias y Barra Lateral"
cd "$(dirname "$0")"
go test -v . -run TestINT06 -tags sqlite_json
