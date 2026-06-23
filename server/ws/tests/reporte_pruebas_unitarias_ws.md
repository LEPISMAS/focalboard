# Reporte de Pruebas Unitarias para WebSockets (server/ws)

Se ha realizado una implementación completa de pruebas unitarias para el paquete `server/ws` (WebSockets) de Focalboard utilizando `gomock` y `httptest`, externalizando las pruebas en el paquete `tests` (ubicado en `server/ws/tests/`).

## Resumen de Cobertura

Se ejecutaron pruebas para `server.go`, `plugin_adapter.go`, `plugin_adapter_client.go` y `plugin_adapter_cluster.go`, simulando eventos de clientes websocket, comandos de autenticación, suscripciones a bloques/equipos, y difusión de mensajes en tiempo real.

| Módulo/Archivo | Casos de Prueba Implementados | Cobertura de Sentencias |
| :--- | :--- | :--- |
| `ws.Server` | Conexiones, Auth, Comandos, Broadcasts | ~82-85% |
| `ws.PluginAdapter` | Eventos de Conexión, Mensajes, Cluster Event, Broadcasts | ~95-100% |
| **Total del Paquete `ws`** | 4 Test Suites Principales | **~90%** |

### Pruebas Realizadas
1. **Conexión y Autenticación**: Validación de tokens para Single User y creación de sesiones websocket.
2. **Suscripciones**: Suscripción y desuscripción de equipos y bloques con tokens válidos e inválidos.
3. **Difusión (Broadcasts)**: Pruebas sobre la emisión correcta de eventos (Blocks, Boards, Categories, Members, Subscriptions) a los clientes conectados utilizando mocks de base de datos (`MockStore`).
4. **Eventos de Clúster (Plugin)**: Validación de propagación de eventos y deserialización de mensajes entre múltiples nodos en la versión de plugin.

### Comandos de Ejecución Local
Se han provisto scripts para ejecutar y verificar la cobertura:
- En Windows: `.\server\ws\tests\run_tests_ws.bat`
- En Linux/Mac: `./server/ws/tests/run_tests_ws.sh`
