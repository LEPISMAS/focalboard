# Argumentacion de pruebas de integracion: Suscripciones y Notificaciones

## 1. Objetivo del flujo

El flujo INT-08 tiene como objetivo comprobar que el sistema de suscripciones a bloques se integra correctamente con la API REST, la capa de negocio, el Store y el punto de entrada del servicio de notificaciones (`NotifyService`). Las suscripciones son la base de las notificaciones que Focalboard envia a los usuarios cuando algo cambia en un tablero al que estan atentos.

El objetivo principal es demostrar que crear, listar y eliminar suscripciones se comportan de forma correcta y persistente, y que el contrato de datos que consume el servicio de notificaciones (el `NotificationHint`) puede registrarse y recuperarse desde el Store.

## 2. Capas integradas

Las pruebas integran las siguientes capas:

- API REST: endpoints `POST /subscriptions`, `GET /subscriptions/{subscriberID}`, `DELETE /subscriptions/{blockID}/{subscriberID}`.
- Capa App: `CreateSubscription`, `GetSubscriptions`, `DeleteSubscription` (`server/app`).
- Store: persistencia de suscripciones y de `NotificationHint`.
- Servicio de notificaciones (`server/services/notify/notifysubscriptions`): consumidor de los `NotificationHint` generados a partir de eventos `BlockChanged`.
- Modelo: estructuras `Subscription`, `Block`, `NotificationHint`.

## 3. Estandares aplicados

- Trazabilidad por identificador: cada funcion conserva el codigo INT-08-01 a INT-08-04.
- Reutilizacion del arnes oficial: se emplean `SetupTestHelper`, clientes reales (`th.Client`, `th.Client2`) y metodos publicos ya usados en el repositorio (`CreateSubscription`, `GetSubscriptions`, `DeleteSubscription`).
- Pruebas contra API real: la creacion, consulta y eliminacion de suscripciones se ejecutan por HTTP.
- Verificacion de estado persistido: se confirma el resultado final consultando `App`/`Store`, no solo el codigo de respuesta HTTP.
- Aislamiento entre usuarios: INT-08-02 crea suscripciones con dos usuarios distintos para comprobar que el listado no mezcla datos entre ellos.
- Transparencia en limitaciones: INT-08-04 declara explicitamente que el arnes de `integrationtests` no registra backends de `notify`/`ws`, siguiendo el mismo criterio de honestidad tecnica que el equipo aplico previamente en el flujo INT-02-07 (WebSocket end-to-end en Gestion de Tableros).

## 4. Importancia desde risk-based testing

### Identificacion de riesgos

Riesgos de producto:

- Suscripciones creadas sin `subscriberType` o `subscriberID` validos, quedando huerfanas o inutilizables.
- Un usuario que ve o recibe suscripciones de otro usuario (fuga de informacion).
- Suscripciones que no se eliminan correctamente, generando notificaciones para bloques que el usuario ya no sigue.
- `NotificationHint` no generado o no persistido tras modificar un bloque, dejando al usuario sin notificacion de un cambio relevante.

Riesgos de proyecto:

- Cambios en el modelo `Subscription` o `NotificationHint` que rompan el contrato entre `App` y el servicio `notify` sin que ninguna prueba lo detecte.
- Ausencia de un arnes de integracion estable para el servicio de notificaciones, dificultando validar el flujo end-to-end (mismo riesgo ya identificado por el equipo para WebSocket en INT-02).

### Evaluacion probabilidad por impacto

| Caso | Riesgo principal | Probabilidad | Impacto | Nivel |
|---|---|---:|---:|---|
| INT-08-01 | Suscripcion creada sin persistencia correcta | Media | Alta | Alto |
| INT-08-02 | Fuga de suscripciones entre usuarios | Media | Critico | Critico |
| INT-08-03 | Suscripcion eliminada pero sigue notificando | Media | Alta | Alto |
| INT-08-04 | Contrato NotificationHint roto entre App y Store | Baja-Media | Alta | Alto |

### Priorizacion

El caso INT-08-02 tiene prioridad critica porque afecta confidencialidad y consistencia: un usuario nunca debe ver las suscripciones de otro. INT-08-01 y INT-08-03 son de prioridad alta porque sostienen el ciclo de vida completo de una suscripcion (creacion y eliminacion efectiva). INT-08-04 es de prioridad alta porque protege el contrato de datos que hace posible que el servicio de notificaciones funcione, aunque su validacion end-to-end completa (incluyendo el envio real por WebSocket) queda fuera del alcance del arnes de integracion actual.

### Mitigacion

La mitigacion se realiza ejercitando las operaciones reales de la API de suscripciones y verificando el resultado en `App`/`Store`, incluyendo un escenario con dos usuarios para descartar fugas de datos. Para el punto de contacto con el servicio de notificaciones, se verifica el contrato de persistencia del `NotificationHint`, que es la pieza que el backend `notifysubscriptions` usa realmente para decidir a quien notificar.

