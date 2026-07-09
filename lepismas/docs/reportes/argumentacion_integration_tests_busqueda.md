# Argumentación: Pruebas de Integración - Flujo Búsqueda (INT-09)

## 1. Estándares Aplicados

Las pruebas de integración para el flujo de búsqueda (INT-09) se han desarrollado siguiendo los siguientes estándares y mejores prácticas:

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

### 1.3 Estándar de Documentación
- **TDD (Test-Driven Documentation)**: Las pruebas sirven como documentación viva
- **Comentarios en código**: Explicación de flujo de integración verificado
- **Mensajes de aserción**: Claros y específicos sobre lo que se verifica

## 2. Risk Based Testing

### 2.1 Identificación de Riesgos

#### Riesgos de Producto
| ID Riesgo | Descripción | Probabilidad | Impacto | Nivel de Riesgo |
|-----------|-------------|--------------|---------|-----------------|
| RP-01 | Búsqueda retorna tableros privados de otros usuarios (violación de seguridad) | Media | Alta | **Alto** |
| RP-02 | Inyección SQL a través de términos de búsqueda | Baja | Crítica | **Alto** |
| RP-03 | Búsqueda no respeta permisos de usuario | Media | Alta | **Alto** |
| RP-04 | Búsqueda con términos especiales causa errores de servidor | Media | Media | **Medio** |
| RP-05 | Búsqueda sin resultados retorna error en lugar de array vacío | Baja | Baja | **Bajo** |

#### Riesgos de Proyecto
| ID Riesgo | Descripción | Probabilidad | Impacto | Nivel de Riesgo |
|-----------|-------------|--------------|---------|-----------------|
| RPR-01 | Cambios en capa de Store rompen integración de búsqueda | Media | Alta | **Alto** |
| RPR-02 | Refactorización de API afecta endpoints de búsqueda | Media | Media | **Medio** |
| RPR-03 | Cambios en lógica de permisos afectan filtrado de resultados | Alta | Alta | **Alto** |

### 2.2 Evaluación de Riesgos

**Matriz de Evaluación:**
- **Probabilidad × Impacto = Nivel de Riesgo**
- Alto: ≥ 6 (requiere cobertura máxima)
- Medio: 3-5 (requiere cobertura moderada)
- Bajo: ≤ 2 (cobertura mínima)

### 2.3 Priorización de Casos de Prueba

| Caso de Prueba | Riesgo Mitigado | Prioridad | Justificación |
|----------------|-----------------|-----------|---------------|
| INT-09-01 | RP-01, RP-03 | Alta | Verifica seguridad y permisos (riesgo alto) |
| INT-09-02 | RP-05 | Media | Verifica comportamiento esperado (riesgo bajo) |
| INT-09-03 | RP-01, RP-03 | Alta | Verifica aislamiento de permisos (riesgo alto) |
| INT-09-04 | RP-02, RP-04 | Media | Verifica robustez contra inyección (riesgo alto) |

### 2.4 Mitigación mediante Cobertura de Prueba

**Áreas de Alto Riesgo - Cobertura Máxima:**
- **Permisos y Seguridad** (INT-09-01, INT-09-03): Verificación exhaustiva de que los usuarios solo ven tableros accesibles
- **Inyección SQL** (INT-09-04): Pruebas con múltiples patrones de ataque para validar sanitización

**Áreas de Riesgo Medio - Cobertura Moderada:**
- **Manejo de Casos Edge** (INT-09-02): Verificación de comportamiento con búsquedas sin resultados

## 3. Justificación de Uso de Herramientas

### 3.1 Herramientas Utilizadas

| Herramienta | Por qué se eligió | Alternativas consideradas | Limitación conocida |
|-------------|------------------|---------------------------|---------------------|
| **Go Testing Framework** | Framework nativo, integrado en el lenguaje, sin dependencias externas, soporte de concurrencia | TestNG, JUnit, PyTest | Sintaxis más verbosa que frameworks de otros lenguajes |
| **Testify/Require** | Aserciones claras con mensajes descriptivos, integración perfecta con Go testing, amplia adopción en comunidad Go | Gomega, assert package nativo | Requiere dependencia externa |
| **TestHelper (custom)** | Abstracción reutilizable para setup/teardown, reduce duplicación, mantiene consistencia entre pruebas | Setup/teardown inline en cada prueba | Requiere mantenimiento del helper |
| **Client API (OctoClient)** | Simula calls reales de frontend, valida integración completa API-App-Store, usa misma infraestructura que producción | Mocks directos de Store, HTTP test server | Dependencia de infraestructura de servidor |

