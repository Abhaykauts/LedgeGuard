Feature: User Authentication and Authorization
  As a registered user
  I want to log in and manage my session
  So that I can access secure financial data

  Scenario: Successful Login
    Given a user exists with username "admin" and password "password123" and role "ADMIN"
    When I login with username "admin" and password "password123"
    Then I should receive a valid access token
    And I should receive a valid refresh token
    And my role should be "ADMIN"

  Scenario: Failed Login with Wrong Password
    Given a user exists with username "analyst" and password "pass123"
    When I login with username "analyst" and password "wrong-pass"
    Then I should receive an "Unauthorized" error

  Scenario: Token Rotation (Refresh)
    Given a user exists with username "admin" and password "password123" and role "ADMIN"
    And I login with username "admin" and password "password123"
    When I request a new access token using the refresh token
    Then I should receive a new valid access token
    And I should receive a new valid refresh token
