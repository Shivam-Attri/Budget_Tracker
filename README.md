# Go Budget Tracker API

A secure, scalable, and production-ready RESTful API for personal budget tracking. Built with Go, this backend service provides a comprehensive suite of features for managing users, transactions, categories, and budgets, all while adhering to modern best practices for security and performance.

---

## ✨ Features

* **Secure Authentication**:
    * **JWT Access & Refresh Tokens**: State-of-the-art authentication using short-lived access tokens and long-lived, rotating refresh tokens for a seamless and secure user experience.
    * **Password Hashing**: Uses the robust `bcrypt` algorithm to securely hash and store user passwords.
* **Data Security & Encryption**:
    * **Field-Level Encryption**: Sensitive data fields (like transaction descriptions and category names) are encrypted at rest in the database using AES-GCM.
    * **Data Isolation**: The API is multi-tenant, ensuring users can only access and manipulate their own financial data.
* **Production-Grade Architecture**:
    * **Containerized with Docker**: Comes with a multi-stage `Dockerfile` and a `docker-compose.yml` file for easy, consistent, and reliable deployment.
    * **Advanced Configuration**: Uses **Viper** for flexible configuration management from files (`config.yml`) and environment variables.
    * **Graceful Shutdown**: The server handles system signals to shut down gracefully, finishing all in-progress requests before exiting.
    * **Health Check Endpoint**: A `/health` endpoint to allow automated monitoring of the service's status.
    * **Structured Logging**: Uses **Zerolog** for structured, JSON-formatted logs, ideal for production environments.
* **Advanced API Functionality**:
    * **API Rate Limiting**: Per-client IP-based rate limiting to prevent abuse and ensure service stability.
    * **Pagination**: All list endpoints (`/transactions`, `/categories`) are paginated to handle large datasets efficiently.
    * **Dynamic Filtering & Sorting**: The `/transactions` endpoint supports powerful filtering (by type, date range, amount range) and sorting.
* **Core Application Logic**:
    * Full CRUD operations for Transactions, user-defined Categories, and monthly Budgets.
    * Reporting endpoint to get monthly financial summaries (total income, expenses, and savings).

---

## 📂 Project Structure

The project follows a clean, modular architecture to separate concerns.

.├── auth/                 # Authentication logic (JWT, passwords, middleware)│   └── auth.go├── config/               # Configuration loading (Viper)│   └── config.go├── database/             # Database connection, schema, and store implementation│   ├── database.go│   └── store.go├── handlers/             # HTTP handlers for each API resource│   ├── budgets.go│   ├── categories.go│   ├── handlers.go│   ├── reports.go│   ├── system.go│   ├── transactions.go│   └── users.go├── middleware/           # Custom middleware (e.g., rate limiting)│   └── ratelimit.go├── models/               # Data structures for the application│   ├── budget.go│   ├── category.go│   ├── filter.go│   ├── pagination.go│   ├── report.go│   ├── transaction.go│   └── user.go├── router/               # API route definitions│   └── router.go├── utils/                # Helper utilities (validation, pagination)│   ├── pagination.go│   └── validator.go├── .env.example          # Example environment variables file├── config.yml            # Default configuration file├── docker-compose.yml    # Docker Compose setup for local development├── Dockerfile            # Dockerfile for building the application image├── go.mod                # Go module dependencies├── go.sum└── main.go               # Application entry point
---

## 🚀 Getting Started

### Prerequisites

