# Plan de pruebas de sistema E2E de Focalboard

## 1. Introducción

### 1.1 Objetivo

Este documento define el plan y la evidencia de dos pruebas de sistema end-to-end (E2E) ya implementadas en Focalboard. Las pruebas ejercitan el sistema desde la interfaz del navegador hasta la persistencia y validan dos atributos de calidad:

1. **Seguridad**, mediante el flujo E2E de autenticación, sesión, cambio de contraseña y registro con invitación.
2. **Usabilidad**, en sus dimensiones automatizables de efectividad y eficiencia, mediante el flujo E2E de creación y gestión de tableros, tarjetas y vistas.

Las pruebas no son propuestas hipotéticas. Se encuentran implementadas con Cypress y se ejecutan automáticamente en GitHub Actions.

### 1.2 Sistema bajo prueba

El recorrido verificado por ambas pruebas es:

```text
Usuario simulado por Cypress
    -> navegador Electron
    -> interfaz React/TypeScript
    -> API HTTP y WebSocket
    -> servidor Go
    -> SQLite
    -> respuesta observable en la interfaz
```

Al recorrer todas estas capas en una instancia real de la aplicación, los escenarios se clasifican como pruebas de sistema E2E y no únicamente como pruebas unitarias o de integración.

### 1.3 Alcance

| ID | Atributo | Escenario automatizado | Archivo |
|---|---|---|---|
| SYS-SEC-01 | Seguridad | Registro, login, logout, cambio de contraseña y control de invitaciones | `webapp/cypress/integration/loginActions.ts` |
| SYS-USA-01 | Usabilidad: efectividad y eficiencia | Crear, modificar y eliminar un tablero; crear y visualizar tarjetas y vistas | `webapp/cypress/integration/createBoard.ts` |

Quedan fuera de alcance las pruebas de penetración, carga, disponibilidad prolongada, satisfacción subjetiva y accesibilidad completa. Estas requieren herramientas y métodos adicionales.

## 2. Atributos de calidad

### 2.1 Seguridad

La prueba `loginActions.ts` evalúa controles de seguridad observables desde el exterior del sistema:

- redirección de visitantes no autenticados hacia `/login`;
- creación de la primera cuenta;
- autenticación con credenciales válidas;
- cierre de sesión y revocación del acceso anterior;
- cambio de contraseña y acceso con la contraseña nueva;
- rechazo del registro de un segundo usuario sin invitación;
- registro permitido mediante una invitación válida.

La prueba verifica seguridad funcional. No sustituye un análisis de vulnerabilidades ni una prueba de penetración.

### 2.2 Usabilidad: efectividad y eficiencia

La prueba `createBoard.ts` evalúa aspectos automatizables de usabilidad durante una tarea principal de Focalboard:

- presencia y visibilidad de controles necesarios;
- capacidad de completar la creación y eliminación de un tablero;
- creación y edición de tarjetas y vistas;
- foco automático en el título de una tarjeta nueva;
- retroalimentación visual después de las acciones;
- ausencia de bloqueos y timeouts durante el flujo;
- duración total informada por Cypress.

La automatización permite medir efectividad y eficiencia operativa, pero no satisfacción. La satisfacción requeriría pruebas con usuarios y un instrumento como SUS (System Usability Scale).

## 3. Entorno y automatización

### 3.1 Componentes

- Frontend React y TypeScript compilado con Webpack.
- Backend Go.
- Base de datos SQLite en memoria para Cypress.
- Cypress 9.5.2 con navegador Electron.
- GitHub Actions sobre Ubuntu.

### 3.2 Ejecución en GitHub Actions

El workflow `.github/workflows/ci.yml`, job `ci-ubuntu-webapp`, realiza las siguientes acciones:

1. instala dependencias con `npm ci`;
2. configura Go y Node.js;
3. compila el servidor Linux y el Webapp;
4. ejecuta `make webapp-ci`;
5. ejecuta lint, pruebas Jest y `npm run cypress:ci`;
6. Cypress inicia el servidor real y ejecuta los specs configurados, incluidos `loginActions.ts` y `createBoard.ts`.

### 3.3 Criterios de entrada

- El repositorio compila correctamente.
- Las dependencias Go y Node.js están disponibles.
- El binario `bin/focalboard-server` fue generado.
- El puerto 8088 está disponible.
- Cypress puede iniciar el navegador.

