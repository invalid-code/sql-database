package main

import (
	"encoding/binary"
	"fmt"
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
		case Leaf:
			for i, key := range curNode.keys {
				err = binary.Write(file, binary.LittleEndian, uint8(key))
				if err != nil {
					panic(err)
				}
				curName := curNode.data[i].name
				err = binary.Write(file, binary.LittleEndian, uint8(len(curName)))
				if err != nil {
					panic(err)
				}
				for _, character := range curName {
					err = binary.Write(file, binary.LittleEndian, uint32(character))
					if err != nil {
						panic(err)
					}
				}
				curEmail := curNode.data[i].email
				err = binary.Write(file, binary.LittleEndian, uint8(len(curEmail)))
				if err != nil {
					panic(err)
				}
				for _, character := range curEmail {
					err = binary.Write(file, binary.LittleEndian, uint32(character))
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
		var tableLength, key, nameLen, emailLen uint8
		err = binary.Read(file, binary.LittleEndian, &tableLength)
		if err != nil {
			panic(err)
		}
		table.length = int(tableLength)
		table.rows.nodeType = Leaf
		table.rows.isRoot = true
		table.rows.keys = make([]int, table.length)
		table.rows.data = make([]Row, table.length)
		for i := 0; i < table.length; i++ {
			err = binary.Read(file, binary.LittleEndian, &key)
			if err != nil {
				panic(err)
			}
			table.rows.keys[i] = int(key)
			err = binary.Read(file, binary.LittleEndian, &nameLen)
			if err != nil {
				panic(err)
			}
			for j := 0; j < int(nameLen); j++ {
				var nameCharacter uint32
				err = binary.Read(file, binary.LittleEndian, &nameCharacter)
				if err != nil {
					panic(err)
				}
				table.rows.data[i].name += string(rune(nameCharacter))
			}
			err = binary.Read(file, binary.LittleEndian, &emailLen)
			if err != nil {
				panic(err)
			}
			for j := 0; j < int(emailLen); j++ {
				var emailCharacter uint32
				err = binary.Read(file, binary.LittleEndian, &emailCharacter)
				if err != nil {
					panic(err)
				}
				table.rows.data[i].email += string(rune(emailCharacter))
			}
		}
		if len(table.rows.keys) > MAX_KEYS {
			table.rows.split([]int{0})
		}
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
