# User Story: Financial Dashboard Analytics

**As a** business owner,
**I want** to see an aggregated summary of my financial state,
**So that** I can make quick, data-driven decisions.

## Business Requirements
1. **Financial Overview**: Calculate `Total Income`, `Total Expenses`, and `Net Balance`.
2. **Categorical Breakdown**: Show a summary of totals per `Category` (e.g., "Software", "Consulting").
3. **Trends**: Display monthly or weekly totals for trend analysis.
4. **Recent Activity**: List the most recent transaction entries.
5. **Real-time Data**: Access cached or real-time views of financial performance.

## Acceptance Criteria
- [ ] Gherkin features pass for calculating totals and trends and mapping categories.
- [ ] Analytics correctly exclude deleted records.
- [ ] Aggregations respond within performance SLAs (e.g., < 200ms).
- [ ] Swagger documentation clearly defines input parameters (e.g., date ranges).
