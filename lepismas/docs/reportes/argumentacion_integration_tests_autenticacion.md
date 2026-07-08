# Argumentacion de pruebas de integracion: Autenticacion

## 1. Objetivo del flujo

El flujo INT-01 tiene como objetivo comprobar que la autenticacion de Focalboard funciona correctamente cuando se ejecuta como una cadena integrada y no como piezas aisladas. El interes principal no es validar solamente que una funcion acepte parametros validos, sino demostrar que una solicitud HTTP real llega a la capa API, invoca la logica de aplicacion, utiliza el servicio de autenticacion y termina persistiendo o consultando datos en el Store.

Este flujo cubre operaciones criticas para la seguridad y continuidad del sistema: registro de usuarios, inicio de sesion, validacion de tokens, rechazo de accesos no autenticados, cambio de contrasena, cierre de sesion y rechazo de registros duplicados.

## 2. Capas integradas

Las pruebas integran las siguientes capas del backend:

- API REST: endpoints como `/api/v2/register`, `/api/v2/login`, `/api/v2/logout`, `/api/v2/users/me` y `/api/v2/users/{userID}/changepassword`.
- Capa App: metodos como `RegisterUser`, `Login`, `GetSession`, `ChangePassword` y `Logout`.
- Servicio Auth: validacion de sesiones, tokens y credenciales.
- Store: persistencia de usuarios y sesiones mediante los metodos publicos disponibles.
- Modelo: estructuras `RegisterRequest`, `LoginRequest`, `LoginResponse`, `ChangePasswordRequest`, `User` y `Session`.
- Cliente de pruebas: cliente HTTP existente del proyecto, que reproduce el consumo real de la API.

La decision de integrar estas capas es relevante porque los defectos de autenticacion suelen aparecer en los bordes: serializacion JSON, cabeceras HTTP, tokens mal propagados, errores de permisos, sesiones no eliminadas o validaciones inconsistentes entre API y Store.

## 3. Estandares aplicados

Las pruebas siguen criterios compatibles con buenas practicas de pruebas de software:

- Pruebas de integracion por flujo: cada caso verifica la comunicacion entre varias capas y no solo una unidad funcional.
- AAA: cada prueba mantiene una estructura clara de preparacion, accion y verificacion.
- Independencia de casos: cada prueba levanta su propio helper y base de prueba, evitando dependencia entre ejecuciones.
- Trazabilidad: los nombres `TestINT0101...` a `TestINT0107...` conservan la relacion directa con los casos definidos.
- Verificacion observable: se comprueban respuestas HTTP, objetos devueltos y persistencia en Store/App cuando el helper lo permite.
- Uso de aserciones explicitas: `require` detiene la prueba en el primer fallo critico y evita falsos positivos.
- Evidencia en salida estandar: cada prueba registra con `t.Log` la utilidad del caso y las capas verificadas.

## 4. Importancia desde risk-based testing

### Identificacion de riesgos

Riesgos de producto:

- Acceso no autorizado a endpoints protegidos.
- Creacion de sesiones invalidas o no persistidas.
- Tokens activos despues de cerrar sesion.
- Usuarios duplicados que rompan unicidad o consistencia.
- Cambio de contrasena que no actualice realmente la credencial persistida.
- Diferencias entre respuesta HTTP exitosa y estado real del Store.

Riesgos de proyecto:

- Confianza excesiva en pruebas unitarias que no cubren comunicacion entre capas.
- Cambios futuros en API, App o Store que mantengan pruebas unitarias verdes, pero rompan el flujo real.
- Dificultad para defender seguridad de autenticacion si no existe evidencia de integracion.

### Evaluacion probabilidad por impacto

| Caso | Riesgo principal | Probabilidad | Impacto | Nivel |
|---|---|---:|---:|---|
| INT-01-01 | Registro no persistido o mal validado | Media | Alta | Alto |
| INT-01-02 | Login sin sesion valida | Media | Alta | Alto |
| INT-01-03 | Token valido no reconocido por endpoint protegido | Media | Alta | Alto |
| INT-01-04 | Endpoint protegido acepta peticiones anonimas | Baja-Media | Critico | Alto |
| INT-01-05 | Cambio de contrasena inconsistente | Media | Media-Alta | Medio-Alto |
| INT-01-06 | Logout no invalida token | Media | Critico | Critico |
| INT-01-07 | Registro duplicado corrompe datos | Media | Media | Medio |

### Priorizacion

Los casos de mayor prioridad son INT-01-02, INT-01-03, INT-01-04 e INT-01-06, porque comprometen directamente el control de acceso. INT-01-01 tambien es de prioridad alta, ya que el registro correcto es la entrada al sistema. INT-01-05 e INT-01-07 se priorizan como controles de consistencia: no siempre bloquean el uso inmediato, pero pueden generar riesgos de seguridad y datos.

### Mitigacion

