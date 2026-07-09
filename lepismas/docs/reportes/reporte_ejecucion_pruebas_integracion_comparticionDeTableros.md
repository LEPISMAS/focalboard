# Reporte de Ejecución — INT-05: Compartición de Tableros

**Proyecto:** Focalboard (fork Lepismas)
**Flujo:** INT-05 — API ↔ Sharing ↔ Store
**Fecha de ejecución:** 2026-07-09
**Entorno:** Windows 11 / Go 1.21 / SQLite (json1)
**Comando ejecutado:**
```
go test -v . -run "TestINT05" -tags sqlite_json
```
**Directorio:** `server/integrationtests/`

---

## Resumen Ejecutivo

| Métrica | Valor |
|---------|-------|
| Total de casos | 4 |
| Casos PASS | 4 |
| Casos FAIL | 0 |
| Tiempo total | ~6.5 s |
| Resultado | **PASS** |

---

## Resultados por Caso

### INT-05-01 — TestINT0501HabilitarComparticion

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.5 s |

**Log de pasos:**
```
INT-05-01: Habilitar comparticion publica de un tablero y verificar token en Store
  Paso 1: Creando tablero de prueba...
  Paso 2: Habilitando comparticion publica del tablero...
  Paso 3: Verificando que la comparticion esta habilitada y tiene token...
  INT-05-01 COMPLETADO: Token generado correctamente
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `UpsertSharing`
- Campo `Enabled == true` en el objeto `Sharing` recuperado del Store
- Campo `Token` no vacío

---

### INT-05-02 — TestINT0502AccederTableroCompartidoValido

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.5 s |

**Log de pasos:**
```
INT-05-02: Acceder a tablero compartido con token valido sin autenticacion
  Paso 1: Creando tablero y habilitando comparticion...
  Paso 2: Accediendo al tablero con token valido sin autenticacion...
  INT-05-02 COMPLETADO: Acceso anonimo con token valido exitoso
```

**Verificaciones superadas:**
- Respuesta HTTP 200 desde cliente anónimo con `read_token` correcto
- ID del tablero recibido coincide con el tablero compartido

---

### INT-05-03 — TestINT0503AccederTableroCompartidoInvalido

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.4 s |

**Log de pasos:**
```
INT-05-03: Acceder a tablero compartido con token invalido debe ser rechazado
  Paso 1: Creando tablero y habilitando comparticion...
  Paso 2: Intentando acceder con token invalido...
  INT-05-03 COMPLETADO: Token invalido correctamente rechazado con 404
```

**Verificaciones superadas:**
- Respuesta HTTP 404 al usar `read_token = "invalid-token-xyz"`
- El sistema no filtra información del tablero en el error

---

### INT-05-04 — TestINT0504DeshabilitarComparticion

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.7 s |

**Log de pasos:**
```
INT-05-04: Deshabilitar comparticion y verificar que el acceso publico queda revocado
  Paso 1: Creando tablero y habilitando comparticion...
  Paso 2: Deshabilitando comparticion...
  Paso 3: Verificando que el acceso con el token anterior es denegado...
  INT-05-04 COMPLETADO: Acceso correctamente revocado tras deshabilitar
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en el `UpsertSharing` de deshabilitación
- Respuesta HTTP 404 al intentar acceder con el token previamente válido

---

## Salida Completa del Test Runner (extracto relevante)

```
=== RUN   TestINT0501HabilitarComparticion
--- PASS: TestINT0501HabilitarComparticion (0.51s)
=== RUN   TestINT0502AccederTableroCompartidoValido
--- PASS: TestINT0502AccederTableroCompartidoValido (0.52s)
=== RUN   TestINT0503AccederTableroCompartidoInvalido
--- PASS: TestINT0503AccederTableroCompartidoInvalido (0.43s)
=== RUN   TestINT0504DeshabilitarComparticion
--- PASS: TestINT0504DeshabilitarComparticion (0.74s)
PASS
ok      github.com/mattermost/focalboard/server/integrationtests        6.479s
```

---

## Observaciones

- Todos los tests inicializan su propio servidor SQLite en memoria mediante `SetupTestHelper`.
- El warning `Unable to serve the index.html file` es esperado en el entorno de tests (no hay frontend compilado) y no afecta el resultado.
- El status code 500 en el primer ping del servidor es el comportamiento esperado del health-check cuando el frontend no está presente; los endpoints de API responden correctamente.

---

## Conclusión

El flujo INT-05 pasa íntegramente. La funcionalidad de compartición pública de tableros —
habilitación, acceso autorizado, rechazo de acceso inválido y revocación — opera correctamente
a través de todas las capas del sistema (API → App → Store).
