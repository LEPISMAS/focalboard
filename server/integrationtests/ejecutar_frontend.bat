@echo off
echo Ejecutando flujo INT-10 Frontend
cd /d "%~dp0"
go test -v . -run TestINT10
