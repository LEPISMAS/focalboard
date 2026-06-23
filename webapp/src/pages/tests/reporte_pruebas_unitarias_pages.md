# Reporte de Pruebas Unitarias - Componente "pages"

Este reporte detalla las pruebas unitarias creadas para los archivos del componente "pages" del módulo "aplicación web" (ubicados en `webapp/src/pages`). Se logró una cobertura de sentencias general del **93.99%**, superando el objetivo establecido de 90%.

## Cuadro de Cobertura y Pruebas Realizadas

La siguiente tabla resume la cobertura de código por archivo e indica los casos de prueba implementados.

| Archivo | Cobertura de Sentencias (%) | Cobertura de Ramas (%) | Cobertura de Funciones (%) | Pruebas y Casos Cubiertos |
| :--- | :---: | :---: | :---: | :--- |
| **changePasswordPage.tsx** | 100.0% | 92.3% | 100.0% | - Renderizado básico y redirección cuando no hay usuario.<br>- Renderizado del formulario con usuario.<br>- Cambios en inputs.<br>- Envío exitoso con respuesta del cliente.<br>- Manejo de errores al fallar el cambio. |
| **errorPage.tsx** | 100.0% | 80.0% | 100.0% | - Renderizado por defecto.<br>- Renderizado de error `team-undefined` y redirección.<br>- Renderizado de error `board-not-found` y retorno a home.<br>- Redirección inmediata para `not-logged-in`. |
| **loginPage.tsx** | 97.05% | 88.88% | 100.0% | - Redirección si ya está autenticado.<br>- Renderizado de formulario de login.<br>- Cambios en inputs.<br>- Login exitoso con redirección (por defecto o con query `r`).<br>- Login fallido con mensaje de error. |
| **registerPage.tsx** | 100.0% | 92.85% | 100.0% | - Redirección si ya está autenticado.<br>- Renderizado de formulario de registro.<br>- Cambios en inputs.<br>- Registro y login exitosos con redirección.<br>- Manejo de token de registro (`signupToken`).<br>- Manejo de error 401 y otros errores devueltos por el cliente. |
| **welcomePage.tsx** | 85.71% | 68.88% | 83.33% | - Visualización de la página de onboarding.<br>- Proceso tras hacer clic en "Take a tour" u omitirlo.<br>- Omisión automática al volver a visitar.<br>- Redirecciones en presencia de parámetros query `r`. |
| **backwardCompatibilityQueryParamsRedirect.tsx** | 100.0% | 100.0% | 100.0% | - Verificación de retorno de valor `null`. |
| **boardPage.tsx** | 88.60% | 79.0% | 90.32% | - Renderizado en modo lectura/escritura y lectura de datos.<br>- Registro y manejo de eventos websocket (actualización incremental de bloques, tableros, miembros, reconexión y seguimiento).<br>- Visualización del diálogo de confirmación para unirse a tablero privado.<br>- Visualización y cierre de advertencia web móvil.<br>- Renderizado en modo sólo lectura. |
| **setWindowTitleAndIcon.tsx** | 100.0% | 100.0% | 100.0% | - Configuración del título a "Focalboard" y favicon indefinido si no hay tablero.<br>- Configuración de título e ícono de tablero activo.<br>- Combinación de título de tablero y vista activa en el documento. |
| **teamToBoardAndViewRedirect.tsx** | 100.0% | 86.95% | 100.0% | - Redirección al último tablero visitado o primer tablero de categorías.<br>- Redirección a la última vista visitada o primera vista de tablero. |
| **undoRedoHotKeys.tsx** | 100.0% | 100.0% | 100.0% | - Retorno de renderizado nulo.<br>- Desencadenamiento de deshacer (con y sin descripción, y mensaje cuando no se puede deshacer).<br>- Desencadenamiento de rehacer (con y sin descripción, y mensaje cuando no se puede rehacer). |
| **websocketConnection.tsx** | 100.0% | 88.88% | 100.0% | - Registro de escuchadores websocket al montar y remover al desmontar.<br>- Visualización del banner de error de conexión tras un retardo de 5s en estado "close".<br>- Ocultamiento correcto si el estado vuelve a ser "open" antes de cumplirse el tiempo. |
| **Total Global Pages** | **93.99%** | **82.67%** | **94.52%** | **Cobertura general del componente pages.** |

## Ejecución de Pruebas

Para ejecutar estas pruebas de manera local y ver los resultados detallados de la cobertura, puede utilizar cualquiera de los scripts provistos:
- **En Linux / macOS:** `bash webapp/src/pages/tests/run_tests_pages.sh`
- **En Windows:** `webapp/src/pages/tests/run_tests_pages.bat`
