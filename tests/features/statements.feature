Feature: Saving to the db

  Background: 
    Given a persistant database
    And a database
    And a table

  Scenario: insert 1 row
    Given a insert command
    When I execute the command
    Then the table should have 1 row
