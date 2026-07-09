# Argumentación de pruebas de integración: Permisos (INT-04)

## 1. Objetivo del flujo

El flujo INT-04 tiene como objetivo validar que el sistema de control de acceso y permisos de Focalboard funcione correctamente en una arquitectura integrada. Los permisos determinan si una operación solicitada a través de la API REST es permitida o rechazada basándose en el rol del usuario (BoardMember) registrado en el Store.

El objetivo primordial es demostrar que la capa API, el PermissionsService y el Store colaboran estrechamente para restringir a usuarios sin permisos (como `viewer` o no miembros) y permitir a usuarios autorizados (como `editor` o administradores) realizar modificaciones en tableros privados, así como garantizar el libre acceso a tableros públicos (`Open`).

## 2. Capas integradas

Las pruebas de este flujo integran y ejercitan las siguientes capas del sistema:
- **API REST (server/api):** Endpoints como `/api/v2/boards/{boardID}`, `/api/v2/boards/{boardID}/members` y `/api/v2/boards/{boardID}/members/{userID}`.
- **PermissionsService (server/services/permissions):** Validación dinámica de los permisos (`view_board`, `manage_board_cards`, `manage_board_properties`, etc.).
- **Capa App (server/app):** Métodos como `GetBoard`, `PatchBoard`, `AddMemberToBoard`, `UpdateBoardMember` y `DeleteBoardMember`.
- **Store (server/store / sqlstore):** Persistencia física de tableros y membresías de tableros (`BoardMember`).

## 3. Estándares aplicados

Para garantizar el rigor metodológico, las pruebas se diseñaron bajo los siguientes estándares:
- **Trazabilidad estricta (IEEE 829 / ISO 29119):** Cada caso de prueba está identificado con la nomenclatura del plan de pruebas (`INT-04-01` a `INT-04-07`).
- **Validación de fronteras (Boundary testing) en seguridad:** No solo se valida la ruta positiva (usuarios que sí pueden), sino que se pone un énfasis especial en las rutas de rechazo, verificando que los códigos de error HTTP 403 Forbidden y la inalterabilidad de los recursos sean respetados.
- **Pruebas multi-usuario concurrentes:** Se emplean múltiples clientes HTTP autenticados de forma simultánea (`th.Client` y `th.Client2`) con tokens independientes para simular el comportamiento real de usuarios colaborando en el sistema.
- **Transparencia en el estado persistido:** Tras un rechazo de edición, se comprueba la base de datos para confirmar que realmente ningún cambio silencioso se haya aplicado en las tablas del Store.

## 4. Importancia desde Risk-Based Testing

### Identificación de Riesgos

**Riesgos de Producto:**
- **RP-04-01:** Usuario con rol de lectura (`viewer`) que logre evadir las validaciones y modifique el título o propiedades de un tablero.
- **RP-04-02:** Usuario con rol de edición (`editor`) bloqueado erróneamente por la capa de permisos, impidiendo la colaboración del equipo.
- **RP-04-03:** Falla de persistencia al registrar a un miembro en el tablero, dejándolo sin permisos operacionales o con un rol desalineado.
- **RP-04-04:** Fuga de información confidencial al permitir que un no miembro acceda, liste o lea un tablero marcado como `Private`.
- **RP-04-05:** Falla en la actualización del rol en caliente, donde el usuario mantenga permisos antiguos hasta reiniciar su sesión.
- **RP-04-06:** Latencia en la revocación de acceso. Un miembro eliminado que siga leyendo o editando el tablero después de ser expulsado.
- **RP-04-07:** Tablero público (`Open`) que imponga validaciones de membresía restrictivas, bloqueando el acceso general de los miembros del equipo.

**Riesgos de Proyecto:**
- **RProj-04-01:** Regresiones en la lógica de middleware de autenticación/autorización que expongan tableros privados globales.
- **RProj-04-02:** Defectos en el mapeo de roles de Mattermost que rompan el sistema de permisos en la versión Standalone de Focalboard.

### Evaluación: Probabilidad × Impacto/Severidad

