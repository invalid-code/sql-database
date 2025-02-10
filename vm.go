package main

import (
	"fmt"
	"os"
)

func executeCommand(command CommandType, table BTreeNode) {
	switch command {
	case Exit:
		saveToFile(table, DB_FILENAME)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
}

func executeStatement(statement StatementType, id int, row Row, table *BTreeNode) {
	switch statement {
	case Insert:
		executeInsert(table, id, row, 0)
	case Select:
		executeSelect(table)
	}
}
