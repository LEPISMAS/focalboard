# Plan de Pruebas de Sistema

## 1. Introducción

### 1.1 Propósito

Este documento define el plan de pruebas de sistema para Focalboard, con el objetivo de evaluar el comportamiento global de la aplicación mediante la validación de flujos completos de extremo a extremo. Las pruebas verifican la integración de todas las capas del sistema: interfaz de usuario, frontend, API, lógica del backend, persistencia en SQLite y el resultado visible para el usuario.

### 1.2 Alcance

El plan se centra en dos características críticas del sistema:

1. **Integridad del flujo de trabajo**: Validación del ciclo de vida de tableros, tarjetas y bloques, incluyendo creación, consulta, edición, movimiento, cambio de estado, eliminación, persistencia y consistencia tras recargas o reinicios.

2. **Seguridad y control de acceso**: Validación de reglas de acceso y visibilidad, incluyendo tableros privados y abiertos, usuarios autenticados y no autorizados, roles (Admin, Editor, Viewer), operaciones de lectura, modificación y eliminación, y prevención de escalamiento de privilegios.

El alcance excluye otras características de calidad como rendimiento, usabilidad, compatibilidad o migración de datos, salvo cuando sean necesarias para validar las dos características principales.

### 1.3 Sistema bajo prueba

**Nombre del sistema**: Focalboard (fork Lepismas)

**Arquitectura**:
- Frontend: React con Redux para gestión de estado
- Backend: Go con API REST
- Persistencia: SQLite (extensión json1 requerida)
- Comunicación en tiempo real: WebSocket
- Capas: UI → Frontend → API → App → Store → SQLite

**Componentes principales**:
- `webapp/src/`: Interfaz de usuario React
- `server/api/`: Endpoints REST
- `server/app/`: Lógica de negocio
- `server/store/`: Capa de persistencia
- `server/services/`: Servicios auxiliares (auth, permissions, notify)

### 1.4 Fuentes utilizadas

**Documentación de pruebas de integración**:
- `lepismas/docs/reportes/argumentacion_integration_tests_autenticacion.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_gestionDeTableros.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_gestionDeTarjetasYBloques.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_permisos.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_comparticionDeTableros.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_categoriasYBarraLateral.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_busqueda.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_frontend.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_importacionYExportacion.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_suscripcionesYNotificaciones.md`

