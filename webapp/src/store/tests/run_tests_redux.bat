@echo off
REM Batch script to run all unit tests for the store (redux) component on Windows

cd /d "%~dp0..\.."
npx.cmd jest src/store/tests --collectCoverageFrom=src/store/**/*.ts --collectCoverageFrom=src/store/**/*.tsx --collectCoverage --watchAll=false
pause