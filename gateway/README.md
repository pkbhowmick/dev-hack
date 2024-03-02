# Auth Service

## Local URL (when run via `docker-compose`)

- `localhost:8084/`

## API Endpoints with examples

- `GET /` (Service's HealthCheck endpoint):
  - Request URL: `localhost:8084/`
  - Response:
    ```
    {
        "ok": true,
        "serviceName": "GatewayService"
    }
    ```
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
