package main

import (
	"encoding/binary"
	"os"
)

func saveToFile(table Table, path string) {
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
		switch curNode.NodeType {
		case Internal:
			for _, childBTreeNode := range curNode.Children {
				nodesQueue.data = nodesQueue.push(nodesQueue.data, childBTreeNode)
			}
		case Leaf:
			for i, key := range curNode.Keys {
				err = binary.Write(file, binary.LittleEndian, uint8(key))
				if err != nil {
					panic(err)
				}
				curName := curNode.Data[i].Name
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
				curEmail := curNode.Data[i].Email
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
		table.rows.NodeType = Leaf
		table.rows.IsRoot = true
		table.rows.Keys = make([]int, table.length)
		table.rows.Data = make([]Row, table.length)
		for i := 0; i < table.length; i++ {
			err = binary.Read(file, binary.LittleEndian, &key)
			if err != nil {
				panic(err)
			}
			table.rows.Keys[i] = int(key)
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
				table.rows.Data[i].Name += string(rune(nameCharacter))
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
				table.rows.Data[i].Email += string(rune(emailCharacter))
			}
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
