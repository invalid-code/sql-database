package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func saveToFile(table Table, path string) {
	if table.length == 0 {
		fmt.Println("No data to save")
		return
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = binary.Write(file, binary.LittleEndian, uint8(table.length))
	if err != nil {
		panic(err)
	}
	nodesQueue := Queue[*BTreeNode]{data: []*BTreeNode{}}
	nodesQueue.data = nodesQueue.push(nodesQueue.data, &(table.rows))
	for !nodesQueue.isEmpty() {
		queueData, err := nodesQueue.pop()
		if err != nil {
			panic(err)
		}
		curNode := *queueData
		switch curNode.nodeType {
		case Internal:
			for _, childBTreeNode := range curNode.children {
				nodesQueue.data = nodesQueue.push(nodesQueue.data, childBTreeNode)
			}
			err = binary.Write(file, binary.LittleEndian, curNode.isRoot)
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(0))
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.keys)))
			if err != nil {
				panic(err)
			}
			for _, key := range curNode.keys {
				err = binary.Write(file, binary.LittleEndian, uint8(key))
				if err != nil {
					panic(err)
				}
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.children)))
			if err != nil {
				panic(err)
			}
		case Leaf:
			err = binary.Write(file, binary.LittleEndian, curNode.isRoot)
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(1))
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.keys)))
			if err != nil {
				panic(err)
			}
			for _, key := range curNode.keys {
				err = binary.Write(file, binary.LittleEndian, uint8(key))
				if err != nil {
					panic(err)
				}
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.data)))
			if err != nil {
				panic(err)
			}
			for _, data := range curNode.data {
				err = binary.Write(file, binary.LittleEndian, uint8(len(data.name)))
				if err != nil {
					panic(err)
				}
				_, err = file.Write([]byte(data.name))
				if err != nil {
					panic(err)
				}
				err = binary.Write(file, binary.LittleEndian, uint8(len(data.email)))
				if err != nil {
					panic(err)
				}
				_, err = file.Write([]byte(data.email))
				if err != nil {
					panic(err)
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
		var tableLength uint8
		err = binary.Read(file, binary.LittleEndian, &tableLength)
		if err != nil {
			panic(err)
		}
		table.length = int(tableLength)
	} else if os.IsNotExist(err) {
		table = Table{
			rows: BTreeNode{
				isRoot:   true,
				nodeType: Leaf,
				parent:   nil,
				keys:     []int{},
				data:     []Row{},
				children: []*BTreeNode{},
			},
			length: 0,
		}
	} else {
		panic(err)
	}
	return table
}
