### Explicacion de pruebas de integracion: Importacion/Exportacion y Suscripciones/Notificaciones

Se realizaron las pruebas de integracion para dos flujos del proyecto Focalboard: Importacion y Exportacion de tableros (INT-07) y Suscripciones y Notificaciones (INT-08).

**Flujo INT-07 (Importacion y Exportacion)**

Se probo el ciclo completo API -> App -> FileSystem -> Store: exportar un tablero a un archivo `.boardarchive`, importarlo de vuelta, hacer un roundtrip completo verificando que no se pierde informacion, rechazar un archivo con formato invalido sin dejar datos parciales, y subir un archivo adjunto a una tarjeta verificando que quede almacenado y recuperable.

Las pruebas quedaron en `server/integrationtests/flujo_importacionYExportacion_int_tests.go`, con scripts de ejecucion en `server/integrationtests/ejecutar_importacionYExportacion.sh` y `.bat`.

**Flujo INT-08 (Suscripciones y Notificaciones)**

Se probo el ciclo API -> App -> Store para suscripciones: crear una suscripcion a un bloque, listar las suscripciones de un usuario (verificando que no se mezclan con las de otro usuario), y eliminar una suscripcion. Tambien se valido el contrato de datos `NotificationHint`, que es la pieza que el servicio de notificaciones usa para decidir a quien avisar cuando cambia un bloque.

Una limitacion importante que documentamos: el arnes de pruebas de integracion del proyecto (`server/integrationtests`) no tiene registrado un backend real de notificaciones ni de websocket, asi que no se puede probar de punta a punta que modificar un bloque dispare automaticamente una notificacion. Esto ya lo habia identificado el equipo antes para websocket en el flujo de Gestion de Tableros, asi que seguimos el mismo criterio: documentar la limitacion en vez de simular algo que el entorno no soporta.

Las pruebas quedaron en `server/integrationtests/flujo_suscripcionesYNotificaciones_int_tests.go`, con scripts de ejecucion en `server/integrationtests/ejecutar_suscripcionesYNotificaciones.sh` y `.bat`.

**Documentacion**

Para cada flujo se genero un documento argumentativo (estandares aplicados, risk-based testing, justificacion de herramientas y de estrategia por caso) y un reporte de ejecucion, ambos en `lepismas/docs/reportes/`.