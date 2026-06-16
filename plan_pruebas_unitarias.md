# Plan de Pruebas Unitarias - Focalboard

Este documento define el **Plan de Pruebas Unitarias** para la aplicación **Focalboard**. Su objetivo es establecer un marco estructurado que asegure la calidad del software, minimice la regresión de errores y garantice la robustez de los componentes del servidor y la interfaz web. Incluye la estrategia general, la especificación detallada de casos de prueba y la matriz de trazabilidad requerida.

---

# 1. Introducción y Objetivos

Las pruebas unitarias en Focalboard tienen como propósito validar el correcto funcionamiento de cada función, clase y módulo de forma aislada.

Los objetivos principales son:

- **Garantizar la Correctitud**: Validar que el código cumpla con los requisitos lógicos definidos.
- **Prevenir Regresiones**: Asegurar que los cambios futuros o la refactorización no rompan la funcionalidad existente.
- **Documentación Viva**: Servir como especificaciones claras del comportamiento esperado de la API y los componentes.
- **Integración Continua**: Facilitar el despliegue seguro mediante la automatización de la validación del código antes de cada mezcla (merge).

---

# 2. Alcance de las Pruebas

El alcance de las pruebas unitarias abarca los dos principales componentes del repositorio.

## 2.1 Backend (Servidor Go)

- **API y Enrutamiento**: Validación de controladores de peticiones, códigos de respuesta HTTP y serialización JSON.
- **Lógica de Negocio (App)**: Reglas relacionadas con bloques, tableros, usuarios, categorías e importación de datos.
- **Capa de Almacenamiento (Store/Database)**: Pruebas unitarias de consultas SQL usando SQLite en memoria para simular transacciones reales rápidamente.
- **Servicios de Terceros**: Mockeo de notificaciones y telemetría.

## 2.2 Frontend (Webapp React/TypeScript)

- **Utilidades y Helpers**: Funciones de procesamiento de fechas, manipulación de URLs y formateo de texto.
- **Componentes de Interfaz**: Renderizado lógico de botones, barras de navegación, diálogos de confirmación y formularios.
- **Gestión de Estado**: Reductores, acciones y manejadores de historial (Undo/Redo).
- **Mocks de API**: Intercepción de peticiones HTTP para probar los clientes API (`octoClient`).

---

# 3. Estrategia y Herramientas

La ejecución de pruebas unitarias utiliza diferentes pilas tecnológicas optimizadas para cada componente.

| Capa / Componente | Herramientas Utilizadas | Detalles de Implementación |
|-------------------|-------------------------|---------------------------|
| Backend (Go) | `testing`, `golang/mock` (mockgen), `testify` | Mocks para servicios externos. Base de datos SQLite local para pruebas rápidas. |
| Frontend (JS/TS) | `Jest`, `React Testing Library`, `fetch-mock` | Pruebas unitarias aisladas. Renderizado de componentes en DOM virtual (`jsdom`). |

---

# 4. Métricas de Éxito y Criterios de Aceptación

Para asegurar que las pruebas aporten valor real, se establecen las siguientes directrices.

## 4.1 Métricas de Calidad

- **Porcentaje de Aprobación**: **100%** de las pruebas deben pasar con éxito para permitir un despliegue o fusión de rama.
- **Cobertura Mínima (Coverage)**: Se aspira a una cobertura general de código de al menos **70%** en componentes críticos (lógica de negocio y utilidades).
- **Tiempo de Ejecución**: La suite rápida de pruebas de backend (SQLite) no debe superar los 3 minutos en entornos de desarrollo local.

## 4.2 Criterios para Escribir Nuevas Pruebas

1. **Aislamiento**: Ninguna prueba debe depender del estado dejado por otra prueba anterior.
2. **Determinismo**: Las pruebas deben dar el mismo resultado independientemente del entorno o la hora del sistema.
3. **Legibilidad**: El nombre de la prueba debe reflejar el escenario evaluado y el resultado esperado.
4. **Manejo de Datos de Prueba**: Toda creación de registros debe limpiarse al finalizar el caso de prueba (`defer` o `afterEach`).

