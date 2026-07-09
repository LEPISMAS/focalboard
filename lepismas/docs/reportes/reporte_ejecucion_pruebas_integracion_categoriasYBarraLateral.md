# Reporte de Ejecución — INT-06: Categorías y Barra Lateral

**Proyecto:** Focalboard (fork Lepismas)
**Flujo:** INT-06 — API ↔ Categories ↔ Store
**Fecha de ejecución:** 2026-07-09
**Entorno:** Windows 11 / Go 1.21 / SQLite (json1)
**Comando ejecutado:**
```
go test -v . -run "TestINT06" -tags sqlite_json
```
**Directorio:** `server/integrationtests/`

---

## Resumen Ejecutivo

| Métrica | Valor |
|---------|-------|
| Total de casos | 5 |
| Casos PASS | 5 |
| Casos FAIL | 0 |
| Tiempo total | ~4.5 s |
| Resultado | **PASS** |

---

## Resultados por Caso

### INT-06-01 — TestINT0601CrearCategoriaYVerificarPersistencia

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.6 s |

**Log de pasos:**
```
INT-06-01: Crear categoria via API y verificar persistencia en Store
  Paso 1: Creando categoria de prueba...
  Paso 2: Verificando que la categoria aparece en GetUserCategoryBoards...
  INT-06-01 COMPLETADO: Categoria creada y persistida correctamente
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `CreateCategory`
- La categoría creada aparece en la lista devuelta por `GetUserCategoryBoards`
- El ID de la categoría coincide

---

### INT-06-02 — TestINT0602MoverTableroACategoriaPersonalizada

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.6 s |

**Log de pasos:**
```
INT-06-02: Mover tablero a categoria personalizada y verificar asociacion en Store
  Paso 1: Creando tablero y categoria personalizada...
  Paso 2: Moviendo tablero a la categoria personalizada...
  Paso 3: Verificando que el tablero aparece en la categoria correcta...
  INT-06-02 COMPLETADO: Tablero movido correctamente a categoria personalizada
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `AddBoardToCategory`
- El tablero aparece en los `CategoryBoards` de la categoría destino

---

### INT-06-03 — TestINT0603ObtenerCategoriasConTableros

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.6 s |

**Log de pasos:**
```
INT-06-03: Obtener categorias de barra lateral y verificar que incluyen tableros asociados
  Paso 1: Creando tablero y categoria personalizada...
  Paso 2: Moviendo tablero a la categoria personalizada...
  Paso 3: Obteniendo categorias de la barra lateral...
  INT-06-03 COMPLETADO: Barra lateral devuelve categorias con tableros
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `GetUserCategoryBoards`
- La respuesta contiene al menos una categoría con al menos un tablero

---

### INT-06-04 — TestINT0604EliminarCategoriaYVerificarTablerosVuelvenADefault

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.7 s |

**Log de pasos:**
```
INT-06-04: Eliminar categoria y verificar que los tableros vuelven a la categoria por defecto
  Paso 1: Creando tablero y categoria personalizada...
  Paso 2: Moviendo tablero a la categoria personalizada...
  Paso 3: Eliminando la categoria personalizada...
  Paso 4: Verificando que el tablero ahora aparece en la categoria por defecto...
  INT-06-04 COMPLETADO: Tableros volvieron correctamente a la categoria por defecto
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `DeleteCategory`
- El tablero aparece en la categoría del sistema (tipo `system`) tras la eliminación

---

### INT-06-05 — TestINT0605ReordenarCategoriasPersistencia

| Campo | Valor |
|-------|-------|
| Estado | **PASS** |
| Duración aprox. | ~0.6 s |

**Log de pasos:**
```
INT-06-05: Reordenar categorias y verificar que el nuevo orden se persista
  Paso 1: Creando dos categorias personalizadas...
  Paso 2: Construyendo lista completa de IDs (custom + sistema)...
  Paso 3: Reordenando categorias (cat2 antes que cat1)...
  Paso 4: Verificando persistencia del nuevo orden...
  INT-06-05 COMPLETADO: Nuevo orden de categorias persistido correctamente
```

**Verificaciones superadas:**
- Respuesta HTTP 200 en `ReorderCategories` (lista incluye TODOS los IDs)
- En `GetUserCategoryBoards` posterior, `indexCat2 < indexCat1` confirma el orden esperado

**Nota de corrección aplicada:**
El test originalmente pasaba solo los IDs de las categorías personalizadas. El backend requiere
que la lista contenga TODOS los IDs (incluyendo la categoría del sistema "Boards"). Se corrigió
recopilando primero todos los IDs de `GetUserCategoryBoards` y construyendo la lista completa
antes de llamar a `ReorderCategories`.

---

## Salida Completa del Test Runner (extracto relevante)

```
=== RUN   TestINT0601CrearCategoriaYVerificarPersistencia
--- PASS: TestINT0601CrearCategoriaYVerificarPersistencia (0.62s)
=== RUN   TestINT0602MoverTableroACategoriaPersonalizada
--- PASS: TestINT0602MoverTableroACategoriaPersonalizada (0.64s)
=== RUN   TestINT0603ObtenerCategoriasConTableros
--- PASS: TestINT0603ObtenerCategoriasConTableros (0.60s)
=== RUN   TestINT0604EliminarCategoriaYVerificarTablerosVuelvenADefault
--- PASS: TestINT0604EliminarCategoriaYVerificarTablerosVuelvenADefault (0.72s)
=== RUN   TestINT0605ReordenarCategoriasPersistencia
--- PASS: TestINT0605ReordenarCategoriasPersistencia (0.61s)
PASS
ok      github.com/mattermost/focalboard/server/integrationtests        4.499s
```

---

## Observaciones

- El warning `Unable to serve the index.html file` y el status 500 en el ping inicial son
  comportamientos esperados en el entorno de tests sin frontend compilado.
- El fix en INT-06-05 fue necesario porque la documentación interna del endpoint no especifica
  explícitamente que se deben incluir TODOS los IDs. Se documentó en la argumentación para
  referencia futura.

---

## Conclusión

El flujo INT-06 pasa íntegramente con los 5 casos de prueba. La gestión de categorías —
creación, asignación de tableros, consulta de barra lateral, eliminación con integridad
referencial y reordenamiento persistente — opera correctamente a través de todas las capas
del sistema (API → App → Store).
