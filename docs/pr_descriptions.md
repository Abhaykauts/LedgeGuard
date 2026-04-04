# Pull Request Descriptions: LedgeGuard Phase 1

This document provides the structured descriptions for the three initial feature branches. You can copy these directly into your GitHub Pull Request descriptions.

---

## PR 1: Authentication Core & RBAC (`feat/auth-core`)

### Description
Implement the foundational security layer using Domain-Driven Design (DDD). This includes JWT-based session management, refresh token rotation, and Role-Based Access Control (RBAC).

### Key Changes
- **Domain**: Defined `User` entity and `Role` types (ADMIN, ANALYST, VIEWER).
- **Application**: Implemented `AuthService` with Login and Refresh token logic.
- **Infrastructure**: Configured GORM with SQLite (Pure Go driver) for user persistence.
- **Security**: Added `Bcrypt` password hashing and `JWT` generation/validation.
- **API**: Created Gin handlers for login and token rotation.
- **Middleware**: Implemented `AuthMiddleware` and `RoleMiddleware` for endpoint protection.

### Testing Results
- **Unit**: 4/4 specs passed for JWT and Password utilities.
- **Functional**: 3 Godog scenarios passed (Successful Login, Failed Login, Token Rotation).

---

## PR 2: Financial Records Management (`feat/records-core`)

### Description
Implement core financial transaction management allowing users to track incomes and expenses with category-based filtering.

### Key Changes
- **Domain**: Defined `Record` entity and `RecordType` (INCOME, EXPENSE).
- **Application**: Implemented `RecordService` with basic validation (e.g., amount > 0).
- **Infrastructure**: Added SQLite repository implementation for Records with auto-migration.
- **API**: Created Gin handlers for CRUD operations on financial records.
- **Access Control**: Integrated RBAC (ADMIN/ANALYST can create, ADMIN can delete, all can list).

### Testing Results
- **Functional**: 3 Godog scenarios passed (Create Income, Access Denied for Viewer, List with Filters).

---

## PR 3: Dashboard Analytics & Swagger (`feat/dashboard-core`)

### Description
Implement high-level financial aggregations for the business dashboard and generate interactive Swagger API documentation.

### Key Changes
- **Application**: Implemented `DashboardService` for calculating Total Income, Total Expenses, and Net Balance.
- **Analytics**: Added category-wise total aggregations.
- **API**: Created Dashboard summary endpoint.
- **Documentation**: Generated full Swagger 2.0 specs (OpenAPI) using `swag`.

### Testing Results
- **Functional**: Godog scenario passed for aggregation accuracy (Calculating totals and category-specific balances).
