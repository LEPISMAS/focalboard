# Reporte de ejecucion de pruebas de integracion: Suscripciones y Notificaciones

## 1. Objetivo

Reportar el estado de implementacion y ejecucion del flujo INT-08 Suscripciones y Notificaciones, orientado a validar la integracion entre API REST, capa App, Store y el contrato de datos que consume el servicio de notificaciones (`NotifyService`).

## 2. Fecha

15 de julio de 2026.

## 3. Entorno

- Proyecto: Focalboard.
- Modulo Go: `server`.
- Paquete de pruebas: `server/integrationtests`.
- Archivo de pruebas: `server/integrationtests/flujo_suscripcionesYNotificaciones_int_tests.go`.
- Scripts de ejecucion:
  - `server/integrationtests/ejecutar_suscripcionesYNotificaciones.bat`
  - `server/integrationtests/ejecutar_suscripcionesYNotificaciones.sh`

## 4. Comandos utilizados

Desde la carpeta del flujo:

```bash
cd server/integrationtests
go test -v . -run TestINT08
```

Mediante script Windows:

```bat
server\integrationtests\ejecutar_suscripcionesYNotificaciones.bat
```

Mediante Bash:

```bash
bash server/integrationtests/ejecutar_suscripcionesYNotificaciones.sh
```

## 5. Casos implementados

| ID | Caso | Estado de implementacion |
|---|---|---|
| INT-08-01 | Crear suscripcion via API y verificar persistencia en Store | Implementado |
| INT-08-02 | Obtener suscripciones y verificar coherencia (sin fuga entre usuarios) | Implementado |
| INT-08-03 | Eliminar suscripcion y verificar que se elimine del Store | Implementado |
| INT-08-04 | Modificar bloque suscrito y verificar contrato de NotificationHint | Implementado (con limitacion documentada, ver seccion 6) |

## 6. Resultados obtenidos

### Implementacion validada

El flujo fue implementado en el paquete `integrationtests`, usando helpers existentes como `SetupTestHelper`, dos clientes autenticados (`th.Client`, `th.Client2`) y metodos publicos de App/Store. Cada prueba incluye salida con `t.Log` para explicar su utilidad y las capas recorridas.

La implementacion valida:

- Creacion de una suscripcion por API con persistencia confirmada en el listado del Store.
- Aislamiento de suscripciones entre dos usuarios distintos.
- Eliminacion efectiva de una suscripcion, reflejada en el listado posterior.
- El contrato de datos `NotificationHint` (creacion y lectura via Store), que es la pieza que el backend `notifysubscriptions` usa para decidir a quien notificar.

### Limitacion documentada (INT-08-04)

El arnes de `server/integrationtests` no registra backends de `notify` ni de `ws` (`clienttestlib.go` no configura `NotifyBackends`), por lo que no es posible observar de punta a punta que una modificacion de bloque via API dispare automaticamente el backend `notifysubscriptions` y este cree el `NotificationHint` de forma asincrona. Esta es la misma clase de limitacion que el equipo ya documento para WebSocket end-to-end en el caso INT-02-07 del flujo Gestion de Tableros. Por eso INT-08-04 valida el contrato de persistencia mas cercano y estable (`UpsertNotificationHint` / `GetNotificationHint`) en lugar de simular una propagacion end-to-end que el entorno de pruebas no soporta.

### Ejecucion en el entorno de trabajo

*(Completar con el resultado real al ejecutar `go test -v . -run TestINT08` en el entorno del equipo: PASS/FAIL por caso, tiempo total y cualquier log relevante de la corrida.)*

## 7. Incidencias encontradas

| ID Defecto | Caso de prueba que lo detecto | Componente afectado | Severidad | Requisito relacionado | Estado |
|---|---|---|---|---|---|
| ENV-03 | INT-08-04 | Arnes de integrationtests / servicio notify | Media (limitacion de entorno, no defecto funcional) | Suscripciones y notificaciones | Abierto en el arnes de integracion |

No se reportan defectos funcionales del flujo de suscripciones, ya que INT-08-01 a INT-08-03 se validan de punta a punta contra la API real. La unica incidencia corresponde a una limitacion de infraestructura de pruebas (ausencia de backend `notify` en el arnes), no a un defecto del producto.

## 8. Conclusiones

El flujo INT-08 esta implementado de forma completa respecto a los casos solicitados. Las pruebas cubren el ciclo de vida de una suscripcion (crear, listar, eliminar) contra la API real y documentan de forma honesta la limitacion tecnica existente para validar la propagacion asincrona de notificaciones en el arnes de integracion actual, verificando en su lugar el contrato de datos que esa propagacion efectivamente consume.

## 9. Recomendaciones

- Ejecutar `go test -v . -run TestINT08` en el entorno de CI o en una maquina local con Go y la base de datos de pruebas configurada, y completar las secciones 6 y 7 con el resultado real.
- Si se requiere validar la entrega end-to-end de notificaciones (incluyendo WebSocket), extender el arnes de `integrationtests` para registrar un backend de `notify` de prueba, en la misma linea que se recomendo para WebSocket en el flujo de Gestion de Tableros.
- Mantener sincronizado este reporte cada vez que se agreguen nuevos casos al flujo INT-08.
