- nombre: pruebas_unitarias_server_api.md

- issue: generar las pruebas unitarias del componente api del modulo del servidor

- leer: Vas a realizar las pruebas unitarias para el componente "api" del modulo "server", los archivos del componente estan en la carpeta server/api/, las pruebas prexistentes contienen una cobertura de sentencias del 36.36%, vas a generar test para subir la cobertura de sentencias al 90%.

- prompt:
<primer_prompt>
<contexto>
Estamos realizando el trabajo final del curso de Pruebas de software.
Esta es la carpeta del proyecto sobre el cual estamos realizando las pruebas de software.
</contexto>
Vas a realizar pruebas unitarias sobre los archivos en la carpeta server/api/ (los cuales estan relacionados a los servicios que el modulo de aplicacion web va a llamar) hasta llegar a una cobertura de sentencias al 90%.
Vas a usar gomock y otras herramientas que el proyecto ya este usando para pruebas en el modulo de server (todo por terminal)
Todas las pruebas deben estar en la siguiente direccion server/api/tests/
Luego vas a realizar:
* un script para ejecutar todas las pruebas generadas de forma local (server/api/tests/run_tests_api.sh & server/api/tests/run_tests_api.bat)
* un reporte de las pruebas realizadas (server/api/tests/reporte_pruebas_unitarias_api.md) en el cual se muestre en un cuadro que indica cobertura y pruebas realizadas.
<primer_prompt>

- necesario:
[server/api/tests/]
[server/api/tests/reporte_pruebas_unitarias_api.md]
[server/api/tests/run_tests_api.sh]
[server/api/tests/run_tests_api.bat]
[lepismas/gabriela/explicacion1.md]
