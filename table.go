package main

import (
	"fmt"
	"slices"
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
	Data     []Row
}

func (bTreeNode *BTreeNode) insertKey(key int, data Row, pathIndex int) {
	switch bTreeNode.NodeType {
	case Internal:
		var foundIndex int
		for i, childKey := range bTreeNode.Keys {
			if key <= childKey {
				foundIndex = i
				bTreeNode.Children[i].insertKey(key, data, i)
				break
			} else if i == len(bTreeNode.Keys)-1 {
				foundIndex = i + 1
				bTreeNode.Children[i+1].insertKey(key, data, i+1)
			}
		}
		bTreeNode = bTreeNode.Children[foundIndex].Parent
		if !bTreeNode.IsRoot {
			bTreeNode.Parent.Children[pathIndex] = bTreeNode
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
			bTreeNode.Data = append(bTreeNode.Data, data)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if key <= nodeKey {
					bTreeNode.Keys = slices.Insert(bTreeNode.Keys, i, key)
					bTreeNode.Data = slices.Insert(bTreeNode.Data, i, data)
					break
				} else if i == len(bTreeNode.Keys)-1 {
					bTreeNode.Keys = append(bTreeNode.Keys, key)
					bTreeNode.Data = append(bTreeNode.Data, data)
				}
			}
			if len(bTreeNode.Keys) > MAX_KEYS {
				bTreeNode.split(pathIndex)
			}
		}
	}
}

func (bTreeNode *BTreeNode) split(pathIndex int) {
	leftKeys, rightKeys := bTreeNode.Keys[:3], bTreeNode.Keys[3:]
	leftData, rightData := bTreeNode.Data[:3], bTreeNode.Data[3:]
	middleKey := bTreeNode.Keys[2]
	switch bTreeNode.NodeType {
	case Internal:
		leftChildren, rightChildren := bTreeNode.Children[:3], bTreeNode.Children[3:]
		if bTreeNode.IsRoot {
			bTreeNode.Children = []*BTreeNode{}
			bTreeNode.Keys = []int{middleKey}
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Internal
				childBTreeNode.Parent = bTreeNode
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Data = leftData
					childBTreeNode.Children = leftChildren
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Data = rightData
					childBTreeNode.Children = rightChildren
				}
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			}
		} else {
			childBTreeNode := new(BTreeNode)
			childBTreeNode.IsRoot = false
			childBTreeNode.NodeType = Internal
			childBTreeNode.Parent = bTreeNode.Parent
			childBTreeNode.Keys = rightKeys
			childBTreeNode.Data = rightData
			childBTreeNode.Children = rightChildren
			bTreeNode.Keys = leftKeys
			bTreeNode.Data = leftData
			bTreeNode.Children = leftChildren
			bTreeNode.Parent.Children = slices.Insert(bTreeNode.Parent.Children, pathIndex+1, childBTreeNode)
		}
	case Leaf:
		if bTreeNode.IsRoot {
			bTreeNode.NodeType = Internal
			bTreeNode.Keys = []int{middleKey}
			bTreeNode.Data = []Row{}
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Leaf
				childBTreeNode.Parent = bTreeNode
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Data = leftData
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Data = rightData
				}
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			}
		} else {
			childBTreeNode := new(BTreeNode)
			childBTreeNode.IsRoot = false
			childBTreeNode.NodeType = Leaf
			childBTreeNode.Parent = bTreeNode.Parent
			childBTreeNode.Keys = rightKeys
			childBTreeNode.Data = rightData
			bTreeNode.Keys = leftKeys
			bTreeNode.Data = leftData
			bTreeNode.Parent.Children = slices.Insert(bTreeNode.Parent.Children, pathIndex+1, childBTreeNode)
		}
	}
}

func (bTreeNode *BTreeNode) printRows() {
	switch bTreeNode.NodeType {
	case Internal:
		for _, childBTreeNode := range bTreeNode.Children {
			childBTreeNode.printRows()
		}
	case Leaf:
		for i, key := range bTreeNode.Keys {
			fmt.Printf("id: %v, name: %v, email: %v\n", key, bTreeNode.Data[i].Name, bTreeNode.Data[i].Email)
		}
	}
}

func (bTreeNode *BTreeNode) printTree(level int) {
	for i := 0; i < level; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("%v\n", bTreeNode)
	for _, childBTreeNode := range bTreeNode.Children {
		childBTreeNode.printTree(level + 1)
	}
}

func (bTreeNode *BTreeNode) Equals(other *BTreeNode) bool {
	isRootEq := bTreeNode.IsRoot == other.IsRoot
	nodeTypeEq := bTreeNode.NodeType == other.NodeType
	keysEq := true
	for i, key := range bTreeNode.Keys {
		if key != other.Keys[i] {
			keysEq = false
			break
		}
	}
	childrenEq := true
	for i, childBTreeNode := range bTreeNode.Children {
		if !childBTreeNode.Equals(other.Children[i]) {
			childrenEq = false
			break
		}
	}
	return (isRootEq && nodeTypeEq && keysEq && childrenEq)
}

const DB_FILENAME = "persistant.db"

type Table struct {
	rows   BTreeNode
	length int
}

func (table *Table) executeInsert(id int, data Row) {
	table.length += 1
	table.rows = *table.rows.insertKey(id, data, 0)
}

func (table *Table) executeSelect() {
	table.rows.printRows()
}

type Row struct {
	Name  string
	Email string
}