---

# 5. Entornos y Comandos de Ejecución

Las pruebas se ejecutan de manera nativa mediante la consola de comandos de desarrollo.

## Ejecutar pruebas del Backend (Go + SQLite)

```bash
cd server && go test -tags 'json1 sqlite3' -race -v -count=1 ./...
```

## Ejecutar pruebas del Frontend (Jest)

```bash
cd webapp && npm run test
```

---

# 6. Catálogo Detallado de Casos de Prueba Unitarios

A continuación se enlistan y describen los casos de prueba unitarios más representativos implementados en el proyecto.

## 6.1 Casos de Prueba del Backend (Go)

### 6.1.1 Capa de Entrada / API (server/api)

| ID de Prueba | Nombre Técnico | Descripción / Qué Hace |
|--------------|----------------|------------------------|
| API-001 | TestErrorResponse | Verifica que el servidor serialice y devuelva correctamente las respuestas de error en formato JSON al cliente con el código HTTP apropiado. |
| API-002 | TestPing | Valida que la ruta `/ping` responda exitosamente, confirmando que la API está levantada y lista para recibir tráfico. |

### 6.1.2 Capa de Negocio / Lógica (server/app)

| ID de Prueba | Nombre Técnico | Descripción / Qué Hace |
|--------------|----------------|------------------------|
| APP-001 | TestLogin | Valida el proceso de login con credenciales correctas, la generación de la sesión de usuario y el rechazo ante credenciales incorrectas. |
| APP-002 | TestRegisterUser | Verifica la creación de nuevos usuarios asegurando que no se permitan nombres de usuario duplicados ni contraseñas débiles. |
| APP-003 | TestUpdateUserPassword | Valida el cambio de contraseña del usuario verificando la fortaleza de la nueva contraseña e invalidando las sesiones previas. |
| APP-004 | TestAddMemberToBoard | Valida la adición de miembros a un tablero verificando permisos y roles. |
| APP-005 | TestDuplicateBoard | Verifica la lógica de duplicación de tableros asegurando la copia de bloques, vistas y propiedades, pero no de membresías. |
| APP-006 | TestInsertBlock | Comprueba que la creación de bloques se realice correctamente respetando la estructura jerárquica. |
| APP-007 | TestPatchBlocks | Valida la actualización parcial de propiedades de bloques. |
| APP-008 | TestDeleteBlock | Valida la baja lógica de bloques y el mantenimiento del historial de cambios. |

## 6.2 Casos de Prueba del Frontend (TypeScript / React)

### 6.2.1 Módulo de Utilidades (webapp/src)

| ID de Prueba | Nombre Técnico | Descripción / Qué Hace |
|--------------|----------------|------------------------|
| WEB-001 | assureProtocol | Verifica que las URLs tengan un protocolo seguro asignado por defecto (`https://`). |
| WEB-002 | createGuid | Valida la generación pseudoaleatoria de IDs únicos. |
| WEB-003 | htmlFromMarkdown | Asegura que la conversión de Markdown a HTML evite ataques XSS. |
| WEB-004 | compareVersions | Valida el ordenamiento lógico de versiones semánticas. |
| WEB-005 | getUserDisplayName | Verifica la prioridad correcta para mostrar nombres de usuario. |

### 6.2.2 Gestión del Estado e Historial

| ID de Prueba | Nombre Técnico | Descripción / Qué Hace |
|--------------|----------------|------------------------|
| WEB-006 | UndoManager | Valida la pila de deshacer/rehacer y la reversibilidad de acciones. |

### 6.2.3 Componentes Visuales Reutilizables