### 3.2 Decisiones de Arquitectura de Pruebas

**Enfoque de Integración vs Unit Testing:**
- Se eligió **integración completa** (API → App → Store) en lugar de mocking parcial
- **Justificación**: El flujo de búsqueda involucra múltiples capas que deben coordinarse correctamente; mocking podría ocultar errores de integración
- **Trade-off**: Pruebas más lentas pero mayor confianza en corrección del sistema

## 4. Justificación de Estrategia por Caso de Prueba

### 4.1 INT-09-01: Buscar tableros por título y verificar accesibilidad

**Estrategia:**
- Crear múltiples tableros con diferentes niveles de acceso (privados de user1, privados de user2, públicos)
- Ejecutar búsqueda como user1
- Verificar que resultados incluyen solo tableros accesibles a user1

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-01 y RP-03 (violación de permisos)
- **Escenario realista**: Simula caso de uso común donde usuarios buscan sus tableros
- **Validación de múltiples capas**: Verifica que API, App y Store coordinan correctamente el filtrado por permisos
- **Prioridad Alta**: Seguridad es crítica; cualquier fuga de datos es inaceptable

### 4.2 INT-09-02: Buscar término sin coincidencias

**Estrategia:**
- Crear tableros con títulos específicos
- Ejecutar búsqueda con término que no coincide con ningún título
- Verificar que retorna array vacío con HTTP 200 (no error)

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-05 (comportamiento incorrecto en casos edge)
- **Validación de contrato API**: Asegura que el endpoint maneja correctamente ausencia de resultados
- **Experiencia de usuario**: Array vacío es más manejable por frontend que error
- **Prioridad Media**: No es crítico para seguridad pero afecta UX

### 4.3 INT-09-03: Verificar respeto de permisos

**Estrategia:**
- Crear tableros privados para user1 y user2
- Crear tablero público
- Ejecutar búsquedas como ambos usuarios
- Verificar aislamiento: cada usuario solo ve sus privados + públicos, no privados ajenos

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-01 y RP-03 (violación de permisos)
- **Validación bidireccional**: Prueba desde perspectiva de ambos usuarios
- **Verificación de lógica de permisos**: Asegura que Store filtra correctamente según membresía
- **Prioridad Alta**: Seguridad y privacidad son fundamentales

### 4.4 INT-09-04: Caracteres especiales sin errores SQL

**Estrategia:**
- Crear tablero de contexto
- Ejecutar búsquedas con patrones de inyección SQL y caracteres especiales
- Verificar que todas retornan respuesta válida (200 OK) sin errores de servidor

**Justificación:**
- **Cobertura de riesgo**: Mitiga RP-02 (inyección SQL) y RP-04 (errores con caracteres especiales)
- **Validación de sanitización**: Asegura que Store sanitiza correctamente inputs
- **Robustez**: Sistema debe ser resistente a inputs maliciosos
- **Prioridad Media-Alta**: Seguridad crítica, pero probabilidad de ataque directo es baja

## 5. Integración con Pruebas Unitarias

Las pruebas de integración INT-09 complementan las pruebas unitarias existentes:

- **Pruebas unitarias**: Validan lógica individual de componentes (Store.SearchBoardsForUser, App.SearchBoardsForUser)
- **Pruebas de integración INT-09**: Validan coordinación entre componentes y flujo end-to-end
- **Sin redundancia**: Las pruebas de integración no duplican validaciones unitarias, sino que verifican integración

## 6. Mantenibilidad y Evolución

### 6.1 Diseño para Mantenimiento
- **TestHelper reutilizable**: Setup/teardown centralizado reduce duplicación
- **Nombres descriptivos**: Facilita identificación de propósito y fallas
- **Logging claro**: Ayuda en debugging cuando pruebas fallan

### 6.2 Escalabilidad
- **Patrón extensible**: Nuevos casos de prueba pueden seguir mismo patrón
- **Independencia**: Pruebas pueden ejecutarse en paralelo sin interferencia
- **Aislamiento de datos**: Cada prueba usa datos únicos para evitar colisiones

## 7. Conclusión

Las pruebas de integración INT-09 proporcionan cobertura robusta para el flujo de búsqueda, priorizando áreas de alto riesgo (seguridad, permisos, inyección SQL) mientras mantienen mantenibilidad y escalabilidad. La estrategia de integración completa (sin mocking) asegura confianza en corrección del sistema, y el diseño siguiendo estándares facilita mantenimiento y evolución futura.
