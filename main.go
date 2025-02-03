package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	MAX_KEYS = 5
)

func insert(a []int, index int, value int) []int {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

type BTreeNodeType int

const (
	Internal BTreeNodeType = iota
	Leaf
)

type BTreeNode struct {
	IsRoot   bool
	NodeType BTreeNodeType
	Children []*BTreeNode
	Keys     []int
}

func (bTreeNode *BTreeNode) insert(key int) {
	switch bTreeNode.NodeType {
	case Internal:
		// for _, childBTreeNode := range bTreeNode.children {
		// 	childBTreeNode.insert(key)
		// }
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if nodeKey < key {
					if i == len(bTreeNode.Keys)-1 {
						bTreeNode.Keys = append(bTreeNode.Keys, key)
					} else {
						bTreeNode.Keys = insert(bTreeNode.Keys, i+1, key)
					}
				}
			}
		}
		if len(bTreeNode.Keys) > MAX_KEYS {
			bTreeNode.split()
		}
	}
}

func (bTreeNode *BTreeNode) split() {
	bTreeNode.NodeType = Internal
	for i := 0; i < 2; i++ {
		childBTreeNode := new(BTreeNode)
		childBTreeNode.IsRoot = false
		childBTreeNode.NodeType = Leaf
		if i == 0 {
			childBTreeNode.Keys = bTreeNode.Keys[0:2]
		} else {
			childBTreeNode.Keys = bTreeNode.Keys[3:]
		}
		bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
	}
	bTreeNode.Keys = bTreeNode.Keys[2:3] // wrong
}

const DB_FILENAME = "persistant.db"

type Row struct {
	Id    int    `gob:"id"`
	Name  string `gob:"name"`
	Email string `gob:"email"`
}

func saveToFile(table BTreeNode) {
	file, err := os.Create(DB_FILENAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(table)
	if err != nil {
		panic(err)
	}
}

func readFile() BTreeNode {
	_, err := os.Stat(DB_FILENAME)
	var table BTreeNode
	var file *os.File
	if err == nil {
		file, err = os.Open(DB_FILENAME)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&table)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		table = BTreeNode{
			IsRoot:   true,
			NodeType: Leaf,
			Keys:     []int{},
			Children: []*BTreeNode{},
		}
	} else {
		panic(err)
	}
	return table
}

func executeCommand(command CommandType, table BTreeNode) {
	switch command {
	case Exit:
		saveToFile(table)
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
}

func executeStatement(statement StatementType, row *Row, table *BTreeNode) {
	switch statement {
	case Insert:
		table.insert(row.Id)
		fmt.Println(table)
	case Select:
	}
}

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

func parseRow(input string) (Row, error) {
	row := Row{Id: 0, Name: "", Email: ""}
	if input[0] != '(' {
		return row, errors.New("invalid array given")
	} else if input[len(input)-1] != ')' {
		return row, errors.New("invalid array given")
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
	row.Id = id
	row.Name = parseQuotedString(inputRow[1])
	row.Email = parseQuotedString(inputRow[2])
	return row, nil
}

func parseStatement(input string) (StatementType, *Row, error) {
	statement, _ := parseWord(input)
	if statement == "insert" {
		row, err := parseRow(input[7:])
		if err != nil {
			panic(err)
		}
		return Insert, &row, nil
	} else if statement == "select" {
		return Select, nil, nil
	}
	return 0, nil, errors.New("unknown statement given")
}

func main() {
	table := readFile()
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

		if input[0] == '.' {
			command, err := parseCommand(input[1:])
			if err != nil {
				fmt.Errorf("%v", err)
				continue
			}
			executeCommand(command, table)
		} else {
			statement, row, err := parseStatement(input)
			if err != nil {
				fmt.Errorf("%v", err)
				continue
			}
			executeStatement(statement, row, &table)
			fmt.Println("Executed!")
		}
	}
}
