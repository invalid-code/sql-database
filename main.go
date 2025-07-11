package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	table := readFile(DB_FILENAME)
	for {
		fmt.Printf("input> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Errorf("%v", err)
			continue
		}
		input := scanner.Text()

		if len(input) < 1 {
			continue
		} else if input[0] == '.' {
			command, err := parseCommand(input[1:])
			if err != nil {
				fmt.Println(err)
				continue
			}
			executeCommand(command, table)
		} else {
			statement, id, row, err := parseStatement(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = executeStatement(statement, id, row, &table)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Executed!")
			table.rows.printTree(0)
		}
	}
}
