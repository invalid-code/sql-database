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

func (bTreeNode *BTreeNode) insertKey(key int, data Row, pathIndex []int) *BTreeNode {
	switch bTreeNode.NodeType {
	case Internal:
		var childBTreeNode *BTreeNode
		for i, childKey := range bTreeNode.Keys {
			if key <= childKey {
				childBTreeNode = bTreeNode.Children[i].insertKey(key, data, append(pathIndex, i))
				break
			} else if i == len(bTreeNode.Keys)-1 {
				childBTreeNode = bTreeNode.Children[i+1].insertKey(key, data, append(pathIndex, i+1))
			}
		}
		bTreeNode = childBTreeNode.Parent
		if !bTreeNode.IsRoot {
			bTreeNode.Parent.Children[pathIndex[len(pathIndex)-1]] = bTreeNode
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
			bTreeNode.Data = append(bTreeNode.Data, data)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if key <= nodeKey {
					destKeys := make([]int, len(bTreeNode.Keys))
					copy(destKeys, bTreeNode.Keys)
					destData := make([]Row, len(bTreeNode.Data))
					copy(destData, bTreeNode.Data)
					bTreeNode.Keys = slices.Insert(destKeys, i, key)
					bTreeNode.Data = slices.Insert(destData, i, data)
					break
				} else if i == len(bTreeNode.Keys)-1 {
					destKeys := make([]int, len(bTreeNode.Keys))
					copy(destKeys, bTreeNode.Keys)
					destData := make([]Row, len(bTreeNode.Data))
					copy(destData, bTreeNode.Data)
					bTreeNode.Keys = append(destKeys, key)
					bTreeNode.Data = append(destData, data)
				}
			}
			if len(bTreeNode.Keys) > MAX_KEYS {
				bTreeNode.split(pathIndex)
			}
		}
	}
	return bTreeNode
}

func (bTreeNode *BTreeNode) split(pathIndex []int) {
	curPathIndex := pathIndex[len(pathIndex)-1]
	leftKeys, rightKeys := bTreeNode.Keys[:3], bTreeNode.Keys[3:]
	middleKey := bTreeNode.Keys[2]
	var leftChildren, rightChildren []*BTreeNode
	var leftData, rightData []Row
	switch bTreeNode.NodeType {
	case Internal:
		leftChildren, rightChildren = bTreeNode.Children[:3], bTreeNode.Children[3:]
	case Leaf:
		leftData, rightData = bTreeNode.Data[:3], bTreeNode.Data[3:]
	}
	if bTreeNode.IsRoot {
		bTreeNode.Keys = []int{middleKey}
		switch bTreeNode.NodeType {
		case Internal:
			bTreeNode.Children = []*BTreeNode{}
		case Leaf:
			bTreeNode.Data = []Row{}
		}
		for i := 0; i < 2; i++ {
			childBTreeNode := new(BTreeNode)
			childBTreeNode.IsRoot = false
			switch bTreeNode.NodeType {
			case Internal:
				childBTreeNode.NodeType = Internal
				if i == 0 {
					childBTreeNode.Children = leftChildren
				} else {
					childBTreeNode.Children = rightChildren
				}
			case Leaf:
				childBTreeNode.NodeType = Leaf
			}
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
		if bTreeNode.NodeType == Leaf {
			bTreeNode.NodeType = Internal
		}
	} else {
		bTreeNode.Parent.Keys = slices.Insert(bTreeNode.Parent.Keys, pathIndex[len(pathIndex)-1], middleKey)
		childBTreeNode := new(BTreeNode)
		childBTreeNode.IsRoot = false
		childBTreeNode.NodeType = Internal
		childBTreeNode.Parent = bTreeNode.Parent
		childBTreeNode.Keys = rightKeys
		childBTreeNode.Data = rightData
		bTreeNode.Keys = leftKeys
		bTreeNode.Data = leftData
		switch bTreeNode.NodeType {
		case Internal:
			childBTreeNode.NodeType = Internal
			childBTreeNode.Children = rightChildren
			bTreeNode.Children = leftChildren
		case Leaf:
			childBTreeNode.NodeType = Leaf
		}
		bTreeNode.Parent.Children = slices.Insert(bTreeNode.Parent.Children, curPathIndex+1, childBTreeNode)
		if len(bTreeNode.Parent.Keys) > MAX_KEYS {
			bTreeNode.Parent.split(pathIndex[:len(pathIndex)-1])
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
	table.rows = *table.rows.insertKey(id, data, []int{0})
}

func (table *Table) executeSelect() {
	if table.length == 0 {
		fmt.Println("You have no rows to print")
	} else {
		table.rows.printRows()
	}
}

type Row struct {
	Name  string
	Email string
}
