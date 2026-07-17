#!/usr/bin/env bash
set -e

echo "Ejecutando flujo INT-08 Suscripciones y Notificaciones"
cd "$(dirname "$0")"
go test -v . -run TestINT08 -tags sqlite_json
