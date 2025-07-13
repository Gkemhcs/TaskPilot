# 🚀 TaskPilot — Scalable Task and Project Management Backend

[![Test and Build](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml/badge.svg)](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml)

> TaskPilot is a clean, modular, and production-grade backend system designed for managing tasks and projects, built with Go, PostgreSQL, JWT Authentication, and REST APIs. It supports background job processing (import/export), RabbitMQ messaging, and secure file handling through cloud/local storage.

---

## 📌 Table of Contents

* [⚙️ Features](#️-features)
* [📆 Tech Stack](#-tech-stack)
* [🧱 Architecture](#-architecture)
* [🔐 Authentication](#-authentication)
* [📄 API Documentation](#-api-documentation)
* [🧪 Testing Strategy](#-testing-strategy)
* [📁 Project Structure](#-project-structure)
* [🚀 Getting Started](#-getting-started)
* [📆 Docker Compose Setup](#-docker-compose-setup)
* [📈 Future Roadmap](#-future-roadmap)
* [👤 Author](#-author)

---

## ⚙️ Features

* 🔐 **Secure JWT Auth**: Access + refresh token rotation with context-based auth middleware
* ⏱ **Per-IP + Route-Based Rate Limiting**: Prevent abuse using `github.com/ulule/limiter/v3`
* 📊 **Prometheus Metrics**: Per-route request counts, error tracking & latency histograms
* 🧼 **Clean Hexagonal Architecture**: Domain-specific handlers, services, and types
* 📇 **Typed DB Access with `sqlc`**: Go code is generated from raw SQL queries, scoped per domain (`user`, `task`, `project`)
* 🐳 **One-Command Docker Compose**: Boots app, migrations, Prometheus, and PostgreSQL
* 📚 **Auto Swagger Docs**: Try-it-out UI + Bearer auth support
* 🧪 **Layered Unit Testing**: Service logic and HTTP handlers tested with mocks & assertions
* 📤 **Async Import/Export with RabbitMQ**: Background job workers for Excel import/export of projects/tasks
* ☁️ **Pluggable Cloud/Local File Storage**: Unified interface to support GCP and local processing
* ⚙️ **GitHub Actions CI**: Automated test and build pipeline

---

## 📆 Tech Stack

| Layer       | Tech                                  |
|------------|---------------------------------------|
| Language    | Go (Golang)                           |
| Framework   | Gin-Gonic (HTTP Routing)              |
| Database    | PostgreSQL                            |
| DB Access   | `sqlc` (type-safe query generator)    |
| Messaging   | RabbitMQ (AMQP workers)               |
| Auth        | JWT (Bearer Token)                    |
| API Docs    | Swagger (Swaggo)                      |
| Testing     | Testify, Mock                         |
| Storage     | Local FS / GCP (via interface)        |
| DevOps      | GitHub Actions                        |

---

## 🧱 Architecture

```mermaid
graph TD
  A[Client] -->|HTTP| B[GIN HTTP Handlers]
  B --> C[Middleware\n(JWT + Rate Limit + Prometheus)]
  C --> D[Services Layer]
  D --> E[Repository Layer\n(sqlc per domain)]
  E --> F[(PostgreSQL)]
  D --> G[RabbitMQ Producers]
  G --> H[Async Workers (Import/Export)]
  H --> I[Excel Engine + File Storage]
```

---

## 🔐 Authentication

* **JWT Bearer Tokens**: Used to secure all `/api/v1/*` routes
* **Token Types**:
  * Access Token (short-lived)
  * Refresh Token (long-lived)
* **Authorization**:
  * Passed via `Authorization: Bearer <token>` in headers
  * Middleware parses and injects `userID` into context

---

## 📄 API Documentation

> Auto-generated using swaggo.

📚 [Live Swagger UI](http://localhost:8080/docs/index.html)

### Try Auth-Protected Endpoints

1. Click the 🔒 “Authorize” button in Swagger UI
2. Paste: `Bearer <your-access-token>`
3. Call secure endpoints like `/api/v1/projects` or `/api/v1/tasks`

---

## 🧪 Testing Strategy

✅ Unit Tests for:

* Handlers (using real service + mocked repo)
* Services (mocked repo)
* Edge case validations

🛠 Test Frameworks:

* `testify`
* `testify/mock`
* `httptest` (for HTTP handlers)

🟢 **CI Integration:**

* Automated tests run on every push and pull request via GitHub Actions.
* Test results are published and displayed directly in the GitHub UI for easy review.

---

## 📁 Project Structure

```bash
.
├── main.go                    # Loads config and calls cmd/server/main.go
├── cmd/
│   └── server/
│       └── main.go            # Entry point: initializes and runs the server
├── internal/
│   ├── auth/                  # JWT handling and generation logic
│   ├── task/                  # Task domain logic
│   ├── project/               # Project domain logic
│   ├── user/                  # User domain logic
│   ├── middleware/            # JWT, metrics, and rate-limiting middleware
│   ├── importer/              # Excel importers with row-level validation
│   ├── exporter/              # Excel exporters + RabbitMQ consumers
│   ├── storage/               # Cloud/Local file storage abstraction
│   ├── db/
│   │   └── migrations/        # SQL schema migrations
│   ├── errors/                # Custom error definitions
│   └── utils/                 # Helper utilities
├── config/                    # Application configuration and env handling
├── docs/                      # Swagger docs (autogenerated)
└── go.mod
```

---

## 🚀 Getting Started

### 1. Clone and setup

```bash
git clone https://github.com/Gkemhcs/taskpilot.git
cd taskpilot
go mod tidy

PORT=5000
DATABASE_URL=postgres://...
JWT_SECRET=your-super-secret-key

go run cmd/server/main.go
```

---

## 📆 Docker Compose Setup (Recommended for local development)

1. Make sure Docker and Docker Compose are installed.
2. Copy or set environment variables in your `.env` file or edit `docker-compose.yaml`.
3. Start all services:

```bash
docker-compose up --build
```

This boots the backend app, PostgreSQL, Swagger docs, Prometheus metrics, and RabbitMQ workers.

---


---

## 📡 Ports

| Service         | Port | Description                            |
|-----------------|------|----------------------------------------|
| Backend API     | 8080 | Main application server                |
| PostgreSQL      | 5432 | Database used for persistence          |
| RabbitMQ (AMQP) | 5672 | Message queue protocol (used by workers)|
| RabbitMQ UI     | 15672| Web-based management interface         |
| Redis           | 6379 | In-memory cache store                  |
| Prometheus      | 9090 | Metrics monitoring and visualization   |
| Postgres Exporter | 9187 | Exports DB metrics for Prometheus     |

> ℹ️ Ensure the listed ports are not blocked by firewalls and not used by other local services.



## 📈 Future Roadmap

* [ ] Task Export per Project (done) ✅
* [ ] Global Task Export for User
* [ ] Role-based access control (RBAC)
* [ ] WebSocket support for live task updates
* [ ] Redis caching for popular projects
* [ ] Admin dashboard & analytics
* [ ] Tracing support (OpenTelemetry)

---

## 👤 Author

**Koti Eswar Mani Gudi**
📧 [gudikotieswarmani@gmail.com](mailto:gudikotieswarmani@gmail.com)
🌐 [GitHub](https://github.com/Gkemhcs) | [LinkedIn](https://www.linkedin.com/in/gkemhcs/) | [Portfolio](https://gkemhcs.github.io)
