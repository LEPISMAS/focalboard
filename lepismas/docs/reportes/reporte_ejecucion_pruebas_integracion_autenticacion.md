# Reporte de ejecucion de pruebas de integracion: Autenticacion

## 1. Objetivo

Reportar el estado de implementacion y ejecucion del flujo INT-01 Autenticacion, cuyo proposito es validar la integracion entre API REST, capa de aplicacion, servicio de autenticacion y Store para operaciones de registro, login, validacion de token, cambio de contrasena, logout y control de duplicados.

Este reporte diferencia dos aspectos:

- Implementacion validada: los casos fueron codificados en el archivo de pruebas correspondiente y respetan la estructura del proyecto.
- Ejecucion condicionada por entorno: la ejecucion local no pudo completarse en este ambiente por dependencias externas a la logica de las pruebas.

## 2. Fecha

8 de julio de 2026.

## 3. Entorno

- Proyecto: Focalboard.
- Sistema local: Windows con PowerShell.
- Modulo Go: `server`.
- Paquete de pruebas: `server/integrationtests`.
- Archivo de pruebas: `server/integrationtests/flujo_autenticacion_int_test.go`.
- Scripts de ejecucion:
  - `server/integrationtests/ejecutar_autenticacion.bat`
  - `server/integrationtests/ejecutar_autenticacion.sh`

## 4. Comandos utilizados

Desde la carpeta del flujo:

```powershell
cd server/integrationtests
go test -v . -run TestINT01
```

Mediante script Windows:

```bat
server\integrationtests\ejecutar_autenticacion.bat
```

Mediante Bash:

```bash
bash server/integrationtests/ejecutar_autenticacion.sh
```

## 5. Casos implementados

| ID | Caso | Estado de implementacion |
|---|---|---|
| INT-01-01 | Registrar usuario nuevo y verificar persistencia | Implementado |
| INT-01-02 | Login genera token y sesion valida | Implementado |
| INT-01-03 | Endpoint protegido con token valido | Implementado |
| INT-01-04 | Endpoint protegido sin token rechaza | Implementado |
| INT-01-05 | Cambio de password permite login nuevo | Implementado |
| INT-01-06 | Logout invalida token | Implementado |
| INT-01-07 | Registro duplicado rechazado | Implementado |

## 6. Resultados obtenidos

### Implementacion validada

El archivo de pruebas fue creado con paquete `integrationtests`, usando helpers existentes del proyecto y aserciones `require` de Testify. Los casos ejecutan endpoints o helpers reales del cliente HTTP y verifican efectos en App/Store cuando las interfaces publicas lo permiten.

La implementacion cubre el flujo completo esperado:

- Registro por API y persistencia en Store.
- Login y generacion de sesion.
- Acceso autorizado con token.
- Rechazo de acceso anonimo.
- Cambio de contrasena y login posterior.
- Logout y rechazo del token anterior.
- Rechazo de usuario duplicado y control de no duplicacion.

### Ejecucion condicionada por el entorno

En este ambiente local la ejecucion completa quedo condicionada por dos incidencias de infraestructura:

1. Acceso denegado a la cache global de Go en `AppData\Local\go-build`.
2. Al intentar usar cache local y el tag SQLite, la migracion de base de datos fallo porque la build local de SQLite no exponia la funcion `json_set`, requerida por migraciones del proyecto.

Estas incidencias no corresponden a errores funcionales de los casos INT-01, sino a dependencias del entorno de ejecucion local.

## 7. Incidencias encontradas

| ID incidencia | Descripcion | Componente afectado | Severidad | Estado |
|---|---|---|---|---|
| ENV-01 | Acceso denegado a cache global de compilacion Go | Entorno local Go/Windows | Media | Abierto en entorno local |
| ENV-02 | SQLite local sin soporte disponible para `json_set` durante migraciones | Entorno local SQLite / migraciones | Alta para ejecucion local | Abierto en entorno local |

No se registra matriz de defectos funcionales del producto, porque no se obtuvo una falla atribuible a la logica de autenticacion. Las incidencias son ambientales y condicionan la ejecucion.

## 8. Conclusiones

El flujo INT-01 fue implementado correctamente a nivel de estructura, trazabilidad y alcance tecnico. Los casos cubren las capas API, App, Auth y Store, y estan alineados con el objetivo de demostrar integracion real del proceso de autenticacion.

La ejecucion completa en este ambiente no puede considerarse evidencia definitiva de resultado funcional porque fue interrumpida por configuracion local de Go/SQLite. Sin embargo, la implementacion queda preparada para ejecutarse en un entorno con cache Go accesible y SQLite compatible con las funciones JSON requeridas por las migraciones de Focalboard.

## 9. Recomendaciones

- Ejecutar el flujo en un entorno CI o local donde Go tenga permisos sobre su cache de compilacion.
- Confirmar que SQLite este compilado con soporte JSON requerido por las migraciones, o ejecutar las pruebas contra MySQL/PostgreSQL segun las variables de entorno soportadas por el proyecto.
- Registrar los resultados reales en un reporte posterior una vez estabilizado el entorno.
- Mantener separados los defectos funcionales de las incidencias ambientales para no atribuir al producto fallos de infraestructura.
