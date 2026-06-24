@echo off
rem run_tests_api.bat - Ejecuta las pruebas unitarias y muestra la cobertura del modulo api en Windows

cd %~dp0\..\..

echo === Ejecutando pruebas unitarias para server/api/tests ===
go test -v -coverprofile=coverage.out -coverpkg=github.com/mattermost/focalboard/server/api ./api/tests

echo.
echo === Resumen de cobertura de sentencias ===
go tool cover -func=coverage.out
