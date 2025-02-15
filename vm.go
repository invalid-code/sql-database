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

func executeStatement(statement StatementType, id int, row Row, table *Table) {
	switch statement {
	case Insert:
		table.executeInsert(id, row)
	case Select:
		table.executeSelect()
	}
}
