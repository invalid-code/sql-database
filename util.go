package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func insert[T any](a []T, index int, value T) []T {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func saveToFile(table BTreeNode) {
	file, err := os.Create(DB_FILENAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var bufferData []byte
	buffer := bytes.NewBuffer(bufferData)
	serializeTable(buffer, &table)
	_, err = buffer.WriteTo(file)
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
		var bufferData []byte
		buffer := bytes.NewBuffer(bufferData)
		_, err := buffer.ReadFrom(file)
		if err != nil {
			panic(err)
		}
		err = binary.Read(buffer, binary.LittleEndian, &table)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		table = BTreeNode{
			IsRoot:   true,
			NodeType: Leaf,
			Parent:   nil,
			Keys:     []int{},
			Children: []*BTreeNode{},
		}
	} else {
		panic(err)
	}
	return table
}

func convToSlice[T any, U any](in []T) ([]U, error) {
	var res []U
	for _, elem := range in {
		if val, ok := interface{}(elem).(U); ok {
			res = append(res, val)
		} else {
			return nil, fmt.Errorf("%v is not of type %v", val, *new(T))
		}
	}
	return res, nil
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
