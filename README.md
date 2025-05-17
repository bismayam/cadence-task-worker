# cadence-task-worker

This project is a **Temporal (formerly Cadence)** worker written in Go. It runs a workflow that:

1. Authenticates a user using provided credentials.
2. Uses the received token to call an external **Task API** and create a task.

This allows you to orchestrate task creation reliably using Temporal workflows.

Cadence UI address : http://localhost:8233/

---


## 🚀 Features

- Connects to Temporal server on `localhost:7233`
- Workflow: Authenticate + Create Task via external API
- Activities:
  - `LoginActivity`: Authenticates and returns token
  - `CallTaskAPI`: Uses token to create a task
- Integrated with external Task API (running on `http://localhost:8080`)

---

## ⚙️ Prerequisites

- [Go](https://golang.org/dl/) 1.20+
- [Docker & Docker Compose](https://docs.docker.com/compose/)
- Temporal running locally (via Docker)
- External Task API (JWT auth enabled) running locally on port `8080`

---

## 🐳 Running Temporal with PostgreSQL (optional)

To start a Temporal + PostgreSQL dev environment with Docker:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose up