## 5. Justificacion de herramientas

| Herramienta | Por que se eligio | Alternativas consideradas | Limitacion conocida |
|---|---|---|---|
| Go test | Mecanismo nativo de pruebas del backend Go, permite ejecutar solo el flujo con `-run TestINT08`. | Herramientas HTTP externas como Postman/Newman. | Depende de que el entorno de base de datos de pruebas este disponible. |
| Testify require | Ya se usa en el repositorio, ofrece aserciones claras (`require.Len`, `require.Equal`). | `testing` puro. | Detiene el caso en el primer fallo critico. |
| TestHelper de integrationtests (`th.Client`, `th.Client2`) | Provee dos clientes autenticados distintos, necesarios para probar aislamiento entre usuarios. | Simular dos usuarios con un unico cliente reautenticado. | Hereda las limitaciones del arnes existente, en particular la ausencia de backends `notify`/`ws`. |
| Cliente HTTP del proyecto | Garantiza que las pruebas pasen por los endpoints reales de suscripciones. | Llamar directamente a `App.CreateSubscription`. | No cubre por si solo la entrega final de la notificacion. |
| App/Store publicos (incluyendo `UpsertNotificationHint`/`GetNotificationHint`) | Permiten verificar el contrato de datos que consume el servicio de notificaciones sin depender de un backend `notify` completo. | Levantar un backend `notify` real con canal de salida simulado. | Solo confirma el contrato Store, no la entrega final por WebSocket u otro canal. |
| Logs con `t.Log` | Explican la intencion de cada prueba en la salida estandar, utiles para el debate. | Comentarios internos solamente. | No sustituyen un reporte de ejecucion formal. |

## 6. Justificacion estrategica por caso

### INT-08-01 Crear suscripcion vía API y verificar persistencia en el Store

Es el caso base del flujo: sin una suscripcion correctamente persistida no existe ningun mecanismo de notificacion posterior. Se verifica tanto la respuesta de la API como el resultado del listado posterior.

### INT-08-02 Obtener suscripciones del usuario y verificar coherencia con las creadas

Este caso es estrategico por confidencialidad: se crean suscripciones con dos usuarios distintos (`th.Client` y `th.Client2`) y se confirma que el listado de un usuario no incluye suscripciones del otro, replicando el mismo criterio de aislamiento por membresia usado en el flujo de Gestion de Tableros.

### INT-08-03 Eliminar una suscripcion y verificar que se elimine del Store

Una suscripcion que no se elimina correctamente puede seguir generando notificaciones no deseadas. La prueba confirma que, tras eliminar, el listado de suscripciones del usuario queda vacio.

### INT-08-04 Modificar un bloque en un tablero suscrito y verificar que se genere un notification hint

El arnes de `server/integrationtests` no registra backends de `notify` ni de `ws` (misma limitacion que el equipo documento en INT-02-07), por lo que no es posible observar de forma end-to-end que una modificacion via API dispare automaticamente el backend `notifysubscriptions` y este cree el `NotificationHint`. Por eso se prueba el contrato mas cercano y estable: que el Store puede registrar (`UpsertNotificationHint`) y recuperar (`GetNotificationHint`) el hint que ese backend genera internamente a partir de un evento `BlockChanged`. Esta decision evita introducir una prueba fragil que dependa de infraestructura no disponible en el entorno de test, y deja la limitacion explicita en lugar de simularla como si funcionara de punta a punta.

## 7. Relacion con pruebas unitarias previas del componente ws

Las pruebas unitarias previas del componente `server/ws` validaban en aislamiento la conexion, autenticacion y difusion de eventos websocket usando mocks (`gomock`) para el Store y para clientes conectados. Esas pruebas confirman que el componente de transporte en tiempo real funciona correctamente de forma unitaria.

Las pruebas INT-08 amplian el alcance hacia el origen de esas notificaciones: verifican que el sistema de suscripciones (quien debe ser notificado) y el contrato de `NotificationHint` (que dispara la notificacion) funcionan correctamente sobre datos reales y persistentes. Ambos conjuntos de pruebas son complementarios: las unitarias de `ws` reducen el riesgo de errores en el transporte de la notificacion, mientras que las de integracion de este flujo reducen el riesgo de que la notificacion nunca llegue a generarse porque la suscripcion o el hint no se persistieron correctamente.

## 8. Conclusion argumentativa

El flujo INT-08 fue seleccionado porque las suscripciones son el mecanismo que decide quien debe ser notificado ante un cambio, y un defecto aqui puede significar que un usuario deje de enterarse de cambios relevantes en un tablero, o que reciba notificaciones de contenido que ya no sigue. La estrategia de prueba cubre el ciclo completo de vida de una suscripcion (crear, listar, eliminar) contra la API real, protege el aislamiento entre usuarios, y documenta de forma honesta la limitacion tecnica del arnes de integracion frente al servicio de notificaciones asincrono, verificando en su lugar el contrato de datos que ese servicio efectivamente consume.
