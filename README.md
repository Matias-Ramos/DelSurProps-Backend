# Nombre del proyecto
DelSurProps - backend

## Sobre el proyecto

### Motivo
La estructura de los datos requeridos por el proyecto generaron la necesidad de una base de datos SQL. Por lo tanto, opté implementar una API que comunique el front end hecho en React con una base de datos hecha en PostgreSQL.

### Link
tbd

### Hosting
tbd

## Código

### Estructura del programa
- func main()
  - getDBdata()
    - generateSQLquery()
    - initBuildingType() 

### Lenguaje
Golang

### Endpoints
- http://localhost:8080/venta-inmuebles
- http://localhost:8080/alquiler-inmuebles
- http://localhost:8080/emprendimientos

### Queries
...

### Ejecución del proyecto
Ejecutar "npm run main.go models.go" en la terminal del IDE. Acto siguiente, generar una http request al host "http://127.0.0.1:8080" con cURL, Postman, Insomnia (u otro) del tipo "get" a los endpoints detallados en el punto anterior.
!!!!!!!!!!!!!!!!! -aclarar configuracion del archivo .env  !!!!!!!!!!!!!!!!!