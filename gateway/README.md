# Auth Service

## Local URL (when run via `docker-compose`)

- `localhost:8084/`

## API Endpoints with examples also with Observability details

- `GET /` (Service's HealthCheck endpoint):

  - Request URL: `localhost:8084/`
  - Response:
    ```
    {
        "ok": true,
        "serviceName": "GatewayService"
    }
    ```
  - Observability:
    - it doesn't have any tracing or etc. because it's not necessary and taking decision which things to trace is most important

- `POST /signup` (SignUp endpoint):

  - Request URL: `localhost:8084/signup`
  - Request Body:
    ```
    {
        "name": "demo",
        "email": "demo@gmail.com",
        "password": "demopass"
    }
    ```
  - Response:
    ```
    {
        "ok": true,
        "message": "Successfully Signed Up!"
    }
    ```
  - Observability:
    - It stores the successful user creation to mongodb for future analytics purpose
    - we added tracing to all the important componenet
    - it uses mysql db for user data storing, as we assumed our userbase is gonna be small but the read operation (to verify whether there's a user or not) is more thus we used mysql for optimal read operation heavy db

- `GET /users/:userId/products` (Get the products of that user):
  - Request URL: `localhost:8085/users/640cebcb-cd27-4dff-86cb-c951e2d65828/products`
  - Response:
    ```
    {
        "ok": true,
        "message": "User is fetched",
        "user": {
            "id": "640cebcb-cd27-4dff-86cb-c951e2d65828",
            "name": "demo",
            "email": "demo@gmail.com",
            "password": "$2a$10$lMRv0ZvYVCmhac4DkiLqQuItC1XEmqclPyEhLavVHaai9XJ0QcG6G",
            "createdAt": "2023-02-18T20:21:40.404Z",
            "updatedAt": "2023-02-18T20:21:40.404Z"
        }
    }
    ```
  - Observability:
    - It usages postgresql db as it can have both write heavy operations
    - Here we used distributed tracing by context propagation and called another different service
- `POST /users/:userId/products` (Create a product by that user):
  - Request URL: `localhost:8085/users/:userId/products`
  - Request Body:
    ```
    {
        "name": "sample",
        "price": 100,
        "description": "sample product",
        "stock": 1
    }
    ```
  - Response:
    ```
    {
        "ok": true,
        "message": "Product is created successfully",
        "productID": "9df4ae40-9cf7-4acd-8d2b-df8be230eaf7"
    }
    ```