**Reportes de ejecución**:
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_autenticacion.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_gestionDeTableros.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_categoriasYBarraLateral.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_comparticionDeTableros.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_busqueda.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_frontend.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_importacionYExportacion.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_suscripcionesYNotificaciones.md`

**Arquitectura y contexto**:
- `lepismas/docs/contexto/modulos_componentes.md`
- `lepismas/ISTQB_Pruebas_de_Software_Dinamicas.md`

**Casos de integración implementados**:
- `server/integrationtests/flujo_gestionDeTableros_int_test.go`
- `server/integrationtests/flujo_gestionDeTarjetasYBloques_int_tests.go`
- `server/integrationtests/flujo_permisos_int_tests.go`
- `server/integrationtests/flujo_comparticionDeTableros_int_tests.go`
- `server/integrationtests/flujo_categoriasYBarraLateral_int_tests.go`

## 2. Características evaluadas

### 2.1 Integridad del flujo de trabajo

**Descripción**: Validar que el ciclo de vida de tableros, tarjetas y bloques funcione correctamente de extremo a extremo, desde la interacción del usuario en la interfaz hasta la persistencia en SQLite y el resultado visible.

**Operaciones cubiertas**:
- Creación de tableros, tarjetas y bloques
- Consulta y visualización en diferentes vistas (Board, Table, Calendar, Gallery)
- Edición de propiedades y contenido
- Movimiento entre categorías y estados
- Cambio de visibilidad (privado ↔ público)
- Eliminación (soft delete) y restauración
- Persistencia en SQLite
- Actualización visible en la interfaz
- Consistencia después de recargar o reiniciar la aplicación

**Componentes involucrados**:
- Interfaz de usuario React
- Frontend (Redux, OctoClient)
- API REST
- Capa App (lógica de negocio)
- Store (SQLite)
- Sistema de categorías

### 2.2 Seguridad y control de acceso

**Descripción**: Validar que las reglas de acceso y visibilidad se cumplan en todo el sistema, garantizando que los usuarios solo puedan acceder a los recursos para los que tienen permisos.

**Operaciones cubiertas**:
- Acceso a tableros privados vs abiertos
- Validación de usuarios autenticados vs no autorizados
- Aplicación de roles (Admin, Editor, Viewer)
- Operaciones de lectura según permisos
- Operaciones de modificación según permisos
- Operaciones de eliminación según permisos
- Prevención de acceso a información restringida
- Prevención de escalamiento de privilegios
- Compartición pública mediante tokens
- Revocación de acceso

**Componentes involucrados**:
- Middleware de autenticación
- PermissionsService
- Capa API
- Capa App
- Store (tablas de membresías)
- Sistema de compartición (Sharing)

### 2.3 Justificación de la selección

**Integridad del flujo de trabajo**: Esta característica es fundamental porque representa el valor core de Focalboard. Los usuarios interactúan principalmente con tableros, tarjetas y bloques. Un fallo en el ciclo de vida de estas entidades afecta directamente la funcionalidad principal del sistema. Las pruebas de integración existentes (INT-02, INT-03, INT-06) demuestran que las capas individuales funcionan, pero se requiere validación de sistema para confirmar que el flujo completo desde la UI funciona correctamente.

**Seguridad y control de acceso**: Esta característica es crítica para la confidencialidad e integridad de los datos. Focalboard maneja información potencialmente sensible en tableros privados. Las pruebas de integración (INT-04, INT-05) validan que el backend aplique permisos correctamente, pero las pruebas de sistema son necesarias para verificar que un usuario no pueda evadir estas restricciones mediante la interfaz o manipulación del cliente.

## 3. Estrategia de prueba

### 3.1 Nivel de prueba

**Nivel**: Sistema (System Testing)

Las pruebas se ejecutan sobre el sistema completo integrado, validando el comportamiento end-to-end desde la perspectiva del usuario. Se diferencia de las pruebas de integración en que estas últimas validan la comunicación entre componentes, mientras que las pruebas de sistema validan el comportamiento observable del sistema como un todo.

### 3.2 Enfoque de caja negra

Las pruebas siguen un enfoque de caja negra, enfocándose en las entradas y salidas del sistema sin conocimiento detallado de la implementación interna. Se valida:

- El comportamiento observable desde la interfaz de usuario
- Las respuestas del sistema ante diferentes entradas
- El cumplimiento de contratos funcionales
- La consistencia de los datos persistidos

### 3.3 Pruebas manuales y automatizadas

**Pruebas manuales**: Recomendadas para la ejecución inicial y validación de flujos críticos, especialmente aquellos que requieren interacción compleja con la UI.

**Pruebas automatizadas**: Se propone el uso de Cypress para automatización de pruebas de UI, complementando las pruebas de integración existentes. Cypress permite validar la interacción real con el navegador y el frontend React.

**Propuesta de automatización**:
- [Pendiente: confirmar instalación y configuración de Cypress]
- Los casos de sistema pueden implementarse como specs Cypress
- Reutilización de datos de prueba y setup de las pruebas de integración

### 3.4 Criterios de entrada

- Pruebas de integración de los flujos base ejecutadas y pasadas (INT-02, INT-03, INT-04, INT-05, INT-06)
- Entorno de ejecución configurado con SQLite compatible (soporte json1)
- Frontend compilado y servidor backend ejecutándose
- Base de datos de prueba limpia o restaurada
- [Pendiente: confirmar disponibilidad de Cypress si se usa automatización]

### 3.5 Criterios de salida

- Todos los casos de prueba de sistema definidos ejecutados
- Tasa de passed ≥ 90% para considerar el ciclo exitoso
- Defectos críticos (bloqueantes) resueltos
- Evidencias de ejecución documentadas (capturas, logs, resultados)
- Matriz de trazabilidad completada
- Limitaciones y riesgos documentados

## 4. Técnicas de diseño de pruebas

### 4.1 Partición de Equivalencia

**Aplicación a Focalboard**:

**Para integridad del flujo de trabajo**:
- **Particiones válidas**: Tableros con datos válidos, tarjetas con propiedades completas, bloques con contenido válido
- **Particiones inválidas**: Tableros sin título, tarjetas sin propiedades requeridas, bloques con parentesco inválido

**Para seguridad y control de acceso**:
- **Particiones válidas**: Usuarios con permisos de lectura, usuarios con permisos de edición, usuarios con permisos de administración
- **Particiones inválidas**: Usuarios sin membresía, usuarios con rol inapropiado para la operación, usuarios no autenticados

**Ejemplo de aplicación**:
- Usuarios autenticados vs no autenticados
- Tableros privados vs tableros abiertos
- Roles Admin, Editor, Viewer

### 4.2 Análisis de Valores Límite

**Aplicación a Focalboard**:

**Límites documentados**:
- [Pendiente: verificar límite permitido para longitud de nombres de tableros]
- [Pendiente: verificar límite permitido para longitud de nombres de tarjetas]
- [Pendiente: verificar cantidad máxima de bloques por tarjeta]
- [Pendiente: verificar cantidad máxima de tarjetas por tablero]

**Nota**: No se inventarán límites que no estén documentados en el repositorio. Los casos de valores límite se agregarán cuando se identifiquen límites verificables en la documentación o código.

### 4.3 Tablas de Decisión

**Aplicación a Focalboard**:

**Para seguridad y control de acceso**:

| Condición | Usuario autenticado | Tiene membresía | Rol | Tablero privado | Tablero abierto | Operación | Resultado esperado |
|-----------|---------------------|-----------------|-----|-----------------|-----------------|-----------|--------------------|
| C1 | Sí | Sí | Admin | Sí | - | Leer | Permitido |
| C2 | Sí | Sí | Editor | Sí | - | Leer | Permitido |
| C3 | Sí | Sí | Viewer | Sí | - | Leer | Permitido |
| C4 | Sí | Sí | Admin | Sí | - | Editar | Permitido |
| C5 | Sí | Sí | Editor | Sí | - | Editar | Permitido |
| C6 | Sí | Sí | Viewer | Sí | - | Editar | Rechazado (403) |
| C7 | Sí | No | - | Sí | - | Leer | Rechazado (403) |
| C8 | Sí | Sí | - | - | Sí | Leer | Permitido |
| C9 | No | - | - | - | Sí | Leer | Permitido (con token) |
| C10 | No | - | - | Sí | - | Leer | Rechazado (403/404) |

**Para integridad del flujo de trabajo**:

| Condición | Tipo de entidad | Estado | Operación | Datos válidos | Resultado esperado |
|-----------|------------------|--------|-----------|---------------|--------------------|
| C1 | Tablero | Activo | Crear | Sí | Creado y persistido |
| C2 | Tablero | Activo | Crear | No (sin título) | Rechazado (400) |
| C3 | Tablero | Activo | Eliminar | - | Soft delete (deleteAt > 0) |
| C4 | Tarjeta | Activo | Crear | Sí | Creada y persistida |
| C5 | Tarjeta | Eliminada | Restaurar | - | Restaurada (deleteAt = 0) |
| C6 | Bloque | Activo | Mover | - | Parentesco actualizado |

### 4.4 Transición de Estados

**Aplicación a Focalboard**:

**Estados de tableros**:
- Activo → Eliminado (soft delete)
- Eliminado → Restaurado
- Privado → Público
- Público → Privado

**Estados de tarjetas**:
- Creada → Editada
- Activa → Eliminada (soft delete)
- Eliminada → Restaurada

**Estados de membresías**:
- No miembro → Viewer
- Viewer → Editor
- Editor → Admin
- Admin → Eliminado

**Diagrama de transición de tableros**:
```
[Activo] → (Eliminar) → [Eliminado]
[Eliminado] → (Restaurar) → [Activo]
[Activo] → (Cambiar visibilidad) → [Público/Privado]
[Público/Privado] → (Cambiar visibilidad) → [Privado/Público]
```

**Casos de prueba basados en transiciones**:
- SYS-01: Crear tablero → Verificar estado activo
- SYS-02: Eliminar tablero → Verificar soft delete
- SYS-03: Restaurar tablero → Verificar estado activo
- SYS-04: Cambiar visibilidad → Verificar nuevo estado

### 4.5 Error Guessing

**Aplicación a Focalboard**:

Basado en problemas observados en pruebas de integración y riesgos identificados:

**Escenarios de error guessing**:

1. **Pérdida de persistencia**: Crear tarjeta, cerrar navegador sin guardar, verificar que los datos persistieron
2. **Datos inconsistentes**: Editar tarjeta mientras se elimina el tablero padre, verificar integridad referencial
3. **Errores de SQLite**: Operaciones con campos JSON que fallen por falta de soporte json1
4. **Fallo con campos JSON**: Propiedades personalizadas con caracteres especiales o estructuras inválidas
5. **Acceso mediante URL directa**: Intentar acceder a un tablero privado mediante URL directa sin autenticación
6. **Sesión expirada**: Operar con una sesión expirada, verificar manejo de error
7. **Acciones repetidas rápidamente**: Crear múltiples tarjetas rápidamente, verificar que todas persistan
8. **Datos incompletos**: Crear tarjeta con campos obligatorios faltantes, verificar validación
9. **Concurrente**: Dos usuarios editando la misma tarjeta simultáneamente, verificar manejo de conflictos
10. **Token inválido tras revocación**: Usar token de compartición después de deshabilitar compartición

## 5. Trazabilidad con pruebas de integración

### 5.1 Flujos de integración utilizados como base

**Para integridad del flujo de trabajo**:
- **INT-02**: Gestión de Tableros (7 casos)
  - INT-02-01: Crear tablero persiste en Store
  - INT-02-02: Crear tablero crea membresía admin
  - INT-02-03: Listar tableros filtra por membresía
  - INT-02-04: Actualizar título persiste
  - INT-02-05: Eliminar tablero por soft delete
  - INT-02-06: Duplicar tablero copia bloques, propiedades y membresías
  - INT-02-07: Crear tablero y notificación WebSocket

- **INT-03**: Gestión de Tarjetas y Bloques (7 casos)
  - INT-03-01: Crear tarjeta vía API y tipo `card`
  - INT-03-02: Obtener tarjetas con propiedades personalizadas
  - INT-03-03: Actualizar propiedades vía PATCH
  - INT-03-04: Insertar bloques de contenido y verificar parentesco
  - INT-03-05: Eliminar tarjeta y bloques hijos (Borrado lógico)
  - INT-03-06: Restaurar tarjeta y bloques hijos
  - INT-03-07: Atomicidad de creación en lote

- **INT-06**: Categorías y Barra Lateral (5 casos)
  - INT-06-01: Crear categoría y verificar persistencia en el Store
  - INT-06-02: Mover tablero a categoría personalizada
  - INT-06-03: Obtener categorías de la barra lateral con tableros asociados
  - INT-06-04: Eliminar categoría y verificar que los tableros vuelven a la categoría por defecto
  - INT-06-05: Reordenar categorías y verificar persistencia del nuevo orden

**Para seguridad y control de acceso**:
- **INT-04**: Permisos (7 casos)
  - INT-04-01: Viewer intenta editar tablero
  - INT-04-02: Editor edita tablero
  - INT-04-03: Admin agrega miembro
  - INT-04-04: No miembro en tablero privado
  - INT-04-05: Cambiar rol de Viewer a Editor
  - INT-04-06: Eliminar membresía y perder acceso
  - INT-04-07: Tablero Open accesible por equipo

- **INT-05**: Compartición de Tableros (4 casos)
  - INT-05-01: Habilitar compartición y verificar generación del token
  - INT-05-02: Acceder al tablero con token válido sin autenticación
  - INT-05-03: Acceder al tablero con token inválido
  - INT-05-04: Deshabilitar compartición y verificar revocación

### 5.2 Evolución de integración a sistema

**Ejemplo de evolución**:

**Prueba de integración (INT-02-01)**:
- Objetivo: Verificar que crear un tablero vía API persiste en Store
- Alcance: API → App → Store
- Verificación: Consulta directa en Store/App

**Prueba de sistema correspondiente (SYS-01)**:
- Objetivo: Verificar que un usuario crea un tablero desde la UI, lo visualiza en el tablero, recarga la aplicación y confirma que el tablero continúa disponible
- Alcance: UI → Frontend → API → App → Store → UI
- Verificación: Observación directa en interfaz de usuario

**Diferencias clave**:
1. Las pruebas de sistema incluyen la interacción real del usuario con la UI
2. Validan la actualización visible en la interfaz
3. Verifican la consistencia tras recargas o reinicios
4. Ejercitan el flujo completo de extremo a extremo

### 5.3 Matriz de trazabilidad

| Caso de sistema | Característica | Requisito funcional | Caso de integración base | Componentes involucrados |
|-----------------|----------------|---------------------|--------------------------|--------------------------|
| SYS-01 | Integridad del flujo de trabajo | RF-004.1 (Crear tablero) | INT-02-01 | UI, Frontend, API, Backend, SQLite |
| SYS-02 | Integridad del flujo de trabajo | RF-004.1 (Crear tablero) | INT-02-02 | UI, Frontend, API, Backend, SQLite |
| SYS-03 | Integridad del flujo de trabajo | RF-004.2 (Listar tableros) | INT-02-03 | UI, Frontend, API, Backend, SQLite |
| SYS-04 | Integridad del flujo de trabajo | RF-004.3 (Editar tablero) | INT-02-04 | UI, Frontend, API, Backend, SQLite |
| SYS-05 | Integridad del flujo de trabajo | RF-004.4 (Eliminar tablero) | INT-02-05 | UI, Frontend, API, Backend, SQLite |
| SYS-06 | Integridad del flujo de trabajo | RF-008.1 (Crear tarjeta) | INT-03-01 | UI, Frontend, API, Backend, SQLite |
| SYS-07 | Integridad del flujo de trabajo | RF-008.2 (Editar tarjeta) | INT-03-03 | UI, Frontend, API, Backend, SQLite |
| SYS-08 | Integridad del flujo de trabajo | RF-008.3 (Eliminar tarjeta) | INT-03-05 | UI, Frontend, API, Backend, SQLite |
| SYS-09 | Integridad del flujo de trabajo | RF-008.3 (Restaurar tarjeta) | INT-03-06 | UI, Frontend, API, Backend, SQLite |
| SYS-10 | Integridad del flujo de trabajo | RF-012.1 (Organización de categorías) | INT-06-02 | UI, Frontend, API, Backend, SQLite |
| SYS-11 | Integridad del flujo de trabajo | RF-012.1 (Organización de categorías) | INT-06-04 | UI, Frontend, API, Backend, SQLite |
| SYS-12 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-01 | UI, Frontend, API, Permissions, SQLite |
| SYS-13 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-02 | UI, Frontend, API, Permissions, SQLite |
| SYS-14 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-04 | UI, Frontend, API, Permissions, SQLite |
| SYS-15 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-05 | UI, Frontend, API, Permissions, SQLite |
| SYS-16 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-06 | UI, Frontend, API, Permissions, SQLite |
| SYS-17 | Seguridad y control de acceso | RF-006.2 (Gestión de permisos) | INT-04-07 | UI, Frontend, API, Permissions, SQLite |
| SYS-18 | Seguridad y control de acceso | RF-007.1 (Compartir tableros) | INT-05-02 | UI, Frontend, API, Sharing, SQLite |
| SYS-19 | Seguridad y control de acceso | RF-007.1 (Compartir tableros) | INT-05-03 | UI, Frontend, API, Sharing, SQLite |
| SYS-20 | Seguridad y control de acceso | RF-007.1 (Compartir tableros) | INT-05-04 | UI, Frontend, API, Sharing, SQLite |

## 6. Casos de prueba de sistema

### 6.1 Casos de integridad del flujo de trabajo

| ID | Característica | Técnica | Flujo de integración base | Precondiciones | Descripción | Datos de prueba | Resultado esperado |
|----|----------------|---------|---------------------------|----------------|-------------|-----------------|--------------------|
| SYS-01 | Integridad del flujo de trabajo | Transición de Estados | INT-02-01 | Usuario autenticado, sesión activa | Usuario crea un tablero desde la UI, verifica que aparece en la barra lateral, recarga la página y confirma que el tablero continúa disponible | Título: "Tablero de prueba", Descripción: "Descripción de prueba" | Tablero visible en sidebar antes y después de recarga, datos persistidos en SQLite |
| SYS-02 | Integridad del flujo de trabajo | Partición de Equivalencia | INT-02-02 | Usuario autenticado | Usuario crea un tablero y verifica que automáticamente se asigna como administrador del mismo | Título: "Tablero admin" | Usuario aparece como admin en membresías del tablero |
| SYS-03 | Integridad del flujo de trabajo | Partición de Equivalencia | INT-02-03 | Dos usuarios autenticados con tableros privados distintos | Usuario 1 lista sus tableros y verifica que solo aparecen sus tableros privados, no los del usuario 2 | Tablero privado user1, Tablero privado user2 | Listado muestra solo tableros del usuario 1 |
| SYS-04 | Integridad del flujo de trabajo | Transición de Estados | INT-02-04 | Tablero existente | Usuario edita el título de un tablero desde la UI, guarda cambios y verifica que el título se actualiza en la interfaz y persiste tras recarga | Título original: "Old", Título nuevo: "New" | Título actualizado visible en UI y persistido en SQLite |
| SYS-05 | Integridad del flujo de trabajo | Transición de Estados | INT-02-05 | Tablero existente | Usuario elimina un tablero desde la UI, verifica que desaparece de la lista y que al recargar no aparece (soft delete) | Tablero ID existente | Tablero no visible en listado, marcado como eliminado en SQLite (deleteAt > 0) |
| SYS-06 | Integridad del flujo de trabajo | Transición de Estados | INT-03-01 | Tablero existente con vista Kanban | Usuario crea una tarjeta desde la UI, verifica que aparece en la columna correspondiente y persiste tras recarga | Título: "Tarjeta prueba", Propiedades: {status: "To Do"} | Tarjeta visible en columna, persistida en SQLite |
| SYS-07 | Integridad del flujo de trabajo | Partición de Equivalencia | INT-03-03 | Tarjeta existente con propiedades | Usuario edita una propiedad de la tarjeta (ej. cambiar estado), guarda y verifica que el cambio se refleja en la UI y persiste | Estado original: "To Do", Estado nuevo: "In Progress" | Propiedad actualizada visible en UI y persistida en SQLite |
| SYS-08 | Integridad del flujo de trabajo | Transición de Estados | INT-03-05 | Tarjeta existente con bloques hijos | Usuario elimina una tarjeta desde la UI, verifica que desaparece y que sus bloques hijos también se marcan como eliminados | Tarjeta con bloques de texto e imagen | Tarjeta y bloques marcados como eliminados en SQLite (deleteAt > 0) |
| SYS-09 | Integridad del flujo de trabajo | Transición de Estados | INT-03-06 | Tarjeta eliminada | Usuario restaura una tarjeta eliminada (si la funcionalidad existe en UI), verifica que reaparece con sus bloques hijos | Tarjeta previamente eliminada | Tarjeta y bloques restaurados (deleteAt = 0) |
| SYS-10 | Integridad del flujo de trabajo | Transición de Estados | INT-06-02 | Tablero y categoría personalizada existentes | Usuario mueve un tablero a una categoría personalizada desde la barra lateral, verifica que aparece en la nueva categoría | Tablero ID, Categoría ID | Tablero aparece en categoría destino en UI y SQLite |
| SYS-11 | Integridad del flujo de trabajo | Transición de Estados | INT-06-04 | Tablero en categoría personalizada | Usuario elimina una categoría personalizada, verifica que los tableros vuelven a la categoría por defecto "Boards" | Categoría personalizada con tableros | Tableros aparecen en categoría "Boards" en UI y SQLite |

### 6.2 Casos de seguridad y control de acceso

| ID | Característica | Técnica | Flujo de integración base | Precondiciones | Descripción | Datos de prueba | Resultado esperado |
|----|----------------|---------|---------------------------|----------------|-------------|-----------------|--------------------|
| SYS-12 | Seguridad y control de acceso | Tablas de Decisión | INT-04-01 | Usuario con rol Viewer en tablero privado | Usuario con rol Viewer intenta editar el título del tablero desde la UI | Rol: Viewer, Operación: Editar título | Edición rechazada, mensaje de error en UI, sin cambios en SQLite |
| SYS-13 | Seguridad y control de acceso | Tablas de Decisión | INT-04-02 | Usuario con rol Editor en tablero privado | Usuario con rol Editor edita el título del tablero desde la UI | Rol: Editor, Operación: Editar título | Edición exitosa, cambios persistidos en SQLite |
| SYS-14 | Seguridad y control de acceso | Partición de Equivalencia | INT-04-04 | Usuario sin membresía en tablero privado | Usuario sin membresía intenta acceder a un tablero privado mediante URL directa | URL directa a tablero privado, No autenticado | Acceso denegado (403/404), sin exposición de datos |
| SYS-15 | Seguridad y control de acceso | Transición de Estados | INT-04-05 | Usuario con rol Viewer en tablero privado | Admin cambia el rol del usuario de Viewer a Editor, usuario intenta editar y verifica que ahora puede | Rol inicial: Viewer, Rol final: Editor | Usuario puede editar tras cambio de rol |
| SYS-16 | Seguridad y control de acceso | Transición de Estados | INT-04-06 | Usuario con membresía en tablero privado | Admin elimina la membresía del usuario, usuario intenta acceder y verifica que es rechazado | Membresía eliminada | Acceso denegado tras eliminación de membresía |
| SYS-17 | Seguridad y control de acceso | Partición de Equivalencia | INT-04-07 | Tablero con visibilidad Open | Usuario del equipo accede a un tablero Open sin tener membresía explícita | Visibilidad: Open, Usuario del equipo | Acceso permitido sin membresía explícita |
| SYS-18 | Seguridad y control de acceso | Error Guessing | INT-05-02 | Tablero con compartición habilitada | Usuario no autenticado accede al tablero usando el token de compartición en la URL | Token válido compartido | Acceso permitido, tablero visible en UI |
| SYS-19 | Seguridad y control de acceso | Error Guessing | INT-05-03 | Tablero con compartición habilitada | Usuario no autenticado intenta acceder con un token inválido | Token inválido: "fake-token-123" | Acceso denegado (404), sin exposición de datos |
| SYS-20 | Seguridad y control de acceso | Transición de Estados | INT-05-04 | Tablero con compartición habilitada | Admin deshabilita la compartición, usuario intenta acceder con el token anteriormente válido | Token previamente válido, ahora deshabilitado | Acceso denegado (404), token revocado |

## 7. Entorno de ejecución

### 7.1 Componentes del entorno

**Requisitos del entorno**:
- **Sistema operativo**: Windows, Linux o macOS
- **Go**: Versión 1.21 o superior (para backend)
- **Node.js**: [Pendiente: verificar versión requerida para frontend]
- **SQLite**: Versión 3.38.0 o superior con soporte para extensión json1
- **Navegador**: Chrome, Firefox, Edge o Safari (para pruebas manuales)
- **Docker**: Opcional, recomendado para reproducibilidad

**Componentes a ejecutar**:
1. **Backend Go**: Servidor API (`server/bin/focalboard-server`)
2. **Frontend React**: Aplicación web compilada (`webapp/dist/`)
3. **Base de datos SQLite**: Base de datos de prueba con migraciones aplicadas
4. **Navegador**: Para ejecución de pruebas manuales o Cypress

### 7.2 Configuración de datos

**Base de datos de prueba**:
- Usar SQLite en memoria o archivo temporal
- Aplicar todas las migraciones del proyecto
- Limpiar o restaurar la base entre ejecuciones de pruebas
- Datos de prueba controlados: usuarios predefinidos, tableros de ejemplo

**Usuarios de prueba**:
- `user_admin`: Usuario con rol administrador del equipo
- `user_editor`: Usuario con permisos de edición
- `user_viewer`: Usuario con permisos de solo lectura
- `user_unauthorized`: Usuario sin membresías específicas

**Datos de ejemplo**:
- Tableros privados para cada usuario
- Tablero público (Open)
- Tablero con compartición habilitada
- Tarjetas con diferentes estados y propiedades
- Categorías personalizadas

### 7.3 Uso recomendado de Docker

**Propuesta de configuración Docker**:

```dockerfile
# Dockerfile para entorno de pruebas de sistema
FROM golang:1.21 AS backend
WORKDIR /app
COPY server/ ./server/
RUN cd server && go build -o bin/focalboard-server

