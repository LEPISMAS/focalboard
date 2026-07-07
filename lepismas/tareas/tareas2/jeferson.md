* nombre: Documentacion argumentativa (pruebas de componente) y pruebas de integracion.

* issue: JEFERSON Documentacion argumentativa (pruebas de componente)

* leer: 
Las proximas exposiciones del trabajo final van a ser debates, por lo que debemos argumentar todas las decisiones que hemos realizado desde el principio y las que realizaremos a partir de ahora.
La tarea tiene 2 partes:
- realizar pruebas de integracion para 2 flujos especificos.
- documentar de forma argumentativa las pruebas de integracion de esta tarea y las pruebas unitarias de la prueba anterior.

* prompts:
<tarea1>
Vas a realizar las pruebas de integracion para los siguientes flujos y casos de prueba
<flujo1>
### 5.7. INT-07: Importación y Exportación — API ↔ App ↔ FileSystem ↔ Store

**Objetivo:** Verificar que los flujos de importación y exportación de tableros funcionen correctamente integrando la API, la lógica de negocio, el sistema de archivos y el Store.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                       | Precondiciones                          | Resultado Esperado                                                      | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|------------------------------------------------|-----------------------------------------|-------------------------------------------------------------------------|-----------|
| INT-07-01  | Exportar un tablero vía API y verificar que se genere un archivo con el formato correcto                   | API → App (`ExportBoardArchive`) → Store → FS  | Tablero con tarjetas y bloques          | Archivo generado, contiene JSON con tablero, bloques y metadatos        | Alta      |
| INT-07-02  | Importar un archivo de respaldo previamente exportado y verificar que se recreen las entidades en el Store | API → App (`ImportBoardArchive`) → Store       | Archivo de exportación válido           | Tablero y bloques recreados en BD con nuevos IDs                        | Alta      |
| INT-07-03  | Exportar e importar un tablero completo (roundtrip) y verificar integridad de datos                        | API → App → Store (Export + Import)            | Tablero con propiedades y tarjetas      | Tablero importado tiene mismo contenido que el original                  | Alta      |
| INT-07-04  | Intentar importar un archivo con formato inválido y verificar que el sistema lo rechace sin corromper datos | API → App → Store                              | Archivo con formato corrupto            | Error descriptivo, ninguna entidad parcialmente creada en BD            | Media     |
| INT-07-05  | Subir un archivo adjunto a una tarjeta y verificar que se almacene y sea recuperable                       | API (`/files`) → App → FileSystem + Store      | Tarjeta existente                       | Archivo almacenado en `filespath`, referencia en BD, descargable vía API | Media     |
</flujo1>
<flujo2>
### 5.8. INT-08: Suscripciones y Notificaciones — API ↔ NotifyService ↔ Store ↔ WebSocket

**Objetivo:** Verificar la integración del sistema de suscripciones con el servicio de notificaciones y la persistencia.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                        | Precondiciones                               | Resultado Esperado                                                | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|------------------------------------------------|----------------------------------------------|-------------------------------------------------------------------|-----------|
| INT-08-01  | Crear una suscripción a un tablero vía API y verificar persistencia en el Store                            | API → App (`CreateSubscription`) → Store        | Tablero existente, `user1` autenticado       | Suscripción creada en BD con subscriberType y subscriberID        | Alta      |
| INT-08-02  | Obtener suscripciones del usuario y verificar coherencia con las creadas                                   | API → App → Store (GetSubscriptions)            | Suscripciones previamente creadas            | Array con suscripciones activas del usuario                       | Media     |
| INT-08-03  | Eliminar una suscripción y verificar que se elimine del Store                                              | API → App (`DeleteSubscription`) → Store        | Suscripción activa                           | Suscripción eliminada de BD                                       | Media     |
| INT-08-04  | Modificar un bloque en un tablero suscrito y verificar que se genere un notification hint                  | App (`InsertBlock`) → Store → NotifyService     | `user2` suscrito al tablero                  | NotificationHint creado en BD tras la modificación                | Media     |
</flujo2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Por cada flujo realiza un archivo "flujo_[NOMBRE DE FLUJO CAMELCASE]_int_tests.go" con las pruebas de integracion (en este caso: flujo_importacionYExportacion_int_tests.go, flujo_suscripcionesYNotificaciones_int_tests.go)
- Por cada caso de prueba va haber una funcion que realiza la prueba.
- La ejecucion de la prueba deja en clara su utilidad a travez de la salida estandar.
- Vas a dejar estos archivos en la ubicacion server/integrationtests/
- Crea un archivo bash y bat para ejecutar las pruebas para cada flujo en la ubicacion server/integrationtests/ejecutar_[NOMBRE FLUJO CAMELCASE].[EXTENCION]
- Realiza un documento de caracter argumentativo para cada flujo en la ubicacion lepismas/docs/reportes/ en el cual se mencione:
-- estandares aplicados en la realizacion de las pruebas
-- importancia del test a raiz del risk based testing:
--- Identificación de riesgos (producto y proyecto)
--- Evaluación: probabilidad de ocurrencia × impacto/severidad
--- Priorización de casos de prueba según el nivel de riesgo resultante
--- Mitigación mediante mayor cobertura de prueba en las áreas de alto riesgo
-- justificacion de uso de herramientas (por cada herramienta realizar un cuadro con estos campos: Herramienta, Por qué se eligió, Alternativas consideradas, Limitación conocida)
-- para cada caso de prueba: (justificacion de estrategia)
-- el nombre del archivo es argumentacion_integration_tests_[NOMBRE FLUJO CAMELCASE].md
- Ejecuta las pruebas y realiza un reporte de ejecucion (en la direccion lepismas/docs/reportes/[NOMBRE DE FLUJO CAMELCASE]_reporte_ejecucion_pruebas_integracion.md) que contenga:
-- Reporte de pruebas exitosas o fallidas
-- Si es que hay pruebas fallidas incluir matriz de trazabilidad de defectos: ID Defecto | Caso de prueba que lo detectó | Componente afectado | Severidad | Requisito relacionado | Estado (abierto/cerrado)
</instrucciones>
Toma en cuenta el contexto de la planeacion de pruebas de integracion
</tarea1>
considera:
Los componentes del proyecto estan establecidos en el documento docs/contexto/modulos_componentes.md

* necesario:
[lepisma/gabriela/explicacion2.md] (explicacion de lo realizado - esto va en el guion cuando expongamos - debes grabarte lo que escribas)
[server/integrationtests/flujo_importacionYExportacion_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 1)
[server/integrationtests/flujo_suscripcionesYNotificaciones_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 2)
[server/integrationtests/ejecutar_importacionYExportacion.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 bash)
[server/integrationtests/ejecutar_importacionYExportacion.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 windows)
[server/integrationtests/ejecutar_suscripcionesYNotificaciones.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 bash)
[server/integrationtests/ejecutar_suscripcionesYNotificaciones.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 windows)
[lepismas/docs/reportes/argumentacion_integration_tests_importacionYExportacion.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 1)
[lepismas/docs/reportes/argumentacion_integration_tests_suscripcionesYNotificaciones.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 2)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_importacionYExportacion.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 1)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_suscripcionesYNotificaciones.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 2)