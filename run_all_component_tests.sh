#!/bin/bash
# run_all_component_tests.sh
# Ejecuta las pruebas de todos los componentes del proyecto Focalboard.
# Debe ejecutarse desde la raíz del repositorio.

ROOT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

PASS_COUNT=0
FAIL_COUNT=0
FAILED_COMPONENTS=()

run_component() {
    local name="$1"
    local script="$2"

    echo ""
    echo "╔══════════════════════════════════════════════════════╗"
    echo "║  Componente: $name"
    echo "╚══════════════════════════════════════════════════════╝"
    bash "$ROOT_DIR/$script"
    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        echo "✔  $name: PASSED"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo "✘  $name: FAILED (exit code $exit_code)"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        FAILED_COMPONENTS+=("$name")
    fi

    return $exit_code
}

echo "╔══════════════════════════════════════════════════════╗"
echo "║      FOCALBOARD - Suite de Pruebas de Componentes   ║"
echo "╚══════════════════════════════════════════════════════╝"

run_component "webapp/src/store (Redux)"  "webapp/src/store/tests/run_tests_redux.sh"
run_component "webapp/src/pages"          "webapp/src/pages/tests/run_tests_pages.sh"
run_component "server/api"               "server/api/tests/run_tests_api.sh"
run_component "server/ws"                "server/ws/tests/run_tests_ws.sh"

echo ""
echo "╔══════════════════════════════════════════════════════╗"
echo "║                   RESUMEN FINAL                     ║"
echo "╚══════════════════════════════════════════════════════╝"
echo "  Componentes ejecutados : $((PASS_COUNT + FAIL_COUNT))"
echo "  ✔ Exitosos             : $PASS_COUNT"
echo "  ✘ Fallidos             : $FAIL_COUNT"

if [ ${#FAILED_COMPONENTS[@]} -gt 0 ]; then
    echo ""
    echo "  Componentes con fallos:"
    for comp in "${FAILED_COMPONENTS[@]}"; do
        echo "    - $comp"
    done
    echo ""
    exit 1
fi

echo ""
echo "  Todas las pruebas de componentes pasaron correctamente."
exit 0