### 3.4 Criterios de salida

- Los escenarios `SYS-SEC-01` y `SYS-USA-01` terminan con estado `Passing`.
- No existen aserciones fallidas ni timeouts.
- El proceso Cypress devuelve código de salida 0.
- GitHub Actions registra la ejecución y conserva el resultado del check.

### 3.5 Ejecución local con Docker

Los siguientes comandos reproducen las pruebas E2E desde Windows PowerShell sin instalar Go, Node.js o Cypress directamente en el equipo. Deben ejecutarse desde la raíz del repositorio:

```powershell
cd C:\Users\marco\Desktop\focalboard
```

#### Paso 1: compilar el servidor Linux compatible con Cypress

Focalboard utiliza CGO y SQLite. Se usa la imagen Bullseye porque su versión de GLIBC es compatible con `cypress/included:9.5.2`.

```powershell
docker run --rm -v "${PWD}:/work" -w /work golang:1.21-bullseye bash -c "mkdir -p bin/linux bin && cd server && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags 'json1 sqlite3' -o ../bin/linux/focalboard-server ./main && cp ../bin/linux/focalboard-server ../bin/focalboard-server"
```

El comando debe finalizar sin errores y generar `bin/focalboard-server`.

#### Paso 2: ejecutar las dos pruebas de sistema E2E

```powershell
docker run --rm --entrypoint bash -v "${PWD}:/work" -v focalboard-webapp-node-modules:/work/webapp/node_modules -w /work/webapp cypress/included:9.5.2 -c "npm ci && npm run pack && { npm run runserver-test & npx --no-install wait-on http://localhost:8088 && npx cypress run --spec 'cypress/integration/loginActions.ts,cypress/integration/createBoard.ts'; }"
```

Este comando instala las dependencias dentro de un volumen Docker, compila el Webapp, inicia Focalboard con SQLite, espera que el servidor responda y ejecuta los dos specs seleccionados.

#### Ejecución individual de seguridad

```powershell
docker run --rm --entrypoint bash -v "${PWD}:/work" -v focalboard-webapp-node-modules:/work/webapp/node_modules -w /work/webapp cypress/included:9.5.2 -c "npm ci && npm run pack && { npm run runserver-test & npx --no-install wait-on http://localhost:8088 && npx cypress run --spec cypress/integration/loginActions.ts; }"
```

#### Ejecución individual de usabilidad

```powershell
docker run --rm --entrypoint bash -v "${PWD}:/work" -v focalboard-webapp-node-modules:/work/webapp/node_modules -w /work/webapp cypress/included:9.5.2 -c "npm ci && npm run pack && { npm run runserver-test & npx --no-install wait-on http://localhost:8088 && npx cypress run --spec cypress/integration/createBoard.ts; }"
```

#### Resultado esperado

Al finalizar deben observarse estas condiciones:

- `All specs passed!`;
- cantidad de pruebas fallidas igual a 0;
- `loginActions.ts` y `createBoard.ts` marcados como aprobados;
- código de salida 0.

Las advertencias de `npm audit`, tamaño del bundle, `tput` o WebSocket `close 1001 (going away)` al cerrarse el navegador no invalidan la ejecución si Cypress informa cero pruebas fallidas.

## 4. Metodología de diseño de caja negra

Las técnicas se aplican sobre entradas y resultados observables, sin utilizar el funcionamiento interno del código como criterio de validación.

### 4.1 Partición de equivalencia

Se dividen las entradas en clases cuyo comportamiento esperado es equivalente. Se selecciona un representante de cada clase cubierta por los escenarios existentes.

### 4.2 Tabla de decisión

Se modelan reglas donde el resultado depende de varias condiciones. Esta técnica se aplica especialmente al registro y al acceso al espacio de trabajo.

### 4.3 Transición de estados

Se modelan estados observables y eventos realizados por el usuario. Se aplica a la sesión y al ciclo de vida del tablero.

### 4.4 Valores límite

El análisis de valores límite forma parte de la metodología general, pero los dos specs seleccionados no prueban longitudes mínimas o máximas. Por tanto, no se declara cobertura de valores límite en esta ejecución. Para aplicarla sería necesario agregar casos explícitos para longitud de usuario, contraseña y título.

Esta delimitación evita atribuir a las pruebas una cobertura que actualmente no ejecutan.

## 5. SYS-SEC-01: seguridad de autenticación y registro