| ID Caso | Probabilidad | Impacto | Nivel de RiesgoResultante | Prioridad |
|---------|:---:|:---:|:---:|:---:|
| **INT-04-01** (Viewer intenta editar) | Media | Alta | **Alto** | Alta |
| **INT-04-02** (Editor edita tablero) | Media | Media | **Medio** | Alta |
| **INT-04-03** (Admin agrega miembro) | Media | Alta | **Alto** | Alta |
| **INT-04-04** (No miembro en privado) | Baja | Crítica | **Alto (Confidencialidad)** | Alta |
| **INT-04-05** (Actualizar rol Viewer->Editor) | Media | Media | **Medio** | Media |
| **INT-04-06** (Eliminar miembro y revocación) | Media | Alta | **Alto** | Media |
| **INT-04-07** (Tablero Open accesible) | Media | Media | **Medio** | Media |

### Priorización y Mitigación
Se priorizan críticamente los casos de acceso no autorizado (**INT-04-04** e **INT-04-01**) para mitigar el riesgo de fuga de datos en tableros privados. La mitigación se implementa mediante la validación rigorosa de respuestas HTTP 403 Forbidden y la comprobación del estado de los registros en la base de datos tras las peticiones no autorizadas.

## 5. Justificación de herramientas

| Herramienta | Por qué se eligió | Alternativas consideradas | Limitación conocida |
|---|---|---|---|
| **Go test toolchain** | Integrado en el entorno de desarrollo de Focalboard. Permite automatizar ejecuciones rápidas y selectivas. | Newman/Postman. | Requiere compilación y base de datos activa para el arnés. |
| **Testify (require)** | Interrumpe la prueba ante fallos de autorización críticos, previniendo aserciones falsas sobre datos corrompidos. | Native Go assertions. | Menor nivel de detalle en fallas acumulativas. |
| **TestHelper / API Clients** | Provee dos clientes HTTP autenticados (`th.Client` y `th.Client2`) de forma nativa para simular interacciones multi-usuario. | Mocks manuales. | No permite validar llamadas directas a middlewares HTTP aislados. |
| **App / Store Services** | Permite registrar roles precondicionales en la base de datos de manera directa y segura sin depender de la UI. | Inserción directa por SQL. | Acoplado al modelo relacional de datos interno de Focalboard. |

## 6. Justificación de estrategia por caso de prueba

### INT-04-01: Viewer intenta editar tablero
Se añade al segundo usuario como `viewer` directamente a nivel de base de datos. Luego, el usuario intenta editar el tablero mediante PATCH. La prueba valida que la API retorne 403 y que el Store no refleje ningún cambio, confirmando la efectividad del control de acceso.

### INT-04-02: Editor edita tablero
Valida la ruta positiva. Un usuario con rol `editor` debe poder editar el título del tablero con éxito, obteniendo una respuesta 200 OK y persistencia confirmada en base de datos.

### INT-04-03: Admin agrega miembro
Verifica que un administrador del tablero pueda conceder acceso a un nuevo miembro mediante la API REST y que la membresía sea persistida en el Store con el rol esperado.

### INT-04-04: No miembro en tablero privado
Valida la confidencialidad más básica de la plataforma. Si un usuario no es miembro de un tablero privado, la API debe denegar cualquier intento de obtener (`GET`) el recurso con un código 403.

### INT-04-05: Cambiar rol de Viewer a Editor
Comprueba la escalación dinámica de permisos. El usuario inicialmente es rechazado al intentar editar, luego se le asciende de rol mediante API, y finalmente la edición se realiza de manera exitosa.

### INT-04-06: Eliminar membresía y perder acceso
Valida la revocación de permisos en tiempo real. Un usuario que tenía acceso deja de tenerlo inmediatamente en la siguiente llamada a la API tras ser eliminado del tablero por un administrador.

### INT-04-07: Tablero Open accesible por equipo
Comprueba el comportamiento por defecto de los tableros públicos (`Open`). Cualquier usuario que pertenezca al equipo del tablero debe poder leerlo sin necesidad de estar registrado explícitamente como miembro en las tablas de permisos de ese tablero.

## 7. Relación con pruebas unitarias del store Redux

Las pruebas del cliente React/Redux validaban únicamente la actualización visual de los estados locales del frontend ante un cambio de roles (como deshabilitar botones de edición en la interfaz de un `viewer`). No obstante, un atacante o error en el cliente podría omitir la interfaz web y llamar a la API directamente. Las pruebas de integración INT-04 aseguran que el servidor proteja los datos y aplique las restricciones de forma autónoma y robusta.
