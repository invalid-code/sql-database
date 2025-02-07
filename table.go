package main

import (
	"bytes"
	"encoding/binary"
)

const (
	MAX_KEYS = 4
)

type BTreeNodeType int

const (
	Internal BTreeNodeType = iota
	Leaf
)

type BTreeNode struct {
	IsRoot   bool
	NodeType BTreeNodeType
	Parent   *BTreeNode
	Children []*BTreeNode
	Keys     []int
}

func (bTreeNode *BTreeNode) insert(key int, pathIndex int) {
	switch bTreeNode.NodeType {
	case Internal:
		for i, nodeKey := range bTreeNode.Keys {
			if nodeKey < key {
				bTreeNode.Children[i].insert(key, i)
			}
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if nodeKey < key {
					if i == len(bTreeNode.Keys)-1 {
						bTreeNode.Keys = append(bTreeNode.Keys, key)
					} else {
						keys, _ := convToSlice[int, int](insert(bTreeNode.Keys, i+1, key))
						bTreeNode.Keys = keys
					}
				}
			}
		}
		if len(bTreeNode.Keys) > MAX_KEYS {
			bTreeNode.split(pathIndex)
		}
	}
}

func (bTreeNode *BTreeNode) split(pathIndex int) {
	childBTreeNodes := []*BTreeNode{}
	middleKey := bTreeNode.Keys[2:3]
	for i := 0; i < 2; i++ {
		childBTreeNode := new(BTreeNode)
		childBTreeNode.IsRoot = false
		childBTreeNode.NodeType = Leaf
		if i == 0 {
			childBTreeNode.Keys = bTreeNode.Keys[0:3]
		} else {
			childBTreeNode.Keys = bTreeNode.Keys[3:]
		}
		childBTreeNodes = append(childBTreeNodes, childBTreeNode)
	}
	if bTreeNode.IsRoot {
		bTreeNode.NodeType = Internal
		bTreeNode.Keys = middleKey
	} else if bTreeNode.NodeType == Leaf {
		bTreeNode.Parent.Children = remove(bTreeNode.Parent.Children, pathIndex)
	}
	for i, childBTreeNode := range childBTreeNodes {
		if bTreeNode.IsRoot {
			childBTreeNode.Parent = bTreeNode
			bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
		} else if bTreeNode.NodeType == Leaf {
			childBTreeNode.Parent = bTreeNode.Parent
			bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex+i, childBTreeNode)
		}
	}
	if len(bTreeNode.Parent.Keys) > MAX_KEYS {
		bTreeNode.Parent.split(pathIndex)
	}
}

const DB_FILENAME = "persistant.db"

type Row struct {
	Id    int    `gob:"id"`
	Name  string `gob:"name"`
	Email string `gob:"email"`
}

func serializeTable(buffer *bytes.Buffer, table *BTreeNode) {
	err := binary.Write(buffer, binary.LittleEndian, table.IsRoot)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buffer, binary.LittleEndian, table.NodeType)
	if err != nil {
		panic(err)
	}
	err = binary.Write(buffer, binary.LittleEndian, table.Keys)
	if err != nil {
		panic(err)
	}
	if table.Parent != nil {
		err = binary.Write(buffer, binary.LittleEndian, table.Keys)
		if err != nil {
			panic(err)
		}
	}
	for _, childBTreeNode := range table.Children {
		serializeTable(buffer, childBTreeNode)
	}
}
