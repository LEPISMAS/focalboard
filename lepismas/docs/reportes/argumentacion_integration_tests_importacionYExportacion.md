# Argumentacion de pruebas de integracion: Importacion y Exportacion

## 1. Objetivo del flujo

El flujo INT-07 tiene como objetivo comprobar que el ciclo de exportacion e importacion de tableros (`.boardarchive`) funciona de forma integrada entre la API REST, la capa de negocio, el sistema de archivos y el Store. Este flujo es critico porque es el unico mecanismo con el que un usuario puede respaldar, migrar o compartir un tablero completo fuera de la base de datos activa.

El objetivo principal es demostrar que exportar, importar, hacer un ciclo completo (roundtrip), rechazar archivos invalidos y adjuntar archivos a una tarjeta se comportan de forma correcta y segura cuando atraviesan las capas reales del backend.

## 2. Capas integradas

Las pruebas integran las siguientes capas:

- API REST: endpoints de exportacion (`ExportBoardArchive`), importacion (`ImportArchive`) y adjuntos (`/files`).
- Capa App: `server/app/export.go` y `server/app/import.go`, encargadas de serializar/deserializar tableros, bloques y metadatos.
- Store: persistencia de los tableros y bloques recreados tras la importacion.
- FileSystem: almacenamiento de los archivos adjuntos subidos a una tarjeta.
- Modelo: estructuras `Board`, `Block`, `BoardsAndBlocks`.

## 3. Estandares aplicados

- Trazabilidad por identificador: cada funcion conserva el codigo INT-07-01 a INT-07-05, igual que en los flujos previos del equipo.
- Reutilizacion del arnes oficial: se emplean `SetupTestHelper`, clientes reales (`th.Client`) y metodos publicos ya usados en el repositorio (`ExportBoardArchive`, `ImportArchive`).
- Pruebas contra API real: la exportacion, importacion y subida de archivos se ejecutan por HTTP, no por llamadas directas a funciones internas.
- Verificacion de estado persistido: se confirma el resultado final consultando `App`/`Store`, no solo el codigo de respuesta HTTP.
- Prueba de rutas positivas y negativas: se valida tanto el camino feliz (exportar, importar, roundtrip, adjuntar) como el camino de error (archivo invalido).
- Independencia de datos entre pruebas: cada test crea su propio tablero con ID generado por `utils.NewID`, evitando colisiones entre casos.

## 4. Importancia desde risk-based testing

### Identificacion de riesgos

Riesgos de producto:

- Un archivo exportado con formato incorrecto o incompleto, inutilizable para respaldo o migracion.
- Perdida o duplicacion de bloques al importar (integridad de datos).
- Aceptar un archivo corrupto y dejar el sistema en un estado parcialmente creado (tableros huerfanos).
- Archivos adjuntos que se suben pero no quedan referenciados correctamente en el Store, volviendose inaccesibles.

Riesgos de proyecto:

- Cambios en el formato del archivo `.boardarchive` que rompan la compatibilidad entre versiones sin que las pruebas lo detecten.
- Falta de pruebas automatizadas para importacion/exportacion, dejando esta funcionalidad sin cobertura de regresion.

### Evaluacion probabilidad por impacto

| Caso | Riesgo principal | Probabilidad | Impacto | Nivel |
|---|---|---:|---:|---|
| INT-07-01 | Archivo exportado con formato invalido | Baja | Alta | Medio-Alto |
| INT-07-02 | Entidades no recreadas o con IDs colisionando | Media | Alta | Alto |
| INT-07-03 | Perdida de integridad de datos en el roundtrip | Media | Critico | Critico |
| INT-07-04 | Importacion de archivo corrupto crea datos parciales | Baja-Media | Alto | Alto |
| INT-07-05 | Adjunto subido pero no recuperable via API | Media | Media | Medio |

### Priorizacion

El caso INT-07-03 (roundtrip) tiene la prioridad mas alta porque es la evidencia definitiva de que exportar e importar preserva el contenido: si este caso falla, ningun respaldo es confiable. INT-07-02 y INT-07-04 son de prioridad alta porque protegen, respectivamente, la correcta recreacion de datos y la integridad ante entradas invalidas. INT-07-01 valida la base tecnica del formato de archivo, e INT-07-05 protege una funcionalidad relacionada (adjuntos) que comparte la capa de FileSystem.

### Mitigacion

La mitigacion se realiza ejercitando el flujo real de exportacion e importacion contra un tablero con contenido, y verificando en `App`/`Store` que los datos exportados coinciden con los importados. Para el caso de error se verifica explicitamente que no queden tableros parcialmente creados, protegiendo la consistencia de la base de datos.

