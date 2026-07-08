# Reporte de ejecucion de pruebas de integracion: Gestion de Tableros

## 1. Objetivo

Reportar el estado de implementacion y ejecucion del flujo INT-02 Gestion de Tableros, orientado a validar la integracion entre API REST, capa App, Store, permisos, membresias, bloques y el tramo observable del flujo relacionado con notificaciones.

El reporte distingue entre:

- Implementacion validada: las pruebas fueron creadas con los helpers reales del proyecto y cubren los casos solicitados.
- Ejecucion condicionada por entorno: la ejecucion local no pudo completarse por dependencias del ambiente, no por una falla demostrada del flujo de tableros.

## 2. Fecha

8 de julio de 2026.

## 3. Entorno

- Proyecto: Focalboard.
- Sistema local: Windows con PowerShell.
- Modulo Go: `server`.
- Paquete de pruebas: `server/integrationtests`.
- Archivo de pruebas: `server/integrationtests/flujo_gestionDeTableros_int_test.go`.
- Scripts de ejecucion:
  - `server/integrationtests/ejecutar_gestionDeTableros.bat`
  - `server/integrationtests/ejecutar_gestionDeTableros.sh`

## 4. Comandos utilizados

Desde la carpeta del flujo:

```powershell
cd server/integrationtests
go test -v . -run TestINT02
```

Mediante script Windows:

```bat
server\integrationtests\ejecutar_gestionDeTableros.bat
```

Mediante Bash:

```bash
bash server/integrationtests/ejecutar_gestionDeTableros.sh
```

## 5. Casos implementados

| ID | Caso | Estado de implementacion |
|---|---|---|
| INT-02-01 | Crear tablero y verificar persistencia | Implementado |
| INT-02-02 | Crear tablero y verificar membresia admin | Implementado |
| INT-02-03 | Listar tableros filtrando por membresia | Implementado |
| INT-02-04 | Actualizar titulo y verificar persistencia | Implementado |
| INT-02-05 | Eliminar tablero y verificar soft delete | Implementado |
| INT-02-06 | Duplicar tablero y verificar bloques, propiedades y membresias | Implementado |
| INT-02-07 | Verificar tramo API-App-Store y documentar limitacion WebSocket end-to-end | Implementado |

## 6. Resultados obtenidos

### Implementacion validada

El flujo fue implementado en el paquete `integrationtests`, usando helpers existentes como `SetupTestHelper`, clientes HTTP reales y metodos publicos de App/Store. Cada prueba incluye salida con `t.Log` para explicar su utilidad y las capas recorridas.

La implementacion valida:

- Creacion de tableros por API.
- Persistencia del tablero y campos principales.
- Creacion automatica de membresia administradora para el creador.
- Filtrado de tableros privados por membresia.
- Actualizacion de titulo mediante PATCH.
- Eliminacion logica con verificacion de historial y ausencia en listados activos.
- Duplicacion de tablero con bloques, propiedades y membresia del creador.
- Reconocimiento de la limitacion de WebSocket end-to-end en el arnes actual de integracion.

### Ejecucion condicionada por el entorno

La ejecucion local quedo condicionada por una incidencia de infraestructura:

1. Acceso denegado a la cache global de Go en `AppData\Local\go-build`, lo que impidio completar el comando `go test` en el entorno Windows local.

Ademas, por antecedentes del flujo INT-01 en este mismo entorno, se identifico que las migraciones pueden depender de una build de SQLite con soporte JSON, en particular funciones como `json_set`. Esa condicion tambien debe controlarse para ejecutar de forma estable los tests de integracion del servidor.

## 7. Incidencias encontradas

| ID incidencia | Descripcion | Componente afectado | Severidad | Estado |
|---|---|---|---|---|
| ENV-01 | Acceso denegado a cache global de compilacion Go | Entorno local Go/Windows | Media | Abierto en entorno local |
| ENV-02 | Riesgo de incompatibilidad SQLite si no esta disponible soporte JSON requerido por migraciones | Entorno local SQLite / migraciones | Alta para ejecucion local | Condicionado a configuracion |

No se reportan defectos funcionales del flujo de tableros, porque la ejecucion no alcanzo un fallo atribuible a la logica del producto. Las incidencias detectadas corresponden al ambiente.

## 8. Conclusiones

El flujo INT-02 esta implementado de forma completa respecto a los casos solicitados. Las pruebas cubren operaciones de ciclo de vida de tablero y verifican estado persistido mediante App/Store, lo cual permite defender que no son pruebas unitarias sino pruebas de integracion entre capas.

La ejecucion local no puede presentarse como resultado funcional definitivo debido a restricciones del entorno. Por ello, este reporte documenta la implementacion validada y deja explicita la condicion ambiental pendiente para obtener evidencia de ejecucion completa.

## 9. Recomendaciones

- Ejecutar las pruebas en un entorno CI o maquina local con permisos correctos sobre la cache de Go.
- Preparar una base de datos de prueba compatible con las migraciones del proyecto.
- Si se requiere validar WebSocket end-to-end, crear posteriormente un arnes estable de integracion para suscripcion y lectura de eventos, en lugar de improvisar clientes fragiles.
- Completar un reporte de resultados finales cuando las pruebas se ejecuten en un entorno estable.