### 5.1 Ficha técnica

| Campo | Valor |
|---|---|
| ID | SYS-SEC-01 |
| Atributo | Seguridad |
| Funcionalidad | Autenticación, sesión, contraseña e invitaciones |
| Implementación | `webapp/cypress/integration/loginActions.ts` |
| Caso Cypress | `Can perform login/register actions` |
| Precondiciones | Servidor limpio; navegador sin sesión activa; registro inicial habilitado |
| Datos principales | Usuario, correo y contraseña definidos en `webapp/cypress.json`; contraseña nueva; enlace de invitación |
| Técnicas | Partición de equivalencia, tabla de decisión y transición de estados |
| Prioridad | Alta |

### 5.2 Partición de equivalencia

| Código | Entrada o condición | Clase válida cubierta | Clase inválida cubierta | Resultado esperado |
|---|---|---|---|---|
| SEC-PE-01 | Sesión | Usuario autenticado | Visitante sin sesión | Autenticado accede al workspace; visitante es redirigido a `/login` |
| SEC-PE-02 | Contraseña después del cambio | Contraseña nueva | Contraseña anterior | La nueva permite autenticarse; la sesión anterior no permanece activa |
| SEC-PE-03 | Registro adicional | Enlace de invitación válido | Registro directo sin invitación | Con invitación se crea la cuenta; sin invitación aparece `Invalid registration link` |

La implementación actual cubre directamente las clases de sesión y registro. El spec confirma el uso de la contraseña nueva, pero no intenta iniciar sesión con la contraseña anterior; por ello no se contabiliza esa comprobación negativa como cubierta.

### 5.3 Tabla de decisión del registro

#### Condiciones

- **C1:** el servidor todavía permite crear la cuenta inicial.
- **C2:** existe una invitación válida.
- **C3:** los datos requeridos son válidos.

#### Acciones

- **A1:** crear la cuenta y mostrar el workspace.
- **A2:** rechazar el registro con `Invalid registration link`.

| Reglas | R1 | R2 | R3 |
|---|:---:|:---:|:---:|
| C1: cuenta inicial permitida | V | F | F |
| C2: invitación válida | - | F | V |
| C3: datos válidos | V | V | V |
| A1: registro exitoso | V | F | V |
| A2: registro rechazado | F | V | F |
| Cobertura en `loginActions.ts` | Sí | Sí | Sí |

El símbolo `-` significa que la condición no modifica el resultado de esa regla.

### 5.4 Transición de estados de sesión

#### Estados

- **S0: No autenticada:** se muestra la página de login y no existe acceso al workspace.
- **S1: Autenticada:** existe una sesión válida y el workspace está visible.
- **S2: Contraseña actualizada:** la cuenta conserva acceso mediante la contraseña nueva.

#### Eventos

- **E1:** registrar o iniciar sesión correctamente.
- **E2:** cerrar sesión.
- **E3:** cambiar contraseña.
- **E4:** iniciar sesión con la contraseña nueva.

| Estado actual | E1: registrar/login | E2: logout | E3: cambiar contraseña | E4: login nueva contraseña |
|---|:---:|:---:|:---:|:---:|
| S0: No autenticada | S1 | - | - | S1 |
| S1: Autenticada | - | S0 | S2 | - |
| S2: Contraseña actualizada | - | S0 | - | S1 |

#### Secuencia automatizada

```text
S0 -> registro -> S1 -> logout -> S0
S0 -> login -> S1 -> cambio de contraseña -> S2
S2 -> logout -> S0 -> login con contraseña nueva -> S1
S1 -> logout -> S0
```

#### Cobertura

- Estados visitados: 3 de 3 = **100%**.
- Transiciones únicas definidas y ejecutadas: 5 de 5 = **100%**.
- Eventos de transición ejecutados en la secuencia completa: 7 (incluye repeticiones de login y logout).

La cobertura anterior corresponde exclusivamente al modelo reducido documentado y no a todos los estados de seguridad posibles de Focalboard.

### 5.5 Pasos y resultados esperados

