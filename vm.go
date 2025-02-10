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

func executeStatement(statement StatementType, row Row, table *BTreeNode) {
	switch statement {
	case Insert:
		executeInsert(table, row, 0)
	case Select:
		executeSelect(table)
	}
}
