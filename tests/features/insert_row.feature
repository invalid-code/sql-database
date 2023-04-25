Feature: Insert Row

  Background: 
    Given a permanent database
    And a database

  Scenario: insert 1 row
    Given a command
    And a table
    When I execute the command
    Then the command should be executed
