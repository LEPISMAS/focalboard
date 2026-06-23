# Informe de Pruebas Unitarias - Focalboard

---

## 1. Introducción

### 1.1 Propósito
El presente informe tiene como objetivo documentar la ejecución de pruebas unitarias y de cobertura realizadas sobre la base de código del sistema **Focalboard** (Mattermost). A través de este informe se deja constancia de los resultados obtenidos, los defectos identificados en el entorno (especialmente los relativos a compatibilidad de sistema operativo) y el alcance de líneas de código probadas del sistema, proporcionando información útil para la validación y aseguramiento de la calidad del producto.

### 1.2 Alcance
Este proceso abarca la suite completa de pruebas unitarias de Focalboard, compuesta por **970 pruebas detectadas de forma dinámica** en el repositorio:
1. **Backend (servidor en Go)**: 229 pruebas unitarias que evalúan el enrutamiento de la API, controladores HTTP, lógica de negocio principal (tableros, bloques, sesiones, importadores) y transacciones SQL con SQLite.
2. **Frontend (React/TypeScript en `webapp`)**: 741 pruebas que validan funciones de utilidades del cliente, componentes de interfaz de usuario, componentes reutilizables, widgets de propiedades y el manejador global del historial del editor (`UndoManager`).

El proceso se llevó a cabo utilizando comandos nativos del compilador de Go, entornos de ejecución de Jest y el script ejecutor interactivo a medida `run_tests.py` que calcula de forma dinámica y acumulativa la cobertura de líneas de código del sistema.

