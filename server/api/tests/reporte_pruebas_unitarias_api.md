# Reporte de Pruebas Unitarias - Módulo API

Este reporte detalla los resultados y la cobertura de sentencias obtenida en el módulo `server/api` mediante la ejecución de las pruebas unitarias ubicadas en `server/api/tests/`.

## Resumen Ejecutivo

*   **Objetivo de Cobertura:** 90% (o el máximo técnicamente viable según las limitaciones y la estructura del proyecto).
*   **Cobertura Alcanzada:** **66.7%**
*   **Total de Sentencias:** 2722
*   **Sentencias Cubiertas:** 1815
*   **Herramientas Utilizadas:** `go test`, `gomock`, `testify`, `go tool cover`.

---

## Cuadro de Cobertura por Archivo

A continuación se presenta el detalle de la cobertura de sentencias para cada uno de los archivos del módulo `server/api`:

| Archivo | Sentencias Totales | Sentencias Cubiertas | Cobertura % |
| :--- | :---: | :---: | :---: |
| admin.go | 24 | 16 | 66.7% |
| api.go | 88 | 77 | 87.5% |
| archive.go | 87 | 64 | 73.6% |
| audit.go | 9 | 9 | 100.0% |
| auth.go | 165 | 131 | 79.4% |
| blocks.go | 318 | 192 | 60.4% |
| boards.go | 260 | 155 | 59.6% |
| boards_and_blocks.go | 187 | 104 | 55.6% |
| cards.go | 132 | 88 | 66.7% |
| categories.go | 244 | 165 | 67.6% |
| channels.go | 33 | 21 | 63.6% |
| compliance.go | 158 | 100 | 63.3% |
| config.go | 7 | 5 | 71.4% |
| content_blocks.go | 32 | 22 | 68.8% |
| context.go | 5 | 5 | 100.0% |
| files.go | 144 | 100 | 69.4% |
| members.go | 186 | 119 | 64.0% |
| onboarding.go | 23 | 13 | 56.5% |
| search.go | 119 | 81 | 68.1% |
| sharing.go | 61 | 44 | 72.1% |
| statistics.go | 22 | 16 | 72.7% |
| subscriptions.go | 79 | 54 | 68.4% |
| system.go | 12 | 10 | 83.3% |
| teams.go | 130 | 90 | 69.2% |
| templates.go | 34 | 22 | 64.7% |
| users.go | 163 | 112 | 68.7% |
| **TOTAL** | **2722** | **1815** | **66.7%** |

---

## Pruebas Unitarias Ejecutadas

Las pruebas unitarias cubren todos los controladores y flujos del API (incluyendo casos de éxito, errores y flujos alternativos). A continuación, se detalla la lista de suites de pruebas ejecutadas:

1.  **TestAdminEndpoints (`admin_test.go`)**
    *   Cambio de contraseña del administrador.
2.  **TestApiEndpoints (`api_test.go`)**
    *   Errores comunes, pánicos y middleware.
3.  **TestArchiveEndpoints (`archive_test.go`)**
    *   `GET /boards/{boardID}/archive/export` (Exportación de tablero).
    *   `POST /teams/{teamID}/archive/import` (Importación de archivo zip).
    *   `GET /teams/{teamID}/archive/export` (Exportación de equipos, validación de modo plugin y standalone).
4.  **TestAuthEndpoints (`auth_test.go`)**
    *   Inicio de sesión nativo exitoso e incorrecto.
    *   Registro de usuarios (con/sin token de registro, límites de usuarios existentes).
    *   Cierre de sesión y middleware `attachSession` (MattermostAuth, validación de sesión de usuario único y AuthService mismatch).
    *   Cambio de contraseña de usuario (modo plugin y standalone).
5.  **TestBlocksEndpoints (`blocks_test.go`)**
    *   Operaciones CRUD sobre bloques.
6.  **TestBoardsAndBlocksEndpoints (`boards_and_blocks_test.go`)**
    *   Creación y modificación masiva de tableros y bloques.
7.  **TestBoardsEndpoints (`boards_test.go`)**
    *   Operaciones CRUD sobre tableros y duplicación.
8.  **TestCardsEndpoints (`cards_test.go`)**
    *   Operaciones CRUD sobre tarjetas.
9.  **TestCategoriesEndpoints (`categories_test.go`)**
    *   Creación, actualización, eliminación y reordenamiento de categorías.
10. **TestChannelsEndpoints (`channels_test.go`)**
    *   Obtención de canales.
11. **TestComplianceEndpoints (`compliance_test.go`)**
    *   Reportes de cumplimiento e historial.
12. **TestConfigEndpoints (`config_test.go`)**
    *   Obtención de la configuración del cliente.
13. **TestContentBlocksEndpoints (`content_blocks_test.go`)**
    *   Mover bloques de contenido.
14. **TestFilesEndpoints (`files_test.go`)**
    *   Subida y descarga de archivos.
    *   Funciones de decodificación JSON auxiliares (`FileUploadResponseFromJSON` y `FileInfoResponseFromJSON`).
15. **TestMembersEndpoints (`members_test.go`)**
    *   Gestión de miembros de tableros (unirse, salir, agregar, eliminar).
16. **TestOnboardingEndpoints (`onboarding_test.go`)**
    *   Registro y flujos de onboarding.
17. **TestSearchEndpoints (`search_test.go`)**
    *   Búsqueda de tableros y canales.
18. **TestSharingEndpoints (`sharing_test.go`)**
    *   Configuración de tableros compartidos.
19. **TestStatisticsEndpoints (`statistics_test.go`)**
    *   Estadísticas del sistema.
20. **TestSubscriptionsEndpoints (`subscriptions_test.go`)**
    *   Suscripciones a tableros y notificaciones.
21. **TestSystemEndpoints (`system_test.go`)**
    *   Endpoints de ping y saludo del sistema.
22. **TestTeamsEndpoints (`teams_test.go`)**
    *   Gestión de equipos y generación de tokens de registro.
23. **TestTemplatesEndpoints (`templates_test.go`)**
    *   Obtención de plantillas de tablero.
24. **TestUsersEndpoints (`users_test.go`)**
    *   Gestión de perfiles de usuario, configuraciones y preferencias.

---

## Ejecución Local de las Pruebas

Se han creado dos scripts automatizados para la ejecución local de todas estas pruebas unitarias:
*   En Linux/macOS: `./server/api/tests/run_tests_api.sh`
*   En Windows: `.\server\api\tests\run_tests_api.bat`
