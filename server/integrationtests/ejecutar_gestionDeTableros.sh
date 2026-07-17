#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-02 Gestion de Tableros"
cd "$(dirname "$0")"
go test -v . -run TestINT02 -tags sqlite_json
