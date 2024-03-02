# dev-hack

In this `dev-hack` monorepo we'll demonstrate Golang based microservices following clean architecture for a ecommerce site with proper observability framwork.

## This monorepo's structure

- [gateway](./gateway/README.md) : Gateway service which have all the endpoints related to sign up the user and creating/fetching the products for the respective users, it's gonna use product service to create/fetch the product basically.
- [product](./product/README.md) : product service which have all the endpoints to do CRUD operations for a ecommerce site.
- [schemas](./schemas/README.md) :
  - DB migrations
  - DB schemas
  - SQLC queries
  - code generations

## Tech Stack & Tools

- Golang
- Gin Framework
- Ginkgo
- Gomega
- Golang Mockgen
- PostgreSQL
- MySQL
- MongoDB
- Redis
- SQLC
- Golang Migrate
- Docker
- OpenAPI Schema

## How to run this project locally?

- clone this repo & go to the root directory (`cd ~/<>/dev-hack`)
- `docker compose up`:
  - to run all the services according to `docker-compose.yaml` file use this command
  - if you change anything in the proeject and wanna run with that then use `docker compose up --build`
  - you can run any specific service by `docker compose up <service_name_from_docker_compose_file>
- [optional] you can connect to the DB via any(Like: `table plus`) GUI providing the connection info given in the `docker-compose.yaml` file's `db` service:
  ```
      POSTGRES_USER: devhack
      POSTGRES_PASSWORD: devhack123
      POSTGRES_DB: devhackdb
  ```

### migrate the postgresql and mysql DB

#### Postgresql

- open another terminal to apply the migrations files into DB
- `cd schemas`
- `export POSTGRESQL_URL='postgresql://devhack:devhack123@localhost:5432/devhackdb?sslmode=disable'`
- `migrate -database $POSTGRESQL_URL -path db/postgresql/migrations up` : apply the migrations files

#### Mysql

- `cd schemas`
- `export MYSQL_URL='mysql://devhack:devhack123@localhost:3306/devhackdb?multiStatements=true'`
- `migrate -database ${MYSQL_URL} -path db/mysql/migrations up`

[🎉] now the projects setup is successful, for checking the API endpoints of each service pls have a look into each service's README.md file and test those endpoints from `postman` or etc.

## Gateway Service's RESTful API Endpoints

- check [here](./gateway/openapi.yaml) for OpenAPI schema of gateway service's api endpoints
- also look into the [examples](./gateway/README.md)

## Product Service's RESTful API Endpoints

- check [here](./product/openapi.yaml) for OpenAPI schema of product service's api endpoints
- also look into the [examples](./ecom/README.md)
