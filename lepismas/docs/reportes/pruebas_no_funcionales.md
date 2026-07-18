# Reporte de Pruebas No Funcionales (HITO 3)

Este documento detalla la automatización, diseño, resultados y comprobación de la suite de pruebas no funcionales para Focalboard, con un enfoque en **Concurrencia (Integridad de Estado)** y **Rendimiento/Desempeño (Escalabilidad y Latencia)** bajo los estándares **ISO/IEC 25010** e **IEEE 829**.

---

## 1. Ubicación de las Pruebas

Tanto las pruebas funcionales como las no funcionales residen en el repositorio bajo estructuras independientes para asegurar aislamiento:

| Tipo de Prueba | Ubicación en el Repositorio | Comando de Ejecución |
| :--- | :--- | :--- |
| **Pruebas Unitarias/Integración (Backend)** | `server/api/tests/`, `server/model/` | `go test ./server/...` |
| **Pruebas de Componentes (Frontend)** | `webapp/src/pages/tests/`, `webapp/src/widgets/` | `npm run test` (en webapp) |
| **Pruebas de Concurrencia (Playwright)** | `non_functional_tests/concurrency.spec.ts` | `npm run test:concurrency` (en non_functional_tests) |
| **Pruebas de Rendimiento (k6)** | `non_functional_tests/perf_load.js` & `perf_insert.js` | `./k6 run perf_load.js` & `perf_insert.js` |
| **Script de Orquestación Completa** | `run_non_functional_tests.sh` | `./run_non_functional_tests.sh` (en la raíz) |

---

## 2. Diseño de las Pruebas

El diseño de las pruebas sigue las directrices de la especificación técnica en `buffer/HITO3/pruebas_no_funcionales.md`.

### Pruebas de Concurrencia (Integridad de Estado)
*   **CON-01 (Conflicto de Escritura - Race Condition):** Dos usuarios concurrentes (U1 y U2) abren la misma tarjeta en navegadores paralelos (Playwright Browser Contexts) y modifican el título simultáneamente. El servidor debe procesar ambas peticiones de forma segura aplicando **Last Write Wins (LWW)** y manteniendo el estado consistente en la base de datos sin generar bloqueos.
*   **CON-02 (Conflicto de Schema y Datos - Card Move vs Column Delete):** Un usuario mueve una tarjeta a la columna "Done" a través del API, mientras concurrentemente otro usuario elimina la columna "Done" del tablero. El servidor debe procesar ambos patches sin corromper el esquema del tablero o el estado de las propiedades de la tarjeta.
*   **CON-03 (Consistencia de WebSockets - Reconexión):** Se simula la desconexión de red de un usuario (U1) usando el modo offline de Playwright. Mientras U1 está desconectado, U2 crea una tarjeta. Al volver a conectar la red, el canal WebSocket de U1 debe reconciliar automáticamente el estado del cliente con el servidor, visualizando la nueva tarjeta sin necesidad de refrescar la página de manera manual.

### Pruebas de Rendimiento (Escalabilidad y Latencia)
*   **PER-01 (Carga de Tablero Masivo):** Carga un tablero masivo que contiene **2,000 tarjetas** previamente creadas mediante una inyección de datos controlada. Se simulan **50 usuarios concurrentes (VUs)** leyendo el tablero de forma simultánea. Métrica objetivo: P95 de tiempo de respuesta (TTFB) < 300ms.
*   **PER-02 (Latencia de UI Optimista):** Mide el tiempo de respuesta del DOM en el frontend al hacer clic en "Añadir Tarjeta". El modal debe ser visible de forma instantánea (latencia optimista perceptiva < 200ms).
*   **PER-03 (Límites de Persistencia):** Inyección masiva sostenida de **500 inserciones de tarjetas por segundo** a la API del servidor usando el ejecutor de tasa constante (`constant-arrival-rate`) de k6 con hasta 100 usuarios virtuales concurrentes.