| ID de Prueba | Nombre Técnico | Descripción / Qué Hace |
|--------------|----------------|------------------------|
| WEB-007 | calculation | Valida cálculos y sumatorias en columnas de tableros. |
| WEB-008 | cardDetail | Comprueba el renderizado correcto de propiedades, comentarios y adjuntos. |
| WEB-009 | sidebarCategory | Valida la expansión y colapso de categorías personalizadas. |

---

# 7. Especificación y Diseño de Pruebas Unitarias (Incremento de Cobertura)

Esta sección define el diseño detallado de las pruebas unitarias requeridas para alcanzar el objetivo de cobertura del 85%.

## 7.1 Mapeo General de Pruebas Unitarias

| ID Prueba | Tipo de Unidad | Componente / Paquete | Unidad de Código a Probar | Objetivo del Test |
|------------|----------------|----------------------|---------------------------|------------------|
| PU-APP-01 | Backend (Go) | server/app/blocks.go | App.CreateBlock | Validar la inserción de un bloque tipo tarjeta. |
| PU-APP-02 | Backend (Go) | server/app/blocks.go | App.DeleteBlock | Comprobar error controlado ante bloque inexistente. |
| PU-APP-03 | Backend (Go) | server/app/members.go | App.AddMember | Verificar asignación de miembros y auditoría. |
| PU-APP-04 | Backend (Go) | server/app/import.go | App.ImportTrello | Validar manejo de JSON corrupto. |
| PU-APP-05 | Backend (Go) | server/app/export.go | App.ExportBoardArchive | Validar generación correcta del archivo ZIP. |
| PU-UI-01 | Frontend (React) | CardDetail | Componente CardDetail | Verificar publicación de comentarios. |
| PU-UI-02 | Frontend (React) | BoardView | Componente BoardViewSelector | Validar cambio de vista. |
| PU-UI-03 | Frontend (React) | Sidebar | Componente Sidebar | Verificar renderizado de tableros. |
| PU-UI-04 | Frontend (React) | Settings | Componente SettingsDialog | Validar cambio de tema visual. |
| PU-UI-05 | Frontend (React) | Properties | Componente CardPropertyDate | Verificar actualización de fecha. |

## 7.2 Especificación de Pruebas Unitarias del Backend (Go)

### PU-APP-01: Creación de Bloque de Tarjeta Exitosa (CreateBlock)

**Unidad Bajo Prueba:** `App.CreateBlock(block *model.Block, userID string)`

**Precondición:** El usuario creador existe y la sesión es válida.

**Mocks Configurados:** `Mock Store.CreateBlock(block)` retorna `nil` y asigna un ID aleatorio.

**Datos de Entrada:**

- `block`: tipo `"card"` con título `"Nueva Tarea"`
- `userID`: `"user-admin-1"`

**Aserciones:**

```go
assert.NoError(t, err)
assert.NotEmpty(t, returnedBlock.ID)
assert.Equal(t, "Nueva Tarea", returnedBlock.Title)
```

### PU-APP-02: Retorno de Error al Eliminar un Bloque Inexistente (DeleteBlock)

**Unidad Bajo Prueba:** `App.DeleteBlock(blockID string, modifiedBy string)`

**Mocks Configurados:** `Mock Store.GetBlock(blockID)` retorna `ErrNotFound`.

**Datos de Entrada:**

- `blockID`: `"invalid-block-uuid-9999"`
- `modifiedBy`: `"user-editor-1"`

**Aserciones:**

```go
assert.Error(t, err)
assert.True(t, model.IsErrNotFound(err))
```

### PU-APP-03: Asociación Exitosa de Miembro a Tablero (AddMember)

**Unidad Bajo Prueba:** `App.AddMember(boardID string, userID string, role string)`

**Precondición:** Existen el tablero y usuario.

**Aserciones:**

```go
assert.NoError(t, err)
```

El servicio de auditoría debe ser invocado exactamente una vez.

### PU-APP-04: Excepción de Formato en Importador de Trello (ImportTrello)

**Unidad Bajo Prueba:** `App.ImportTrello(reader io.Reader, workspaceID string)`

