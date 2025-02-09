package main

import (
	"errors"
	"strconv"
)

func parseWord(input string) (string, int) {
	parsedInput, offset := "", 0
	for i, character := range input {
		if character == ' ' {
			offset = i
			break
		}
		parsedInput += string(character)
	}
	return parsedInput, offset
}

type CommandType int

const (
	Exit CommandType = iota
)

func parseCommand(input string) (CommandType, error) {
	cmd, _ := parseWord(input)
	if cmd == "exit" {
		return Exit, nil
	}
	return 0, errors.New("unknown command given")
}

type StatementType int

const (
	Insert StatementType = iota
	Select
)

func parseQuotedString(input string) string {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		return input[1 : len(input)-1]
	}
	return input
}

func parseRow(input string) (int, Row, error) {
	row := Row{Name: "a", Email: "a"}
	if input[0] != '(' {
		return 0, row, errors.New("invalid array given")
	} else if input[len(input)-1] != ')' {
		return 0, row, errors.New("invalid array given")
	}
	offset := 1
	inputRow := []string{}
	for i := 0; i < 3; i++ {
		parsedInput, parsedOffset := parseWord(input[offset : len(input)-1])
		offset += parsedOffset + 1
		inputRow = append(inputRow, parsedInput)
	}
	id, err := strconv.Atoi(inputRow[0])
	if err != nil {
		panic(err)
	}
	row.Name = parseQuotedString(inputRow[1])
	row.Email = parseQuotedString(inputRow[2])
	return id, row, nil
}

func parseStatement(input string) (StatementType, int, Row, error) {
	statement, _ := parseWord(input)
	if statement == "insert" {
		id, row, err := parseRow(input[7:])
		if err != nil {
			panic(err)
		}
		return Insert, id, row, nil
	} else if statement == "select" {
		return Select, 0, Row{}, nil
	}
	return 0, 0, Row{}, errors.New("unknown statement given")
}
