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
	nodesPathQueue := Queue[int]{data: []int{}}
	nodesPathQueue.data = nodesPathQueue.push(nodesPathQueue.data, 0)
	for !nodesQueue.isEmpty() {
		queueData, err := nodesQueue.pop()
		if err != nil {
			panic(err)
		}
		pathQueueData, err := nodesPathQueue.pop()
		if err != nil {
			panic(err)
		}
		err = binary.Write(file, binary.LittleEndian, uint8(*pathQueueData))
		if err != nil {
			panic(err)
		}
		curNode := *queueData
		switch curNode.nodeType {
		case Internal:
			for i, childBTreeNode := range curNode.children {
				nodesQueue.data = nodesQueue.push(nodesQueue.data, childBTreeNode)
				nodesPathQueue.data = nodesPathQueue.push(nodesPathQueue.data, i)
			}
		case Leaf:
			err = binary.Write(file, binary.LittleEndian, uint8(255))
			if err != nil {
				panic(err)
			}
			err = binary.Write(file, binary.LittleEndian, uint8(len(curNode.keys)))
			if err != nil {
				panic(err)
			}
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
		var tableLength, bTreeIndex, key, nameLen, emailLen, keysLen uint8
		err = binary.Read(file, binary.LittleEndian, &tableLength)
		if err != nil {
			panic(err)
		}
		table.length = int(tableLength)
		isRootIndex := true
		curBTreeNode := &table.rows
		for {
			err = binary.Read(file, binary.LittleEndian, &bTreeIndex)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			if isRootIndex {
				*curBTreeNode = BTreeNode{isRoot: true, nodeType: Leaf, parent: nil, children: []*BTreeNode{}, keys: []int{}, data: []Row{}}
				isRootIndex = false
				continue
			}
			if bTreeIndex == 255 {
				err = binary.Read(file, binary.LittleEndian, &keysLen)
				if err != nil {
					panic(err)
				}
				curBTreeNode.keys = make([]int, keysLen)
				curBTreeNode.data = make([]Row, keysLen)
				for i := 0; i < int(keysLen); i++ {
					err = binary.Read(file, binary.LittleEndian, &key)
					if err != nil {
						panic(err)
					}
					curBTreeNode.keys[i] = int(key)
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
						curBTreeNode.data[i].name = string(rune(nameCharacter))
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
						curBTreeNode.data[i].name = string(rune(emailCharacter))
					}
				}
			} else {
				// if curBTreeNode.nodeType == Leaf {
				// 	(*curBTreeNode).nodeType = Internal
				// }
				// (*curBTreeNode).children[bTreeIndex] = &BTreeNode{isRoot: false, nodeType: Leaf, parent: curBTreeNode, children: []*BTreeNode{}, keys: []int{}, data: []Row{}}
				// curBTreeNode = curBTreeNode.children[bTreeIndex]
				// curBTreeNodeParent := curBTreeNode
				// for {
				// 	if curBTreeNodeParent.parent == nil {
				// 		table.rows = *curBTreeNodeParent
				// 		break
				// 	}
				// 	curBTreeNodeParent = curBTreeNodeParent.parent
				// }
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
