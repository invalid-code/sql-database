Feature: Insert Row

  Background: 
    Given A Table
    And A Execute Result

  Scenario: insert 1 row
    Given A Statement
    When I execute A Statement
    Then A Statement should be executed
