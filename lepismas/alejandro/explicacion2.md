# Explicacion 2: pruebas de integracion

## Guion breve para exposicion

En esta parte del trabajo implementamos pruebas de integracion para dos flujos criticos de Focalboard: Autenticacion y Gestion de Tableros. La finalidad fue pasar de pruebas unitarias, que validan piezas aisladas, a pruebas que recorren varias capas del sistema y permiten comprobar si esas capas se comunican correctamente.

La tarea se implemento en el backend Go del proyecto, dentro de `server/integrationtests`. Se crearon dos archivos principales de pruebas: `flujo_autenticacion_int_test.go` para el flujo INT-01 y `flujo_gestionDeTableros_int_test.go` para el flujo INT-02. Tambien se crearon scripts `.bat` y `.sh` para ejecutar cada flujo por separado usando `go test` con filtros por nombre de prueba.

## Arquitectura recorrida por las pruebas

Las pruebas recorren la arquitectura real del servidor. En lugar de llamar funciones aisladas sin contexto, se usan helpers existentes del proyecto para levantar un servidor de prueba, crear clientes HTTP y ejecutar operaciones contra endpoints reales.

En el flujo de autenticacion se recorren estas capas:

- API REST, con endpoints como registro, login, logout, usuario actual y cambio de contrasena.
- Capa App, donde se ejecutan reglas como registrar usuario, iniciar sesion, cambiar contrasena y cerrar sesion.
- Servicio de autenticacion, encargado de validar tokens y sesiones.
- Store, donde se persisten usuarios y sesiones.

En el flujo de gestion de tableros se recorren:

- API REST de tableros, miembros, bloques y duplicacion.
- Capa App, donde se aplican reglas de negocio.
- Store, donde se persisten tableros, membresias, historial y bloques.
- Sistema de permisos, especialmente importante para tableros privados.
- El componente WebSocket se reconoce como parte de la arquitectura, aunque la prueba end-to-end de suscripcion no se implemento porque no existe un helper estable de WebSocket dentro de `server/integrationtests`.

## Que valida cada flujo

El flujo INT-01 Autenticacion valida siete casos. Primero verifica que un usuario nuevo pueda registrarse y quede persistido. Luego comprueba que el login genere un token y una sesion valida. Despues valida que un endpoint protegido responda correctamente con token valido y rechace una peticion sin token. Tambien se prueba el cambio de contrasena, asegurando que el usuario pueda volver a iniciar sesion con la nueva clave. Finalmente se valida que el logout invalide el token anterior y que un registro duplicado sea rechazado sin duplicar datos en el Store.

El flujo INT-02 Gestion de Tableros valida el ciclo de vida principal de un tablero. Se prueba la creacion y persistencia del tablero, la creacion automatica de membresia administradora para el creador, el filtrado de tableros privados por membresia, la actualizacion de titulo, la eliminacion logica mediante `deleteAt`, la duplicacion con bloques y propiedades, y el tramo API-App-Store asociado a la creacion que posteriormente deberia emitir notificaciones.

## Por que son pruebas de integracion

Son pruebas de integracion porque no se limitan a verificar una funcion aislada. Cada caso prueba la colaboracion entre varias capas. Por ejemplo, en autenticacion, un login exitoso no solo significa que una funcion compare contrasenas, sino que la API reciba el request, App procese la regla, Auth genere la sesion y Store la persista. En tableros ocurre algo similar: crear un tablero implica API, validacion de permisos, App, Store y creacion de membresia.

Esa diferencia es importante para el debate. Una prueba unitaria podria decirnos que un metodo funciona correctamente con mocks, pero una prueba de integracion nos dice si el flujo completo mantiene su contrato cuando las partes reales interactuan.

## Relacion con pruebas unitarias anteriores

Las pruebas unitarias anteriores, especialmente las relacionadas con el store Redux del frontend, validaban unidades aisladas: reducers, acciones, selectores o transformaciones de estado. Esas pruebas siguen siendo necesarias porque reducen errores locales del cliente.

Sin embargo, las pruebas de integracion tienen otro proposito. Buscan validar la comunicacion entre capas. En otras palabras, las unitarias responden: "esta pieza funciona sola?". Las de integracion responden: "estas piezas funcionan correctamente juntas?". Ambas estrategias se complementan. Primero se controla el comportamiento individual y luego se comprueba el comportamiento integrado.

## Dificultades encontradas

La principal dificultad no estuvo en la logica de los casos, sino en el entorno local de ejecucion. En Windows aparecio un problema de acceso denegado a la cache global de Go en `AppData\Local\go-build`. Ademas, al intentar usar una cache local y SQLite, se detecto que las migraciones del proyecto requieren soporte JSON en SQLite, como la funcion `json_set`, que no estaba disponible en esa configuracion local.

Por esa razon, los reportes separan claramente dos ideas: la implementacion de las pruebas esta realizada y organizada segun el proyecto, pero la ejecucion completa queda condicionada por dependencias del entorno. Esta separacion es importante para no confundir un defecto del producto con un problema de infraestructura.

Otra dificultad fue el caso WebSocket. Focalboard tiene pruebas del servidor WebSocket en otro paquete, pero no hay un helper estable de WebSocket dentro de `server/integrationtests`. Por eso se decidio no inventar un cliente fragil. Se verifico el tramo API-App-Store y se documento la limitacion tecnica para una mejora futura.

## Conclusiones

La tarea fortalece la estrategia de pruebas del proyecto porque agrega cobertura sobre flujos reales y de alto riesgo. En autenticacion se cubren riesgos de seguridad como sesiones invalidas, acceso anonimo y tokens reutilizados despues del logout. En tableros se cubren riesgos de integridad y confidencialidad, como tableros sin propietario, filtrado incorrecto de tableros privados y duplicaciones incompletas.

La conclusion principal es que estas pruebas permiten defender con mayor solidez el comportamiento del sistema, porque validan comunicacion entre capas y no solo unidades aisladas. Tambien muestran una practica importante en pruebas profesionales: documentar las limitaciones del entorno y diferenciar problemas de infraestructura de defectos funcionales del producto.
