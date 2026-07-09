@echo off
echo Ejecutando flujo INT-05 Comparticion de Tableros
cd /d "%~dp0"
go test -v . -run TestINT05 -tags sqlite_json
