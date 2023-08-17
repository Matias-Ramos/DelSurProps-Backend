# Nombre del proyecto
DelSurProps - backend

## Sobre el proyecto

### Motivo
La estructura de los datos requeridos por el proyecto generaron la necesidad de una base de datos SQL. Por lo tanto, opté implementar una API que comunique el front end hecho en React con una base de datos hecha en PostgreSQL.

### Link
delsurprops.com.ar

### Hosting
donweb.com

## Código

### Estructura del programa
- func main()
  - getDBdata()
    - generateSQLquery()
    - initBuildingType() 

### Lenguaje
Golang

### Endpoints
- http://localhost:8080/api/venta_inmuebles
- http://localhost:8080/api/alquiler_inmuebles
- http://localhost:8080/api/emprendimientos

### Queries
- ?location=string
- ?price_init=integer / ?price_limit=integer
- ?env_init=integer / ?env_limit=integer
- ?bedroom_init=integer / ?bedroom_limit=integer
- ?bathroom_init=integer / ?bathroom_limit=integer
- ?garage_init=integer / ?garage_limit=integer
- ?total_surface_init=integer / ?total_surface_limit=integer
- ?covered_surface_init=integer / ?covered_surface_limit=integer
- ?building_status= in_progress-or-pozo/pozo/in_progress

### Ejecución del proyecto
1. Crear archivo .env en el directorio root del proyecto con el siguiente contenido:
USER="new_user"
PWD="secure777"
DB_NAME="Inmobiliaria_BD"
2. Ejecutar "go run main.go models.go" en la terminal del IDE. Acto siguiente, generar una http request al host "http://127.0.0.1:8080/api/{category}" con cURL, Postman, Insomnia (u otro) del tipo "get" a los endpoints detallados en el punto anterior.