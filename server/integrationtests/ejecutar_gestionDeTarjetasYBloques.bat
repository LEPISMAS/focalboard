@echo off
echo Ejecutando pruebas de integracion para Flujo 1: Gestion de Tarjetas y Bloques (INT-03)
cd /d "%~dp0"
go test -v . -run TestINT03
