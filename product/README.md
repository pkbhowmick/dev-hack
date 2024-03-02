# Ecom service

## Local URL (when run via `docker-compose`)

- `localhost:8085/`

## API endpoints with examples

- `GET /` (Service's HealthCheck endpoint):
  - Request URL: `localhost:8085/`
  - Response:
    ```
    {
        "ok": true,
        "serviceName": "EcomService"
    }
    ```
- `GET /products/:userId` (Get all the products for that user):
  - Request URL: `localhost:8085/products/<userId>`
  - Response:
    ```
    {
        "ok": true,
        "message": "Products fetching is succeeded",
        "products": [
            {
                "id": "dcab11a7-7c33-4cc8-b177-7e3d00f7af0a",
                "userId": "",
                "name": "demo2",
                "price": 11,
                "description": "demo product2",
                "stock": 3,
                "createdAt": "2023-02-14T07:52:45.809Z",
                "updatedAt": "2023-02-14T07:52:45.809Z"
            }
        ]
    }
    ```
- `POST /products` (Create a product):
  - Request URL: `localhost:8085/products`
  - Request Body:
    ```
    {
        "name": "sample",
        "userId": "",
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
