Feature: Insert Row

  Background: 
    Given A Table
    * A Execute Result

  Scenario: insert 1 row
    Given A Statement
    When I execute A Statement
    Then A Statement should be executed

  Scenario: insert max row
    Given A Max Statement
    When I execute A Statement until table is full
    Then A Table is full
