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

func executeStatement(statement StatementType, row *Row, table *BTreeNode) {
	switch statement {
	case Insert:
		table.insert(row.Id, 0)
		fmt.Println("inside executeStatement")
		table.printTree(0)
	case Select:
	}
}