**Datos de Entrada:**

```json
{ invalid json data ... }
```

**Aserciones:**

```go
assert.Error(t, err)
```

Las llamadas a `Store.SaveBlock` nunca deben ejecutarse.

### PU-APP-05: Validación de Estructura de Exportación (ExportBoardArchive)

**Unidad Bajo Prueba:** `App.ExportBoardArchive(boardID string, writer io.Writer)`

**Aserciones:**

```go
assert.NoError(t, err)
```

Los primeros bytes generados deben corresponder a la firma ZIP:

```text
PK\003\004
```

## 7.3 Especificación de Pruebas Unitarias del Frontend (React/Jest)

### PU-UI-01: Envío de Comentario en CardDetail

**Componente Bajo Prueba:**

```tsx
<CardDetail card={mockCard} />
```

**Interacción:**

- Escribir `"Listo para revisión"`
- Clic en `"Send"`

**Aserciones:**

- `postComment` recibe `"Listo para revisión"`
- El textarea queda vacío.

### PU-UI-02: Cambio de Vista en BoardViewSelector

**Interacción:**

- Clic sobre la pestaña `"Table"`

**Aserciones:**

- Se ejecuta `onChangeView` con el ID de la vista de tabla.

### PU-UI-03: Carga de Tableros en Sidebar

**Store Mock:**

```json
[
  {"id":"b1","title":"Sprint 1"},
  {"id":"b2","title":"Planificación"}
]
```

**Aserciones:**

```javascript
expect(screen.getByText("Sprint 1")).toBeInTheDocument()
expect(screen.getByText("Planificación")).toBeInTheDocument()
```

### PU-UI-04: Seteo de Tema Oscuro en SettingsDialog

**Interacción:**

- Seleccionar `"Theme"`
- Elegir `"Dark"`

**Aserciones:**

- Se dispara `setTheme("dark")`
- El contenedor raíz cambia a `.theme-dark`

### PU-UI-05: Modificación de Fecha Límite en CardPropertyDate

**Interacción:**

- Clic en la propiedad de fecha
- Selección del día `"15"`

**Aserciones:**

- `mockOnChange` recibe un objeto `Date` correspondiente al día seleccionado.

---

# 8. Matriz de Trazabilidad: Requisito → Caso de Prueba

| ID Requisito | Descripción del Requisito Funcional | ID de Caso(s) de Prueba Asociado(s) |
|-------------|--------------------------------------|-------------------------------------|
| HU-01 | Gestión de Tableros (Crear, Duplicar, Eliminar) | APP-005, PU-UI-03 |
| HU-02 | Gestión de Tarjetas (Crear, Editar, Eliminar Bloque) | APP-006, APP-008, PU-APP-01, PU-UI-01 |
| HU-03 | Sistema de Miembros y Roles en Tableros | APP-004, PU-APP-03 |
| HU-04 | Autenticación y Gestión de Usuarios | APP-001, APP-002, APP-003 |
| HU-05 | Persistencia de Datos y Almacenamiento | PU-APP-01, APP-005 |
| HU-06 | Funcionalidad de Deshacer/Rehacer (Undo/Redo) | WEB-006 |
| HU-07 | Importación de Datos desde Trello | PU-APP-04 |
| HU-08 | Exportación de Tableros (Copia de Seguridad) | PU-APP-05 |
| HU-09 | Vistas Alternativas del Tablero | PU-UI-02 |
| HU-10 | Configuración de Preferencias de Usuario | PU-UI-04 |
| HU-11 | Visualización de Detalle y Comentarios en Tarjetas | WEB-008, PU-UI-01 |
| HU-12 | Propiedades Personalizadas en Tarjetas | PU-UI-05, WEB-007 |
| HU-13 | Utilidades de Formato y Seguridad | WEB-001, WEB-003 |
| HU-14 | Identificadores Únicos y Versionado Semántico | WEB-002, WEB-004 |

---
**Fin del Documento**