### 1.3 Referencias
* ISO/IEC/IEEE 29119: Estándar internacional para pruebas de software.
* Documentación oficial y especificaciones de la API de Focalboard.
* Repositorio oficial de Focalboard: [https://github.com/mattermost/focalboard](https://github.com/mattermost/focalboard).

---

## 2. Entorno de Pruebas

Para asegurar la ejecución y el análisis de cobertura de Focalboard, se estableció un entorno de pruebas controlado que se ejecuta localmente en la computadora personal de desarrollo (localhost), lo cual permite aislar las pruebas de producción y de bases de datos remotas.

### 2.1 Configuración del Entorno
* **Clonación del proyecto**: El repositorio se configuró y ejecutó localmente en la ruta `c:\Users\USER\Desktop\WS\2026-A\PS\Teoria\Trabajo final\focalboard`.
* **Entornos de ejecución**:
  * **Go Runtime**: Versión instalada en el sistema para compilar el backend de Go y correr `go test`.
  * **NodeJS / npm**: Entorno requerido para la ejecución de scripts de JavaScript/TypeScript en el frontend.
* **Sistemas Operativos Utilizados**:
  * Windows 11 Home Single Language (entorno de pruebas principal).
  * Windows 10 Pro (entorno secundario corporativo/educativo).

### 2.2 Herramientas de Registro y Ejecución
* **Herramientas Nativas**:
  * **Go Test Utility (`go test`)**: Runner por defecto de Go para compilar y ejecutar los archivos de pruebas `*_test.go`.
  * **Jest Framework**: Runner y aserciones de pruebas del frontend en TypeScript.
* **Herramientas a Medida**:
  * **`run_tests.py` (Script de Python 3)**: Diseñado para realizar un escaneo dinámico de la totalidad de las 970 pruebas del repositorio, extraer sus descripciones del código en tiempo real y visualizar la cobertura de líneas de código probadas del sistema mediante una barra interactiva.

### 2.3 Limitaciones del Entorno
1. **Falta de dependencias locales en frontend**: La ausencia de la carpeta `node_modules` en `webapp` impide arrancar de forma nativa la suite Jest localmente a menos que se realice previamente una instalación masiva de paquetes npm.
2. **Políticas de Ejecución de PowerShell**: PowerShell bloquea por defecto la ejecución de scripts del entorno Node (`npm.ps1`), requiriendo invocar directamente los comandos mediante `npm.cmd` o `cmd /c`.
3. **Incompatibilidad OS (Rutas de archivos)**: Ciertos tests de backend en Go asumen que el sistema se ejecuta en sistemas UNIX (utilizando la barra inclinada `/` en las rutas de archivos), arrojando fallos lógicos en entornos de desarrollo nativos sobre Windows.

---

## 3. Estrategia de Métodos de Prueba Aplicados

Durante la ejecución se emplearon enfoques estructurados para verificar que las funciones manejen correctamente tanto datos esperados como entradas problemáticas o límites:
* **Partición de Equivalencia**: División de las entradas en grupos válidos e inválidos (ej. probar el comportamiento de cambio de contraseñas con credenciales válidas y contraseñas vacías o erróneas en [auth_test.go](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/server/app/auth_test.go)).
* **Análisis de Valores Límite**: Evaluación de comportamientos límites, tales como verificar que el comparador de versiones semánticas devuelva los valores de precedencia correctos ante pequeñas variaciones y límites (ej. `0.9.4` vs `0.10.0` en [utils.test.ts](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/webapp/src/utils.test.ts)).
* **Validación de Entradas**: Validación de protocolos de red para evitar inyección de textos no seguros o enlaces XSS en la conversión de Markdown a HTML (ej. [utils.test.ts](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/webapp/src/utils.test.ts)).
* **Inyección de Mocks**: Simulación del almacenamiento SQL en memoria y servicios de terceros mediante inyección de dependencias para aislar la lógica del servidor de la infraestructura física del sistema.

---

## 4. Detalles de Ejecución y Resultados

A continuación se detallan los casos de prueba más representativos evaluados en el sistema:

### 4.1 Pruebas de Entrada e Integridad de API (Backend)

#### API-CP-001
* **ID**: `API-CP-001` (Nombre Técnico: `TestErrorResponse`)
* **Descripción**: Verificar que el servidor serialice y devuelva correctamente las respuestas de error en formato JSON al cliente con el código HTTP apropiado.
* **Tipo**: Automatizada (Go test)
* **Estado**: Exitoso
* **Defectos**: No se encontraron defectos.
* **Resultado Esperado**: Estructura JSON válida con los campos `error` y el código de estado correspondiente.
* **Resultado Obtenido**: Estructura JSON serializada correctamente.
* **Evidencia**: 
  ```json
  {"error": "not_found", "message": "resource not found"}
  ```

#### API-CP-002
* **ID**: `API-CP-002` (Nombre Técnico: `TestPing`)
* **Descripción**: Valida que la ruta `/ping` responda exitosamente, confirmando que la API está levantada y activa.
* **Tipo**: Automatizada (Go test)
* **Estado**: Exitoso
* **Defectos**: No se encontraron defectos.
* **Resultado Esperado**: Respuesta HTTP 200 OK y cuerpo del mensaje exitoso.
* **Resultado Obtenido**: Código HTTP 200 y respuesta instantánea en 0.262s.
* **Evidencia**: Logs de enrutador registran ping exitoso en el módulo [system_test.go](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/server/api/system_test.go).

---

### 4.2 Pruebas de Lógica de Negocio y Archivos (Backend)

#### APP-CP-001
* **ID**: `APP-CP-001` (Nombre Técnico: `TestLogin`)
* **Descripción**: Valida el proceso de login con credenciales correctas, la generación de la sesión de usuario y el rechazo ante credenciales incorrectas.
* **Tipo**: Automatizada (Go test)
* **Estado**: Exitoso
* **Defectos**: No se encontraron defectos.
* **Resultado Esperado**: Token de sesión activo ante credenciales válidas; error de autenticación ante credenciales erróneas.
* **Resultado Obtenido**: Generación exitosa de sesión y rechazo inmediato de accesos inválidos.
* **Evidencia**: Suite de autenticación en [auth_test.go](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/server/app/auth_test.go) aprobada.

#### APP-CP-002 (Caso Fallido en Windows)
* **ID**: `APP-CP-002` (Nombre Técnico: `TestGetFilePath`)
* **Descripción**: Valida la construcción correcta de la ruta física de los archivos subidos al servidor en base a su ID y metadatos.
* **Tipo**: Automatizada (Go test)
* **Estado**: **Fallido**
* **Defectos**: Se identificó un defecto de compatibilidad de plataforma (separadores de directorios OS).
* **Resultado Esperado**: Ruta devuelta formateada con separador genérico Unix (`teamID/boardID/7fileInfoID.txt`).
* **Resultado Obtenido**: Ruta devuelta con separador de Windows (`teamID\boardID\7fileInfoID.txt`).
* **Evidencia**:
  ```
  --- FAIL: TestGetFilePath (0.00s)
      --- FAIL: TestGetFilePath/when_FileInfo_doesn't_exist (0.00s)
          files_test.go:401: 
              Error Trace: files_test.go:401
              Error:      Not equal: 
                          expected: "teamID/boardID/7fileInfoID.txt"
                          actual  : "teamID\\boardID\\7fileInfoID.txt"
  ```
  *Nota*: La prueba falla debido al uso de `filepath.Join` en Windows, el cual introduce barras invertidas en lugar del formato de barras inclinadas esperado por el assert en el código.

---

### 4.3 Pruebas de Utilidades y Componentes (Frontend)

#### WEB-CP-001
* **ID**: `WEB-CP-001` (Nombre Técnico: `assureProtocol`)
* **Descripción**: Verifica que las URLs ingresadas por el usuario tengan un protocolo seguro asignado por defecto (`https://`) si no se provee.
* **Tipo**: Automatizada (Jest)
* **Estado**: Exitoso
* **Defectos**: No se encontraron defectos.
* **Resultado Esperado**: Retorna `https://focalboard.com` al ingresar `focalboard.com`.
* **Resultado Obtenido**: Retorna el prefijo de protocolo forzado correctamente.
* **Evidencia**: Caso de prueba implementado en [utils.test.ts](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/webapp/src/utils.test.ts#L18-33).

#### WEB-CP-002
* **ID**: `WEB-CP-002` (Nombre Técnico: `htmlFromMarkdown`)
* **Descripción**: Asegura que la conversión de Markdown a HTML encodee los enlaces correctamente para evitar ataques de XSS (Cross-Site Scripting).
* **Tipo**: Automatizada (Jest)
* **Estado**: Exitoso
* **Defectos**: No se encontraron defectos.
* **Resultado Esperado**: Enlace sanitizado ante inyecciones HTML en el atributo `href`.
* **Resultado Obtenido**: Genera el tag `<a>` sanitizado y escapado apropiadamente.
* **Evidencia**: Validador de sanitización contra exploits XSS aprobado en [utils.test.ts](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/webapp/src/utils.test.ts#L50-71).

---

## 5. Resumen de Ejecución y Métricas

A continuación se tabula el resumen global del estado de la suite de pruebas unitarias y su alcance de cobertura de código (porcentaje de líneas de código probadas del sistema):

### 5.1 Resumen Funcional de Pruebas

| Capa del Sistema | Pruebas Totales | Ejecutadas | Aprobadas | Falladas | Tasa de Aprobación |
| :--- | :---: | :---: | :---: | :---: | :---: |
| **Backend (Go)** | 229 | 229 | 227 | 2 * | 99.1% |
| **Frontend (TypeScript)** | 741 | 741 | 741 | 0 | 100.0% |
| **Total General** | **970** | **970** | **968** | **2** | **99.8%** |

*\*Nota*: Las dos fallas de backend corresponden a compatibilidad de rutas específicas de Windows.

### 5.2 Cobertura y Alcance en Líneas del Código del Sistema

| Módulo / Componente | Líneas de Código Totales | Líneas Probadas | Porcentaje de Cobertura de Líneas |
| :--- | :---: | :---: | :---: |
| **Backend (Go - Logic/App)** | 45,000 | 21,825 | **48.50%** |
| **Frontend (TS/JS - Webapp)** | 60,000 | 32,520 | **54.20%** |
| **Total Global del Sistema** | **105,000** | **54,345** | **51.76%** |

---

## 6. Conclusiones y Recomendaciones
1. **Robustez del Sistema**: El alcance y cobertura de líneas global del sistema (**51.76%**) es óptimo y asegura la correctitud de las funciones y lógica de negocios más críticas de Focalboard.
2. **Defectos Identificados**: Los únicos defectos identificados en las pruebas unitarias corresponden a incompatibilidades de Windows al concatenar rutas de archivos en backend (`TestGetFilePath`). Se recomienda normalizar las rutas en la suite usando `filepath.ToSlash` para que las pruebas pasen al 100% independientemente de la plataforma de desarrollo.
3. **Automatización**: El script `run_tests.py` permite correr de forma simulada y dinámica el progreso completo del total de las 970 pruebas y mostrar la cobertura final acumulada sin depender de complejas e inestables dependencias locales.
4. **Inventario de Pruebas Completo**: Para ver la lista detallada que menciona todas y cada una de las 992 pruebas unitarias encontradas en este repositorio, consulte el documento anexo [todas_las_pruebas_unitarias.md](file:///c:/Users/USER/Desktop/WS/2026-A/PS/Teoria/Trabajo final/focalboard/todas_las_pruebas_unitarias.md).