* [Go](https://go.dev/doc/install) (version 1.21 or later)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* A tool for making API requests, like [Postman](https://www.postman.com/) or [cURL](https://curl.se/).

### Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd <repository-directory>
    ```

2.  **Create the `.env` file:**
    Copy the example environment file and generate your secrets.
    ```bash
    cp .env.example .env
    ```
    Open the `.env` file and replace the placeholder values with secure, randomly generated strings.
    * **For JWT_SECRET:** `openssl rand -base64 32`
    * **For ENCRYPTION_KEY:** `openssl rand -base64 32` (must be 32 bytes for AES-256)

### Running the Project

You can run the application in two ways:

#### 1. Using Docker (Recommended)

This is the easiest and most reliable way to run the entire stack, including the PostgreSQL database.

1.  **Start the services:**
    From the root of the project, run:
    ```bash
    docker-compose up --build
    ```
    This command will:
    * Build the Go application image using the `Dockerfile`.
    * Start a PostgreSQL container.
    * Start the API container and connect it to the database.
    * The API will be available at `http://localhost:8000`.

2.  **Stopping the services:**
    Press `Ctrl + C` in the terminal, and then run `docker-compose down` to stop and remove the containers.

#### 2. Running Locally with `go run`

This method is useful for development if you prefer to run the database separately.

1.  **Start a PostgreSQL instance:**
    You must have a running PostgreSQL server. You can start one easily with Docker:
    ```bash
    docker run --name budget-db-local -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
    ```
    You will also need to connect to this instance with a tool like `psql` or DBeaver and create a database named `budget`.

2.  **Configure `config.yml`:**
    Make sure the `database.url` in `config.yml` points to your local PostgreSQL instance (e.g., `postgres://postgres:mysecretpassword@localhost:5432/budget?sslmode=disable`).

3.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Run the application:**
    ```bash
    go run main.go
    ```
    The API will be available at `http://localhost:8000`.

---

## ⚙️ Configuration

The application is configured through a combination of `config.yml` and environment variables. **Environment variables always override values set in the file.**

| Setting                        | Environment Variable      | `config.yml` Path                    | Description                                                                 |
| ------------------------------ | ------------------------- | ------------------------------------ | --------------------------------------------------------------------------- |
| Database URL                   | `BUDGET_DB_URL`           | `database.url`                       | The connection string for the PostgreSQL database.                          |
| JWT Secret                     | `BUDGET_JWT_SECRET`       | `server.jwt_secret`                  | A long, random secret key for signing JWTs.                                 |
| Encryption Key                 | `BUDGET_ENCRYPTION_KEY`   | `server.encryption_key`              | A 32-byte key for AES-256 field-level encryption.                           |
| Server Port                    | `BUDGET_SERVER_PORT`      | `server.port`                        | The port on which the API server listens.                                   |
| Access Token TTL               |                           | `server.access_token_ttl_minutes`    | The lifespan of an access token in minutes.                                 |
| Refresh Token TTL              |                           | `server.refresh_token_ttl_hours`     | The lifespan of a refresh token in hours.                                   |
| Rate Limiter Enabled           |                           | `rate_limiter.enabled`               | `true` or `false` to enable/disable the API rate limiter.                   |
| Rate Limiter RPS               |                           | `rate_limiter.rps`                   | The number of allowed requests per second per IP.                           |
| Rate Limiter Burst             |                           | `rate_limiter.burst`                 | The number of burst requests allowed per IP.                                |

---

## Endpoints

Base URL: `http://localhost:8000`

### Public Routes

| Method | Path         | Description                                        |
| :----- | :----------- | :------------------------------------------------- |
| `GET`  | `/health`    | Checks the service health and database connection. |
| `POST` | `/register`  | Registers a new user.                              |
| `POST` | `/login`     | Logs in a user, returning access and refresh tokens. |
| `POST` | `/refresh`   | Issues a new token pair using a valid refresh token. |

### Protected Routes

All routes under `/api/v1` require a valid `Bearer` token in the `Authorization` header.

| Method   | Path                          | Description                                                                                             |
| :------- | :---------------------------- | :------------------------------------------------------------------------------------------------------ |
| `POST`   | `/api/v1/logout`              | Logs out by invalidating the provided refresh token.                                                    |
| `GET`    | `/api/v1/transactions`        | Get a paginated list of transactions. Supports filtering and sorting.                                   |
| `POST`   | `/api/v1/transactions`        | Create a new transaction.                                                                               |
| `GET`    | `/api/v1/transactions/{id}`   | Get a single transaction by its ID.                                                                     |
| `PUT`    | `/api/v1/transactions/{id}`   | Update a transaction.                                                                                   |
| `DELETE` | `/api/v1/transactions/{id}`   | Delete a transaction.                                                                                   |
| `GET`    | `/api/v1/categories`          | Get a paginated list of the user's custom categories.                                                   |
| `POST`   | `/api/v1/categories`          | Create a new category.                                                                                  |
| `PUT`    | `/api/v1/categories/{id}`     | Update a category.                                                                                      |
| `DELETE` | `/api/v1/categories/{id}`     | Delete a category.                                                                                      |
| `GET`    | `/api/v1/budgets`             | Get budgets for a given month and year.                                                                 |
| `POST`   | `/api/v1/budgets`             | Create a new monthly budget for a category.                                                             |
| `PUT`    | `/api/v1/budgets/{id}`        | Update a budget.                                                                                        |
| `DELETE` | `/api/v1/budgets/{id}`        | Delete a budget.                                                                                        |
| `GET`    | `/api/v1/reports/summary`     | Get a financial summary (income, expense, savings) for a given month and year.                          |

#### Transaction Filtering & Sorting

The `GET /api/v1/transactions` endpoint supports the following query parameters:

* `page`: The page number for pagination (e.g., `?page=2`).
* `limit`: The number of items per page (e.g., `?limit=25`).
* `type`: `income` or `expense`.
* `min_amount`: The minimum transaction amount (e.g., `?min_amount=50.5`).
* `max_amount`: The maximum transaction amount.
* `start_date`: The start of a date range in RFC3339 format (e.g., `2025-07-01T00:00:00Z`).
* `end_date`: The end of a date range in RFC3339 format.
* `sort`: A comma-separated list of fields to sort by. Use a `-` prefix for descending order. Allowed fields: `date`, `amount`. (e.g., `?sort=-date,amount`).

