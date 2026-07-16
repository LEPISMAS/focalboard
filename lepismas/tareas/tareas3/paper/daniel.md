* nombre: Redacción de secciones del paper IEEE - Parte 2

* issue: DANIEL Redacción de secciones del paper IEEE - Parte 2

* leer: 
El paper IEEE en lepismas/paper/main.tex ha sido estructurado con todas las secciones necesarias para el artículo técnico sobre pruebas de software en Focalboard. Ahora debemos llenar el contenido de cada sección con información real del proyecto.

* prompts:
<tarea1>
Vas a redactar el contenido de las siguientes secciones del paper IEEE:
<seccion1>
### IV. Estrategia y planificación de pruebas
#### A. Alcance
#### B. Plan de pruebas unitarias
#### C. Plan de pruebas funcionales
#### D. Plan de pruebas de integración
#### E. Entorno y herramientas
</seccion1>
<seccion2>
### V. Pruebas unitarias y cobertura
#### A. Backend Go
#### B. Frontend TypeScript
#### C. Criterios de cobertura
#### D. Resultados de cobertura
</seccion2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Para cada sección, reemplaza los placeholders con contenido real basado en el repositorio de Focalboard
- Usa la información del Makefile para entender los comandos de testing
- Usa la información de webapp/package.json para las pruebas del frontend
- Usa la información de server/go.mod para las pruebas del backend
- Incluye los resultados reales de cobertura si están disponibles
- Mantén el formato LaTeX y los comandos existentes
- No elimines la tabla de cobertura, actualiza su contenido con datos reales
- Incluye citas bibliográficas apropiadas usando \cite{}
- Describe los frameworks utilizados: Go testing, Jest, Cypress
- Menciona las bases de datos de prueba: SQLite, MySQL, MariaDB, PostgreSQL
</instrucciones>
Toma en cuenta el contexto de las pruebas unitarias realizadas en el proyecto
</tarea1>

* necesario:
[lepismas/paper/main.tex] (archivo principal del paper - editar secciones IV, V)
[Makefile] (comandos de testing)
[webapp/package.json] (scripts de testing del frontend)
[server/go.mod] (dependencias del backend)
[.github/workflows/ci.yml] (pipeline de CI/CD)
