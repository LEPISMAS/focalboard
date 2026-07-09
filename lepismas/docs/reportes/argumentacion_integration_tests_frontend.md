# Argumentación: Pruebas de Integración - Flujo Frontend (INT-10)

## 1. Estándares Aplicados

Las pruebas de integración para el flujo frontend-backend (INT-10) se han desarrollado siguiendo los siguientes estándares y mejores prácticas:

### 1.1 Estándar de Pruebas de Integración
- **ISO/IEC/IEEE 29119**: Metodología estructurada para diseño y ejecución de pruebas
- **Test Pyramid**: Enfoque equilibrado entre pruebas unitarias, de integración y E2E
- **AAA Pattern (Arrange-Act-Assert)**: Estructura clara de cada caso de prueba
- **Independencia de pruebas**: Cada prueba es autónoma y no depende del estado de otras

### 1.2 Estándar de Código
- **Go Testing Framework**: Uso de `testing` package de Go estándar
- **Testify/Require**: Aserciones claras y descriptivas
- **Nomenclatura descriptiva**: Nombres de funciones que indican propósito y caso de prueba
- **Logging informativo**: Uso de `t.Log()` para documentar propósito y resultados

### 1.3 Estándar de Simulación de Frontend
- **OctoClient Simulation**: Uso del mismo cliente HTTP que usa el frontend React
- **API Contract Validation**: Verificación de contratos REST entre frontend y backend
- **State Validation**: Verificación de persistencia de estado en Store

## 2. Risk Based Testing

### 2.1 Identificación de Riesgos

#### Riesgos de Producto
| ID Riesgo | Descripción | Probabilidad | Impacto | Nivel de Riesgo |
|-----------|-------------|--------------|---------|-----------------|
| RP-10-01 | Registro de usuario falla pero no reporta error al frontend | Media | Alta | **Alto** |
| RP-10-02 | Login exitoso pero sesión no persiste en Store | Media | Alta | **Alto** |
| RP-10-03 | Creación de tablero no se refleja en sidebar (caching) | Alta | Alta | **Alto** |
| RP-10-04 | Tarjeta creada no persiste correctamente | Media | Alta | **Alto** |
| RP-10-05 | Cambio de vista pierde datos o inconsistencia | Baja | Alta | **Medio** |
| RP-10-06 | Filtros no aplican correctamente o muestran datos incorrectos | Media | Media | **Medio** |
| RP-10-07 | Logout no invalida sesión correctamente (security breach) | Baja | Crítica | **Alto** |

#### Riesgos de Proyecto
| ID Riesgo | Descripción | Probabilidad | Impacto | Nivel de Riesgo |
|-----------|-------------|--------------|---------|-----------------|
| RPR-10-01 | Cambios en API REST rompen integración con frontend | Alta | Alta | **Alto** |
| RPR-10-02 | Cambios en schema de Store afectan persistencia de datos | Media | Alta | **Alto** |
| RPR-10-03 | Refactorización de OctoClient afecta comunicación | Media | Media | **Medio** |

### 2.2 Evaluación de Riesgos

**Matriz de Evaluación:**
- **Probabilidad × Impacto = Nivel de Riesgo**
- Alto: ≥ 6 (requiere cobertura máxima)
- Medio: 3-5 (requiere cobertura moderada)
- Bajo: ≤ 2 (cobertura mínima)

### 2.3 Priorización de Casos de Prueba

| Caso de Prueba | Riesgo Mitigado | Prioridad | Justificación |
|----------------|-----------------|-----------|---------------|
| INT-10-01 | RP-10-01, RP-10-02 | Alta | Flujo crítico de autenticación (riesgo alto) |
| INT-10-02 | RP-10-03 | Alta | Creación de tableros es core (riesgo alto) |
| INT-10-03 | RP-10-04 | Alta | Persistencia de tarjetas es fundamental (riesgo alto) |
| INT-10-04 | RP-10-05 | Media | Cambio de vista es importante pero edge case (riesgo medio) |
| INT-10-05 | RP-10-06 | Media | Filtros son feature importante (riesgo medio) |
| INT-10-06 | RP-10-07 | Alta | Logout es crítico para seguridad (riesgo alto) |

### 2.4 Mitigación mediante Cobertura de Prueba

**Áreas de Alto Riesgo - Cobertura Máxima:**
- **Autenticación** (INT-10-01, INT-10-06): Verificación completa de flujo de registro, login y logout
- **Persistencia de Datos** (INT-10-02, INT-10-03): Validación de que datos creados persisten correctamente
- **Seguridad de Sesión** (INT-10-06): Verificación de invalidación de tokens

