# Project Name
DelSurProps - backend

## About the project

### What it is 
An API that connects the React frontend with a PostgreSQL database.

### What it does
1. GET: Filters real state objects based on the input the API receives and returns the output back to the frontend.
2. DELETE - POST: Deletes and Creates real state objects based on the frontend input.

### Link
- Front end demo: [https://endearing-melomakarona-8df72e.netlify.app](https://endearing-melomakarona-8df72e.netlify.app)
- Video of the working app w/ backend and frontend: [https://youtu.be/kzu_LFJki7s](https://youtu.be/kzu_LFJki7s)
- delsurprops.com.ar (site shut down as per client request)

### Front end
[https://github.com/Matias-Ramos/Inmobiliaria](https://github.com/Matias-Ramos/DelSurProps-FrontEnd)

### Full stack diagram
![Full stack architecture diagram](architecture.png)

### Challenges on this project:
- First time writting a Golang API.
- Learnt and implemented JWT for session authentication.
- SQL query generation on the backend.
- REST API architecture.


## Code

### Language
Golang

### Program structure
I won't go through all the program structure, but only through the three files that have the biggest complexity, as I find the overall perspective quite helpful to understand the logic in them. Remember that further details are written as comments in the respective files:
#### post.go
- PostData()
  - generateInsertQuery()
    - convertToLowerCase()
    - parseToSqlSyntax()

#### get.go
- GetDBdata()
  - generateGetQuery()
    - handleLocationField() - due to query syntax uniqueness, a specific method was required.
    - handleBuildingStatusField() - due to query syntax uniqueness, a specific method was required.
    - handleCommonField()
      - includeNullValues()
  - initBuildingType()

#### auth

![auth_diagram](auth_diagram.png)



### Endpoints
#### GET
- http://localhost:8080/api/venta_inmuebles
- http://localhost:8080/api/alquiler_inmuebles
- http://localhost:8080/api/emprendimientos

#### POST
- http://localhost:8080/admin/post/{category}

#### DELETE
- http://localhost:8080/admin/delete/{category}

#### AUTH
- http://localhost:8080/admin/jwt

### Queries
- ?location=string
- ?price_init=integer / ?price_limit=integer
- ?env_init=nullInteger / ?env_limit=nullInteger
- ?bedroom_init=nullInteger / ?bedroom_limit=nullInteger
- ?bathroom_init=nullInteger / ?bathroom_limit=nullInteger
- ?garage_init=nullInteger / ?garage_limit=nullInteger
- ?total_surface_init=nullInteger / ?total_surface_limit=nullInteger
- ?covered_surface_init=nullInteger / ?covered_surface_limit=nullInteger
- ?building_status= in_progress-or-pozo/pozo/in_progress

*nullInteger = accepts both null and integer.

## Running the project

1. Clone the project and run "go run main.go" or "go run ." in the IDE terminal.

### Environment setup
2. Create a .env file in the root directory of the project with the following environment variables:
- PGS_USER="new_user"
- PGS_PWD="secure777"
- PGS_DB_NAME="Inmobiliaria_BD"
- API_KEY="1234"
- JWT_SECRET="super-secret-auth-key"

### Database setup

3. Create a database with a user, password and database name as shown in the previous point. As for the tables:

#### common fields for all tables:
- id - bigint notnull
- location - text notnull
- price - integer notnull
- env - smallint
- bedrooms - smallint
- bathrooms - smallint
- garages - smallint
- images - text[] notnull
- link_ml - text
- link_zonaprop - text
- link_argenprop - text

#### venta_inmuebles:
- currency - text

#### alquiler_inmuebles: 
- total_surface - smallint
- covered_surface - smallint

#### emprendimientos: 
- total_surface - smallint
- covered_surface - smallint
- pozo - boolean notnull
- in_progress - boolean notnull


## POST

4. Hit "localhost:8080/admin/jwt" with a get request. Inside the header, add a key value pair "Access":1234. You will receive a token, save it. 

5. Hit "localhost:8080/admin/alquiler_inmueble" (note that the category is singular, unlike the database name that is plural) with a post request. <br>
Inside the header, add a key value pair "Authentication": jwtToken, being jwtToken the value you saved in the previous step.<br>
Inside the body, this is an example of a valid value to send in JSON format: 
*{"location":"1","price":"2","currency":"usd","env":"3","bedrooms":"4","bathrooms":"5","garages":"6","link_ml":"7","link_zonaprop":"8","link_argenprop":"9","image_links":["10","11"]}*

## DELETE

6. Hit "localhost:8080/admin/alquiler_inmuebles". <br>
Inside the header, add a key value pair "Authentication": jwtToken, being jwtToken the value you saved in the previous step. <br>
Inside the body, this is an example of a valid value to send in JSON format: 
*{"buildingId": 1628372893729}* <br>
Note that in order to know the ID you can either check the database, or use the Get request that will be shown next.

## GET

7. Hit "localhost:8080/admin/alquiler_inmuebles". <br>
No query params will retrieve all the records.  <br>
The possible queries may be find in the previous section "Queries", but here is an example: <br>
/alquiler_inmuebles?location=Villa+Crespo&price_init=10000&price_limit=20000&env_init=1&env_limit=7&bedroom_init=1&bedroom_limit=7&bathroom_init=1&bathroom_limit=7&garage_init=0&garage_limit=7
