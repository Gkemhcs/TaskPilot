# 🚀 TaskPilot — Scalable Task and Project Management Backend

[![Test and Build](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml/badge.svg)](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml)

> TaskPilot is a clean, modular, and production-grade backend system designed for managing tasks and projects, built with Go, PostgreSQL, JWT Authentication, and REST APIs.

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
* [🌐 Deployment](#-deployment)
* [📈 Future Roadmap](#-future-roadmap)
* [👤 Author](#-author)

---

## ⚙️ Features

```mermaid
graph TD
  A[Client] -->|HTTP| B[API Gateway (Gin)]
  B --> C[JWT Middleware]
  C --> D[Business Logic (Services)]
  D --> E[Data Access (Repositories via sqlc)]
  E --> F[(PostgreSQL)]
  G[Docker Compose] --> B
  H[CI/CD (GitHub Actions)] --> G
  I[Monitoring (Prometheus)] --> B
  I --> F
```

* 🔐 **Secure JWT Auth**: Access + refresh token rotation with context-based auth middleware
* ⏱ **Per-IP + Route-Based Rate Limiting**: Prevent abuse using `golang.org/x/time/rate`
* 📊 **Prometheus Metrics**: Per-route request counts, error tracking & latency histograms
* 🧼 **Clean Hexagonal Architecture**: Domain-specific handlers, services, and types
* 📇 **Typed DB Access with `sqlc`**: Go code is generated from raw SQL queries, scoped per domain (`user`, `task`, `project`)
* 🐳 **One-Command Docker Compose**: Boots app, migrations, Prometheus, and PostgreSQL
* 📚 **Auto Swagger Docs**: Try-it-out UI + Bearer auth support
* 🧪 **Layered Unit Testing**: Service logic and HTTP handlers tested with mocks & assertions
* ⚙️ **GitHub Actions CI/CD**: Test, build, and deploy pipeline configured

---

## 📆 Tech Stack

| Layer     | Tech                                  |
| --------- | ------------------------------------- |
| Language  | Go (Golang)                           |
| Framework | Gin-Gonic (HTTP Routing)              |
| Database  | PostgreSQL                            |
| DB Access | `sqlc` (type-safe query generator)    |
| Auth      | JWT (Bearer Token)                    |
| API Docs  | Swagger (Swaggo)                      |
| Testing   | Testify, Mock                         |
| DevOps    | GitHub Actions, Azure PostgreSQL, GCP |

---

## 🧱 Architecture

```mermaid
graph TD
  A[Client] -->|HTTP| B[GIN HTTP Handlers]
  B --> C[Middleware (JWT, Rate Limit, Prometheus)]
  C --> D[Domain Services (Business Logic)]
  D --> E[Repositories (sqlc per domain)]
  E --> F[(PostgreSQL)]
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
* Test results are published and displayed directly in the GitHub UI for easy review (see the "Checks" tab on your PRs and commits).

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
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── types.go
│   │   └── gen/               # sqlc-generated DB access code
│   ├── project/               # Project domain logic
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── types.go
│   │   └── gen/               # sqlc-generated DB access code
│   ├── user/                  # User domain logic
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── types.go
│   │   └── gen/               # sqlc-generated DB access code
│   ├── middleware/            # JWT, metrics, and rate-limiting middleware
│   ├── db/
│   │   └── migrations/        # SQL schema migrations
│   ├── errors/                # Custom error definitions
│   └── utils/                 # Helper utilities for response formatting etc.
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
2. Copy or set environment variables in your `.env` file or edit `docker-compose.yaml` as needed.
3. Start all services:

```bash
docker-compose up --build
```

This will start the backend app, PostgreSQL database, run DB migrations using `migrate`, and expose Prometheus for metrics collection. Swagger will be available at `http://localhost:8080/docs/index.html`.

---

## 🌐 Deployment

* GCP VM, Azure PostgreSQL, Docker Compose
* CI/CD using GitHub Actions

---

## 📈 Future Roadmap

* [ ] Role-based access control
* [ ] Redis caching for frequently accessed tasks
* [ ] WebSocket support for live task updates
* [ ] OpenTelemetry tracing
* [ ] Admin panel / dashboard integration

---

## 👤 Author

**Koti Eswar Mani Gudi**
📧 [gudikotieswarmani@gmail.com](mailto:gudikotieswarmani@gmail.com)
🌐 [GitHub](https://github.com/Gkemhcs) | [LinkedIn](https://www.linkedin.com/in/gkemhcs/) | [Portfolio](https://gkemhcs.dev)