**Áreas de Riesgo Medio - Cobertura Moderada:**
- **Funcionalidades de Vista** (INT-10-04, INT-10-05): Verificación de cambio de vista y filtros

## 3. Justificación de Uso de Herramientas

### 3.1 Herramientas Utilizadas

| Herramienta | Por qué se eligió | Alternativas consideradas | Limitación conocida |
|-------------|------------------|---------------------------|---------------------|
| **Go Testing Framework** | Framework nativo, integrado en el lenguaje, sin dependencias externas, soporte de concurrencia | TestNG, JUnit, PyTest | Sintaxis más verbosa que frameworks de otros lenguajes |
| **Testify/Require** | Aserciones claras con mensajes descriptivos, integración perfecta con Go testing, amplia adopción en comunidad Go | Gomega, assert package nativo | Requiere dependencia externa |
| **TestHelper (custom)** | Abstracción reutilizable para setup/teardown, reduce duplicación, mantiene consistencia entre pruebas | Setup/teardown inline en cada prueba | Requiere mantenimiento del helper |
| **OctoClient** | Simula calls reales de frontend React, valida integración completa Frontend→API→App→Store, usa misma infraestructura que producción | Cypress E2E, Puppeteer, HTTP test server | No prueba UI React real, solo capa de comunicación |
| **Store Direct Access** | Verificación directa de persistencia en base de datos, asegura que datos realmente se guardaron | Solo verificar respuesta API, mocks de Store | Acoplamiento a implementación de Store |

### 3.2 Decisiones de Arquitectura de Pruebas

**Enfoque de Integración Backend vs E2E Cypress:**
- Se eligió **integración backend con OctoClient** en lugar de Cypress E2E completo
- **Justificación**: 
  - Las pruebas Cypress existentes ya cubren flujos UI reales
  - Estas pruebas de integración complementan validando la capa de comunicación API
  - Más rápidas y estables que pruebas E2E completas
  - Permiten verificar persistencia directa en Store
- **Trade-off**: No prueban UI React real, pero validan integración completa de datos

**Simulación de Frontend:**
- Se usa **OctoClient** (mismo cliente que usa React) en lugar de mocks
- **Justificación**: Valida contrato real entre frontend y backend
- **Beneficio**: Detecta cambios en API que romperían frontend

## 4. Justificación de Estrategia por Caso de Prueba

### 4.1 INT-10-01: Registrar usuario y verificar login exitoso

**Estrategia:**
- Simular registro desde UI usando OctoClient.Register
- Verificar persistencia en Store (GetUserByUsername)
- Simular login desde UI usando OctoClient.Login
- Verificar token generado y acceso a workspace (GetMe)

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-01 y RP-10-02 (registro/login)
- **Validación de flujo completo**: Verifica que registro → persistencia → login funciona end-to-end
- **Contrato API**: Valida que endpoints /register y /login funcionan correctamente
- **Prioridad Alta**: Autenticación es flujo crítico; cualquier falla bloquea uso del sistema

### 4.2 INT-10-02: Crear tablero y verificar aparición en sidebar

**Estrategia:**
- Simular creación de tablero desde UI usando OctoClient.CreateBoard
- Verificar respuesta API con ID y título
- Simular consulta de tableros para sidebar (GetBoardsForTeam)
- Verificar que tablero creado aparece en listado
- Verificar persistencia en Store (GetBoard)

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-03 (creación y visibilidad)
- **Validación de caching**: Asegura que sidebar refleja cambios sin recarga
- **Contrato API**: Valida endpoints /boards y /teams/{teamID}/boards
- **Prioridad Alta**: Creación de tableros es funcionalidad core

### 4.3 INT-10-03: Crear tarjeta Kanban, verificar persistencia y edición

**Estrategia:**
- Crear tablero con vista Board (Kanban)
- Crear tarjeta usando OctoClient.CreateCard
- Verificar persistencia en Store (GetCard)
- Simular edición usando OctoClient.PatchCard
- Verificar que edición persistió en Store

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-04 (persistencia de tarjetas)
- **Validación de CRUD completo**: Verifica Create y Update de tarjetas
- **Integración con vistas**: Asegura que tarjetas funcionan en contexto de vista Kanban
- **Prioridad Alta**: Tarjetas son unidad fundamental de datos en Focalboard

### 4.4 INT-10-04: Cambiar de vista (Board → Table) con mismas tarjetas

