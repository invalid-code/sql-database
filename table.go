package main

import (
	"errors"
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
	isRoot   bool
	nodeType BTreeNodeType
	parent   *BTreeNode
	children []*BTreeNode
	keys     []int
	data     []Row
}

func (bTreeNode *BTreeNode) insertKey(key int, data Row, pathIndex []int) *BTreeNode {
	switch bTreeNode.nodeType {
	case Internal:
		var childBTreeNode *BTreeNode
		for i, childKey := range bTreeNode.keys {
			if key <= childKey {
				childBTreeNode = bTreeNode.children[i].insertKey(key, data, append(pathIndex, i))
				break
			} else if i == len(bTreeNode.keys)-1 {
				childBTreeNode = bTreeNode.children[i+1].insertKey(key, data, append(pathIndex, i+1))
			}
		}
		if childBTreeNode.parent != nil {
			bTreeNode = childBTreeNode.parent
			if bTreeNode.parent != nil {
				childI := 0
				if pathIndex[len(pathIndex)-1] > len(bTreeNode.parent.children) {
					childI = len(bTreeNode.parent.children) - 1
				} else {
					childI = pathIndex[len(pathIndex)-1]
				}
				bTreeNode.parent.children[childI] = bTreeNode
			}
		}
	case Leaf:
		if len(bTreeNode.keys) == 0 {
			bTreeNode.keys = append(bTreeNode.keys, key)
			bTreeNode.data = append(bTreeNode.data, data)
		} else {
			for i, nodeKey := range bTreeNode.keys {
				if key <= nodeKey {
					destKeys := make([]int, len(bTreeNode.keys))
					copy(destKeys, bTreeNode.keys)
					destData := make([]Row, len(bTreeNode.data))
					copy(destData, bTreeNode.data)
					bTreeNode.keys = slices.Insert(destKeys, i, key)
					bTreeNode.data = slices.Insert(destData, i, data)
					break
				} else if i == len(bTreeNode.keys)-1 {
					destKeys := make([]int, len(bTreeNode.keys))
					copy(destKeys, bTreeNode.keys)
					destData := make([]Row, len(bTreeNode.data))
					copy(destData, bTreeNode.data)
					bTreeNode.keys = append(destKeys, key)
					bTreeNode.data = append(destData, data)
				}
			}
			if len(bTreeNode.keys) > MAX_KEYS {
				bTreeNode = bTreeNode.split(pathIndex)
			}
		}
	}
	return bTreeNode
}

