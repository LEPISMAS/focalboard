# Explicación: Pruebas de Integración INT-05 e INT-06

Hola Gabriela! Este documento explica qué se hizo, por qué y cómo funciona cada parte de
las pruebas de integración que implementamos para los flujos de Compartición de Tableros (INT-05)
y Categorías / Barra Lateral (INT-06).

---

## ¿Qué son las pruebas de integración y por qué son importantes?

Las **pruebas unitarias** verifican piezas pequeñas de código de forma aislada. Son rápidas
pero no detectan problemas de comunicación entre componentes.

Las **pruebas de integración** verifican que varios componentes funcionan correctamente
*juntos*. En Focalboard, eso significa que:

1. La **API REST** recibe la petición correctamente,
2. La **lógica de negocio (App layer)** la procesa bien, y
3. El **Store (base de datos)** guarda y recupera los datos correctamente.

Si cualquiera de esas tres capas falla o se comunica mal, la prueba lo detecta.

---

## Flujo INT-05: Compartición de Tableros

### ¿Qué verifica?

Cuando un usuario marca un tablero como "público", Focalboard genera un **token secreto**.
Cualquiera que tenga ese token puede ver el tablero *sin* estar logueado.

Los 4 tests verifican:

| Test | ¿Qué hace? | ¿Por qué importa? |
|------|-----------|------------------|
| INT-05-01 | Habilita la compartición y verifica que el token se guardó | Sin token no hay acceso público posible |
| INT-05-02 | Accede al tablero con el token correcto, sin login | Es la experiencia del usuario final que recibe el link |
| INT-05-03 | Intenta acceder con un token inventado | Verifica que la seguridad funciona |
| INT-05-04 | Deshabilita la compartición y verifica que el token ya no funciona | Garantiza que el usuario puede revocar el acceso |

### ¿Cómo funciona el código?

```go
// INT-05-01: Habilitar compartición
th.Client.UpsertSharing(model.Sharing{
    ID:      board.ID,
    Enabled: true,
    Token:   utils.NewID(utils.IDTypeNone),
})
// Luego verificamos que el Store lo guardó:
sharing, _ := th.Client.GetSharingBoard(board.ID)
require.True(t, sharing.Enabled)
require.NotEmpty(t, sharing.Token)
```

```go
// INT-05-02: Acceso anónimo con token válido
anonClient := th.CreateAnon()
board, resp := anonClient.GetBoard(board.ID, sharing.Token)
th.CheckOK(resp)  // debe ser HTTP 200
```

---

## Flujo INT-06: Categorías y Barra Lateral

### ¿Qué verifica?

La barra lateral de Focalboard agrupa los tableros en categorías. Los usuarios pueden crear
sus propias categorías, mover tableros entre ellas, y cambiar el orden.

Los 5 tests verifican:

| Test | ¿Qué hace? | ¿Por qué importa? |
|------|-----------|------------------|
| INT-06-01 | Crea una categoría y verifica que quedó guardada | Base de todo lo demás |
| INT-06-02 | Mueve un tablero a la categoría creada | La organización de tableros es la función principal |
| INT-06-03 | Consulta la barra lateral y verifica que tiene tableros | Valida la vista del usuario |
| INT-06-04 | Elimina la categoría y verifica que los tableros van a "Boards" | Integridad de datos: los tableros no se pierden |
| INT-06-05 | Reordena las categorías y verifica que el orden persiste | El reorden debe sobrevivir a la siguiente carga |

### Un problema que encontramos y cómo lo resolvimos

En INT-06-05, la primera versión del test fallaba con este error:

```
cannot update category order, passed list of categories different size than in database
length new categories: 2, length existing categories: 3
```

**¿Por qué?** Focalboard siempre tiene una categoría del sistema llamada "Boards" que no
se puede eliminar. Cuando llamamos a `ReorderCategories`, el backend exige que le pasemos
*todos* los IDs — no solo los de las categorías que creamos nosotros.

**La solución:**
```go
// Primero obtenemos TODOS los IDs que ya existen
initialCats, _ := th.Client.GetUserCategoryBoards(testTeamID)

// Construimos la lista completa: primero nuestras categorías en el orden deseado,
// y al final las del sistema que no son nuestras
customIDs := map[string]bool{cat1.ID: true, cat2.ID: true}
newOrder := []string{cat2.ID, cat1.ID}  // cat2 antes que cat1
for _, cb := range initialCats {
    if !customIDs[cb.ID] {
        newOrder = append(newOrder, cb.ID)  // añadir las del sistema
    }
}
```

---

## Estructura de archivos

```
server/integrationtests/
├── flujo_comparticionDeTableros_int_tests.go   <- Tests INT-05 (fuente)
├── flujo_categoriasYBarraLateral_int_tests.go  <- Tests INT-06 (fuente)
├── ejecutar_comparticionDeTableros.bat         <- Runner Windows INT-05
├── ejecutar_comparticionDeTableros.sh          <- Runner Linux/Mac INT-05
├── ejecutar_categoriasYBarraLateral.bat        <- Runner Windows INT-06
└── ejecutar_categoriasYBarraLateral.sh         <- Runner Linux/Mac INT-06

lepismas/docs/reportes/
├── argumentacion_integration_tests_comparticionDeTableros.md
├── argumentacion_integration_tests_categoriasYBarraLateral.md
├── reporte_ejecucion_pruebas_integracion_comparticionDeTableros.md
└── reporte_ejecucion_pruebas_integracion_categoriasYBarraLateral.md
```

---

## ¿Por qué los archivos se llaman `_int_tests.go` y no `_int_test.go`?

Go tiene una regla estricta: **solo compila como tests los archivos que terminan en `_test.go`**.

El nombre que se pidió termina en `_tests.go` (con "s" al final). Para evitar que Go se confunda,
añadimos esta línea al inicio de esos archivos:

```go
//go:build ignore
```

Eso le dice a Go "ignora este archivo durante la compilación normal". Antes de ejecutar los
tests, el script runner hace una copia del archivo con el nombre correcto (`_test.go`), ejecuta
los tests, y listo.

---

## Cómo ejecutar las pruebas

### En Windows:
```batch
cd server\integrationtests
ejecutar_comparticionDeTableros.bat
ejecutar_categoriasYBarraLateral.bat
```

### En Linux/Mac:
```bash
cd server/integrationtests
bash ejecutar_comparticionDeTableros.sh
bash ejecutar_categoriasYBarraLateral.sh
```

### Directamente con Go:
```bash
# Copiar primero el archivo fuente al nombre que Go acepta
cat flujo_comparticionDeTableros_int_tests.go | tail -n +2 > flujo_comparticionDeTableros_int_test.go
go test -v . -run "TestINT05" -tags sqlite_json
```

> **Importante:** Siempre usar `-tags sqlite_json`. Sin ese tag, la base de datos SQLite
> no tiene la extensión JSON necesaria y los tests fallan al inicializar.

---

## Resultados finales

| Flujo | Tests | Resultado |
|-------|-------|-----------|
| INT-05: Compartición de Tableros | 4/4 | ✅ PASS |
| INT-06: Categorías y Barra Lateral | 5/5 | ✅ PASS |
| **Total** | **9/9** | **✅ PASS** |

---

¡Cualquier duda me avisas! 🙌
