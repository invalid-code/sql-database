Feature: running meta commands

  Scenario: open a database
    Given a open command
    When I execute the meta command
    Then there should be a persistent database
