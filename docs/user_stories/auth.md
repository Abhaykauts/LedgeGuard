# User Story: Secure Access and Role Control

**As a** system user,
**I want** to securely authenticate and access features based on my role,
**So that** financial data remains confidential and protected.

## Business Requirements
1. **Login**: Users can log in with `username` and `password`.
2. **Tokens**: Successful login returns an `Access Token` (short-lived) and a `Refresh Token` (long-lived).
3. **Session Management**: Users can use a `Refresh Token` to rotate their `Access Token` without logging in again.
4. **Role Enforcement**:
   - **Viewer**: Read-only access to dashboard and records.
   - **Analyst**: Access to dashboard, records, and summaries/insights.
   - **Admin**: Full control (Full management of users and records).
5. **Security**: Password hashing (Bcrypt) and secure JWT signing.

## Acceptance Criteria
- [ ] Gherkin features pass for login, token refresh, and role-based rejection.
- [ ] Unauthorized access to Admin endpoints returns `403 Forbidden`.
- [ ] Passwords are never stored or transmitted in plain text.
- [ ] Swagger documentation clearly shows restricted endpoints.
