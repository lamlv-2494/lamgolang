# Swagger API Documentation

## Overview

This project includes comprehensive Swagger/OpenAPI documentation for the Food Delivery API. The API follows RESTful principles and uses JWT tokens for authentication.

## Accessing Swagger UI

To view the interactive API documentation:

1. **Install Swagger UI** (if not already installed):
   ```bash
   go get github.com/swaggo/swag/cmd/swag
   go get github.com/swaggo/gin-swagger
   go get github.com/swaggo/files
   ```

2. **Generate Swagger files**:
   ```bash
   swag init -g cmd/app/main.go
   ```

3. **Add Swagger to your main.go**:
   ```go
   import (
       "github.com/swaggo/gin-swagger"
       "github.com/swaggo/files"
   )
   
   // In your route setup
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   ```

4. **Access the documentation**:
   - Navigate to `http://localhost:8080/swagger/index.html` in your browser

## API Base URL

- Production: `https://api.example.com/api`
- Development: `http://localhost:8080/api`

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

To get a token:
1. Register a new user: `POST /api/user/register`
2. Login: `POST /api/user/login`
3. Use the returned token in subsequent requests

## API Endpoints Summary

### Auth (Public)
- `POST /user/register` - Register new user
- `POST /user/login` - Login user

### User
- `GET /user` - Get current user profile
- `PUT /user` - Update user profile

### Products
- `GET /products` - List products (with pagination)
- `GET /products/{id}` - Get product details

### Categories
- `GET /categories` - List categories

### Cart (Authenticated)
- `GET /cart` - Get user's cart
- `POST /cart` - Add item to cart
- `PUT /cart/{id}` - Update cart item quantity
- `DELETE /cart/{id}` - Remove item from cart

### Orders (Authenticated)
- `POST /orders` - Create order (checkout)
- `GET /orders` - Get user's order history

### Ratings (Authenticated)
- `POST /products/{id}/rating` - Create rating for product
- `GET /products/{id}/rating` - Get ratings for product

### Suggestions (Authenticated)
- `POST /suggestions` - Create a suggestion

### Admin Routes (Admin Only)

**Users**:
- `GET /admin/users` - Get all users
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user

**Categories**:
- `POST /admin/categories` - Create category
- `PUT /admin/categories/{id}` - Update category
- `DELETE /admin/categories/{id}` - Delete category

**Products**:
- `POST /admin/products` - Create product
- `PUT /admin/products/{id}` - Update product
- `DELETE /admin/products/{id}` - Delete product

**Orders**:
- `GET /admin/orders` - Get all orders
- `PUT /admin/orders/{id}/status` - Update order status

**Suggestions**:
- `GET /admin/suggestions` - Get all suggestions

## Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "username": "john_doe",
      "email": "john@example.com",
      "password": "password123"
    }
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Products
```bash
curl -X GET "http://localhost:8080/api/products?page=1&limit=10" \
  -H "Content-Type: application/json"
```

### Add to Cart (Authenticated)
```bash
curl -X POST http://localhost:8080/api/cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

### Checkout
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "delivery_address": "123 Main St, City, Country",
    "payment_method": "credit_card"
  }'
```

## Response Format

All API responses follow a consistent format:

**Success Response**:
```json
{
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2026-06-01T10:00:00Z",
    "updated_at": "2026-06-01T10:00:00Z"
  }
}
```

**Error Response**:
```json
{
  "error": "Invalid credentials"
}
```

## Pagination

List endpoints support pagination with these query parameters:
- `page` (default: 1) - Page number
- `limit` (default: 10) - Items per page

Example:
```
GET /api/products?page=2&limit=20
```

## Status Codes

- `200 OK` - Successful GET, PUT requests
- `201 Created` - Successful POST requests
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - User lacks required permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Rate Limiting

Currently, no rate limiting is implemented. This should be added for production deployments.

## Support

For API issues or questions, please refer to the handler implementations in `internal/handlers/` or contact the development team.
