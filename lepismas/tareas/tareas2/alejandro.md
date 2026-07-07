* nombre: Documentacion argumentativa (pruebas de componente) y pruebas de integracion.

* issue: ALEJANDRO Documentacion argumentativa (pruebas de integracion)

* leer: 
Las proximas exposiciones del trabajo final van a ser debates, por lo que debemos argumentar todas las decisiones que hemos realizado desde el principio y las que realizaremos a partir de ahora.
La tarea tiene 2 partes:
- realizar pruebas de integracion para 2 flujos especificos.
- documentar de forma argumentativa las pruebas de integracion de esta tarea y las pruebas unitarias de la prueba anterior.

* prompts:
<tarea1>
Vas a realizar las pruebas de integracion para los siguientes flujos y casos de prueba
<flujo1>
### INT-01: Autenticación — API REST ↔ Servicio Auth ↔ Store

**Objetivo:** Verificar que el flujo completo de autenticación funcione correctamente a través de las capas API, Auth Service y Store de sesiones.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                | Precondiciones                        | Resultado Esperado                                                    | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|-----------------------------------------|---------------------------------------|-----------------------------------------------------------------------|-----------|
| INT-01-01  | Registrar un usuario nuevo vía API y verificar que se persista en el Store                                 | API (`/api/v2/register`) → App → Store  | Servidor levantado, DB vacía          | Usuario creado en BD, respuesta HTTP 200 con datos del usuario        | Alta      |
| INT-01-02  | Iniciar sesión con el usuario registrado y verificar que se genere un token de sesión válido                | API (`/api/v2/login`) → Auth → Store    | Usuario `user1` registrado            | Token de sesión válido, sesión creada en BD, cookie configurada       | Alta      |
| INT-01-03  | Acceder a un endpoint protegido con token válido y verificar acceso concedido                              | API → Auth (validación de token)        | Sesión activa de `user1`              | Respuesta HTTP 200, datos retornados correctamente                    | Alta      |
| INT-01-04  | Acceder a un endpoint protegido sin token y verificar rechazo                                              | API → Auth (validación de token)        | Sin sesión                            | Respuesta HTTP 401 Unauthorized                                       | Alta      |
| INT-01-05  | Cambiar contraseña vía API y verificar que el nuevo password funcione en el siguiente login                 | API (`/api/v2/users/.../changepassword`) → Auth → Store | Sesión activa de `user1` | Contraseña actualizada en BD, login exitoso con nueva contraseña      | Media     |
| INT-01-06  | Cerrar sesión vía API y verificar que el token quede invalidado para peticiones subsiguientes               | API → Auth → Store (DeleteSession)      | Sesión activa de `user1`              | Sesión eliminada de BD, siguiente petición con ese token retorna 401  | Alta      |
| INT-01-07  | Registrar usuario con datos duplicados y verificar que Store rechace la operación y API retorne error      | API → App → Store (constraint unique)   | Usuario `user1` ya existe             | Respuesta HTTP 200 con error o HTTP 409, usuario no duplicado en BD   | Media     |
</flujo1>
<flujo2>
### INT-02: Gestión de Tableros — API REST ↔ App ↔ Store ↔ WebSocket

**Objetivo:** Verificar el ciclo de vida completo de los tableros a través de las capas del sistema, incluyendo notificación en tiempo real.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                          | Precondiciones                         | Resultado Esperado                                                          | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------------|----------------------------------------|-----------------------------------------------------------------------------|-----------|
| INT-02-01  | Crear un tablero vía API y verificar que se persista en la BD con todos sus campos                         | API → App (`CreateBoard`) → Store → DB            | `user1` autenticado                    | Tablero con ID generado, `createdBy` = user1, tipo correcto, persistido     | Alta      |
| INT-02-02  | Crear un tablero y verificar que se cree automáticamente la membresía del propietario                      | API → App → Store (Board + BoardMember)           | `user1` autenticado                    | BoardMember con rol `admin` creado para `user1`                             | Alta      |
| INT-02-03  | Obtener tableros del equipo y verificar que solo retorne tableros donde el usuario es miembro               | API → App → Store (GetBoardsForUserAndTeam)       | `user1` con 2 tableros, `user2` con 1 | `user1` obtiene solo sus tableros, no los privados de `user2`               | Alta      |
| INT-02-04  | Actualizar el título de un tablero vía PATCH y verificar persistencia y respuesta                          | API → App (`PatchBoard`) → Store                  | Tablero existente de `user1`           | Título actualizado en BD, respuesta con tablero actualizado                 | Media     |
| INT-02-05  | Eliminar un tablero vía API y verificar eliminación lógica (soft delete) en la BD                          | API → App (`DeleteBoard`) → Store                 | Tablero existente de `user1`           | Campo `deleteAt` > 0 en BD, tablero no aparece en listados                 | Alta      |
| INT-02-06  | Duplicar un tablero y verificar que se copien bloques, propiedades y membresías                            | API → App (`DuplicateBoard`) → Store              | Tablero con tarjetas y propiedades     | Nuevo tablero con copia de bloques, IDs distintos, membresía del creador    | Media     |
| INT-02-07  | Crear un tablero y verificar que el cambio se notifique vía WebSocket a clientes conectados                | API → App → Store + App → WS (Broadcast)          | `user2` suscrito al equipo vía WS      | `user2` recibe evento de creación de tablero por WebSocket                  | Media     |

</flujo2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Por cada flujo realiza un archivo "flujo_[NOMBRE DE FLUJO CAMELCASE]_int_tests.go" con las pruebas de integracion (en este caso: flujo_autenticacion_int_tests.go, flujo_gestionDeTableros_int_tests.go)
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
[lepisma/alejandro/explicacion2.md] (explicacion de lo realizado - esto va en el guion cuando expongamos - debes grabarte lo que escribas)
[server/integrationtests/flujo_autenticacion_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 1)
[server/integrationtests/flujo_gestionDeTableros_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 2)
[server/integrationtests/ejecutar_autenticacion.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 bash)
[server/integrationtests/ejecutar_autenticacion.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 windows)
[server/integrationtests/ejecutar_gestionDeTableros.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 bash)
[server/integrationtests/ejecutar_gestionDeTableros.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 windows)
[lepismas/docs/reportes/argumentacion_integration_tests_autenticacion.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 1)
[lepismas/docs/reportes/argumentacion_integration_tests_gestionDeTableros.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 2)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_autenticacion.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 1)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_gestionDeTableros.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 2)