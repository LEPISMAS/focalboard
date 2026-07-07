* nombre: Documentacion argumentativa (pruebas de componente) y pruebas de integracion.

* issue: DANIEL Documentacion argumentativa (pruebas de integracion)

* leer: 
Las proximas exposiciones del trabajo final van a ser debates, por lo que debemos argumentar todas las decisiones que hemos realizado desde el principio y las que realizaremos a partir de ahora.
La tarea tiene 2 partes:
- realizar pruebas de integracion para 2 flujos especificos.
- documentar de forma argumentativa las pruebas de integracion de esta tarea y las pruebas unitarias de la prueba anterior.

* prompts:
<tarea1>
Vas a realizar las pruebas de integracion para los siguientes flujos y casos de prueba
<flujo1>
### INT-03: Gestión de Tarjetas y Bloques — API ↔ App ↔ Store

**Objetivo:** Verificar la integración en la gestión de tarjetas (cards) y bloques de contenido, que son las entidades centrales de Focalboard.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                      | Precondiciones                              | Resultado Esperado                                                        | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|-----------------------------------------------|---------------------------------------------|---------------------------------------------------------------------------|-----------|
| INT-03-01  | Crear una tarjeta dentro de un tablero vía API y verificar que se persista como Block de tipo `card`       | API → App (`CreateCard`) → Store              | Tablero existente, `user1` autenticado      | Block con `type=card`, `boardId` correcto, campos generados               | Alta      |
| INT-03-02  | Obtener tarjetas de un tablero y verificar que incluyan propiedades personalizadas                         | API → App → Store (GetCards + propiedades)    | Tablero con tarjetas y propiedades          | Array de tarjetas con campo `properties` poblado correctamente            | Alta      |
| INT-03-03  | Actualizar propiedades de una tarjeta vía PATCH y verificar que el Store refleje los cambios               | API → App (`PatchCard`) → Store               | Tarjeta existente con propiedades           | Propiedades actualizadas en BD, respuesta refleja cambios                 | Media     |
| INT-03-04  | Insertar un bloque de contenido (texto, imagen, checkbox) dentro de una tarjeta y verificar relación padre | API → App (`InsertBlock`) → Store             | Tarjeta existente                           | Bloque con `parentId` = ID de la tarjeta, tipo correcto                   | Alta      |
| INT-03-05  | Eliminar una tarjeta y verificar que sus bloques hijos también se marquen como eliminados                  | API → App (`DeleteBlock`) → Store             | Tarjeta con bloques hijos                   | Tarjeta y bloques hijos con `deleteAt` > 0                               | Alta      |
| INT-03-06  | Restaurar una tarjeta eliminada y verificar que se recupere junto con sus bloques hijos                    | API → App (`UndeleteBlock`) → Store           | Tarjeta previamente eliminada               | Tarjeta y bloques hijos con `deleteAt` = 0                               | Media     |
| INT-03-07  | Crear un tablero con bloques en lote (`CreateBoardsAndBlocks`) y verificar atomicidad de la operación      | API → App → Store (transacción)               | `user1` autenticado                         | Tablero y bloques creados atómicamente o ninguno (rollback si falla)      | Alta      |
</flujo1>
<flujo2>
### INT-04: Permisos — API ↔ PermissionsService ↔ Store

**Objetivo:** Verificar que el sistema de permisos funcione correctamente integrando la capa API con el servicio de permisos y la información de membresías almacenada en el Store.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                            | Precondiciones                                    | Resultado Esperado                                                  | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|---------------------------------------------------|---------------------------------------------------------------------|-----------|
| INT-04-01  | Un usuario `viewer` intenta editar un tablero y verificar que el sistema rechace la operación              | API → Permissions → Store (BoardMember role)        | `user2` agregado como `viewer` al tablero         | HTTP 403 Forbidden, tablero sin cambios en BD                       | Alta      |
| INT-04-02  | Un usuario `editor` edita un tablero y verificar que el sistema permita la operación                       | API → Permissions → Store                           | `user2` agregado como `editor` al tablero         | HTTP 200, cambios persistidos correctamente                         | Alta      |
| INT-04-03  | Un usuario `admin` del tablero agrega un nuevo miembro y verificar la persistencia de la membresía         | API → App → Permissions → Store (AddMember)         | `user1` es admin del tablero                      | Nuevo BoardMember creado en BD con rol asignado                     | Alta      |
| INT-04-04  | Un usuario no miembro intenta acceder a un tablero privado y verificar rechazo                             | API → Permissions → Store                           | Tablero privado de `user1`, `user2` sin membresía | HTTP 403 Forbidden                                                  | Alta      |
| INT-04-05  | Cambiar el rol de un miembro de `viewer` a `editor` y verificar que los permisos se actualicen             | API → App → Store (UpdateBoardMember)               | `user2` es `viewer`                               | Rol actualizado en BD, `user2` ahora puede editar                   | Media     |
| INT-04-06  | Eliminar membresía de un usuario y verificar que pierde acceso inmediatamente                              | API → App → Store (DeleteMember) + Permissions      | `user2` es miembro del tablero                    | Membresía eliminada, siguiente GET al tablero retorna 403           | Media     |
| INT-04-07  | Verificar que un tablero público (tipo `Open`) sea accesible para cualquier miembro del equipo             | API → Permissions → Store                           | Tablero tipo `Open` en el equipo                  | Cualquier usuario del equipo puede ver el tablero (HTTP 200)        | Media     |
</flujo2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Por cada flujo realiza un archivo "flujo_[NOMBRE DE FLUJO CAMELCASE]_int_tests.go" con las pruebas de integracion (en este caso: flujo_gestionDeTarjetasYBloques_int_tests.go, flujo_permisos_int_tests.go)
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
[lepisma/daniel/explicacion2.md] (explicacion de lo realizado - esto va en el guion cuando expongamos - debes grabarte lo que escribas)
[server/integrationtests/flujo_gestionDeTarjetasYBloques_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 1)
[server/integrationtests/flujo_permisos_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 2)
[server/integrationtests/ejecutar_gestionDeTarjetasYBloques.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 bash)
[server/integrationtests/ejecutar_gestionDeTarjetasYBloques.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 windows)
[server/integrationtests/ejecutar_permisos.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 bash)
[server/integrationtests/ejecutar_permisos.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 windows)
[lepismas/docs/reportes/argumentacion_integration_tests_gestionDeTarjetasYBloques.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 1)
[lepismas/docs/reportes/argumentacion_integration_tests_permisos.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 2)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_gestionDeTarjetasYBloques.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 1)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_permisos.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 2)