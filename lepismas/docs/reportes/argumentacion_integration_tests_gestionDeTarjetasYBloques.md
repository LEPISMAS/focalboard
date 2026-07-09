# Argumentación de pruebas de integración: Gestión de Tarjetas y Bloques (INT-03)

## 1. Objetivo del flujo

El flujo INT-03 tiene como objetivo verificar la correcta integración en el ciclo de vida de las tarjetas (cards) y los bloques de contenido (blocks), que constituyen las entidades centrales del modelo de datos de Focalboard. Las tarjetas representan los registros operacionales, y los bloques hijos (texto, imágenes, checkboxes) estructuran su contenido interno.

El objetivo primordial es demostrar que las operaciones de creación de tarjetas, recuperación de propiedades personalizadas, actualización vía PATCH, e inserción de bloques hijos se coordinan correctamente a través de la API REST, la capa App (lógica de negocio) y el Store (base de datos), garantizando consistencia, integridad referencial y atomicidad.

## 2. Capas integradas

Las pruebas de este flujo integran y ejercitan las siguientes capas del sistema:
- **API REST (server/api):** Endpoints como `/api/v2/boards/{boardID}/cards`, `/api/v2/cards/{cardID}`, `/api/v2/boards/{boardID}/blocks` y `/api/v2/boards-and-blocks`.
- **Capa App (server/app):** Métodos de orquestación como `CreateCard`, `GetCards`, `PatchCard`, `InsertBlocks`, `DeleteBlock`, `UndeleteBlock` y `CreateBoardsAndBlocks`.
- **Store (server/store / sqlstore):** Persistencia física en la base de datos (SQLite/PostgreSQL/MySQL), mapeo de bloques y gestión del historial de versiones.
- **Modelos de datos (server/model):** Estructuras como `Card`, `CardPatch`, `Block` y `BoardsAndBlocks`, además de sus funciones de conversión interna (`Card2Block`, `Block2Card`).

## 3. Estándares aplicados

Para garantizar el rigor metodológico, las pruebas se diseñaron bajo los siguientes estándares de ingeniería de software y control de calidad:
- **Trazabilidad estricta (IEEE 829 / ISO 29119):** Cada caso de prueba está identificado con la nomenclatura del plan de pruebas (`INT-03-01` a `INT-03-07`).
- **Verificación multi-capa (Integración Real):** Las acciones se inician a través del cliente HTTP oficial contra la API REST, pero los resultados de persistencia se validan directamente en el Store o en la capa App (comprobando el estado de la base de datos).
- **Limpieza y aislamiento de datos (Data Sandbox):** Cada ejecución utiliza una base de datos temporal SQLite en memoria que se inicializa y destruye con el ciclo de vida de la prueba (`th.TearDown()`).
- **Validación de efectos colaterales (Cascading deletion):** Se verifican explícitamente los efectos de cascada lógicos (`deleteAt > 0`) en lugar de asumir que la eliminación de un nodo padre es exitosa de forma aislada.

## 4. Importancia desde Risk-Based Testing

### Identificación de Riesgos

**Riesgos de Producto:**
- **RP-03-01:** Creación de tarjetas sin la conversión correspondiente a bloques `card` en la base de datos, corrompiendo la integridad del tablero.
- **RP-03-02:** Pérdida o desalineación de propiedades personalizadas de las tarjetas al guardarse en el mapa genérico de campos de bloques.
- **RP-03-03:** Actualizaciones parciales de propiedades (PATCH) que sobrescriben o eliminan el resto de propiedades preexistentes.
- **RP-03-04:** Pérdida de la relación de parentesco (`parentId`) de bloques de contenido, dejando elementos huérfanos e invisibles en la UI.
- **RP-03-05:** Eliminación lógica de una tarjeta que mantiene activos sus bloques de contenido internos, causando fuga de almacenamiento y basura de datos.
- **RP-03-06:** Restauración incompleta de una tarjeta que deja sus bloques de contenido marcados como eliminados permanentemente.
- **RP-03-07:** Operación parcial fallida en creación de lotes que deja la base de datos en estado inconsistente (infracción de la regla ACID).

**Riesgos de Proyecto:**
- **RProj-03-01:** Refactorizaciones del modelo de bloques que rompan de forma silenciosa el comportamiento de las tarjetas en el cliente web.
- **RProj-03-02:** Modificaciones en el controlador de base de datos que anulen las transacciones en lote y no disparen alertas a nivel unitario.

### Evaluación: Probabilidad × Impacto/Severidad

