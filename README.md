# LedgeGuard Finance Backend

LedgeGuard is a production-grade financial data processing and access control backend built with Go, following **Domain-Driven Design (DDD)** and **Spec-Driven Development** principles.

## Features
- **JWT-based Authentication**: Secure login with refresh token rotation.
- **Role-Based Access Control (RBAC)**: Fine-grained permissions (ADMIN, ANALYST, VIEWER).
- **Financial Records**: Management of income and expense entries.
- **Dashboard Analytics**: Real-time aggregations (Total Income, Expenses, Net Balance).
- **Swagger Documentation**: Interactive API testing at `/swagger/index.html`.
- **BDD Testing**: Ginkgo (Unit) and Godog (Functional/Gherkin) test suites.

## Tech Stack
- **Language**: Go 1.25.x
- **Framework**: Gin (Web), GORM (ORM)
- **Database**: SQLite (Pure-Go driver)
- **Security**: JWT, Bcrypt
- **Testing**: Ginkgo, Gomega, Godog

## Getting Started

### Prerequisites
- Go installed on your system.
- Git.

### Setup
1. Clone the repository: `git clone https://github.com/Abhaykauts/LedgeGuard.git`
2. Navigate to the backend: `cd backend`
3. Install dependencies: `go mod download`
4. Run the API: `go run cmd/api/main.go`

### Testing
- **Unit Tests**: `cd backend/pkg/auth && go test -v .`
- **Functional Tests**: `cd backend/tests && go test -v .`

## API Documentation
Once the server is running, visit:
`http://localhost:8080/swagger/index.html`

### Initial Admin Credentials
- **Username**: `admin`
- **Password**: `admin123`

## Development Workflow
1. Create a feature branch: `feat/your-feature`
2. Define User Stories in `docs/user_stories/`
3. Write Gherkin features in `backend/features/`
4. Implement Domain and Application logic.
5. Run tests and commit with **Conventional Commits**.
