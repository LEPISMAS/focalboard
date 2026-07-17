@echo off
echo Ejecutando flujo INT-08 Suscripciones y Notificaciones
cd /d "%~dp0"
go test -v . -run TestINT08 -tags sqlite_json
