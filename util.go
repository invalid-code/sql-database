package main

import (
	"encoding/binary"
	"os"
)

type Stack[T any] struct {
	data   []T
	length int
}

func (stack *Stack[T]) pop() T {
	poppedItem := stack.data[stack.length]
	stack.data = append(stack.data, stack.data[:stack.length]...)
	return poppedItem
}

func (stack *Stack[T]) push(item T) {
	stack.data = append(stack.data, item)
	stack.length += 1
}

func insert[T any](a []T, index int, value T) []T {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func saveToFile(table Table, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = binary.Write(file, binary.LittleEndian, uint32(table.length))
	if err != nil {
		panic(err)
	}
	rowsStack := Stack[*BTreeNode]{data: []*BTreeNode{&table.rows}, length: 0}
	for len(rowsStack.data) != 0 {
		curBTreeNode := rowsStack.pop()
		switch curBTreeNode.NodeType {
		case Internal:
			for _, childBTreeNode := range curBTreeNode.Children {
				rowsStack.push(childBTreeNode)
			}
		case Leaf:
			for i, key := range curBTreeNode.Keys {
				name, email := curBTreeNode.Data[i].Name, curBTreeNode.Data[i].Email
				err = binary.Write(file, binary.LittleEndian, key)
				if err != nil {
					panic(err)
				}
				err = binary.Write(file, binary.LittleEndian, len(name))
				if err != nil {
					panic(err)
				}
				for _, character := range name {
					err = binary.Write(file, binary.LittleEndian, character)
					if err != nil {
						panic(err)
					}
				}
				err = binary.Write(file, binary.LittleEndian, len(email))
				if err != nil {
					panic(err)
				}
				for _, character := range email {
					err = binary.Write(file, binary.LittleEndian, character)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

func readFile(path string) Table {
	_, err := os.Stat(path)
	var table Table
	var file *os.File
	if err == nil {
		file, err = os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		err = binary.Read(file, binary.LittleEndian, &table.length)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		table = Table{
			rows: BTreeNode{
				IsRoot:   true,
				NodeType: Leaf,
				Parent:   nil,
				Keys:     []int{},
				Data:     []Row{},
				Children: []*BTreeNode{},
			},
			length: 0,
		}
	} else {
		panic(err)
	}
	return table
}
