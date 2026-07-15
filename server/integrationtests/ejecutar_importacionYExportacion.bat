@echo off
echo Ejecutando flujo INT-07 Importacion y Exportacion
cd /d "%~dp0"
go test -v . -run TestINT07
