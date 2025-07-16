# 🚀 TaskPilot — Scalable Task and Project Management Backend

[![Test and Build](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml/badge.svg)](https://github.com/Gkemhcs/TaskPilot/actions/workflows/ci.yml)

> TaskPilot is a clean, modular, and production-grade backend system designed for managing tasks and projects, built with Go, PostgreSQL, JWT Authentication, and REST APIs. It supports background job processing (import/export), RabbitMQ messaging, and secure file handling through cloud/local storage.



---

## 📌 Table of Contents
- [🚀 TaskPilot — Scalable Task and Project Management Backend](#-taskpilot--scalable-task-and-project-management-backend)
  - [📌 Table of Contents](#-table-of-contents)
  - [🎬 Demo Videos](#-demo-videos)
  - [⚙️ Features](#️-features)
  - [📆 Tech Stack](#-tech-stack)
  - [🧱 Architecture](#-architecture)
  - [🔐 Authentication](#-authentication)
  - [📄 API Documentation](#-api-documentation)
    - [Try Auth-Protected Endpoints](#try-auth-protected-endpoints)
  - [🧪 Testing Strategy](#-testing-strategy)
  - [📁 Project Structure](#-project-structure)
  - [🚀 Getting Started](#-getting-started)
    - [1. Clone the Repository and Install Dependencies](#1-clone-the-repository-and-install-dependencies)
    - [2. Set Up Google Cloud Storage (GCP Bucket)](#2-set-up-google-cloud-storage-gcp-bucket)
    - [3. Running Without Docker (Manual Mode)](#3-running-without-docker-manual-mode)
    - [4. Docker Compose Deployment (Recommended for Local Development)](#4-docker-compose-deployment-recommended-for-local-development)
  - [📡 Ports](#-ports)
  - [📈 Future Roadmap](#-future-roadmap)
  - [👤 Author](#-author)

---
## 🎬 Demo Videos

### 📤 Project Import and Task Export Demo


Check out the video demo showcasing both the **Project Import** and **Task Export** features of TaskPilot — built with Go, PostgreSQL, RabbitMQ, and a clean layered architecture.

[![Watch Demo Video](https://img.youtube.com/vi/YOUR_VIDEO_ID/0.jpg)](https://youtu.be/SJZP2r9oeXs)





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

### 1. Clone the Repository and Install Dependencies

1. Clone the repository:

```bash
git clone https://github.com/Gkemhcs/taskpilot.git
cd taskpilot
go mod tidy
```

2. Set up your environment variables (replace with your actual values):

```bash
PORT=5000
DATABASE_URL=postgres://...
JWT_SECRET=your-super-secret-key
REDIS_URL
```

---

### 2. Set Up Google Cloud Storage (GCP Bucket)

1. Enter your Google Cloud project ID (must have billing enabled):

```bash
echo "Enter  your google cloud project id with billing account attached"
read PROJECT_ID
gcloud config set project $PROJECT_ID
```

2. Create the GCP storage bucket:

```bash
echo  "Creating the gcp  storage bucket"
gsutil mb  -l asia-south1 "gs://taskpilot-${PROJECT_ID}"
echo "Bucket creation successful. BucketName:- taskpilot-${PROJECT_ID}"
```

3. Create a service account, grant permissions, and download the key file:

```bash
echo "Creating the service account and adding permissions and downloading the key file"

gcloud iam service-accounts create taskpilot-storage-writer --display-name "Taskpilot backend service account"
gcloud projects add-iam-policy-binding $PROJECT_ID \
--member "serviceAccount:taskpilot-storage-writer@${PROJECT_ID}.iam.gserviceaccount.com" \
--role roles/storage.admin 

gcloud iam service-accounts keys create key.json --iam-account "taskpilot-storage-writer@${PROJECT_ID}.iam.gserviceaccount.com"
```

---

### 3. Running Without Docker (Manual Mode)

> Open **three separate terminals** and run each of the following command sets in its own terminal window:

**Terminal 1: Start the API Backend Server**

```bash
export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/key.json"
export GCP_BUCKET="taskpilot-${PROJECT_ID}"
export GCP_PREFIX="taskpilot-backend-data"
echo "Staring the  TaskPilot API backend server"
go run main.go
```

**Terminal 2: Start the Project Worker**

```bash
export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/key.json"
export GCP_BUCKET="taskpilot-${PROJECT_ID}"
export GCP_PREFIX="taskpilot-backend-data"
echo "Staring the  Taskpilot Project Worker"
go run .cmd/worker/project
```

**Terminal 3: Start the Task Worker**

```bash
export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/key.json"
export GCP_BUCKET="taskpilot-${PROJECT_ID}"
export GCP_PREFIX="taskpilot-backend-data"
echo "Staring the  Taskpilot  Task Worker"
go run .cmd/worker/task
```

---

### 4. Docker Compose Deployment (Recommended for Local Development)

1. Ensure Docker and Docker Compose are installed on your system.
2. Update the `GCP_BUCKET` value in the `docker-compose.yaml` file to match the bucket you created above.

```bash
docker-compose up --build
```

This will start the backend app, PostgreSQL, Swagger docs, Prometheus metrics, and RabbitMQ workers in one go.

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
