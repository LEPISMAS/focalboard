* nombre: Documentacion argumentativa (pruebas de componente) y pruebas de integracion.

* issue: GABRIELA Documentacion argumentativa (pruebas de componente)

* leer: 
Las proximas exposiciones del trabajo final van a ser debates, por lo que debemos argumentar todas las decisiones que hemos realizado desde el principio y las que realizaremos a partir de ahora.
La tarea tiene 2 partes:
- realizar pruebas de integracion para 2 flujos especificos.
- documentar de forma argumentativa las pruebas de integracion de esta tarea y las pruebas unitarias de la prueba anterior.

* prompts:
<tarea1>
Vas a realizar las pruebas de integracion para los siguientes flujos y casos de prueba
<flujo1>
### 5.5. INT-05: Compartición de Tableros — API ↔ Sharing ↔ Store

**Objetivo:** Verificar la integración del sistema de compartición pública de tableros.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados               | Precondiciones                      | Resultado Esperado                                                    | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|----------------------------------------|-------------------------------------|-----------------------------------------------------------------------|-----------|
| INT-05-01  | Habilitar compartición pública de un tablero y verificar que se genere un token de compartición en el Store | API → App (`UpsertSharing`) → Store    | Tablero de `user1`, sin compartir   | Registro `Sharing` creado en BD con token único, `enabled=true`       | Alta      |
| INT-05-02  | Acceder a un tablero compartido con token válido sin autenticación y verificar respuesta                   | API (ruta pública) → Store (Sharing)   | Tablero con sharing habilitado      | HTTP 200 con datos del tablero                                        | Alta      |
| INT-05-03  | Acceder a un tablero compartido con token inválido y verificar rechazo                                     | API (ruta pública) → Store (Sharing)   | Token de sharing incorrecto         | HTTP 401 o HTTP 404                                                   | Alta      |
| INT-05-04  | Deshabilitar compartición y verificar que el acceso público se revoque                                     | API → App → Store (UpsertSharing)      | Tablero previamente compartido      | Sharing `enabled=false` en BD, acceso público retorna error           | Media     |
</flujo1>
<flujo2>
### 5.6. INT-06: Categorías y Barra Lateral — API ↔ App ↔ Store

**Objetivo:** Verificar la integración en la gestión de categorías de la barra lateral, que permiten organizar los tableros por agrupaciones personalizadas.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                    | Precondiciones                           | Resultado Esperado                                                  | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------|------------------------------------------|---------------------------------------------------------------------|-----------|
| INT-06-01  | Crear una categoría vía API y verificar persistencia en Store                                              | API → App (`CreateCategory`) → Store        | `user1` autenticado                      | Categoría creada en BD con nombre, teamID y userID correctos        | Alta      |
| INT-06-02  | Mover un tablero a una categoría personalizada y verificar la asociación en el Store                       | API → App → Store (CategoryBoards)          | Categoría y tablero existentes           | Relación CategoryBoard creada, sidebar muestra tablero en categoría | Media     |
| INT-06-03  | Obtener categorías del sidebar y verificar que incluyan los tableros asociados                             | API → App → Store (GetUserCategoryBoards)   | Categorías con tableros asignados        | Array de categorías con campo `boardIDs` poblado                    | Alta      |
| INT-06-04  | Eliminar una categoría y verificar que los tableros se muevan a la categoría por defecto                   | API → App → Store                           | Categoría personalizada con tableros     | Categoría eliminada, tableros reasignados a categoría default       | Media     |
| INT-06-05  | Reordenar categorías y verificar que el nuevo orden se persista                                            | API → App → Store                           | Múltiples categorías creadas             | Orden actualizado en BD, siguiente GET retorna en nuevo orden       | Baja      |
</flujo2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Por cada flujo realiza un archivo "flujo_[NOMBRE DE FLUJO CAMELCASE]_int_tests.go" con las pruebas de integracion (en este caso: flujo_comparticionDeTableros_int_tests.go, flujo_CategoriasYBarraLateral_int_tests.go)
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
[server/integrationtests/flujo_comparticionDeTableros_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 1)
[server/integrationtests/flujo_categoriasYBarraLateral_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 2)
[server/integrationtests/ejecutar_comparticionDeTableros.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 bash)
[server/integrationtests/ejecutar_comparticionDeTableros.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 windows)
[server/integrationtests/ejecutar_categoriasYBarraLateral.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 bash)
[server/integrationtests/ejecutar_categoriasYBarraLateral.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 windows)
[lepismas/docs/reportes/argumentacion_integration_tests_comparticionDeTableros.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 1)
[lepismas/docs/reportes/argumentacion_integration_tests_categoriasYBarraLateral.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 2)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_comparticionDeTableros.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 1)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_categoriasYBarraLateral.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 2)