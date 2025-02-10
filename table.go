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
	Keys     []Row
}

func executeInsert(bTreeNode *BTreeNode, key Row, pathIndex int) {
	switch bTreeNode.NodeType {
	case Internal:
		for i, nodeKey := range bTreeNode.Keys {
			if nodeKey.Id >= key.Id {
				executeInsert(bTreeNode.Children[i], key, i)
				break
			} else if len(bTreeNode.Keys)-1 == i {
				executeInsert(bTreeNode.Children[i], key, i)
			}
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if nodeKey.Id > key.Id {
					keys := insert(bTreeNode.Keys, i, key)
					bTreeNode.Keys = keys
					break
				}
				if i == len(bTreeNode.Keys)-1 {
					bTreeNode.Keys = append(bTreeNode.Keys, key)
				}
			}
		}
		if len(bTreeNode.Keys) > MAX_KEYS {
			split(bTreeNode, pathIndex)
		}
	}

}

func split(bTreeNode *BTreeNode, pathIndex int) {
	leftKeys, rightKeys, middleKey := bTreeNode.Keys[:3], bTreeNode.Keys[3:], bTreeNode.Keys[2]
	switch bTreeNode.NodeType {
	case Internal:
		leftChildren, rightChildren := bTreeNode.Children[:3], bTreeNode.Children[3:]
		if bTreeNode.IsRoot {
			bTreeNode.Keys = []Row{middleKey}
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
			bTreeNode.Keys = []Row{middleKey}
			bTreeNode.NodeType = Internal
			for i := 0; i < 2; i++ {
				childBTreeNode := new(BTreeNode)
				childBTreeNode.IsRoot = false
				childBTreeNode.NodeType = Leaf
				if i == 0 {
					childBTreeNode.Keys = leftKeys
				} else {
					childBTreeNode.Keys = rightKeys
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
				childBTreeNode.NodeType = Leaf
				if i == 0 {
					childBTreeNode.Keys = leftKeys
				} else {
					childBTreeNode.Keys = rightKeys
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
		printTree(childBTreeNode, level + 1)
	}
}

func executeSelect(bTreeNode *BTreeNode) {
	switch bTreeNode.NodeType {
	case Internal:
		for _, childBTreeNode := range bTreeNode.Children {
			executeSelect(childBTreeNode)
		}
	case Leaf:
		for _, childBTreeNodeKey := range bTreeNode.Keys {
			fmt.Println(childBTreeNodeKey)
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
	Id    int    `gob:"id"`
	Name  string `gob:"name"`
	Email string `gob:"email"`
}
