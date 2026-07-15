# Reporte de ejecucion de pruebas de integracion: Importacion y Exportacion

## 1. Objetivo

Reportar el estado de implementacion y ejecucion del flujo INT-07 Importacion y Exportacion, orientado a validar la integracion entre API REST, capa App, sistema de archivos y Store, para las operaciones de exportar, importar, roundtrip, rechazo de archivos invalidos y adjuntos.

## 2. Fecha

15 de julio de 2026.

## 3. Entorno

- Proyecto: Focalboard.
- Modulo Go: `server`.
- Paquete de pruebas: `server/integrationtests`.
- Archivo de pruebas: `server/integrationtests/flujo_importacionYExportacion_int_tests.go`.
- Scripts de ejecucion:
  - `server/integrationtests/ejecutar_importacionYExportacion.bat`
  - `server/integrationtests/ejecutar_importacionYExportacion.sh`

## 4. Comandos utilizados

Desde la carpeta del flujo:

```bash
cd server/integrationtests
go test -v . -run TestINT07
```

Mediante script Windows:

```bat
server\integrationtests\ejecutar_importacionYExportacion.bat
```

Mediante Bash:

```bash
bash server/integrationtests/ejecutar_importacionYExportacion.sh
```

## 5. Casos implementados

| ID | Caso | Estado de implementacion |
|---|---|---|
| INT-07-01 | Exportar tablero y verificar formato de archivo generado | Implementado |
| INT-07-02 | Importar archivo y verificar recreacion de entidades en Store | Implementado |
| INT-07-03 | Exportar e importar (roundtrip) y verificar integridad de datos | Implementado |
| INT-07-04 | Importar archivo invalido y verificar rechazo sin datos parciales | Implementado |
| INT-07-05 | Subir archivo adjunto y verificar que sea recuperable | Implementado |

## 6. Resultados obtenidos

### Implementacion validada

El flujo fue implementado en el paquete `integrationtests`, usando helpers existentes como `SetupTestHelper`, clientes HTTP reales y metodos publicos de App/Store. Cada prueba incluye salida con `t.Log` para explicar su utilidad y las capas recorridas.

La implementacion valida:

- Generacion de un archivo de exportacion con formato de cabecera valido.
- Recreacion de tablero y bloques tras una importacion, con nuevos IDs.
- Preservacion de contenido (titulo de tablero, cantidad y titulos de bloques) en un ciclo completo de exportar e importar.
- Rechazo de un archivo con formato invalido, sin crear tableros parciales.
- Subida de un archivo adjunto a una tarjeta y recuperacion de su informacion via API.

### Ejecucion en el entorno de trabajo

*(Completar con el resultado real al ejecutar `go test -v . -run TestINT07` en el entorno del equipo: PASS/FAIL por caso, tiempo total y cualquier log relevante de la corrida.)*

## 7. Incidencias encontradas

| ID Defecto | Caso de prueba que lo detecto | Componente afectado | Severidad | Requisito relacionado | Estado |
|---|---|---|---|---|---|
| *(pendiente de completar tras primera ejecucion real)* | | | | | |

Si al ejecutar las pruebas en el entorno del equipo no se detectan fallas, esta seccion debe indicarse explicitamente como "Sin defectos detectados" en lugar de dejarse vacia.

## 8. Conclusiones

El flujo INT-07 esta implementado de forma completa respecto a los casos solicitados por el issue de la tarea. Las pruebas cubren el ciclo de vida de exportacion e importacion de tableros y verifican estado persistido mediante App/Store, lo que permite defender que no son pruebas unitarias sino pruebas de integracion entre capas.

## 9. Recomendaciones

- Ejecutar `go test -v . -run TestINT07` en el entorno de CI o en una maquina local con Go y la base de datos de pruebas configurada, y completar las secciones 6 y 7 con el resultado real.
- Si se detectan defectos, documentarlos en la tabla de trazabilidad de la seccion 7 antes de la exposicion.
- Mantener sincronizado este reporte cada vez que se agreguen nuevos casos al flujo INT-07.
