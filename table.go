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
}

func (bTreeNode *BTreeNode) insert(key int, pathIndex int) {
	switch bTreeNode.NodeType {
	case Internal:
		for i, nodeKey := range bTreeNode.Keys {
			if nodeKey > key {
				bTreeNode.Children[i].insert(key, i)
			} else {
				bTreeNode.Children[i+1].insert(key, i+1)
			}
		}
	case Leaf:
		if len(bTreeNode.Keys) == 0 {
			bTreeNode.Keys = append(bTreeNode.Keys, key)
		} else {
			for i, nodeKey := range bTreeNode.Keys {
				if nodeKey > key {
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
			bTreeNode.split(pathIndex)
            fmt.Println("inside insert")
            bTreeNode.Parent.printTree(0)
		}
	}
}

func (bTreeNode *BTreeNode) split(pathIndex int) {
	leftKeys, rightKeys, middleKey := bTreeNode.Keys[0:3], bTreeNode.Keys[3:], bTreeNode.Keys[2]
	var leftChildren, rightChildren []*BTreeNode
	if bTreeNode.IsRoot {
		bTreeNode.Keys = []int{middleKey}
        switch bTreeNode.NodeType {
        case Leaf:
            bTreeNode.NodeType = Internal
        case Internal:
            bTreeNode.Children = []*BTreeNode{}
        }
	} else {
		bTreeNode.Parent.Children = remove(bTreeNode.Parent.Children, pathIndex)
		bTreeNode.Parent.Keys = insert(bTreeNode.Parent.Keys, pathIndex, middleKey)
	}
    if bTreeNode.NodeType == Internal {
		leftChildren, rightChildren = bTreeNode.Children[:3], bTreeNode.Children[3:]
    }
	for i := 0; i < 2; i++ {
		childBTreeNode := new(BTreeNode)
		childBTreeNode.IsRoot = false
		if i == 0 {
			childBTreeNode.Keys = leftKeys
		} else {
			childBTreeNode.Keys = rightKeys
		}
		switch bTreeNode.NodeType {
		case Internal:
			childBTreeNode.NodeType = Internal
			if i == 0 {
				childBTreeNode.Children = leftChildren
			} else {
				childBTreeNode.Children = rightChildren
			}
			if bTreeNode.IsRoot {
				childBTreeNode.Parent = bTreeNode
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			} else {
				childBTreeNode.Parent = bTreeNode.Parent
				bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex+i, childBTreeNode)
			}
		case Leaf:
			childBTreeNode.NodeType = Leaf
			if bTreeNode.IsRoot {
				childBTreeNode.Parent = bTreeNode
				bTreeNode.Children = append(bTreeNode.Children, childBTreeNode)
			} else {
				childBTreeNode.Parent = bTreeNode.Parent
				bTreeNode.Parent.Children = insert(bTreeNode.Parent.Children, pathIndex+i, childBTreeNode)
			}
		}
	}
	if !bTreeNode.IsRoot {
		if len(bTreeNode.Parent.Keys) > MAX_KEYS {
			bTreeNode.Parent.split(pathIndex)
		}
	}
    fmt.Println("inside split")
    bTreeNode.Parent.printTree(0)
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

type Row struct {
	Id    int    `gob:"id"`
	Name  string `gob:"name"`
	Email string `gob:"email"`
}
