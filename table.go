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

func executeInsert(bTreeNode *BTreeNode, key int, row Row, pathIndex int) {
	switch bTreeNode.NodeType {
	case Internal:
		for i, nodeKey := range bTreeNode.Keys {
			if nodeKey >= key {
				executeInsert(bTreeNode.Children[i], key, row, i)
				break
			} else if len(bTreeNode.Keys)-1 == i {
				executeInsert(bTreeNode.Children[i], key, row, i)
			}
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
			bTreeNode.Data = append(bTreeNode.Data, row)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if nodeKey > key {
					bTreeNode.Keys = insert(bTreeNode.Keys, i, key)
					bTreeNode.Data = insert(bTreeNode.Data, i, row)
					break
				}
				if i == len(bTreeNode.Keys)-1 {
					bTreeNode.Keys = append(bTreeNode.Keys, key)
					bTreeNode.Data = append(bTreeNode.Data, row)
				}
			}
		}
		if len(bTreeNode.Keys) > MAX_KEYS {
			split(bTreeNode, pathIndex)
		}
	}

}

func split(bTreeNode *BTreeNode, pathIndex int) {
	leftKeys, rightKeys := bTreeNode.Keys[:3], bTreeNode.Keys[3:]
	middleKey := bTreeNode.Keys[2]
	leftData, rightData := bTreeNode.Data[:3], bTreeNode.Data[3:]
	switch bTreeNode.NodeType {
	case Internal:
		leftChildren, rightChildren := bTreeNode.Children[:3], bTreeNode.Children[3:]
		if bTreeNode.IsRoot {
			bTreeNode.Keys = []int{middleKey}
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Internal
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Children = leftChildren
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Children = rightChildren
				}
				childBTreeNode.Parent = bTreeNode
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			}
		} else {
			bTreeNode.Parent.Keys = insert(bTreeNode.Parent.Keys, pathIndex, middleKey)
			bTreeNode.Parent.Children = remove(bTreeNode.Parent.Children, pathIndex)
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Internal
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Children = leftChildren
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Children = rightChildren
				}
				childBTreeNode.Parent = bTreeNode.Parent
				bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex+i, childBTreeNode)
			}
			if len(bTreeNode.Parent.Keys) > MAX_KEYS {
				split(bTreeNode.Parent, pathIndex)
			}
		}
	case Leaf:
		if bTreeNode.IsRoot {
			bTreeNode.Keys = []int{middleKey}
			bTreeNode.NodeType = Internal
			bTreeNode.Data = []Row{}
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Leaf
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Data = leftData
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Data = rightData
				}
				childBTreeNode.Parent = bTreeNode
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			}
		} else {
			fmt.Println("hi")
			bTreeNode.Parent.Keys = insert(bTreeNode.Parent.Keys, pathIndex, middleKey)
			bTreeNode.Parent.Children = remove(bTreeNode.Parent.Children, pathIndex)
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Leaf
				if i == 0 {
					childBTreeNode.Keys = leftKeys
					childBTreeNode.Data = leftData
				} else {
					childBTreeNode.Keys = rightKeys
					childBTreeNode.Data = rightData
				}
				childBTreeNode.Parent = bTreeNode.Parent
				bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex+i, childBTreeNode)
			}
			if len(bTreeNode.Parent.Keys) > MAX_KEYS {
				split(bTreeNode.Parent, pathIndex)
			}
		}
	}
}

func printTree(bTreeNode *BTreeNode, level int) {
	for i := 0; i < level; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("%v\n", bTreeNode)
	for _, childBTreeNode := range bTreeNode.Children {
		printTree(childBTreeNode, level+1)
	}
}

func executeSelect(bTreeNode *BTreeNode) {
	switch bTreeNode.NodeType {
	case Internal:
		for _, childBTreeNode := range bTreeNode.Children {
			executeSelect(childBTreeNode)
		}
	case Leaf:
		for i, childBTreeNodeKey := range bTreeNode.Keys {
			row := bTreeNode.Data[i]
			fmt.Printf("id: %v, name: %v, email: %v\n", childBTreeNodeKey, row.Name, row.Email)
		}
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

type Row struct {
	Name  string
	Email string
}
