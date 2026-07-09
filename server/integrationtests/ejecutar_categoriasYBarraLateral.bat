@echo off
echo Ejecutando flujo INT-06 Categorias y Barra Lateral
cd /d "%~dp0"
go test -v . -run TestINT06 -tags sqlite_json
