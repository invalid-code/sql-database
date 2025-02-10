package main

import (
	"bufio"
	"io"
	"os"

	"github.com/Sereal/Sereal/Go/sereal"
)

func insert[T any](a []T, index int, value T) []T {
	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

func saveToFile(table BTreeNode, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := sereal.NewEncoder()
	bytes, err := encoder.Marshal(table)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func readFile(path string) BTreeNode {
	_, err := os.Stat(path)
	var table BTreeNode
	var file *os.File
	if err == nil {
		file, err = os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		bytes, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		decoder := sereal.NewDecoder()
		err = decoder.Unmarshal(bytes, &table)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		table = BTreeNode{
			IsRoot:   true,
			NodeType: Leaf,
			Parent:   nil,
			Keys:     []int{},
			Data:     []Row{},
			Children: []*BTreeNode{},
		}
	} else {
		panic(err)
	}
	return table
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
