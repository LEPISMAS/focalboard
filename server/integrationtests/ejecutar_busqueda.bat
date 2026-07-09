@echo off
echo Ejecutando flujo INT-09 Búsqueda
cd /d "%~dp0"
go test -v . -run TestINT09
