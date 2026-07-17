@echo off
echo Ejecutando flujo INT-02 Gestion de Tableros
cd /d "%~dp0"
go test -v . -run TestINT02 -tags sqlite_json
