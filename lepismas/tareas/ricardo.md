- nombre: pruebas_unitarias_server_model.md

- issue: generar las pruebas unitarias del componente model del modulo del servidor

- leer: Vas a realizar las pruebas unitarias para el componente "model" del modulo "server", los archivos del componente estan en la carpeta server/model, las pruebas prexistentes contienen una cobertura de sentencias del 23%, vas a generar test para subir la cobertura de sentencias al 90%.

- prompt:
<primer_prompt>
<contexto>
Estamos realizando el trabajo final del curso de Pruebas de software.
Esta es la carpeta del proyecto sobre el cual estamos realizando las pruebas de software.
</contexto>
Vas a realizar pruebas unitarias sobre los archivos en la carpeta server/model/ (los cuales estan relacionados a las declaraciones de bloques, tableros y campos) hasta llegar a una cobertura de sentencias al 90%.
Vas a usar gomock y otras herramientas que el proyecto ya este usando para pruebas en el modulo de webapp (todo por terminal)
Todas las pruebas deben estar en la siguiente direccion server/model/tests/
Luego vas a realizar:
* un script para ejecutar todas las pruebas generadas de forma local (server/model/tests/run_tests_model.sh & server/model/tests/run_tests_model.bat)
* un reporte de las pruebas realizadas (server/model/tests/reporte_pruebas_unitarias_model.md) en el cual se muestre en un cuadro que indica cobertura y pruebas realizadas.
<primer_prompt>

- necesario:
[server/model/tests/]
[server/model/tests/reporte_pruebas_unitarias_model.md]
[server/model/tests/run_tests_model.sh]
[server/model/tests/run_tests_model.bat]
[lepismas/ricardo/explicacion1.md]
