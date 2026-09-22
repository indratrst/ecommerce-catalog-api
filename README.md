# Go Fiber & PostgreSQL RESTful API with Clean Architecture

A production-ready, highly scalable RESTful API built using Go (Fiber framework) and PostgreSQL. This project adheres strictly to **Clean Architecture** principles and is fully containerized using Docker and Docker Compose for seamless deployment and development.

---

## 🏗️ Architecture Overview

The project is structured according to **Clean Architecture** to ensure low coupling and high maintainability:

```text
├── config/             # Database connection & AutoMigration configurations
├── domain/             # Entities, Models, DTOs, and Interfaces (Contracts)
├── repository/         # Data Access Layer (GORM / SQL queries)
├── service/            # Core Business Logic Layer
├── handler/            # Delivery Layer (HTTP Request/Response Controllers)
├── middleware/         # HTTP Interceptors (JWT Authentication)
├── utils/              # Helper utilities (Bcrypt, JWT Token generators)
├── Dockerfile          # Multi-Stage Build definition
├── docker-compose.yml  # Docker multi-container orchestration
└── main.go             # Application entry point & Dependency Injection
```

---

## 🚀 Features

- **Authentication & Authorization**: User registration and login using `bcrypt` password hashing and stateless **JWT** bearer tokens.
- **Product & Category Management**: Full CRUD capabilities with a **One-to-Many** relationship (Categories -> Products).
- **Advanced Querying**: Dynamic **Pagination** (`page`, `limit`), case-insensitive **Search** (`search`), and Category Filtering.
- **Nested JSON Responses**: Optimized database preloading (`Preload`) for relational data fetching.
- **Production Containerization**: Lightweight Docker container (~20-30MB) powered by Go **Multi-Stage Builds**.

---

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.23+
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL 16
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker & Docker Compose

---

## 🚦 Getting Started

### Prerequisites

Ensure you have the following installed on your host machine:
- [Docker Engine](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Application

1. **Clone the Repository**
   ```bash
   git clone https://github.com/your-username/my-first-go-api.git
   cd my-first-go-api
   ```

2. **Start the Infrastructure and Application**
   Run the following command to build and start both the Go API and PostgreSQL containers:
   ```bash
   docker compose up --build -d
   ```

3. **Verify Execution**
   The API will be available at `http://localhost:8080`.

---

## 🐳 Docker Commands Reference

Here is a guide to the key Docker commands used in this project and their functions:

| Command | Description |
| :--- | :--- |
| `docker compose up --build -d` | **Builds/Rebuilds** the Docker image using the `Dockerfile`, starts both `api` and `postgres` containers, and runs them in **detached mode** (background). Use this every time you change Go source code. |
| `docker compose up -d` | Starts the existing containers in the background without rebuilding the application code. |
| `docker compose down` | Stops and removes all running containers, networks, and resources defined in `docker-compose.yml`. |
| `docker compose logs -f api` | Stream live console logs from the `go_fiber_api` container (replaces local `go run main.go` output). |
| `docker compose ps` | Displays the current running status and port mappings for all project containers. |
| `docker exec -it go_postgres_db psql -U postgres -d mygodb` | Accesses the interactive PostgreSQL CLI inside the database container directly. |

---

## 📑 API Endpoint Documentation

### Auth Endpoints
- `POST /api/v1/register` — Register a new user
- `POST /api/v1/login` — Login and receive a JWT Token

### Category Endpoints *(Protected)*
- `POST /api/v1/categories` — Create a new Category
- `GET /api/v1/categories` — Retrieve all Categories

### Product Endpoints *(Protected)*
- `POST /api/v1/products` — Create a new Product (Requires `category_id`)
- `GET /api/v1/products` — Retrieve Products with Pagination & Search

#### Query Parameters for `GET /api/v1/products`:
| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `page` | `int` | `1` | Page number to fetch |
| `limit` | `int` | `10` | Number of items per page |
| `search` | `string` | `""` | Search products by name (case-insensitive) |
| `category_id` | `uint` | - | Filter products by specific Category ID |

**Example Search & Pagination URL:**
```text
http://localhost:8080/api/v1/products?search=keyboard&page=1&limit=5&category_id=1
```

---

## 🔐 Authorization Header Example

For protected endpoints, attach the JWT token in the request header:

```text
Authorization: Bearer <YOUR_JWT_TOKEN>
```