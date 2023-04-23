Feature: Data Persistance

  Scenario: data should be permanent
    Given a select statement
    When I execute the statement
    Then the previous data should be saved