## 5. Justificacion de herramientas

| Herramienta | Por que se eligio | Alternativas consideradas | Limitacion conocida |
|---|---|---|---|
| Go test | Mecanismo nativo de pruebas del backend Go, permite ejecutar solo el flujo con `-run TestINT07`. | Herramientas HTTP externas como Postman/Newman. | Depende de que el entorno de base de datos de pruebas este disponible. |
| Testify require | Ya se usa en el repositorio, ofrece aserciones claras (`require.Len`, `require.Equal`). | `testing` puro. | Detiene el caso en el primer fallo critico. |
| TestHelper de integrationtests | Provee servidor, clientes autenticados y limpieza consistente de recursos. | Crear un setup propio para cada prueba. | Hereda limitaciones del arnes existente. |
| Cliente HTTP del proyecto (`th.Client`) | Garantiza que las pruebas pasen por los endpoints reales de exportacion/importacion/adjuntos. | Llamar directamente a `App.ExportArchive`/`App.ImportArchive`. | Algunas verificaciones internas requieren complementarse con `App`/`Store`. |
| App/Store publicos | Permiten confirmar persistencia real de tableros y bloques tras importar, sin SQL directo. | Consultas SQL manuales sobre las tablas. | Solo exponen la informacion disponible via interfaces publicas. |
| Logs con `t.Log` | Explican la intencion de cada prueba en la salida estandar, utiles para el debate. | Comentarios internos solamente. | No sustituyen un reporte de ejecucion formal. |

## 6. Justificacion estrategica por caso

### INT-07-01 Exportar tablero genera archivo con formato correcto

Antes de confiar en el contenido de un archivo exportado hay que confirmar que el formato base es el esperado. La prueba valida que la respuesta no este vacia y que respete la firma de cabecera de un archivo zip, que es el contenedor usado por Focalboard para los `.boardarchive`.

### INT-07-02 Importar archivo recrea entidades en el Store

Este caso es central porque valida el proposito mismo de la importacion: que el tablero y sus bloques existan de nuevo en el sistema, con nuevos IDs (para evitar colisionar con el original) pero con el mismo contenido logico (titulo).

### INT-07-03 Exportar e importar tablero (roundtrip) preserva integridad

Es el caso de mayor valor del flujo: ejecuta exportacion e importacion en secuencia y compara cantidad y contenido de bloques entre el tablero original y el importado, demostrando que el ciclo completo no pierde ni corrompe informacion.

### INT-07-04 Importar archivo invalido es rechazado sin corromper datos

Un sistema de importacion debe ser resiliente ante archivos malformados. La prueba confirma que ante un archivo invalido no se crean tableros nuevos, protegiendo la base de datos de estados parciales o corruptos.

### INT-07-05 Subir archivo adjunto se almacena y es recuperable

Los adjuntos comparten la capa de FileSystem con el mecanismo de exportacion/importacion. La prueba valida que un archivo subido a una tarjeta queda referenciado en el sistema y su informacion es recuperable via API, cerrando el flujo de principio a fin.

## 7. Relacion con pruebas unitarias previas del componente ws

Las pruebas unitarias previas del componente `server/ws` validaban en aislamiento la logica de conexion, autenticacion y difusion de eventos websocket, con dependencias externas simuladas mediante `gomock`. Esas pruebas responden si cada pieza del componente `ws` funciona correctamente de forma unitaria.

Las pruebas INT-07 amplian el alcance: verifican que el contenido que un usuario exporta o importa efectivamente atraviesa la API REST, la logica de negocio, el sistema de archivos y el Store de forma coherente. La diferencia central es que las unitarias comprueban una unidad aislada del servidor de tiempo real, mientras que estas pruebas de integracion comprueban un flujo de negocio completo, de punta a punta, sobre datos persistentes.

La relacion es complementaria: las unitarias del componente `ws` reducen el riesgo de errores en la logica de comunicacion en tiempo real, mientras que las pruebas de integracion de este flujo reducen el riesgo de perdida o corrupcion de datos en operaciones de respaldo y migracion de tableros.

## 8. Conclusion argumentativa

El flujo INT-07 fue seleccionado porque la importacion y exportacion de tableros es la unica via de respaldo y portabilidad de datos que ofrece Focalboard, y un defecto en este flujo puede significar perdida de informacion irreversible para el usuario. La estrategia de prueba prioriza el roundtrip como evidencia de integridad, protege contra archivos invalidos y verifica el flujo de adjuntos relacionado, aportando evidencia solida de que el sistema mantiene la consistencia de los datos exportados e importados.
