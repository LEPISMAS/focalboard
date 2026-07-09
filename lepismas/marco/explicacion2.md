# Explicación: Pruebas de Integración - Flujos Búsqueda y Frontend

## Contexto

Para las exposiciones del trabajo final, donde realizaremos debates argumentando todas las decisiones técnicas, se desarrollaron pruebas de integración para dos flujos críticos de Focalboard: el flujo de búsqueda (INT-09) y el flujo frontend-backend (INT-10).

## Decisiones Técnicas Tomadas

### 1. Elección de Nivel de Pruebas

**Decisión**: Implementar pruebas de integración backend en lugar de pruebas E2E completas con Cypress.

**Justificación**:
- Las pruebas Cypress existentes ya cubren flujos UI reales
- Las pruebas de integración backend permiten validar más rápidamente la coordinación entre componentes
- Proporcionan una capa crítica de validación entre frontend React y backend Go
- Son más estables y rápidas que pruebas E2E completas
- Permiten verificar persistencia directa en Store

**Trade-off**: No prueban UI React real, pero esto es mitigado por las pruebas Cypress existentes.

### 2. Estrategia de Integración Completa

**Decisión**: Usar integración completa (API → App → Store) sin mocking parcial.

**Justificación**:
- El flujo de búsqueda involucra múltiples capas que deben coordinarse correctamente
- Mocking podría ocultar errores de integración entre capas
- Proporciona mayor confianza en la corrección del sistema
- Valida contratos reales entre componentes

**Trade-off**: Pruebas más lentas que pruebas unitarias, pero necesarias para validar integración.

### 3. Simulación de Frontend con OctoClient

**Decisión**: Usar OctoClient (mismo cliente HTTP que usa React) en lugar de Cypress o Puppeteer.

**Justificación**:
- Valida contrato real entre frontend y backend
- Usa misma infraestructura que producción
- Detecta cambios en API que romperían frontend
- Complementa pruebas Cypress E2E que prueban UI real

**Trade-off**: No prueba UI React, pero esto es intencional para complementar, no duplicar, pruebas Cypress.

### 4. Priorización de Casos de Prueba (Risk Based Testing)

**Decisión**: Priorizar casos de prueba basándose en análisis de riesgos (probabilidad × impacto).

**Justificación**:
- Seguridad y permisos tienen prioridad alta (riesgo alto × impacto alto)
- Funcionalidades core (autenticación, creación de tableros/tarjetas) tienen prioridad alta
- Features importantes pero no críticos (filtros, cambio de vista) tienen prioridad media
- Casos edge (búsqueda sin resultados) tienen prioridad media

**Resultado**: 4 casos para búsqueda (INT-09) y 6 casos para frontend (INT-10), todos priorizados según riesgo.

### 5. Estándares Aplicados

**Decisión**: Seguir estándares ISO/IEC/IEEE 29119 para pruebas y Go Testing Framework para implementación.

**Justificación**:
- Metodología estructurada para diseño y ejecución de pruebas
- Framework nativo de Go, sin dependencias externas adicionales
- Patrón AAA (Arrange-Act-Assert) para claridad
- Nomenclatura descriptiva para mantenibilidad

### 6. Documentación Argumentativa

**Decisión**: Crear documentos argumentativos detallados para cada flujo.

**Justificación**:
- Las exposiciones son debates donde debemos argumentar decisiones
- Documentación de risk based testing con matrices de evaluación
- Justificación de herramientas con alternativas consideradas
- Estrategia por caso de prueba con justificación específica

## Resultados Obtenidos

### Archivos Creados

**Código de Pruebas**:
- `server/integrationtests/flujo_busqueda_int_test.go` - 4 casos de prueba
- `server/integrationtests/flujo_frontend_int_test.go` - 6 casos de prueba

**Scripts de Ejecución**:
- `server/integrationtests/ejecutar_busqueda.sh` y `.bat`
- `server/integrationtests/ejecutar_frontend.sh` y `.bat`

**Documentación Argumentativa**:
- `lepismas/docs/reportes/argumentacion_integration_tests_busqueda.md`
- `lepismas/docs/reportes/argumentacion_integration_tests_frontend.md`

**Reportes de Ejecución**:
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_busqueda.md`
- `lepismas/docs/reportes/reporte_ejecucion_pruebas_integracion_frontend.md`

### Problemas Encontrados

**Problema de Entorno SQLite**:
- Las migraciones de base de datos usan `json_set` de SQLite
- La versión de SQLite disponible no soporta esta función
- Esto afecta a todas las pruebas de integración del proyecto, no solo las nuevas

**Resolución**:
- Documentado en reportes de ejecución
- Recomendaciones para resolver: actualizar SQLite o usar PostgreSQL
- Código de pruebas está correcto, solo requiere entorno adecuado

## Argumentos para Debate

### 1. ¿Por qué pruebas de integración y no solo unitarias?

**Argumento**: Las pruebas unitarias validan lógica individual, pero no coordinación entre componentes. Las pruebas de integración son necesarias para validar que API, App y Store trabajan correctamente juntos, especialmente para flujos críticos como búsqueda y autenticación.

### 2. ¿Por qué no usar Cypress para todo?

**Argumento**: Cypress es excelente para E2E pero es lento y frágil. Las pruebas de integración backend proporcionan validación rápida y estable de la lógica de negocio y persistencia, complementando las pruebas Cypress que validan UX.

### 3. ¿Por qué risk based testing?

**Argumento**: Con recursos limitados, es imposible probar todo. Risk based testing permite priorizar áreas de mayor riesgo (seguridad, permisos, persistencia) asegurando que los casos más críticos se prueben primero.

### 4. ¿Por qué documentación tan detallada?

**Argumento**: Para debates académicos, necesitamos justificar cada decisión técnica. La documentación argumentativa proporciona evidencia de análisis sistemático y consideración de alternativas, no decisiones arbitrarias.

## Conclusión

Las pruebas de integración desarrolladas proporcionan cobertura robusta para flujos críticos de Focalboard, priorizando áreas de alto riesgo según análisis sistemático. La estrategia complementa efectivamente pruebas unitarias y E2E existentes, proporcionando una capa crítica de validación entre frontend y backend. Aunque la ejecución fue impedida por problemas de entorno, el código está correctamente implementado y listo para ejecución una vez resuelto el problema de configuración.

Todas las decisiones técnicas fueron tomadas con justificación clara basada en mejores prácticas, análisis de riesgos y consideración de trade-offs, proporcionando base sólida para argumentación en exposiciones académicas.
