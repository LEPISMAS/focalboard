@echo off
echo Ejecutando flujo INT-01 Autenticacion
cd /d "%~dp0"
go test -v . -run TestINT01 -tags sqlite_json