**Estrategia:**
- Crear tablero con vistas Board y Table
- Crear tarjetas en el tablero
- Obtener tarjetas desde Store (simula carga en cualquier vista)
- Verificar que mismas tarjetas están disponibles independientemente de la vista

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-05 (consistencia entre vistas)
- **Validación de arquitectura**: Asegura que vistas son solo presentaciones, no almacenes separados
- **Experiencia de usuario**: Usuario espera ver mismos datos en diferentes layouts
- **Prioridad Media**: Importante pero edge case comparado con creación de datos

### 4.5 INT-10-05: Aplicar filtro y verificar tarjetas filtradas

**Estrategia:**
- Crear tablero con propiedades personalizadas (priority)
- Crear tarjetas con diferentes valores de propiedad
- Obtener todas las tarjetas desde Store
- Simular filtrado por propiedad (en aplicación real esto sería via API)
- Verificar que solo tarjetas que cumplen condición son retornadas

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-06 (correctitud de filtros)
- **Validación de propiedades**: Asegura que sistema de propiedades funciona correctamente
- **Feature importante**: Filtros son funcionalidad clave para gestión de tableros grandes
- **Prioridad Media**: Feature importante pero no crítico como autenticación

### 4.6 INT-10-06: Cerrar sesión y verificar invalidación de acceso

**Estrategia:**
- Registrar y hacer login de usuario
- Obtener token antes de logout
- Verificar que usuario está autenticado (GetMe)
- Simular logout desde UI (POST /logout)
- Intentar usar token anterior para acceder a /users/me
- Verificar que acceso es denegado (401 Unauthorized)
- Verificar que nuevo login funciona correctamente

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-10-07 (invalidación de sesión)
- **Validación de seguridad**: Asegura que logout realmente invalida tokens
- **Contrato API**: Valida endpoint /logout y manejo de tokens
- **Prioridad Alta**: Seguridad crítica; sesión no invalidada es vulnerabilidad

## 5. Integración con Pruebas Unitarias y E2E

Las pruebas de integración INT-10 complementan otras capas de pruebas:

- **Pruebas unitarias**: Validan lógica individual de componentes (App.CreateBoard, Store.InsertBoard)
- **Pruebas de integración INT-10**: Validan coordinación entre componentes y flujo Frontend→API→App→Store
- **Pruebas Cypress E2E**: Validan UI React real y flujos de usuario completos en navegador
- **Sin redundancia**: Las pruebas de integración no duplican validaciones unitarias ni E2E, sino que validan la capa de comunicación API

## 6. Mantenibilidad y Evolución

### 6.1 Diseño para Mantenimiento
- **TestHelper reutilizable**: Setup/teardown centralizado reduce duplicación
- **Nombres descriptivos**: Facilita identificación de propósito y fallas
- **Logging claro**: Ayuda en debugging cuando pruebas fallan
- **Simulación consistente**: Uso de OctoClient mantiene consistencia con frontend real

### 6.2 Escalabilidad
- **Patrón extensible**: Nuevos casos de prueba pueden seguir mismo patrón
- **Independencia**: Pruebas pueden ejecutarse en paralelo sin interferencia
- **Aislamiento de datos**: Cada prueba usa datos únicos para evitar colisiones

### 6.3 Evolución hacia Cypress E2E
- Estas pruebas de integración pueden servir como base para migrar a Cypress E2E
- La lógica de setup/teardown puede reutilizarse en comandos Cypress personalizados
- Los casos de prueba validados aquí pueden implementarse como specs Cypress

## 7. Limitaciones y Consideraciones

### 7.1 Limitaciones Actuales
- **No prueba UI React**: Estas pruebas validan comunicación API pero no interfaz de usuario real
- **No prueba WebSocket**: No validan actualizaciones en tiempo real
- **No prueba navegación**: No validan routing de React Router

### 7.2 Mitigación de Limitaciones
- **Complemento con Cypress**: Pruebas Cypress existentes cubren UI y navegación
- **Pruebas de WebSocket**: Casos específicos de WS se validan en pruebas separadas
- **Enfoque en datos**: Estas pruebas se enfocan en integridad de datos, no en UX

## 8. Conclusión

Las pruebas de integración INT-10 proporcionan cobertura robusta para el flujo frontend-backend, priorizando áreas de alto riesgo (autenticación, persistencia, seguridad) mientras mantienen mantenibilidad y escalabilidad. La estrategia de usar OctoClient para simular frontend asegura validación del contrato API real, y el diseño siguiendo estándares facilita mantenimiento y evolución futura. Estas pruebas complementan efectivamente las pruebas unitarias y E2E Cypress, proporcionando una capa de validación crítica entre frontend React y backend Go.
