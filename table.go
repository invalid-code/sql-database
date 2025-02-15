package main

import "fmt"

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

func (bTreeNode *BTreeNode) insertKey(key int, data Row, pathIndex int) *BTreeNode {
	switch bTreeNode.NodeType {
	case Internal:
		for i, childBTreeNodeKey := range bTreeNode.Keys {
			if key <= childBTreeNodeKey {
				bTreeNode.Children[i] = bTreeNode.Children[i].insertKey(key, data, i)
				break
			} else if i == len(bTreeNode.Keys)-1 {
				panic("todo")
			}
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
			bTreeNode.Data = append(bTreeNode.Data, data)
		} else {
			for i, childBTreeNodeKey := range bTreeNode.Keys {
				if key <= childBTreeNodeKey {
					bTreeNode.Keys = insert(bTreeNode.Keys, i, key)
					bTreeNode.Data = insert(bTreeNode.Data, i, data)
				} else if i == len(bTreeNode.Keys)-1 {
					bTreeNode.Keys = append(bTreeNode.Keys, key)
					bTreeNode.Data = append(bTreeNode.Data, data)
				}
			}
			if len(bTreeNode.Keys) > MAX_KEYS {
				return bTreeNode.split(pathIndex)
			}
		}
	}
	return bTreeNode
}

func (bTreeNode *BTreeNode) split(pathIndex int) *BTreeNode {
	switch bTreeNode.NodeType {
	case Internal:
		panic("todo")
	case Leaf:
		leftKeys, rightKeys := bTreeNode.Keys[:3], bTreeNode.Keys[3:]
		leftData, rightData := bTreeNode.Data[:3], bTreeNode.Data[3:]
		middleKey := bTreeNode.Keys[2]
		bTreeNode.Keys = leftKeys
		bTreeNode.Data = leftData
		if bTreeNode.IsRoot {
			bTreeNode.Parent = new(BTreeNode)
			bTreeNode.Parent.NodeType = Internal
			bTreeNode.Parent.Keys = append(bTreeNode.Parent.Keys, middleKey)
			bTreeNode.Parent.IsRoot = true
			bTreeNode.IsRoot = false
		}
		childBTreeNode := new(BTreeNode)
		childBTreeNode.IsRoot = false
		childBTreeNode.NodeType = Leaf
		childBTreeNode.Keys = rightKeys
		childBTreeNode.Data = rightData
		childBTreeNode.Parent = bTreeNode.Parent
		fmt.Println("before")
		bTreeNode.Parent.Parent.printTree(0)
		bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex, []*BTreeNode{bTreeNode, childBTreeNode}...)
		fmt.Println("after")
		bTreeNode.Parent.Parent.printTree(0)
		return bTreeNode.Parent
	}
	return nil
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
