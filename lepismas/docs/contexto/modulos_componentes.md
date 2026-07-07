## Módulos y Componentes del Sistema

Focalboard está estructurado en torno a conceptos fuertemente cohesionados. Aquí identificamos los módulos principales y los componentes individuales que los conforman.

### A. Módulo del Servidor (Backend Go)

1. **Modelo de Bloques (`server/model` y `server/app/blocks.go`):**
   * *Componentes:* `Block`, `Board`, `Card`.
   * *Detalle:* Todo en Focalboard se abstrae bajo el concepto de **Block** (bloque). Un bloque puede representar un tablero, una tarjeta, una vista, un bloque de texto, o una imagen. Cada bloque cuenta con un `id`, `parent_id`, `type`, `schema`, `fields` (campos dinámicos en formato JSON) e información de auditoría.
2. **Capa de Persistencia (`server/services/store`):**
   * *Componentes:* `Store` (interfaz común), `SQLStore` (SQLite, MySQL, PostgreSQL).
   * *Detalle:* Provee el acceso a datos. Utiliza el generador de consultas fluidas `squirrel`. Cuenta con un sistema interno de migraciones automáticas (`sqlstore/migrate.go` y directorio `migrations/`) que asegura que el esquema esté actualizado al arrancar el servidor.
3. **Módulo de Lógica de Negocio (`server/app`):**
   * *Componentes:* `boards.go`, `cards.go`, `category.go`, `files.go`, `sharing.go`, `import.go`, `export.go`.
   * *Detalle:* Contiene las reglas de negocio esenciales. Valida los permisos de los usuarios antes de ejecutar cualquier acción de creación o modificación, gestiona el flujo de trabajo de plantillas, organiza la agrupación de tableros por categorías y gestiona la lógica de importación y exportación (ej. compatibilidad con archivos `.boardarchive`).
4. **Módulo API REST (`server/api`):**
   * *Componentes:* `api.go`, `api_boards.go`, `api_blocks.go`, `api_users.go`.
   * *Detalle:* Expone la interfaz HTTP que consume el frontend. Se encarga de procesar los parámetros de las solicitudes, deserializar el cuerpo JSON, autenticar las cabeceras de sesión mediante middlewares, invocar la capa de negocio (`app`) y retornar respuestas REST uniformes.
5. **Módulo de Comunicación en Tiempo Real (`server/ws`):**
   * *Componentes:* `ws.go`, `server.go`.
   * *Detalle:* Implementa un servidor de WebSockets. Cuando un cliente modifica un bloque a través del enrutador API, el servidor notifica a todos los demás clientes con sesión activa en el mismo tablero sobre los cambios para actualizar la UI en vivo.
6. **Módulos de Servicios Auxiliares (`server/services/`):**
   * *Componentes:*
     * `auth`: Controla el ciclo de vida de las sesiones y la encriptación de contraseñas.
     * `notify`: Envía notificaciones de cambios en tarjetas basadas en suscripciones.
     * `permissions`: Motor de control de accesos que evalúa si un rol tiene permisos de lectura, escritura o administración sobre un tablero o equipo.
     * `metrics` / `telemetry`: Monitoreo del rendimiento con Prometheus y telemetría de comportamiento de uso.

### B. Módulo de Aplicación Web (Frontend React)

1. **Gestión de Estado y API (`webapp/src/store`, `octoClient.ts`, `wsclient.ts`):**
   * *Componentes:* `octoClient` (llamadas HTTP), `wsClient` (suscripción a WebSockets), `Redux Slices` (User, Boards, Cards, Comments, Views, activeView).
   * *Detalle:* Mantiene el estado en sincronía con el servidor. Cuando se carga la página, `octoClient` solicita los datos necesarios y `Redux` los organiza en el árbol de estado global.
2. **Módulo de Mutaciones e Historial (`webapp/src/mutator.ts`, `webapp/src/undomanager.ts`):**
   * *Componentes:* `Mutator`, `UndoManager`.
   * *Detalle:* El `Mutator` actúa como una capa intermedia para realizar cambios en los bloques. Genera y despacha operaciones de manera optimista en el cliente para dar una experiencia fluida, y si una acción falla, revierte el estado. `UndoManager` almacena una pila de estados anteriores para permitir operaciones "deshacer" (Ctrl+Z) y "rehacer" (Ctrl+Y) de forma nativa.
3. **Componentes Visuales de la Interfaz (`webapp/src/components/`):**
   * *Componentes:*
     * `kanban/`: La vista clásica de columnas basadas en propiedades (ej. Estado).
     * `table/`: Vista tabular que permite organizar propiedades y filtros detallados.
     * `calendar/`: Mapea tarjetas que tienen propiedades de fecha a un calendario mensual interactivo.
     * `gallery/`: Muestra tarjetas como cajas o mosaicos con imágenes de portada.
     * `cardDetail/`: Modal interactivo que se abre al hacer clic en una tarjeta; permite modificar títulos, propiedades personalizadas, agregar comentarios, adjuntar archivos y ver el historial de cambios.
     * `sidebar/`: Panel de navegación lateral para cambiar de tablero, colapsar categorías y gestionar la visibilidad de tableros privados frente a compartidos.