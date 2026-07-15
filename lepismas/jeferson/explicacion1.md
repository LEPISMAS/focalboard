### Explicación de pruebas unitarias del componente ws (websockets) del servidor

Se realizaron pruebas unitarias para el componente `ws` (websockets) del módulo `server` de Focalboard, encargado de la conexión y actualización en tiempo real de los tableros colaborativos.

El trabajo se centró en los archivos `server.go`, `plugin_adapter.go`, `plugin_adapter_client.go` y `plugin_adapter_cluster.go`, usando `gomock` y `httptest`. Las pruebas quedaron externalizadas en `server/ws/tests/`.

Se cubrieron los siguientes escenarios:
- Conexión y autenticación de clientes websocket (tokens válidos e inválidos, Single User).
- Suscripción y desuscripción a equipos y bloques.
- Difusión (broadcast) de eventos como Blocks, Boards, Categories, Members y Subscriptions.
- Eventos de clúster en la versión de plugin (propagación y deserialización de mensajes entre nodos).

La cobertura obtenida fue:

- `ws.Server`: ~82-85%
- `ws.PluginAdapter`: ~95-100%
- **Total del paquete `ws`: ~90%**

Con esto se cumplió el objetivo de cobertura de sentencias solicitado (90%).

Las pruebas se pueden ejecutar localmente con:
- Windows: `.\server\ws\tests\run_tests_ws.bat`
- Linux/Mac: `./server/ws/tests/run_tests_ws.sh`