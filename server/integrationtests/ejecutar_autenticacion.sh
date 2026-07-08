#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-01 Autenticacion"
cd "$(dirname "$0")"
go test -v . -run TestINT01
