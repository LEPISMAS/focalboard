# Argumentacion de pruebas de integracion: Gestion de Tableros

## 1. Objetivo del flujo

El flujo INT-02 tiene como objetivo comprobar el ciclo de vida de los tableros en Focalboard desde una perspectiva integrada. La gestion de tableros no depende de una sola funcion: involucra API REST, reglas de negocio, permisos, persistencia, membresias, bloques asociados y, en el diseno del sistema, notificaciones en tiempo real.

El objetivo principal es demostrar que las operaciones de creacion, membresia inicial, listado filtrado, actualizacion, eliminacion logica, duplicacion y notificacion se comportan de forma coherente cuando atraviesan las capas reales del backend.

## 2. Capas integradas

Las pruebas integran las siguientes capas:

- API REST: endpoints como `/api/v2/boards`, `/api/v2/teams/{teamID}/boards`, `/api/v2/boards/{boardID}`, `/api/v2/boards/{boardID}/duplicate` y rutas de bloques.
- Capa App: metodos como `CreateBoard`, `GetBoardsForUserAndTeam`, `PatchBoard`, `DeleteBoard`, `DuplicateBoard` y `GetMembersForBoard`.
- Store: persistencia de tableros, historial, membresias y bloques.
- Modelo: estructuras `Board`, `BoardPatch`, `BoardMember`, `Block`, `BoardsAndBlocks` y `QueryBoardHistoryOptions`.
- Sistema de permisos: validacion de pertenencia, administracion y visibilidad de tableros privados.
- WebSocket: reconocido como parte del flujo arquitectonico de notificacion, aunque la validacion end-to-end queda limitada por la ausencia de un helper estable en `server/integrationtests`.

## 3. Estandares aplicados

Las pruebas se disenaron siguiendo criterios de calidad aplicables a pruebas de integracion:

- Trazabilidad por identificador: cada funcion conserva el codigo INT-02-01 a INT-02-07.
- Separacion de casos: cada caso prueba un riesgo concreto y evita mezclar responsabilidades innecesarias.
- Reutilizacion del arnes oficial: se emplean `SetupTestHelper`, clientes reales y metodos publicos ya usados por el repositorio.
- Pruebas contra API real: las acciones principales se realizan mediante helpers que llaman endpoints REST, no mediante simulaciones de bajo nivel.
- Verificacion de estado persistido: cuando corresponde, se confirma el resultado mediante `App` o `Store`.
- Prueba de rutas positivas y restricciones: el flujo no solo valida que se pueda crear o actualizar, sino tambien que el listado respete membresias.
- Transparencia en limitaciones: INT-02-07 declara la limitacion tecnica sobre WebSocket end-to-end en integracion, sin inventar un cliente fragil.

## 4. Importancia desde risk-based testing

### Identificacion de riesgos

Riesgos de producto:

- Tableros creados sin ID, sin creador o sin persistencia real.
- Ausencia de membresia admin para el creador, dejando tableros inaccesibles o sin propietario.
- Exposicion de tableros privados de otros usuarios.
- Actualizaciones que responden correctamente pero no persisten.
- Eliminacion fisica accidental en lugar de soft delete.
- Duplicacion incompleta de bloques, propiedades o membresias.
- Ausencia de notificacion a clientes conectados cuando cambia un tablero.

Riesgos de proyecto:

- Refactors en API/App/Store que mantengan pruebas unitarias verdes, pero rompan el flujo real.
- Cobertura insuficiente de permisos y membresias.
- Dificultad para defender el comportamiento del sistema si solo se prueban unidades aisladas.
- Pruebas WebSocket inestables si se implementan sin arnes confiable.

### Evaluacion probabilidad por impacto

| Caso | Riesgo principal | Probabilidad | Impacto | Nivel |
|---|---|---:|---:|---|
| INT-02-01 | Tablero creado sin persistencia correcta | Media | Alta | Alto |
| INT-02-02 | Tablero sin admin propietario | Media | Alta | Alto |
| INT-02-03 | Fuga de tableros privados en listados | Media | Critico | Critico |
| INT-02-04 | Actualizacion no persistida | Media | Media | Medio |
| INT-02-05 | Eliminacion incorrecta o visible tras borrado | Media | Alta | Alto |
| INT-02-06 | Duplicacion incompleta de contenido o permisos | Media | Media-Alta | Medio-Alto |
| INT-02-07 | Cambio no notificado en tiempo real | Baja-Media | Media | Medio |

### Priorizacion

El caso INT-02-03 tiene prioridad critica porque afecta confidencialidad: un tablero privado no debe aparecer en listados de usuarios sin membresia. INT-02-01, INT-02-02 e INT-02-05 son de prioridad alta porque sostienen la existencia, propiedad y ciclo de vida de los tableros. INT-02-06 tiene prioridad media-alta porque duplicar tableros es una operacion compuesta y propensa a inconsistencias. INT-02-07 es importante para experiencia en tiempo real, pero su validacion end-to-end requiere un arnes especializado.

### Mitigacion

La mitigacion se realiza ejercitando operaciones reales de API y verificando resultados en App/Store. Esto permite detectar fallos en contratos entre capas: payloads mal interpretados, permisos aplicados de forma incompleta, membresias omitidas, historial de eliminacion incorrecto o duplicaciones parciales.

