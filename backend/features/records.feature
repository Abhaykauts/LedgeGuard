Feature: Financial Records Management
  As a user with sufficient role
  I want to manage financial transactions
  So that I can track incomes and expenses

  Scenario: Create a valid income record
    Given I am an authenticated "ADMIN"
    When I create a "INCOME" record with amount 1000 and category "Salary"
    Then the record should be saved successfully
    And the total count of records should be 1

  Scenario: Unauthorized Viewer cannot create a record
    Given I am an authenticated "VIEWER"
    When I create a "EXPENSE" record with amount 50 and category "Coffee"
    Then I should receive an "Access Denied" error

  Scenario: List and filter records by type
    Given I am an authenticated "ANALYST"
    And the following records exist:
      | Type    | Amount | Category |
      | INCOME  | 2000   | Sales    |
      | EXPENSE | 500    | Rent     |
    When I list records filtered by type "INCOME"
    Then I should see 1 record
    And the record amount should be 2000
