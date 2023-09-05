# Project Name
DelSurProps - backend

## About the project

### Purpose
The data structure required by the project led to the usage of an SQL database. Therefore, I chose to implement an API that connects the React frontend with a PostgreSQL database.

### Link
delsurprops.com.ar

### Hosting
donweb.com

## Code

### Program structure
- func main()
TO BE COMPLETED!!!!!!!

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

### Running the project
1. Create a .env file in the root directory of the project with the following content:
PGS_USER="new_user"
PGS_PWD="secure777"
PGS_DB_NAME="Inmobiliaria_BD"
API_KEY="1234"
JWT_SECRET="super-secret-auth-key"
2. Run "go run main.go" in the IDE terminal. Next, make an HTTP request to the host "http://127.0.0.1:8080/api/{category}" using cURL, Postman, Insomnia (or another tool) with a "get" request to the endpoints detailed in the previous section.
