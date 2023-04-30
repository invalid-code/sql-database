Feature: testing persistancy

  Scenario: save a database
    Given a open command
    And a create command
    And a persistant database
    When I execute all the commands
    When I save the database
    Then the database should have been saved
