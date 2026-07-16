* nombre: Redacción de secciones del paper IEEE - Parte 4

* issue: JEFERSON Redacción de secciones del paper IEEE - Parte 4

* leer: 
El paper IEEE en lepismas/paper/main.tex ha sido estructurado con todas las secciones necesarias para el artículo técnico sobre pruebas de software en Focalboard. Ahora debemos llenar el contenido de cada sección con información real del proyecto.

* prompts:
<tarea1>
Vas a redactar el contenido de las siguientes secciones del paper IEEE:
<seccion1>
### VIII. Pruebas de sistema
#### A. Selección de atributos de calidad
#### B. Primer atributo evaluado
#### C. Segundo atributo evaluado
#### D. Métricas y criterios de aceptación
#### E. Resultados
</seccion1>
<seccion2>
### IX. Automatización mediante GitHub Actions
#### A. Workflows implementados
#### B. Pruebas automatizadas
#### C. Cobertura y artefactos
#### D. Resultados del pipeline
</seccion2>
Siguiendo las siguientes instrucciones:
<instrucciones>
- Para cada sección, reemplaza los placeholders con contenido real basado en el repositorio de Focalboard
- Usa la información de .github/workflows/ para la automatización
- Describe los atributos de calidad: funcionalidad, usabilidad, eficiencia, etc.
- Selecciona 2 atributos de calidad para evaluar con detalle
- Describe los workflows: ci.yml, component-tests.yml, codeql-analysis.yml, etc.
- Explica las pruebas automatizadas en el pipeline
- Menciona la generación de artefactos y reportes de cobertura
- Incluye citas bibliográficas apropiadas usando \cite{}
- Mantén el formato LaTeX y los comandos existentes
- No elimines la figura del pipeline, actualiza su caption si es necesario
</instrucciones>
Toma en cuenta el contexto de las pruebas de sistema y la automatización CI/CD
</tarea1>

* necesario:
[lepismas/paper/main.tex] (archivo principal del paper - editar secciones VIII, IX)
[.github/workflows/ci.yml] (workflow principal de CI)
[.github/workflows/component-tests.yml] (pruebas de componentes)
[.github/workflows/codeql-analysis.yml] (análisis de seguridad)
[.github/workflows/lint-server.yml] (linting del servidor)