| Paso | Acción E2E | Resultado esperado |
|---:|---|---|
| 1 | Visitar `/` sin sesión | Redirección a `/login` |
| 2 | Registrar la primera cuenta | Workspace y barra lateral visibles |
| 3 | Cerrar sesión y volver a `/` | Se mantiene la redirección al login |
| 4 | Iniciar sesión | Workspace visible |
| 5 | Cambiar contraseña | Mensaje de éxito y continuidad de la cuenta |
| 6 | Iniciar sesión con la contraseña nueva | Acceso concedido |
| 7 | Intentar registrar un segundo usuario sin invitación | Registro rechazado |
| 8 | Copiar una invitación y registrar otro usuario | Registro exitoso y workspace visible |

### 5.6 Métricas y aceptación

| Métrica | Fórmula o medición | Criterio |
|---|---|---|
| Controles de seguridad aprobados | verificaciones aprobadas / verificaciones ejecutadas × 100 | 100% |
| Transiciones cubiertas | transiciones ejecutadas / transiciones definidas × 100 | 100% del modelo reducido |
| Errores Cypress | Conteo del reporte | 0 |
| Timeouts | Conteo del reporte | 0 |

## 6. SYS-USA-01: usabilidad del flujo de tablero y tarjeta

### 6.1 Ficha técnica

| Campo | Valor |
|---|---|
| ID | SYS-USA-01 |
| Atributo | Usabilidad: efectividad y eficiencia |
| Funcionalidad | Crear y gestionar un tablero, tarjetas y vistas |
| Implementación | `webapp/cypress/integration/createBoard.ts` |
| Caso Cypress principal | `Can create and delete a board and a card` |
| Precondiciones | Servidor inicializado; usuario de prueba autenticado; tableros reiniciados; tour omitido |
| Datos principales | Títulos únicos de tablero, tarjeta y vistas construidos con timestamp |
| Técnicas | Partición de equivalencia y transición de estados |
| Prioridad | Alta |

### 6.2 Partición de equivalencia

| Código | Elemento | Clase válida cubierta | Resultado esperado |
|---|---|---|---|
| USA-PE-01 | Tipo de tablero | Tablero vacío | Se crea y aparece el componente del tablero |
| USA-PE-02 | Título | Texto no vacío y único | El título escrito queda visible |
| USA-PE-03 | Tipo de vista | Kanban y Tabla | La tarjeta puede observarse en ambas representaciones |
| USA-PE-04 | Estado de barra lateral | Visible y oculta | El usuario puede ocultarla y volver a mostrarla |

Las clases inválidas de títulos vacíos, excesivamente largos o con caracteres no admitidos no se ejecutan en el spec y, por tanto, quedan fuera de la cobertura declarada.

### 6.3 Transición de estados del tablero

El ejemplo genérico de eliminar y restaurar una tarjeta no representa la implementación seleccionada: este spec no restaura tarjetas. El modelo se adapta a las transiciones que realmente ejecuta Cypress.

#### Estados

- **B0: Inexistente:** el tablero de prueba todavía no existe.
- **B1: Activo sin personalizar:** se creó un tablero vacío.
- **B2: Activo personalizado:** el tablero tiene título, vistas y tarjetas creadas por el usuario.
- **B3: Eliminado:** el tablero desaparece de la barra lateral y deja de estar disponible en el flujo activo.

#### Eventos

- **E1:** crear tablero vacío.
- **E2:** editar título y agregar tarjeta/vista.
- **E3:** eliminar tablero y confirmar la operación.

| Estado actual | E1: crear | E2: personalizar | E3: eliminar |
|---|:---:|:---:|:---:|
| B0: Inexistente | B1 | - | - |
| B1: Activo sin personalizar | - | B2 | B3 |
| B2: Activo personalizado | - | B2 | B3 |
| B3: Eliminado | - | - | - |

#### Secuencia automatizada

```text
B0 -> crear tablero vacío -> B1
B1 -> nombrar, crear tarjeta y agregar vista -> B2
B2 -> eliminar y confirmar -> B3
```

#### Cobertura

- Estados visitados: 4 de 4 = **100%**.
- Transiciones objetivo ejecutadas: 3 de 3 = **100%**.

La cobertura corresponde al flujo seleccionado. Restaurar, archivar y eliminar permanentemente no forman parte del spec.

### 6.4 Pasos y resultados esperados