| ID Caso | Probabilidad | Impacto | Nivel de RiesgoResultante | Prioridad |
|---------|:---:|:---:|:---:|:---:|
| **INT-03-01** (Crear Tarjeta) | Media | Alta | **Alto** | Alta |
| **INT-03-02** (Obtener Propiedades) | Media | Alta | **Alto** | Alta |
| **INT-03-03** (PATCH Propiedades) | Media | Media | **Medio** | Media |
| **INT-03-04** (Insertar Bloque e Hijo) | Alta | Alta | **Alto** | Alta |
| **INT-03-05** (Borrado en Cascada) | Media | Alta | **Alto** | Alta |
| **INT-03-06** (Restaurar Cascada) | Baja | Alta | **Medio-Alto** | Media |
| **INT-03-07** (Atomicidad Lote) | Media | Crítica | **Alto (Crítico)** | Alta |

### Priorización y Mitigación
Se asigna la máxima prioridad a la atomicidad en lote (**INT-03-07**) y al borrado en cascada (**INT-03-05**), dado que fallas en estas áreas destruyen la integridad referencial y corrompen el modelo de datos. La mitigación se logra ejerciendo llamadas HTTP directas y contrastándolas con el historial del Store, confirmando que la lógica transaccional y el `deleteAt` funcionen exactamente como se espera.

## 5. Justificación de herramientas

| Herramienta | Por qué se eligió | Alternativas consideradas | Limitación conocida |
|---|---|---|---|
| **Go test toolchain** | Estándar nativo de Go, excelente rendimiento e integración con el arnés del proyecto. | Postman/Newman, k6. | Requiere compilación y entorno local listo. |
| **Testify (require)** | Proporciona aserciones legibles y detiene la prueba inmediatamente si falla una precondición crítica. | Native Go assertions (`if err != nil`). | Dificulta la obtención de reportes acumulativos en un solo test si hay múltiples aserciones. |
| **TestHelper integrado** | Reutiliza la configuración oficial de base de datos de pruebas (SQLite en memoria) y los clientes HTTP autenticados. | Crear mocks personalizados. | Puede acoplarse estrechamente al comportamiento por defecto de la base de datos de prueba. |
| **Client HTTP del proyecto** | Garantiza que las pruebas de integración pasen exactamente por el mapeo JSON y los routers REST. | Llamar directamente a la estructura `App`. | Agrega latencia de red/sockets simulados durante la ejecución de las pruebas. |

## 6. Justificación de estrategia por caso de prueba

### INT-03-01: Crear tarjeta vía API y tipo `card`
La prueba valida la conversión transparente de `model.Card` a un `model.Block` de tipo `card` en el Store. Se comprueba que los metadatos de autoría (`createdBy`, `modifiedBy`) se inyecten de forma integrada a través de la cabecera de sesión.

### INT-03-02: Obtener tarjetas con propiedades personalizadas
Valida que el Store sea capaz de deserializar el mapa JSON de propiedades en el cliente Go sin alterar tipos de datos fundamentales.

### INT-03-03: Actualizar propiedades vía PATCH
Garantiza que el envío de una propiedad modificada mediante un PATCH actualice los campos correspondientes sin borrar las otras propiedades asignadas previamente a la tarjeta, probando la estrategia de mezcla (`merge`) en la capa App.

### INT-03-04: Insertar bloques de contenido y verificar parentesco
Valida la jerarquía de bloques. Un bloque de texto, imagen o checkbox debe estar atado operativamente a la tarjeta contenedora mediante `parentId`.

### INT-03-05: Eliminar tarjeta y bloques hijos (Borrado lógico)
Valida la propagación en cascada de la eliminación. Al borrar una tarjeta, el Store debe marcar los registros hijos con `deleteAt > 0`. Esto se comprueba a través del historial de bloques del Store.

### INT-03-06: Restaurar tarjeta y bloques hijos
Verifica que la operación inversa (`UndeleteBlock`) limpie el marcador de eliminación (`deleteAt = 0`) tanto en la tarjeta como en los bloques de contenido que dependen de ella.

### INT-03-07: Atomicidad de creación en lote
Comprueba que la inserción de múltiples tableros y bloques se maneje de forma atómica en una única transacción de base de datos. Si un bloque es inválido (por ejemplo, pertenece a un tablero inexistente), la operación debe fallar en su totalidad y el tablero válido que venía en el mismo lote no debe persistirse (Rollback).

## 7. Relación con pruebas unitarias del store Redux

Las pruebas unitarias del store Redux de la webapp validaban únicamente que los estados locales del frontend reaccionaran correctamente ante acciones puras de React/Redux. Sin embargo, no podían garantizar que el backend entendiera los payloads transmitidos o que la base de datos aplicara las restricciones transaccionales, de parentesco y de borrado lógico descritas aquí. Las pruebas INT-03 complementan esa capa asegurando que la API y el Store proporcionen un comportamiento confiable y consistente.
