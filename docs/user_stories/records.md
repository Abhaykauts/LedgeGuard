# User Story: Financial Records Management

**As an** Accountant or Analyst,
**I want** to record and manage financial transactions,
**So that** I can track income, expenses, and overall business health.

## Business Requirements
1. **Record Entries**: Users can create records with `Amount`, `Type` (Income/Expense), `Category`, `Date`, and `Notes`.
2. **Read Records**:
   - **Viewers/Analysts/Admins** can list records.
   - **Filtering**: Support filtering by `Date Range`, `Category`, and `Type`.
3. **Manage Records**:
   - **Admins** can Update or Delete any record.
   - **Analysts** (optional) can Update their own recorded entries.
4. **Data Integrity**: Records must be validated (e.g., non-negative amounts, valid categories).
5. **Soft Delete**: Records shouldn't be permanently purged from the database initially; they should be marked as deleted.

## Acceptance Criteria
- [ ] Gherkin features pass for creating, listing, and filtering records.
- [ ] Unauthorized users (Viewers) are blocked from Create/Update/Delete actions.
- [ ] Summaries correctly reflect only non-deleted records.
- [ ] Swagger documentation shows clear validation rules for record entry.
