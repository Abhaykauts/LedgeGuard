Feature: Financial Dashboard Analytics
  As a business owner
  I want to see my financial summary
  So that I can understand my business health

  Scenario: Calculate Total Summary
    Given the following financial records exist:
      | Type    | Amount | Category   |
      | INCOME  | 5000   | Sales      |
      | INCOME  | 1000   | Consulting |
      | EXPENSE | 2000   | Rent       |
      | EXPENSE | 500    | Software   |
    When I request the dashboard summary
    Then the total income should be 6000
    And the total expenses should be 2500
    And the net balance should be 3500
    And the category "Sales" total should be 5000
    And the category "Rent" total should be 2000
