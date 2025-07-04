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
		err = binary.Write(file, binary.LittleEndian, curNode.isRoot)
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
		switch curNode.nodeType {
		case Internal:
			err = binary.Write(file, binary.LittleEndian, uint8(0))
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.children)))
			if err != nil {
				panic(err)
			}
			for _, childBTreeNode := range curNode.children {
				nodesQueue.data = nodesQueue.push(nodesQueue.data, childBTreeNode)
			}
		case Leaf:
			err = binary.Write(file, binary.LittleEndian, uint8(1))
			if err != nil {
				panic(err)
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
		var tableLength, isRootByte, nodeType, keysLen, curKey, childrenLen, dataLen, nameLen, nameByte, emailLen, emailByte uint8
		err = binary.Read(file, binary.LittleEndian, &tableLength)
		if err != nil {
			panic(err)
		}
		table.length = int(tableLength)
		bTreeNodeQueue := Queue[*BTreeNode]{}
		bTreeNodeQueue.data = bTreeNodeQueue.push(bTreeNodeQueue.data, &(table.rows))
		for {
			poppedBTreeNodeQueue, err := bTreeNodeQueue.pop()
			if err != nil {
				break
			}
			curBTreeNode := *poppedBTreeNodeQueue
			err = binary.Read(file, binary.LittleEndian, &isRootByte)
			if err != nil {
				panic(err)
			}
			curBTreeNode.isRoot = isRootByte == 1
			err = binary.Read(file, binary.LittleEndian, &keysLen)
			if err != nil {
				panic(err)
			}
			for range keysLen {
				err = binary.Read(file, binary.LittleEndian, &curKey)
				if err != nil {
					panic(err)
				}
				curBTreeNode.keys = append(curBTreeNode.keys, int(curKey))
			}
			err = binary.Read(file, binary.LittleEndian, &nodeType)
			if err != nil {
				panic(err)
			}
			switch nodeType {
			case 0:
				curBTreeNode.nodeType = Internal
			case 1:
				curBTreeNode.nodeType = Leaf
			}
			switch curBTreeNode.nodeType {
			case Internal:
				err = binary.Read(file, binary.LittleEndian, &childrenLen)
				if err != nil {
					panic(err)
				}
				for range childrenLen {
					newBTreeNode := BTreeNode{
						isRoot:   false,
						nodeType: Leaf,
						parent:   curBTreeNode,
						children: []*BTreeNode{},
						keys:     []int{},
						data:     []Row{},
					}
					bTreeNodeQueue.data = bTreeNodeQueue.push(bTreeNodeQueue.data, &newBTreeNode)
					curBTreeNode.children = append(curBTreeNode.children, &newBTreeNode)
				}
			case Leaf:
				err = binary.Read(file, binary.LittleEndian, &dataLen)
				if err != nil {
					panic(err)
				}
				for range dataLen {
					err = binary.Read(file, binary.LittleEndian, &nameLen)
					if err != nil {
						panic(err)
					}
					name := ""
					for range nameLen {
						err = binary.Read(file, binary.LittleEndian, &nameByte)
						if err != nil {
							panic(err)
						}
						name += string(nameByte)
					}
					err = binary.Read(file, binary.LittleEndian, &emailLen)
					if err != nil {
						panic(err)
					}
					email := ""
					for range nameLen {
						err = binary.Read(file, binary.LittleEndian, &emailByte)
						if err != nil {
							panic(err)
						}
						email += string(emailByte)
					}
					curBTreeNode.data = append(curBTreeNode.data, Row{name, email})
				}
			}
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
