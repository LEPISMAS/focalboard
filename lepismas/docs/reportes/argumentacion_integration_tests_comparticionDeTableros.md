# Argumentación de Pruebas de Integración — INT-05: Compartición de Tableros

**Proyecto:** Focalboard (fork Lepismas)
**Flujo evaluado:** INT-05 — API ↔ Sharing ↔ Store
**Archivo de pruebas:** `server/integrationtests/flujo_comparticionDeTableros_int_tests.go`
**Fecha de elaboración:** 2026-07-09
**Autor:** Equipo Lepismas

---

## 1. Contexto y Motivación

La compartición pública de tableros es una funcionalidad crítica de Focalboard que permite exponer
el contenido de un tablero sin requerir autenticación por parte del visitante.
El flujo involucra tres capas de la arquitectura:

| Capa | Rol en el flujo |
|------|----------------|
| **API REST** (`server/api/boards.go`) | Recibe y valida las peticiones HTTP, gestiona el token de compartición |
| **App / Lógica de Negocio** (`server/app/sharing.go`) | Crea, persiste y revoca el estado de compartición |
| **Store / Persistencia** (`server/store/`) | Almacena y recupera el registro `Sharing` con su token |

Las pruebas de integración de este flujo verifican que **las tres capas se comunican correctamente
de extremo a extremo**, algo que las pruebas unitarias aisladas no pueden garantizar por sí solas.

---

## 2. Estrategia de Prueba

Se adoptó una estrategia de **caja negra sobre la interfaz HTTP**: los tests interactúan
exclusivamente con el cliente oficial (`server/client/client.go`) de la misma manera en que lo
haría un frontend real, sin acceder directamente a la base de datos ni a la lógica interna.

Esto asegura que:
- Los contratos de la API se cumplen desde la perspectiva del consumidor.
- Los errores de serialización, autorización o enrutamiento se detectan tempranamente.
- El entorno de prueba (SQLite con extensión json1) es reproducible en CI.

---

## 3. Casos de Prueba y su Justificación

### INT-05-01 — Habilitar compartición y verificar generación del token

**Objetivo:** Confirmar que al marcar un tablero como compartido, la API persiste un token
no vacío en el Store y lo devuelve en la respuesta.

**Justificación:**
El token es el mecanismo central de acceso público. Si el token no se genera o no se persiste,
todos los accesos públicos posteriores fallarían silenciosamente. Este caso verifica el "camino
feliz" de la habilitación y la coherencia entre lo devuelto por la API y lo almacenado.

**Pasos clave:**
1. Crear un tablero con usuario autenticado.
2. Llamar a `UpsertSharing` con `Enabled = true`.
3. Recuperar el objeto `Sharing` mediante `GetSharingBoard`.
4. Comprobar que `Enabled == true` y `Token != ""`.

---

### INT-05-02 — Acceder al tablero con token válido sin autenticación

**Objetivo:** Verificar que un visitante no autenticado puede obtener los datos del tablero
usando el token correcto en el parámetro `read_token`.

**Justificación:**
Esta es la funcionalidad principal desde la perspectiva del usuario final. Si el acceso anónimo
falla, la característica de compartición no tiene valor. Además, este caso detecta posibles
errores de middleware de CSRF o autenticación que puedan bloquear peticiones sin sesión.

**Pasos clave:**
1. Habilitar compartición en un tablero (obtener token).
2. Crear un cliente sin autenticación (`th.CreateAnon()`).
3. Llamar a `GetBoard` con el `read_token` como parámetro de query.
4. Confirmar respuesta HTTP 200 y que el ID del tablero recibido coincide.

---

### INT-05-03 — Acceder al tablero con token inválido

**Objetivo:** Verificar que la API rechaza el acceso cuando el token proporcionado no corresponde
al tablero solicitado.

**Justificación:**
La seguridad del mecanismo depende de que únicamente el token correcto conceda acceso. Este caso
es el contrapunto de INT-05-02 y protege contra regresiones que pudieran permitir acceso con
cualquier token o sin token.

**Pasos clave:**
1. Habilitar compartición en un tablero.
2. Intentar acceder con el token `"invalid-token-xyz"`.
3. Verificar que la respuesta es HTTP 404 (tablero no encontrado / acceso denegado).

---

### INT-05-04 — Deshabilitar compartición y verificar revocación

**Objetivo:** Confirmar que al deshabilitar la compartición, el acceso público queda
efectivamente revocado incluso usando el token que anteriormente era válido.

**Justificación:**
La revocación es tan importante como la habilitación. Un usuario puede querer retirar el
acceso público en cualquier momento. Este caso garantiza que la operación de deshabilitar no es
meramente cosmética sino que bloquea peticiones reales.

**Pasos clave:**
1. Habilitar compartición y guardar el token.
2. Llamar a `UpsertSharing` con `Enabled = false`.
3. Intentar acceder con el token original usando cliente anónimo.
4. Verificar que la respuesta es HTTP 404.

---

## 4. Consideraciones Técnicas

### Build tag requerido

```
go test -tags sqlite_json
```

La migración de base de datos versión 18 usa la función SQL `json_set`, que requiere la
extensión json1 de SQLite. Sin el tag `sqlite_json`, el proceso de migración falla y los
tests no pueden inicializar.

### Archivos involucrados

| Archivo | Rol |
|---------|-----|
| `flujo_comparticionDeTableros_int_tests.go` | Definición de los tests (con `//go:build ignore`) |
| `ejecutar_comparticionDeTableros.bat` | Script de ejecución para Windows |
| `ejecutar_comparticionDeTableros.sh` | Script de ejecución para Unix/macOS |

### ¿Por qué `//go:build ignore`?

Go solo compila como tests los archivos que terminan en `_test.go`. El nombre requerido por
el equipo termina en `_tests.go` (plural), por lo que se añade la directiva `//go:build ignore`
para evitar conflictos de compilación, y los tests se copian a un archivo `_test.go` temporal
antes de ejecutarlos.

---

## 5. Cobertura del Flujo

```
Cliente HTTP
    |
    v
API REST (/api/v2/boards/:id/sharing)
    |  <- INT-05-01, INT-05-04: UpsertSharing
    |  <- INT-05-02, INT-05-03: GetBoard + read_token
    v
App Layer (app/sharing.go)
    |
    v
Store (SQLite) — tabla: focalboard_sharing
```

Todos los casos atraviesan las tres capas de extremo a extremo, validando:
- Creacion y persistencia del token (INT-05-01)
- Acceso publico autorizado (INT-05-02)
- Rechazo de tokens invalidos (INT-05-03)
- Revocacion efectiva del acceso (INT-05-04)

---

## 6. Conclusión

Los cuatro casos de prueba del flujo INT-05 proporcionan cobertura completa del ciclo de vida
de la compartición pública de tableros. La estrategia de caja negra sobre la API asegura que
el comportamiento observable desde el exterior es correcto, protegiendo los contratos del
sistema frente a futuras refactorizaciones.
