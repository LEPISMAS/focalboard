- nombre: pruebas_unitarias_server_ws.md

- issue: generar las pruebas unitarias del componente ws(websocket) del modulo del servidor

- leer: Vas a realizar las pruebas unitarias para el componente "ws" del modulo "server", los archivos del componente estan en la carpeta server/ws/, las pruebas prexistentes contienen una cobertura de sentencias del 23%, vas a generar test para subir la cobertura de sentencias al 90%.

- prompt:
<primer_prompt>
<contexto>
Estamos realizando el trabajo final del curso de Pruebas de software.
Esta es la carpeta del proyecto sobre el cual estamos realizando las pruebas de software.
</contexto>
Vas a realizar pruebas unitarias sobre los archivos en la carpeta server/ws/ (los cuales estan relacionados a la conexion de websockets para la actualizacion real de tablas colaborativas) hasta llegar a una cobertura de sentencias al 90%.
Vas a usar gomock y otras herramientas que el proyecto ya este usando para pruebas en el modulo de server (todo por terminal)
Todas las pruebas deben estar en la siguiente direccion server/ws/tests/
Luego vas a realizar:
* un script para ejecutar todas las pruebas generadas de forma local (server/ws/tests/run_tests_ws.sh & server/ws/tests/run_tests_ws.bat)
* un reporte de las pruebas realizadas (server/ws/tests/reporte_pruebas_unitarias_ws.md) en el cual se muestre en un cuadro que indica cobertura y pruebas realizadas.
<primer_prompt>

- necesario:
[server/ws/tests/]
[server/ws/tests/reporte_pruebas_unitarias_ws.md]
[server/ws/tests/run_tests_ws.sh]
[server/ws/tests/run_tests_ws.bat]
[lepismas/jeferson/explicacion1.md]
