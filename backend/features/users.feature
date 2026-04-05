Feature: User Management
  As an administrator
  I want to manage system users
  So that I can control access levels and status

  Scenario: Admin can list all users
    Given I am an authenticated "ADMIN"
    When I list all users
    Then the response status should be 200
    And the response should contain a list of users

  Scenario: Admin can create a new user
    Given I am an authenticated "ADMIN"
    When I create a user with:
      | username  | sarah_analyst |
      | role      | ANALYST       |
      | is_active | true          |
    Then the user should be created successfully
    And the user "sarah_analyst" should have role "ANALYST"

  Scenario: Non-admin cannot manage users
    Given I am an authenticated "ANALYST"
    When I list all users
    Then I should receive a "Forbidden" error