## 5. Justificacion de herramientas

| Herramienta | Por que se eligio | Alternativas consideradas | Limitacion conocida |
|---|---|---|---|
| Go test | Es el mecanismo nativo de pruebas del backend Go y permite ejecutar solo el flujo con `-run TestINT02`. | Herramientas externas HTTP como Postman/Newman. | Depende de que el entorno de base de datos de pruebas este disponible. |
| Testify require | Ya se usa en el repositorio y permite aserciones claras y expresivas. | `testing` puro o librerias distintas. | Detiene el caso en el primer fallo critico. |
| TestHelper de integrationtests | Provee servidor, clientes autenticados y limpieza consistente de recursos. | Crear setup propio para cada prueba. | Hereda las limitaciones del arnes existente, especialmente en aspectos no cubiertos como WebSocket end-to-end. |
| Cliente HTTP del proyecto | Asegura que las pruebas pasen por endpoints reales y cabeceras correctas. | Llamar directamente a App. | Algunas verificaciones internas requieren complementar con App/Store. |
| App/Store publicos | Permiten comprobar persistencia, membresias, historial y bloques sin SQL directo. | Consultas SQL manuales sobre tablas. | Solo permiten verificar informacion expuesta por interfaces publicas. |
| Logs con `t.Log` | Hacen visible la intencion de cada prueba en salida estandar para debate y trazabilidad. | Comentarios internos solamente. | No sustituyen un reporte de ejecucion formal. |

## 6. Justificacion estrategica por caso

### INT-02-01 Crear tablero persiste en Store

Este caso verifica que la creacion de un tablero no termine solo en una respuesta HTTP correcta. La persistencia se valida con App/Store, lo que permite defender que el flujo API-App-Store funciona como unidad integrada.

### INT-02-02 Crear tablero crea membresia admin

Un tablero sin miembro administrador queda funcionalmente inconsistente: existe, pero puede no tener propietario operativo. La prueba valida que el creador recibe `SchemeAdmin`, que es el equivalente real del modelo para rol administrador del tablero.

### INT-02-03 Listar tableros filtra por membresia

Este caso es estrategico por confidencialidad. La prueba crea tableros privados para dos usuarios y verifica que el usuario 1 no vea el tablero privado del usuario 2. Esto demuestra que el listado no es solo una consulta de equipo, sino una consulta filtrada por permisos y membresias.

### INT-02-04 Actualizar titulo persiste

Actualizar titulo es una operacion comun de usuario. La prueba confirma tanto la respuesta del PATCH como la persistencia posterior, mitigando el riesgo de inconsistencias entre API y Store.

### INT-02-05 Eliminar tablero por soft delete

El sistema debe preservar historial mediante eliminacion logica. La prueba verifica que el tablero ya no aparece como activo y que en el historial existe `deleteAt > 0`. Esto protege contra dos riesgos: borrado fisico accidental y tableros eliminados que siguen visibles.

### INT-02-06 Duplicar tablero copia bloques, propiedades y membresias

Duplicar es una operacion compuesta: crea un nuevo tablero, copia bloques, conserva propiedades relevantes y asigna membresia al creador. Es un caso de integracion de alto valor porque combina varias entidades y reglas de negocio en una sola accion.

### INT-02-07 Crear tablero y notificacion WebSocket

La arquitectura contempla que los cambios se comuniquen por WebSocket. Sin embargo, en `server/integrationtests` no existe un helper estable para suscripcion WebSocket end-to-end. Por eso se implementa la comprobacion mas cercana y segura: creacion por API y verificacion en App/Store, dejando explicita la limitacion tecnica. Esta decision evita introducir una prueba fragil que pueda fallar por infraestructura de test y no por defecto real del producto.

## 7. Relacion con pruebas unitarias previas del store Redux

Las pruebas unitarias previas del store Redux validaban unidades aisladas del frontend: transformaciones de estado, reducers, selectores y acciones en condiciones controladas. Esas pruebas son necesarias porque aseguran que el cliente maneje correctamente su estado local.

Las pruebas INT-02 amplian el alcance: verifican que la informacion que el frontend eventualmente consumiria se produce correctamente desde el backend integrado. La diferencia central es que las unitarias responden si una pieza aislada funciona, mientras que estas pruebas responden si las capas API, App, Store y permisos colaboran correctamente en escenarios reales de gestion de tableros.

La relacion es complementaria. Las unitarias del store Redux reducen el riesgo de errores locales en el cliente; las pruebas de integracion reducen el riesgo de errores de contrato y comunicacion entre capas del servidor.

## 8. Conclusion argumentativa

El flujo INT-02 fue seleccionado porque la gestion de tableros concentra reglas esenciales del producto: propiedad, visibilidad, permisos, persistencia, actualizacion, eliminacion y duplicacion. Un defecto en este flujo puede afectar tanto la experiencia del usuario como la confidencialidad y consistencia de los datos.

La estrategia de prueba es defendible porque prioriza escenarios segun riesgo, usa herramientas existentes del proyecto y evita verificaciones artificiales. En lugar de probar funciones aisladas, las pruebas recorren endpoints reales y confirman efectos persistidos. Esto aporta evidencia mas fuerte para sostener que Focalboard mantiene integridad y control de acceso en operaciones centrales de tableros.
