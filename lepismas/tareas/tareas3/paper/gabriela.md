* nombre: Redacción de secciones del paper IEEE - Parte 3

* issue: GABRIELA Redacción de secciones del paper IEEE - Parte 3

* leer: 
El paper IEEE en lepismas/paper/main.tex ha sido estructurado con todas las secciones necesarias para el artículo técnico sobre pruebas de software en Focalboard. Ahora debemos llenar el contenido de cada sección con información real del proyecto.

* prompts:
<tarea1>
Vas a redactar el contenido de las siguientes secciones del paper IEEE:
<seccion1>
### VI. Pruebas funcionales
#### A. Técnicas de caja negra
#### B. Diseño de casos de prueba
#### C. Ejecución manual
#### D. Resultados
</seccion1>
<seccion2>
### VII. Pruebas de integración
#### A. APIs críticas
#### B. Flujos de integración
#### C. Casos implementados
#### D. Resultados
</seccion2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Para cada sección, reemplaza los placeholders con contenido real basado en el repositorio de Focalboard
- Usa la información de webapp/cypress/ para las pruebas funcionales
- Usa la información de server/integrationtests/ para las pruebas de integración
- Usa los documentos en lepismas/docs/reportes/ para los argumentos y resultados
- Describe las técnicas de caja negra: partición de equivalencia, valores límite, etc.
- Lista los flujos de integración implementados: autenticación, tableros, tarjetas, permisos, etc.
- Actualiza la tabla de casos funcionales con datos reales
- Actualiza la tabla de casos de integración con datos reales
- Incluye citas bibliográficas apropiadas usando \cite{}
- Mantén el formato LaTeX y los comandos existentes
</instrucciones>
Toma en cuenta el contexto de las pruebas funcionales y de integración realizadas
</tarea1>

* necesario:
[lepismas/paper/main.tex] (archivo principal del paper - editar secciones VI, VII)
[webapp/cypress/] (pruebas funcionales con Cypress)
[server/integrationtests/] (pruebas de integración)
[lepismas/docs/reportes/] (documentación de pruebas de integración)
[webapp/cypress.json] (configuración de Cypress)
