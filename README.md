# 🍔 Food Delivery - Fullstack Web Application

A comprehensive online Food Delivery platform built with a Client-Server architecture (RESTful API). This project implements a clean Layered Architecture (Repository - Service - Handler) ensuring scalability and maintainability, secured by JWT authentication, and features a lightweight Static Single Page Application (SPA) frontend for optimal performance.

## ✨ Key Features

### 👤 For Customers (Normal Users)
- **Authentication:** Secure Registration and Login (Bcrypt password hashing, JWT tokens).
- **Profile Management:** View and update personal information (username, email, password).
- **Food Discovery:** - Browse food catalog and categories.
  - Advanced filtering by price range, category, minimum rating, and dynamic keyword search.
  - Sorting capabilities (alphabetical, price ascending/descending).
- **Shopping Cart:** Add items to cart, dynamically update quantities, and remove items.
- **Checkout & Orders:** Secure checkout process, view detailed order history, and track real-time order status (Pending, Processing, Delivering, Completed, Cancelled).
- **Reviews & Ratings:** Submit 1-5 star ratings and written feedback for purchased products.
- **Suggestions:** Submit feedback or request new menu items directly to the administration.

### 👑 For Administrators (Admin)
- **User Management:** View the list of registered accounts, update user roles, or delete users.
- **Category Management:** Full CRUD (Create, Read, Update, Delete) operations for food categories.
- **Product Management:** Full CRUD operations for menu items, including **physical image file uploads** stored securely on the server.
- **Order Management:** Monitor all system-wide customer invoices and update order statuses dynamically.
- **Suggestions Inbox:** Read and review customer feedback and feature requests.

---

## 🛠️ Technology Stack

### Backend (Core API)
- **Language:** Golang (Go 1.26+)
- **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
- **ORM:** [GORM v2](https://gorm.io/) (MySQL Driver)
- **Database:** MySQL
- **Authentication:** JSON Web Tokens (JWT) + Bcrypt
- **API Documentation:** Swaggo (Swagger UI)

### Frontend (Client-Side Rendering)
- **Core:** HTML5 / CSS3 / Vanilla JavaScript
- **UI Framework:** Bootstrap 5 (via CDN)
- **API Communication:** Native `fetch()` API with automated JWT Interceptor Pattern.

---

## 📂 Project Structure

The project follows the **Standard Go Project Layout** combined with a static frontend directory.

```text
lamgolang/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── internal/
│   ├── configs/                 # DB connections, JWT config, Error mapping
│   ├── handlers/                # HTTP Controllers (Request/Response parsing)
│   ├── middlewares/             # JWT Auth, Admin RBAC, Logging middlewares
│   ├── models/
│   │   ├── dto/                 # Data Transfer Objects & Generic List Responses
│   │   └── entities/            # GORM Database Models
│   ├── repositories/            # Database abstraction layer (CRUD ops)
│   ├── routes/                  # API routing definitions
│   ├── services/                # Core Business Logic layer
│   └── utils/                   # Shared constants and utilities
├── static/                      # FRONTEND ASSETS (Served globally via Gin)
│   ├── admin/                   # Admin portal HTML views
│   ├── user/                    # Customer HTML views
│   ├── js/                      # Core API engine and Auth scripts
│   └── index.html               # Customer storefront landing page
├── uploads/                     # Server storage for uploaded product images
├── docs/                        # Auto-generated Swagger documentation
├── .env                         # Environment variables configuration
├── go.mod                       # Golang dependencies
└── README.md

## 📡 API Endpoints Summary

Here is a high-level overview of the RESTful API endpoints available in the system. All routes are prefixed with `/api`.

### 🌍 Public Routes (No Auth Required)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/categories` | Get a paginated list of all food categories |
| `GET` | `/products` | Get a paginated list of products (supports search, sort, filters) |
| `GET` | `/products/:id` | Retrieve detailed information of a specific product |
| `POST` | `/user/register` | Register a new customer account |
| `POST` | `/user/login` | Authenticate and receive a JWT token |

### 👤 Customer Routes (Requires Bearer Token)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/user` | Get current user's profile information |
| `PUT` | `/user` | Update current user's profile/password |
| `GET` | `/cart` | View current user's shopping cart |
| `POST` | `/cart` | Add a product to the shopping cart |
| `PUT` | `/cart/:id` | Update quantity of a specific cart item |
| `DELETE` | `/cart/:id` | Remove an item from the cart |
| `POST` | `/orders` | Checkout cart and create a new order |
| `GET` | `/orders` | View user's order history |
| `POST` | `/products/:id/rating`| Submit a rating and review for a product |
| `GET` | `/products/:id/rating` | Get user's existing rating for a specific product |
| `GET` | `/products/ratings` | Get a list of all ratings submitted by the user |
| `POST` | `/suggestions` | Submit a feedback/suggestion to the admin |

### 👑 Admin Routes (Requires Admin Bearer Token)
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/admin/users` | Get a paginated list of all registered users |
| `PUT` | `/admin/users/:id` | Update user roles or information |
| `DELETE` | `/admin/users/:id` | Permanently delete a user account |
| `POST` | `/admin/categories` | Create a new food category |
| `PUT` | `/admin/categories/:id`| Update an existing category |
| `DELETE` | `/admin/categories/:id`| Delete a category |
| `POST` | `/admin/products` | Add a new product (Supports `multipart/form-data` image upload) |
| `PUT` | `/admin/products/:id` | Update product details or replace image |
| `DELETE` | `/admin/products/:id` | Delete a product from the catalog |
| `GET` | `/admin/orders` | View all customer orders across the system |
| `PUT` | `/admin/orders/:id/status`| Update the delivery status of an order |
| `GET` | `/admin/suggestions` | View all user feedback and suggestions |