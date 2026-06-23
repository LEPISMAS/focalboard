@echo off
REM Batch script to run all unit tests for the pages component on Windows

cd /d "%~dp0..\.."
npx.cmd jest src/pages/tests --collectCoverageFrom=src/pages/**/*.ts --collectCoverageFrom=src/pages/**/*.tsx --collectCoverage --watchAll=false