La mitigacion se logra probando el recorrido completo de autenticacion con cliente HTTP real y Store real de pruebas. De esta forma se cubren puntos que una prueba unitaria aislada no detecta: formato de request, propagacion de token, persistencia de sesion, eliminacion de sesion y respuestas HTTP esperadas.

## 5. Justificacion de herramientas

| Herramienta | Por que se eligio | Alternativas consideradas | Limitacion conocida |
|---|---|---|---|
| Go test | Es el runner nativo del proyecto backend y permite ejecutar pruebas integradas junto al codigo Go existente. | Ejecutores externos como Postman/Newman o scripts manuales. | Requiere que el entorno local tenga configurada correctamente la base de datos de prueba y dependencias Go. |
| Testify require | Ya es usado por el proyecto y permite aserciones claras, deteniendo el caso ante fallos criticos. | `testing` puro o `assert`. | Detiene la prueba en el primer fallo, por lo que puede ocultar fallos secundarios dentro del mismo caso. |
| TestHelper de integrationtests | Reutiliza el servidor, cliente HTTP y ciclo de vida de pruebas ya aceptado por el repositorio. | Crear un servidor manual por cada prueba. | Depende de la configuracion interna del helper y de la base de datos temporal. |
| Cliente HTTP del proyecto | Ejecuta endpoints reales y respeta cabeceras como `Authorization` y `X-Requested-With`. | Llamar directamente metodos de App. | No expone todas las verificaciones internas, por eso algunas validaciones se complementan con App/Store. |
| Store/App publicos | Permiten verificar persistencia y sesiones sin SQL directo. | Consultas SQL manuales. | Solo se puede verificar lo que las interfaces publicas exponen. |

## 6. Justificacion estrategica por caso

### INT-01-01 Registrar usuario nuevo

Este caso valida que el registro no sea un exito superficial de API. La prueba confirma que la respuesta HTTP correcta corresponde a un usuario realmente persistido. Es estrategica porque el registro es la puerta de entrada del sistema y una falla aqui afecta todos los flujos posteriores.

### INT-01-02 Login genera sesion valida

El inicio de sesion no se considera correcto solo por devolver un token. La prueba verifica que el token exista y que la sesion asociada sea recuperable desde la capa App/Store. Esto reduce el riesgo de tokens huerfanos, sesiones no persistidas o inconsistencias entre credenciales y estado del servidor.

### INT-01-03 Endpoint protegido con token valido

Este caso prueba la autorizacion positiva. Su valor esta en confirmar que un token emitido por `/login` es aceptado por un endpoint protegido real. Asi se cubre la continuidad entre generacion de token, almacenamiento de sesion y validacion posterior.

### INT-01-04 Endpoint protegido sin token

La prueba valida el control negativo mas importante: una peticion sin sesion no debe acceder a datos protegidos. Es un caso pequeno, pero de alto valor, porque una regresion aqui tendria impacto critico sobre confidencialidad.

### INT-01-05 Cambiar contrasena

El cambio de contrasena se verifica no solo por respuesta exitosa, sino por la posibilidad de iniciar sesion con la nueva credencial. Esto prueba que el cambio llego hasta Store y que la autenticacion posterior usa el dato actualizado.

### INT-01-06 Logout invalida token

Cerrar sesion debe destruir la capacidad de reutilizar el token anterior. Este caso tiene prioridad critica porque los tokens persistentes despues de logout representan riesgo de secuestro de sesion o acceso no autorizado.

### INT-01-07 Registro duplicado rechazado

El registro duplicado se prueba para proteger consistencia e identidad. La verificacion de conteo de usuarios o consulta al Store confirma que el rechazo no solo ocurre en la respuesta HTTP, sino tambien en el estado persistido.

## 7. Relacion con pruebas unitarias previas del store Redux

Las pruebas unitarias previas del store Redux validaban unidades aisladas del frontend: reducers, acciones, selectores o transformaciones de estado en memoria. Su objetivo era comprobar que una pieza del cliente actualiza correctamente el estado bajo condiciones controladas.

Estas pruebas de integracion tienen otro alcance. No sustituyen a las unitarias; las complementan. Mientras las unitarias responden si una unidad aislada se comporta correctamente, las pruebas INT-01 responden si varias capas del backend se comunican correctamente bajo un flujo real de autenticacion. La relacion defendible es que primero se asegura la correccion local de unidades y luego se verifica la colaboracion entre API, App, Auth y Store.

## 8. Conclusion argumentativa

El flujo INT-01 fue seleccionado por su riesgo funcional y de seguridad. La autenticacion es una zona critica: si falla, el sistema puede negar acceso legitimo, permitir acceso indebido o conservar sesiones que deberian estar invalidadas. Por ello, las pruebas no se limitan a llamadas internas, sino que ejercitan endpoints reales y verifican efectos persistidos.

La estrategia es defendible porque combina trazabilidad, priorizacion por riesgo, reutilizacion de herramientas existentes y validacion entre capas. En conjunto, estas pruebas aportan evidencia de que el flujo de autenticacion funciona como contrato integrado del sistema, no solo como suma de funciones unitarias.
