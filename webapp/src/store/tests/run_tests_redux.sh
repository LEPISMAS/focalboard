#!/bin/bash
# Script to run all unit tests for the store (redux) component

# Get the directory of this script and navigate to webapp root
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "$DIR/../.."

# Run jest on the store/redux tests with coverage restricted to the store components
npx jest src/store/tests --collectCoverageFrom=src/store/**/*.ts --collectCoverageFrom=src/store/**/*.tsx --collectCoverage --watchAll=false