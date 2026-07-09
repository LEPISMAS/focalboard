#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-05 Comparticion de Tableros"
cd "$(dirname "$0")"
go test -v . -run TestINT05 -tags sqlite_json
