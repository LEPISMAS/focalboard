* nombre: Documentacion argumentativa (pruebas de componente) y pruebas de integracion.

* issue: MARCO Documentacion argumentativa (pruebas de componente)

* leer: 
Las proximas exposiciones del trabajo final van a ser debates, por lo que debemos argumentar todas las decisiones que hemos realizado desde el principio y las que realizaremos a partir de ahora.
La tarea tiene 2 partes:
- realizar pruebas de integracion para 2 flujos especificos.
- documentar de forma argumentativa las pruebas de integracion de esta tarea y las pruebas unitarias de la prueba anterior.

* prompts:
<tarea1>
Vas a realizar las pruebas de integracion para los siguientes flujos y casos de prueba
<flujo1>
### 5.9. INT-09: Búsqueda — API ↔ App ↔ Store (Full-Text)

**Objetivo:** Verificar que la funcionalidad de búsqueda integre correctamente la API con la lógica de consulta en el Store.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados             | Precondiciones                            | Resultado Esperado                                                   | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|--------------------------------------|-------------------------------------------|----------------------------------------------------------------------|-----------|
| INT-09-01  | Buscar tableros por título y verificar que retorne solo los tableros accesibles al usuario                 | API → App (`SearchBoardsForUser`) → Store | Tableros con títulos diversos          | Solo tableros donde `user1` es miembro o son públicos                | Alta      |
| INT-09-02  | Buscar con un término que no coincide con ningún tablero y verificar respuesta vacía                       | API → App → Store                    | Tableros existentes                       | Array vacío, HTTP 200                                                | Media     |
| INT-09-03  | Buscar tableros y verificar que los resultados respeten los permisos (no mostrar tableros privados ajenos) | API → App → Store + Permissions      | Tableros privados de otro usuario         | Resultados no incluyen tableros privados donde no se es miembro      | Alta      |
| INT-09-04  | Buscar con caracteres especiales y verificar que no cause errores en la capa de Store                      | API → App → Store (SQL query)        | Servidor funcionando                      | Respuesta válida (vacía o con resultados), sin error SQL             | Media     |
</flujo1>
<flujo2>
### 5.10. INT-10: Frontend ↔ Backend — OctoClient ↔ API REST

**Objetivo:** Verificar la integración end-to-end entre el frontend React y el backend Go a través de las pruebas Cypress E2E, validando que los flujos de usuario completen su ciclo a través de todas las capas.

| ID         | Caso de Prueba                                                                                             | Componentes Involucrados                          | Precondiciones                       | Resultado Esperado                                                     | Prioridad |
|------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------------|--------------------------------------|------------------------------------------------------------------------|-----------|
| INT-10-01  | Registrar un usuario desde la interfaz web y verificar que pueda iniciar sesión                            | UI (RegisterPage) → OctoClient → API → Store      | Servidor Focalboard corriendo        | Registro exitoso, login exitoso, espacio de trabajo visible            | Alta      |
| INT-10-02  | Crear un tablero desde la interfaz y verificar que aparezca en la barra lateral                            | UI (AddBoard) → Mutator → OctoClient → API → Store → WS | Usuario autenticado           | Tablero creado, aparece en sidebar sin recargar página                 | Alta      |
| INT-10-03  | Crear una tarjeta en la vista Kanban y verificar que se persista y sea editable                            | UI (KanbanCard) → Mutator → OctoClient → API → Store | Tablero existente con vista Board | Tarjeta creada, editable al hacer clic, datos persistidos              | Alta      |
| INT-10-04  | Cambiar de vista (Board → Table) y verificar que se carguen las mismas tarjetas en diferente layout        | UI (ViewMenu) → Mutator → OctoClient → API → Store | Tablero con tarjetas existentes   | Mismas tarjetas visibles en formato tabla                              | Media     |
| INT-10-05  | Aplicar un filtro en una vista y verificar que solo se muestren las tarjetas que cumplen la condición       | UI (FilterComponent) → Redux → OctoClient → API   | Tarjetas con propiedades diversas    | Solo tarjetas filtradas son visibles, sin errores de consola           | Media     |
| INT-10-06  | Cerrar sesión desde el menú de usuario y verificar redirección a login e invalidación de acceso            | UI (SidebarUserMenu) → OctoClient → API → Auth    | Sesión activa                        | Sesión cerrada, redirección a `/login`, acceso posterior denegado      | Alta      |
</flujo2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Por cada flujo realiza un archivo "flujo_[NOMBRE DE FLUJO CAMELCASE]_int_tests.go" con las pruebas de integracion (en este caso: flujo_busqueda_int_tests.go, flujo_frontend_int_tests.go)
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
[server/integrationtests/flujo_busqueda_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 1)
[server/integrationtests/flujo_frontend_int_tests.go] (UBICACION CODIGO PRUEBAS INT flujo 2)
[server/integrationtests/ejecutar_busqueda.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 bash)
[server/integrationtests/ejecutar_busqueda.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 1 windows)
[server/integrationtests/ejecutar_frontend.sh] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 bash)
[server/integrationtests/ejecutar_frontend.bat] (UBICACION SCRIPT EJECUTAR PRUEBAS flujo 2 windows)
[lepismas/docs/reportes/argumentacion_integration_tests_busqueda.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 1)
[lepismas/docs/reportes/argumentacion_integration_tests_frontend.md] (UBICACION DOCUMENTO ARGUMENTATIVO INT flujo 2)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_busqueda.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 1)
[lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_frontend.md] (UBICACION DOCUMENTO REPORTE PRUEBAS INT flujo 2)