| Paso | Acción E2E | Evidencia de usabilidad esperada |
|---:|---|---|
| 1 | Crear un tablero vacío desde `+ Add board` | Controles visibles y tablero creado |
| 2 | Escribir el título del tablero | Valor actualizado y visible |
| 3 | Ocultar y volver a mostrar la barra lateral | Estado visual cambia según la acción |
| 4 | Renombrar la vista Kanban | Nombre nuevo visible |
| 5 | Crear una tarjeta | Diálogo visible y foco automático en el título |
| 6 | Escribir el título de la tarjeta | Valor visible en la tarjeta |
| 7 | Crear y renombrar una vista de tabla | Vista disponible y tarjeta visible en la tabla |
| 8 | Ordenar la tabla | Acción completada sin bloqueo |
| 9 | Eliminar el tablero y confirmar | El tablero desaparece de la interfaz |

### 6.5 Métricas de usabilidad

| Métrica | Fórmula o fuente | Criterio de aceptación |
|---|---|---|
| Efectividad de tareas | pasos completados / 9 pasos definidos × 100 | 100% |
| Errores de interacción | aserciones Cypress fallidas | 0 |
| Bloqueos/timeouts | timeouts registrados | 0 |
| Foco correcto | comprobaciones `should('have.focus')` aprobadas | 100% |
| Eficiencia temporal | duración informada por Cypress | Registrar como línea base; la ejecución debe finalizar dentro del timeout de CI |

No se fija retrospectivamente un umbral numérico de segundos porque el spec no fue diseñado con un SLA y la duración depende del runner. Después de reunir varias ejecuciones puede definirse una línea base y un percentil aceptable.

## 7. Catálogo consolidado

| ID | Atributo | Técnica principal | Entrada/acción | Resultado esperado | Automatización |
|---|---|---|---|---|---|
| SYS-SEC-01 | Seguridad | Tabla de decisión y transición de estados | Registro, login, logout, cambio de contraseña e invitación | Accesos válidos permitidos y registro sin invitación rechazado | `loginActions.ts` |
| SYS-USA-01 | Usabilidad | Partición de equivalencia y transición de estados | Crear, personalizar y eliminar tablero desde UI | Tarea completada, foco y retroalimentación correctos, cero bloqueos | `createBoard.ts` |

## 8. Trazabilidad

| Caso de sistema | Requisitos relacionados | Capas recorridas | Evidencia |
|---|---|---|---|
| SYS-SEC-01 | Registro, autenticación, sesión y cambio de contraseña | Navegador, React, API Auth, backend, SQLite | Resultado de `loginActions.ts` en `ci-ubuntu-webapp` |
| SYS-USA-01 | Gestión de tableros, tarjetas y vistas | Navegador, React, API Boards/Blocks, backend, SQLite | Resultado de `createBoard.ts` en `ci-ubuntu-webapp` |

## 9. Evidencias que deben conservarse

Para la presentación final se deben guardar:

1. enlace al run de GitHub Actions;
2. captura del job `ci-ubuntu-webapp` en verde;
3. sección `Run Finished` de Cypress;
4. nombre de cada spec y cantidad de casos aprobados/fallidos;
5. duración reportada por Cypress;
6. commit que contiene las pruebas y sus correcciones;
7. versión de Go, Node.js y Cypress usada por el workflow.

## 10. Riesgos y limitaciones

- Los escenarios se ejecutan con SQLite; no demuestran el mismo recorrido E2E con MySQL, MariaDB o PostgreSQL.
- Cypress utiliza un usuario simulado; no mide percepción o satisfacción humana.
- `loginActions.ts` comprueba controles funcionales, no resistencia frente a ataques.
- `createBoard.ts` mide efectividad y eficiencia técnica, no accesibilidad completa.
- El tiempo de CI varía según carga del runner y descarga de dependencias.
- Una advertencia WebSocket `close 1001 (going away)` al cerrar Cypress representa una desconexión normal, no un fallo del escenario.

## 11. Conclusión

El plan se concentra en dos atributos solicitados: seguridad y usabilidad. Cada atributo se sustenta con una prueba E2E existente, ejecutable y automatizada en GitHub Actions. `loginActions.ts` valida controles de autenticación, sesión, contraseña e invitaciones; `createBoard.ts` valida que una persona pueda completar eficazmente el flujo principal de tableros y tarjetas con retroalimentación y foco adecuados.

Las técnicas de partición de equivalencia, tabla de decisión y transición de estados se aplican solamente a comportamientos que los specs ejecutan realmente. Los valores límite, la restauración de tarjetas, las pruebas de penetración y la satisfacción subjetiva se documentan como fuera de alcance para evitar afirmar una cobertura inexistente.