FROM node:18 AS frontend
WORKDIR /app
COPY webapp/ ./webapp/
RUN cd webapp && npm install && npm run build

FROM alpine:latest
RUN apk add --no-cache sqlite
COPY --from=backend /app/server/bin/focalboard-server /usr/local/bin/
COPY --from=frontend /app/webapp/dist /usr/share/nginx/html
EXPOSE 8080
CMD ["focalboard-server"]
```

**Beneficios de Docker**:
- Reproducibilidad del entorno
- Aislamiento de dependencias
- Facilita ejecución en CI/CD
- Control de versiones de SQLite con soporte json1

### 7.4 Automatización de UI

**Propuesta de Cypress**:

**Instalación**:
```bash
cd webapp
npm install --save-dev cypress
npx cypress open
```

**Estructura de specs Cypress**:
```
cypress/
├── e2e/
│   ├── workflow/
│   │   ├── sys01_create_board.cy.js
│   │   ├── sys02_edit_board.cy.js
│   │   └── sys03_delete_board.cy.js
│   └── security/
│       ├── sys12_viewer_edit.cy.js
│       ├── sys14_unauthorized_access.cy.js
│       └── sys18_sharing_token.cy.js
├── support/
│   ├── commands.js
│   └── e2e.js
└── cypress.config.js
```

**Ejemplo de caso Cypress**:
```javascript
describe('SYS-01: Crear tablero desde UI', () => {
  it('Debe crear tablero y persistir tras recarga', () => {
    cy.login('user_admin');
    cy.visit('/boards');
    cy.contains('Create Board').click();
    cy.get('[data-testid="board-title"]').type('Tablero de prueba');
    cy.contains('Create').click();
    cy.contains('Tablero de prueba').should('be.visible');
    cy.reload();
    cy.contains('Tablero de prueba').should('be.visible');
  });
});
```

**Estado actual**:
- [Pendiente: confirmar instalación y configuración de Cypress]
- [Pendiente: verificar si Cypress ya está implementado en el repositorio]

## 8. Riesgos y limitaciones

### 8.1 Riesgos técnicos

| ID Riesgo | Descripción | Probabilidad | Impacto | Mitigación |
|-----------|-------------|--------------|---------|------------|
| RIESGO-01 | SQLite local sin soporte json1 causa fallos en migraciones | Media | Alta | Usar Docker con SQLite compatible o PostgreSQL |
| RIESGO-02 | Frontend no compilado correctamente impide ejecución de pruebas UI | Baja | Alta | Verificar build de frontend antes de ejecutar pruebas |
| RIESGO-03 | Dependencias de Node.js desactualizadas causan errores en Cypress | Baja | Media | Mantener dependencias actualizadas |
| RIESGO-04 | Concurrencia en base de datos causa inconsistencias en pruebas | Baja | Media | Ejecutar pruebas de forma secuencial o con aislamiento |
| RIESGO-05 | Cambios en API rompen contratos con frontend sin detectarse en pruebas de integración | Media | Alta | Mantener pruebas de sistema actualizadas con cambios de API |

### 8.2 Limitaciones del entorno

**Limitaciones identificadas en pruebas de integración**:

1. **SQLite sin soporte json1**: Las migraciones de Focalboard requieren la función `json_set` de SQLite. Versiones antiguas de SQLite no incluyen esta función, causando fallos en la inicialización de la base de datos de pruebas.

2. **Acceso a cache de Go**: En entornos Windows locales, se han reportado problemas de permisos en la cache de compilación de Go (`AppData\Local\go-build`), lo que impide la ejecución de pruebas.

3. **WebSocket end-to-end**: El arnés de `server/integrationtests` no registra backends de WebSocket ni de notificaciones, lo que limita la validación end-to-end de actualizaciones en tiempo real (documentado en INT-02-07 e INT-08-04).

### 8.3 Problemas conocidos de SQLite y JSON

**Problemas documentados en reportes de ejecución**:

1. **Error de migración**: `panic: driver: sqlite, message: failed when applying migration, command: apply_migration, originalError: no such function: json_set`

   **Causa**: La versión de SQLite disponible en el entorno no soporta la función `json_set` requerida por las migraciones de Focalboard.

   **Impacto**: Afecta a todas las pruebas de integración que requieren inicialización de base de datos.

   **Solución**: Usar SQLite 3.38.0+ con soporte json1, o configurar el entorno para usar PostgreSQL/MySQL.

2. **Requisito de build tag**: Las pruebas de integración requieren el tag `sqlite_json` para compilar correctamente con el soporte JSON de SQLite.

   **Comando**: `go test -v . -run TestINTXX -tags sqlite_json`

3. **Dependencia de extensión json1**: Las migraciones de base de datos versión 18 y superiores usan funciones JSON de SQLite. Sin la extensión json1, las migraciones fallan.

### 8.4 Dependencias externas

**Dependencias críticas**:

1. **SQLite con json1**: Requerido para migraciones de base de datos
2. **Go 1.21+**: Requerido para compilar y ejecutar backend
3. **Node.js**: Requerido para compilar frontend y ejecutar Cypress
4. **Navegador moderno**: Requerido para pruebas de UI manuales o automatizadas

**Riesgos de dependencias**:
- Cambios en versiones de SQLite pueden romper migraciones
- Actualizaciones de Go pueden introducir incompatibilidades
- Cambios en React pueden afectar renderizado de UI

## 9. Criterios de aceptación

**Criterios generales**:
- Todos los casos de prueba de sistema definidos han sido ejecutados
- La tasa de passed es ≥ 90%
- No hay defectos críticos (bloqueantes) sin resolver
- Las evidencias de ejecución están documentadas

**Criterios específicos por característica**:

**Integridad del flujo de trabajo**:
- SYS-01 a SYS-11 pasan: El ciclo de vida de tableros, tarjetas y bloques funciona correctamente
- La persistencia en SQLite es consistente
- La UI refleja correctamente el estado persistido
- Las recargas de página no pierden datos

**Seguridad y control de acceso**:
- SYS-12 a SYS-20 pasan: Las reglas de acceso se aplican correctamente
- Usuarios sin permisos no pueden acceder a recursos restringidos
- Los roles se respetan en todas las operaciones
- La compartición y revocación funcionan correctamente

**Criterios de automatización** (si se implementa Cypress):
- Los casos críticos (SYS-01, SYS-05, SYS-12, SYS-14) están automatizados
- Las pruebas automatizadas son estables y reproducibles
- El tiempo de ejecución de la suite es aceptable (< 15 minutos)

## 10. Evidencias esperadas

**Por cada caso de prueba**:
- Capturas de pantalla de la interfaz antes y después de la operación
- Logs de la consola del navegador (si aplica)
- Logs del servidor backend
- Estado de la base de datos antes y después (query SQLite)
- Registro de timestamp y ejecutor

**Documentación de ejecución**:
- Reporte de ejecución con resultados PASS/FAIL por caso
- Tiempo de ejecución por caso y total
- Defectos encontrados con severidad y estado
- Incidencias de entorno (si aplica)

**Evidencias de defectos**:
- Pasos para reproducir
- Comportamiento observado vs esperado
- Capturas del error
- Logs relevantes
- Impacto en el sistema

## 11. Conclusiones del plan

**Resumen**:
Este plan de pruebas de sistema define una estrategia para validar el comportamiento global de Focalboard, enfocándose en dos características críticas: integridad del flujo de trabajo y seguridad y control de acceso. El plan se basa en las pruebas de integración existentes y las evoluciona para validar flujos completos de extremo a extremo.

**Técnicas ISTQB aplicadas**:
- **Partición de Equivalencia**: Para clasificar usuarios, roles, tipos de tableros y estados válidos/inválidos
- **Análisis de Valores Límite**: Pendiente de aplicar cuando se identifiquen límites verificables en el sistema
- **Tablas de Decisión**: Para validar combinaciones de roles, permisos y operaciones en seguridad
- **Transición de Estados**: Para validar el ciclo de vida de tableros, tarjetas y membresías
- **Error Guessing**: Para identificar escenarios de fallo basados en problemas observados

**Flujos de integración utilizados como base**:
- INT-02 (Gestión de Tableros): 7 casos para integridad del flujo de trabajo
- INT-03 (Gestión de Tarjetas y Bloques): 7 casos para integridad del flujo de trabajo
- INT-06 (Categorías y Barra Lateral): 5 casos para integridad del flujo de trabajo
- INT-04 (Permisos): 7 casos para seguridad y control de acceso
- INT-05 (Compartición de Tableros): 4 casos para seguridad y control de acceso

**Cantidad de casos de sistema definidos**:
- Total: 20 casos de prueba de sistema
- Integridad del flujo de trabajo: 11 casos (SYS-01 a SYS-11)
- Seguridad y control de acceso: 9 casos (SYS-12 a SYS-20)

**Archivos consultados en lepismas**:
- `ISTQB_Pruebas_de_Software_Dinamicas.md`

**Archivos consultados en lepismas/docs**:
- `contexto/modulos_componentes.md`
- `reportes/argumentacion_integration_tests_autenticacion.md`
- `reportes/argumentacion_integration_tests_busqueda.md`
- `reportes/argumentacion_integration_tests_categoriasYBarraLateral.md`
- `reportes/argumentacion_integration_tests_comparticionDeTableros.md`
- `reportes/argumentacion_integration_tests_frontend.md`
- `reportes/argumentacion_integration_tests_gestionDeTableros.md`
- `reportes/argumentacion_integration_tests_gestionDeTarjetasYBloques.md`
- `reportes/argumentacion_integration_tests_importacionYExportacion.md`
- `reportes/argumentacion_integration_tests_permisos.md`
- `reportes/argumentacion_integration_tests_suscripcionesYNotificaciones.md`
- `reportes/reporte_ejecucion_pruebas_integracion_autenticacion.md`
- `reportes/reporte_ejecucion_pruebas_integracion_busqueda.md`
- `reportes/reporte_ejecucion_pruebas_integracion_categoriasYBarraLateral.md`
- `reportes/reporte_ejecucion_pruebas_integracion_comparticionDeTableros.md`
- `reportes/reporte_ejecucion_pruebas_integracion_frontend.md`
- `reportes/reporte_ejecucion_pruebas_integracion_gestionDeTableros.md`
- `reportes/reporte_ejecucion_pruebas_integracion_importacionYExportacion.md`
- `reportes/reporte_ejecucion_pruebas_integracion_suscripcionesYNotificaciones.md`

**Limitaciones detectadas**:
1. SQLite requiere versión 3.38.0+ con soporte json1 para ejecutar migraciones
2. Entornos Windows locales pueden tener problemas de permisos en cache de Go
3. El arnés de integración no soporta WebSocket end-to-end ni backends de notificaciones
4. No se identificaron límites verificables para aplicar Análisis de Valores Límite en el catálogo de requisitos proporcionado

**Placeholders pendientes de completar**:
- [Pendiente: verificar límite permitido para longitud de nombres de tableros]
- [Pendiente: verificar límite permitido para longitud de nombres de tarjetas]
- [Pendiente: verificar cantidad máxima de bloques por tarjeta]
- [Pendiente: verificar cantidad máxima de tarjetas por tablero]
- [Pendiente: confirmar instalación y configuración de Cypress]
- [Pendiente: verificar versión requerida para frontend (Node.js)]
- [Pendiente: identificar endpoint específico] (si aplica)

**Confirmación de que no se modificó código fuente**:
Este plan de pruebas de sistema es un documento de planificación y análisis. No se han realizado modificaciones al código fuente de Focalboard. El documento se ha creado exclusivamente en `docs/hito3/plan_pruebas_sistema.md` como entregable de la tarea.
