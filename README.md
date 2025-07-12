# 🚀 TaskPilot — Scalable Task and Project Management Backend

![Go](https://img.shields.io/badge/Go-1.21+-blue?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue?logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-brightgreen?logo=swagger)
![CI/CD](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml/badge.svg)
![GCP](https://img.shields.io/badge/Deployed%20on-GCP-%23039BE5?logo=googlecloud)

> TaskPilot is a clean, modular, and production-grade backend system designed for managing tasks and projects, built with Go, PostgreSQL, JWT Authentication, and REST APIs.

---

## 📌 Table of Contents

- [⚙️ Features](#️-features)
- [📦 Tech Stack](#-tech-stack)
- [🧱 Architecture](#-architecture)
- [🔐 Authentication](#-authentication)
- [📄 API Documentation](#-api-documentation)
- [🧪 Testing Strategy](#-testing-strategy)
- [📁 Project Structure](#-project-structure)
## 📌 Table of Contents
### 2. Run with Docker Compose (Recommended for local development)
- [⚙️ Features](#%EF%B8%8F-features)
1. Make sure Docker and Docker Compose are installed.
2. Copy or set environment variables in your `.env` file or edit `docker-compose.yaml` as needed.
3. Start all services:

```bash
docker-compose up --build
```
### 2. Run with Docker Compose (Recommended for local development)
This will start the backend app, PostgreSQL database, and run migrations automatically. The API will be available at `http://localhost:8080` and Swagger UI at `http://localhost:8080/docs/index.html`.

---
- [🚀 Getting Started](#-getting-started)
- [🌐 Deployment](#-deployment)
- [📈 Future Roadmap](#-future-roadmap)
- [👤 Author](#-author)

---

## ⚙️ Features
  subgraph User
    A[Client]
  end
  subgraph Backend
    B[API Gateway (Gin)]
    C[JWT Middleware]
    D[Business Logic (Services)]
    E[Data Access (Repositories)]
  end
  subgraph Database
    F[(PostgreSQL)]
  end
  subgraph DevOps
    G[Docker Compose]
    H[CI/CD (GitHub Actions)]
    I[Monitoring (Prometheus)]
  end
  A -->|HTTP| B
  B --> C
  C --> D
  D --> E
  E --> F
  G --> B
  H --> G
  I --> B
  I --> F
- 🔄 **Token Refresh Endpoint**: Securely renew access tokens
- 🧪 **Mock-based Unit Tests**: Thoroughly tested using `testify/mock`
- 🗃️ **PostgreSQL Integration**: Relational schema with migrations

---

## 📦 Tech Stack

| Layer          | Tech                                                                 |
|----------------|----------------------------------------------------------------------|
| Language       | Go (Golang)                                                          |
| Framework      | Gin-Gonic (HTTP Routing)                                             |
| Database       | PostgreSQL                                                           |
| Auth           | JWT (Bearer Token)                                                   |
| API Docs       | Swagger (Swaggo)                                                     |
| Testing        | Testify, Mock                                                        |
| DevOps         | GitHub Actions, Azure PostgreSQL, GCP                                |

---

## 🧱 Architecture

```mermaid
graph TD
  A[Client] -->|HTTP| B[GIN HTTP Handlers]
  B --> C[Middleware (JWT)]
  C --> D[Services Layer]
  D --> E[Repository Layer]
  E --> F[(PostgreSQL)]

```

## 🔐 Authentication

- **JWT Bearer Tokens**: Used to secure all `/api/v1/*` routes
- **Token Types**:
  - Access Token (short-lived)
  - Refresh Token (long-lived)
- **Authorization**:
  - Passed via `Authorization: Bearer <token>` in headers
  - Middleware parses and injects `userID` into context


## 📄 API Documentation

> Auto-generated using swaggo.

📚 [Live Swagger UI](http://localhost:8080/docs/index.html)

### Try Auth-Protected Endpoints

1. Click the 🔒 “Authorize” button in Swagger UI
2. Paste: `Bearer <your-access-token>`
3. Call secure endpoints like `/api/v1/projects` or `/api/v1/tasks`



## 🧪 Testing Strategy

✅ Unit Tests for:
- Handlers (using real service + mocked repo)
- Services (mocked repo)
- Edge case validations

🛠 Test Frameworks:
- `testify`
- `testify/mock`
- `httptest` (for HTTP handlers)

🟢 **CI Integration:**
- Automated tests run on every push and pull request via GitHub Actions.
- Test results are published and displayed directly in the GitHub UI for easy review (see the "Checks" tab on your PRs and commits).


## 📁 Project Structure

```bash
.
├── cmd/
├── internal/
│   ├── auth/            # JWT generation and verification
│   ├── task/            # Handlers, services, repo for tasks
│   ├── project/         # Handlers, services, repo for projects
│   ├── user/            # Auth, registration, login
│   ├── middleware/      # JWT middleware
│   └── utils/           # Error, success response wrappers
├── docs/                # Swagger generated docs
├── db/migrations/       # SQL schema migrations
├── config/              # App config
└── go.mod



## 🚀 Getting Started

### 1. Clone and setup
```bash
git clone https://github.com/Gkemhcs/taskpilot.git
cd taskpilot
go mod tidy



PORT=5000
DATABASE_URL=postgres://...
JWT_SECRET=your-super-secret-key

go run main.go

```



## 📈 Future Roadmap

- [ ] Add role-based access control
- [ ] Integrate Redis for caching tasks
- [ ] WebSocket for live task updates
- [ ] Metrics endpoint for Prometheus
- [ ] Full Docker + CI/CD setup