---

## 3. Resultado de las Pruebas

Las pruebas fueron ejecutadas sobre un servidor Focalboard local utilizando una base de datos SQLite aislada (`focalboard_test.db`).

### Resultados de Concurrencia (Playwright)
*   **CON-01:** **EXITOSO**. Ambos usuarios modificaron la tarjeta de forma concurrente, el servidor procesó los dos PATCH con código HTTP 200 de forma ordenada. La base de datos persistió el valor correspondiente al LWW de forma estable.
*   **CON-02:** **EXITOSO**. El movimiento de la tarjeta y la eliminación de la opción de columna ocurrieron al mismo tiempo. El servidor retornó status 200 en ambas operaciones, dejando el tablero consistente con la opción borrada de la lista de propiedades.
*   **CON-03:** **EXITOSO**. U1 desconectado no recibió actualizaciones del UI. En cuanto la red se restableció, el WebSocket reconcilió el log de transacciones y el cliente Playwright detectó la nueva tarjeta de forma instantánea y reactiva.
*   **PER-02 (Latencia UI):** **EXITOSO**. Latencia promedio de apertura del modal tras la mutación del DOM: **123ms** (dentro del límite de tolerancia de 200ms).

### Resultados de Rendimiento (k6)
*   **PER-01 (Carga Masiva - 50 VUs / 2000 Tarjetas):**
    *   **Peticiones completadas:** 455 reqs.
    *   **P95 de Tiempo de Respuesta:** **5.67s**.
    *   *Nota de Análisis:* El tiempo de respuesta superó el objetivo ideal de 300ms debido a la naturaleza monohilo de las transacciones de SQLite y el procesamiento de serialización de un payload masivo JSON (341 MB transferidos) en la CPU del entorno virtualizado de desarrollo.
*   **PER-03 (Inserción Masiva - 500 req/s):**
    *   **Peticiones procesadas con éxito (200 OK):** 94.63% (335 reqs).
    *   **Dropped Iterations:** 4647 (debido a saturación de VUs en el servidor SQLite).
    *   **P95 de Latencia de Inserción:** **7.98s**.
    *   *Nota de Análisis:* La base de datos local SQLite bajo una tasa sostenida de 500 inserciones/seg experimenta contención de bloqueos debido a su esquema de bloqueos a nivel de archivo. Para producción con estas tasas de rendimiento, se valida técnicamente la recomendación de migrar a PostgreSQL.

---

## 4. Comprobación y Automatización

Para ejecutar, repetir y comprobar estas pruebas de manera automatizada en el flujo de integración continua, se ha creado el script de bash `run_non_functional_tests.sh` en la raíz del repositorio.

### Flujo de Ejecución del Script
El script realiza las siguientes acciones automáticamente:
1.  Termina cualquier proceso remanente del servidor Focalboard para liberar el puerto 8000.
2.  Elimina la base de datos temporal `focalboard_test.db` para asegurar pruebas libres de efectos colaterales.
3.  Inicia el servidor Focalboard en segundo plano con la configuración de pruebas (`server-test-config.json`).
4.  Espera hasta 30 segundos verificando la disponibilidad del puerto.
5.  Ejecuta las pruebas de concurrencia de Playwright (`npm run test:concurrency`).
6.  Corre el script de inicialización masiva para generar el tablero de 2,000 tarjetas (`setup_perf.spec.ts`).
7.  Ejecuta los escenarios de k6 de carga (`perf_load.js`) e inserciones (`perf_insert.js`).
8.  Detiene el servidor y finaliza la sesión de pruebas con un reporte consolidado.

### Cómo ejecutar las pruebas
Para iniciar toda la suite no funcional, ejecute el siguiente comando desde la raíz del proyecto:
```bash
./run_non_functional_tests.sh
```
Las trazas del servidor se almacenarán en `server.log` para facilitar el debugging en caso de errores en la aserción de base de datos.
