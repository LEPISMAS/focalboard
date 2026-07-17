#!/usr/bin/env bash
set -e
echo "Ejecutando pruebas de integracion para Flujo 1: Gestion de Tarjetas y Bloques (INT-03)"
cd "$(dirname "$0")"
go test -v . -run TestINT03 -tags sqlite_json
