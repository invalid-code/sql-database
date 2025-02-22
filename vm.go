package main

import (
	"fmt"
	"os"
)

func executeCommand(command CommandType, table Table) {
	switch command {
	case Exit:
		saveToFile(table, DB_FILENAME)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
}

func executeStatement(statement StatementType, id int, row Row, table *Table) error {
	switch statement {
	case Insert:
		err := table.executeInsert(id, row)
		if err != nil {
			return err
		}
	case Select:
		table.executeSelect()
	}
	return nil
}
