@echo off
echo Ejecutando pruebas de integracion para Flujo 2: Permisos (INT-04)
cd /d "%~dp0"
go test -v . -run TestINT04
