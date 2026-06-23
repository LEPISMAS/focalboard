#!/bin/bash
# Script to run all unit tests for the pages component

# Get the directory of this script and navigate to webapp root
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR/../.."

# Run jest on the pages tests with coverage restricted to the pages components
npx jest src/pages/tests --collectCoverageFrom=src/pages/**/*.ts --collectCoverageFrom=src/pages/**/*.tsx --collectCoverage --watchAll=false
