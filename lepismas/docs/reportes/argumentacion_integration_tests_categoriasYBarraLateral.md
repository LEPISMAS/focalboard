# Argumentación de Pruebas de Integración — INT-06: Categorías y Barra Lateral

**Proyecto:** Focalboard (fork Lepismas)
**Flujo evaluado:** INT-06 — API ↔ Categories ↔ Store
**Archivo de pruebas:** `server/integrationtests/flujo_categoriasYBarraLateral_int_tests.go`
**Fecha de elaboración:** 2026-07-09
**Autor:** Equipo Lepismas

---

## 1. Contexto y Motivación

La barra lateral de Focalboard organiza los tableros del usuario en categorías personalizables.
Esta funcionalidad involucra un flujo completo entre la API REST, la capa de lógica de negocio
y el almacenamiento persistente.

El flujo cubre las siguientes operaciones:
- Crear y eliminar categorías personalizadas.
- Mover tableros entre categorías.
- Consultar la barra lateral con su organización actual.
- Reordenar categorías de forma persistente.

Las pruebas de integración validan que estas operaciones son coherentes a través de todas las
capas del sistema, no solo en los componentes individuales.

---

## 2. Estrategia de Prueba

Al igual que en INT-05, se utiliza una estrategia de **caja negra sobre HTTP** mediante el
cliente oficial. Las pruebas se encadenan lógicamente: cada test establece el estado necesario
para verificar un aspecto específico de la gestión de categorías.

Cada test es independiente (usa `SetupTestHelper().InitBasic()` y `defer th.TearDown()`) para
garantizar aislamiento entre casos.

---

## 3. Casos de Prueba y su Justificación

### INT-06-01 — Crear categoría y verificar persistencia en el Store

**Objetivo:** Verificar que al crear una categoría vía API, esta queda persistida y es
recuperable mediante una consulta posterior.

**Justificación:**
Es el caso base del flujo. Sin una creación funcional de categorías, ninguno de los casos
restantes puede ejecutarse. Confirma que la serialización del modelo `Category`, el enrutamiento
y la escritura al Store funcionan correctamente de extremo a extremo.

**Verificación:** Se llama a `GetUserCategoryBoards` tras la creación y se busca la categoría
por ID en la lista devuelta.

---

### INT-06-02 — Mover tablero a categoría personalizada

**Objetivo:** Verificar que al reasignar un tablero a una categoría personalizada, la asociación
se refleja correctamente en el Store.

**Justificación:**
La organización de tableros en categorías es el valor central de la barra lateral. Este caso
garantiza que la operación de movimiento (`AddBoardToCategory`) actualiza la persistencia
correctamente y que la consulta posterior (`GetUserCategoryBoards`) refleja el cambio.

**Verificación:** Se itera sobre los tableros de la categoría destino y se confirma que el
tablero recién movido aparece en ella.

---

### INT-06-03 — Obtener categorías de la barra lateral con tableros asociados

**Objetivo:** Verificar que la respuesta de la barra lateral incluye correctamente las
categorías con sus tableros asociados.

**Justificación:**
La consulta de la barra lateral es el punto de entrada principal del frontend para renderizar
la navegación. Un error aquí afecta directamente la experiencia del usuario. Este test valida
que la agregación de categorías + tableros funciona correctamente en la respuesta de la API.

**Verificación:** Se comprueba que existe al menos una categoría con al menos un tablero
asociado en la respuesta.

---

### INT-06-04 — Eliminar categoría y verificar que los tableros vuelven a la categoría por defecto

**Objetivo:** Confirmar que al eliminar una categoría personalizada, los tableros que contenía
son reasignados automáticamente a la categoría del sistema ("Boards").

**Justificación:**
La eliminación de una categoría no debe provocar pérdida de tableros. Este caso garantiza la
integridad referencial del sistema: los tableros huérfanos son recuperados en la categoría
predeterminada. Es un caso crítico de integridad de datos.

**Verificación:** Tras eliminar la categoría, se consultan los tableros de la categoría del
sistema y se verifica que el tablero movido ahora aparece allí.

---

### INT-06-05 — Reordenar categorías y verificar persistencia del nuevo orden

**Objetivo:** Verificar que al reordenar las categorías, el nuevo orden queda efectivamente
persistido en el Store.

**Justificación:**
El reorden es una funcionalidad de UX importante para usuarios con muchas categorías. Este
test garantiza que la operación `ReorderCategories` no solo devuelve éxito, sino que el orden
es efectivamente almacenado y recuperado en consultas posteriores.

**Nota técnica importante:** La API exige que la lista pasada a `ReorderCategories` incluya
TODOS los IDs de categorías del usuario, incluyendo la categoría del sistema ("Boards"). Si se
omite alguna, el backend devuelve error 500 con el mensaje:
`cannot update category order, passed list of categories different size than in database`.

Por esta razón, el test recopila todos los IDs de `GetUserCategoryBoards` y construye la
lista completa antes de llamar a `ReorderCategories`.

**Verificación:** Se consultan las categorías tras el reorden y se verifica que la categoría
que debería ir primero tiene un índice menor que la que debería ir después.

---

## 4. Consideraciones Técnicas

### Build tag requerido

```
go test -tags sqlite_json
```

Mismo requisito que INT-05: la extensión json1 de SQLite es obligatoria.

### Archivos involucrados

| Archivo | Rol |
|---------|-----|
| `flujo_categoriasYBarraLateral_int_tests.go` | Definición de los tests (con `//go:build ignore`) |
| `ejecutar_categoriasYBarraLateral.bat` | Script de ejecución para Windows |
| `ejecutar_categoriasYBarraLateral.sh` | Script de ejecución para Unix/macOS |

---

## 5. Cobertura del Flujo

```
Cliente HTTP
    |
    v
API REST (/api/v2/teams/:teamID/categories)
    |  <- INT-06-01: CreateCategory
    |  <- INT-06-02: AddBoardToCategory
    |  <- INT-06-03: GetUserCategoryBoards
    |  <- INT-06-04: DeleteCategory
    |  <- INT-06-05: ReorderCategories
    v
App Layer (app/categories.go)
    |
    v
Store (SQLite) — tablas: focalboard_categories, focalboard_category_boards
```

Todos los casos atraviesan las tres capas de extremo a extremo, validando:
- Persistencia de creacion de categorias (INT-06-01)
- Asociacion correcta tablero-categoria (INT-06-02)
- Consulta agregada de barra lateral (INT-06-03)
- Integridad referencial en eliminacion (INT-06-04)
- Persistencia del reorden (INT-06-05)

---

## 6. Conclusión

Los cinco casos de prueba del flujo INT-06 cubren el ciclo de vida completo de la gestión
de categorías en la barra lateral. La combinación de operaciones CRUD y de reorden, junto
con la verificación de integridad referencial en la eliminación, proporciona una cobertura
robusta de los escenarios más críticos para el usuario.
