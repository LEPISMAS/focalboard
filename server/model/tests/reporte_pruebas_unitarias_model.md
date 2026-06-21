# Reporte de Pruebas Unitarias - Módulo Model de Focalboard

Este reporte detalla el plan de pruebas unitarias implementado para el paquete `server/model/` de Focalboard. El objetivo fue alcanzar y superar el 90% de cobertura en las estructuras clave relacionadas a declaraciones de bloques, tableros, campos y sus respectivos métodos de validación y utilidades.

## Resumen del Progreso de Cobertura

A continuación se muestra el porcentaje de cobertura de sentencias por cada archivo clave bajo el alcance:

| Archivo | Área Evaluada / Descripción | Pruebas Unitarias Realizadas | Cobertura de Sentencias |
| :--- | :--- | :--- | :---: |
| `block.go` | Estructuras de bloques, validación, copias de auditoría, parches y límites de almacenamiento. | `TestBlockValidation`, `TestBlockPatch`, `TestBlockLogClone`, `TestBlockGetLimited`, `TestBlockStampModificationMetadata` | **100.0%** |
| `blockid.go` | Generación y asignación de identificadores de bloques, mapeo de referencias de dependencias (`contentOrder`, `defaultTemplateId`). | `TestGenerateBlockIDs`, `TestGenerateBlockIDsEdgeCases` | **100.0%** |
| `blocktype.go` | Manejo de enums de tipo de bloque, conversiones string, mapeos de tipo ID. | `TestBlockTypeString`, `TestBlockTypeFromString`, `TestBlockType2IDType`, `TestIsErrInvalidBlockType` | **100.0%** |
| `board.go` | Estructuras de tableros, validaciones, mapeos de campos de búsqueda, lectura de JSON. | `TestBoardGetPropertyString`, `TestBoardJSONHelpers`, `TestBoardPatch`, `TestBoardValidation`, `TestBoardSearchFieldFromString` | **100.0%** |
| `boards_and_blocks.go` | Operaciones agrupadas sobre tableros y bloques, parches y eliminación. | `TestBoardsAndBlocks` | **100.0%** |
| `card.go` | Estructuras y ciclos de vida de tarjetas (cards), conversiones a bloques planos, validaciones. | `TestCardErrors`, `TestCardPopulateAndValid`, `TestCardPatch`, `TestCard2Block`, `TestBlock2Card`, `TestCardPatch2BlockPatch` | **100.0%** |
| `properties.go` | Resolutor de valores de propiedades, esquemas JSON y serialización de campos del tablero. | `TestPropDefGetValue`, `TestParseDate`, `TestParsePropertySchema`, `TestParseProperties` | **100.0%** |
| `auth.go` | Peticiones de login/registro, validación de contraseñas seguras y deserialización. | `TestAuthParamError`, `TestLoginResponseFromJSON`, `TestRegisterRequestIsValid`, `TestChangePasswordRequestIsValid` | **100.0%** |
| `error.go` | Mapeo y envolturas de errores HTTP, manejadores de estado de error (Bad Request, Forbidden, Not Found). | `TestErrorStructures`, `TestIsErrHandlers` | **100.0%** |
| `file.go` | Metadatos y constructor de archivos multimedia adjuntos. | `TestFile` | **100.0%** |
| `sharing.go` | Configuración e importación JSON de tableros públicos/compartidos. | `TestSharing` | **100.0%** |
| `team.go` | Gestión e importación JSON de equipos/organizaciones. | `TestTeam` | **100.0%** |
| `util.go` | Utilidades de conversión de marcas de tiempo (timestamp milliseconds). | `TestUtils` | **100.0%** |
| `version.go` | Escritura de trazas y metadatos de versión del servidor en logs. | `TestVersion` | **100.0%** |
| `category.go` | Agrupación de categorías de tableros, hidratación de defaults y validación de tipos. | `TestCategory` | **93.2%** |
| **Total Global** | **Módulo `server/model` Completo** | **Todas las pruebas del paquete** | **92.5%** |

## Ejecución Local de Pruebas

Se han creado scripts de automatización para facilitar la ejecución y análisis de cobertura local:

### En Linux / MacOS:
Ejecuta el script:
```bash
./server/model/tests/run_tests_model.sh
```

### En Windows:
Ejecuta el archivo batch:
```cmd
server\model\tests\run_tests_model.bat
```
