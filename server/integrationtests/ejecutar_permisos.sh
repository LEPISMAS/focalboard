#!/usr/bin/env bash
set -e
echo "Ejecutando pruebas de integracion para Flujo 2: Permisos (INT-04)"
cd "$(dirname "$0")"
go test -v . -run TestINT04 -tags sqlite_json
