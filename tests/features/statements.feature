Feature: Saving to the db

  Background: 
    Given a persistant database
    And a database
    And a table

  Scenario: insert 1 row
    Given a insert command
    When I execute the command
    Then the table should have 1 row

  Scenario: create a database
    Given a create database command
    When I execute the command
    Then the persistant database should have 1 database

  Scenario: create a table
    Given a create table command
    When I execute the command
    Then the database should have 1 table