func (bTreeNode *BTreeNode) split(pathIndex []int) *BTreeNode {
	curPathIndex := pathIndex[len(pathIndex)-1]
	leftKeys, rightKeys := bTreeNode.keys[:3], bTreeNode.keys[3:]
	middleKey := bTreeNode.keys[2]
	var leftChildren, rightChildren []*BTreeNode
	var leftData, rightData []Row
	switch bTreeNode.nodeType {
	case Internal:
		leftChildren, rightChildren = bTreeNode.children[:3], bTreeNode.children[3:]
	case Leaf:
		leftData, rightData = bTreeNode.data[:3], bTreeNode.data[3:]
	}
	if bTreeNode.isRoot {
		bTreeNode.keys = []int{middleKey}
		switch bTreeNode.nodeType {
		case Internal:
			bTreeNode.children = []*BTreeNode{}
		case Leaf:
			bTreeNode.data = []Row{}
		}
		for i := range 2 {
			childBTreeNode := new(BTreeNode)
			childBTreeNode.isRoot = false
			switch bTreeNode.nodeType {
			case Internal:
				childBTreeNode.nodeType = Internal
				if i == 0 {
					for leftChildI := range len(leftChildren) {
						leftChildren[leftChildI].parent = childBTreeNode
					}
					childBTreeNode.children = leftChildren
				} else {
					for rightChildI := range len(rightChildren) {
						rightChildren[rightChildI].parent = childBTreeNode
					}
					childBTreeNode.children = rightChildren
				}
			case Leaf:
				childBTreeNode.nodeType = Leaf
			}
			childBTreeNode.parent = bTreeNode
			if i == 0 {
				childBTreeNode.keys = leftKeys
				childBTreeNode.data = leftData
			} else {
				childBTreeNode.keys = rightKeys
				childBTreeNode.data = rightData
			}
			bTreeNode.children = append(bTreeNode.children, childBTreeNode)
		}
		if bTreeNode.nodeType == Leaf {
			bTreeNode.nodeType = Internal
		}
		return bTreeNode
	} else {
		bTreeNode.parent.keys = slices.Insert(bTreeNode.parent.keys, pathIndex[len(pathIndex)-1], middleKey)
		childBTreeNode := new(BTreeNode)
		childBTreeNode.isRoot = false
		childBTreeNode.nodeType = Internal
		childBTreeNode.parent = bTreeNode.parent
		childBTreeNode.keys = rightKeys
		childBTreeNode.data = rightData
		bTreeNode.keys = leftKeys
		bTreeNode.data = leftData
		switch bTreeNode.nodeType {
		case Internal:
			childBTreeNode.nodeType = Internal
			childBTreeNode.children = rightChildren
			for _, childBTreeNodeChild := range childBTreeNode.children {
				childBTreeNodeChild.parent = childBTreeNode
			}
			bTreeNode.children = leftChildren
		case Leaf:
			childBTreeNode.nodeType = Leaf
		}
		bTreeNode.parent.children = slices.Insert(bTreeNode.parent.children, curPathIndex+1, childBTreeNode)
		if len(bTreeNode.parent.keys) > MAX_KEYS {
			bTreeNode.parent = bTreeNode.parent.split(pathIndex[:len(pathIndex)-1])
			bTreeNode.parent = childBTreeNode.parent
		}
		return bTreeNode.parent
	}
}

func (bTreeNode *BTreeNode) printRows() {
	switch bTreeNode.nodeType {
	case Internal:
		for _, childBTreeNode := range bTreeNode.children {
			childBTreeNode.printRows()
		}
	case Leaf:
		for i, key := range bTreeNode.keys {
			fmt.Printf("id: %v, name: %v, email: %v\n", key, bTreeNode.data[i].name, bTreeNode.data[i].email)
		}
	}
}

func (bTreeNode *BTreeNode) printTree(level int) {
	for range level {
		fmt.Printf("\t")
	}
	fmt.Printf("%v\n", bTreeNode)
	for _, childBTreeNode := range bTreeNode.children {
		childBTreeNode.printTree(level + 1)
	}
}

func (bTreeNode *BTreeNode) Equals(other *BTreeNode) bool {
	isRootEq := bTreeNode.isRoot == other.isRoot
	nodeTypeEq := bTreeNode.nodeType == other.nodeType
	keysEq := true
	for i, key := range bTreeNode.keys {
		if key != other.keys[i] {
			keysEq = false
			break
		}
	}
	childrenEq := true
	for i, childBTreeNode := range bTreeNode.children {
		if !childBTreeNode.Equals(other.children[i]) {
			childrenEq = false
			break
		}
	}
	// todo
	// parentEq := true
	// if bTreeNode.parent != nil && other.parent != nil {
	// 	parentEq = bTreeNode.parent.Equals(other.parent)
	// }
	return (isRootEq && nodeTypeEq && keysEq && childrenEq)
}

const DB_FILENAME = "persistant.db"

type Table struct {
	rows   BTreeNode
	length int
}

func (table *Table) executeInsert(id int, data Row) error {
	if len(data.email) > 255 || len(data.email) < 1 {
		return errors.New("email length must be between 1 and 255")
	}
	if len(data.name) > 33 || len(data.name) < 1 {
		return errors.New("name length must be between 1 and 33")
	}
	table.length += 1
	table.rows = *table.rows.insertKey(id, data, []int{0})
	return nil
}

func (table *Table) executeSelect() {
	if table.length == 0 {
		fmt.Println("You have no rows to print")
	} else {
		table.rows.printRows()
	}
}

type Row struct {
	name  string
	email string
}
