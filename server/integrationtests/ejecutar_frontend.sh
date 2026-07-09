#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-10 Frontend"
cd "$(dirname "$0")"
go test -v . -run TestINT10
