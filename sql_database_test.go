package main

import "testing"

const TEST_DB_FILENAME = "testing.db"

func TestParseCommand(t *testing.T) {
	_, res := parseCommand("exit")
	if res != nil {
		t.Errorf("%v", res)
	}
}

func TestParseStatement(t *testing.T) {
	_, _, _, res := parseStatement("select")
	if res != nil {
		t.Errorf("%v", res)
	}
}

func TestBTreeInitialRootSplit(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	// setup
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(3, Row{name: "a", email: "a"})
	table.executeInsert(4, Row{name: "a", email: "a"})
	// test
	table.executeInsert(5, Row{name: "a", email: "a"})
	if table.rows.keys[0] != 3 {
		t.Error("didnt propogate the middle key up to parent")
	}
	if table.rows.children[0].keys[0] != 1 {
		t.Error("didn't save the left keys correctly")
	}
	if table.rows.children[1].keys[0] != 4 {
		t.Error("didn't save the right keys correctly")
	}
}

func TestBTreeNonRootLeafSplit(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			keys:     []int{},
			data:     []Row{},
			children: []*BTreeNode{},
		},
	}
	// setup
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(10, Row{name: "a", email: "a"})
	table.executeInsert(11, Row{name: "a", email: "a"})
	table.executeInsert(12, Row{name: "a", email: "a"})
	// test
	table.executeInsert(3, Row{name: "a", email: "a"})
	table.executeInsert(4, Row{name: "a", email: "a"})
	if len(table.rows.children) != 3 {
		t.Error("didn't create 2 new children")
	}
	if table.rows.children[0].keys[0] != 1 {
		t.Error("didn't split left keys correctly")
	}
	if table.rows.children[1].keys[0] != 4 {
		t.Error("didn't split right keys correctly")
	}
	if table.rows.keys[0] != 3 {
		t.Error("did't propogate the middle key correctly")
	}
}

func TestInsert(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	table.executeInsert(1, Row{name: "a", email: "a"})
	if table.rows.keys[0] != 1 {
		t.Error("didn't insert key")
	}
}

func TestPersistance(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	table.executeInsert(1, Row{name: "a", email: "a"})
	saveToFile(table, TEST_DB_FILENAME)
	readTable := readFile(TEST_DB_FILENAME)
	if !table.rows.Equals(&readTable.rows) {
		t.Error("persistance failed")
	}
}

func TestBTreeNonRootLeafToInternalRootParentSplit(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	// setup
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(5, Row{name: "a", email: "a"})
	table.executeInsert(6, Row{name: "a", email: "a"})
	table.executeInsert(7, Row{name: "a", email: "a"})
	table.executeInsert(10, Row{name: "a", email: "a"})
	table.executeInsert(11, Row{name: "a", email: "a"})
	table.executeInsert(12, Row{name: "a", email: "a"})
	table.executeInsert(15, Row{name: "a", email: "a"})
	table.executeInsert(16, Row{name: "a", email: "a"})
	table.executeInsert(17, Row{name: "a", email: "a"})
	table.executeInsert(20, Row{name: "a", email: "a"})
	table.executeInsert(21, Row{name: "a", email: "a"})
	table.executeInsert(22, Row{name: "a", email: "a"})
	table.executeInsert(25, Row{name: "a", email: "a"})
	table.executeInsert(26, Row{name: "a", email: "a"})
	// test
	table.executeInsert(27, Row{name: "a", email: "a"})
	if len(table.rows.children) != 2 {
		t.Error("didn't create 2 children")
	}
	for _, rootNodeChild := range table.rows.children {
		if rootNodeChild.nodeType != Internal {
			t.Error("Supposed to be 2 internal children")
		}
		if !table.rows.Equals(rootNodeChild.parent) {
			t.Error("Parent should be root")
		}
		for _, leafNode := range rootNodeChild.children {
			if leafNode.nodeType != Leaf {
				t.Error("Nodetype should be Leaf")
			}
			if !rootNodeChild.Equals(leafNode.parent) {
				t.Errorf("Parent should be %p: %p %v", leafNode, rootNodeChild.parent, rootNodeChild.parent)
			}
		}
	}
}

func TestBTreeNonRootLeafToInternalNonRootParentSplit(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	// setup
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(5, Row{name: "a", email: "a"})
	table.executeInsert(6, Row{name: "a", email: "a"})
	table.executeInsert(7, Row{name: "a", email: "a"})
	table.executeInsert(10, Row{name: "a", email: "a"})
	table.executeInsert(11, Row{name: "a", email: "a"})
	table.executeInsert(12, Row{name: "a", email: "a"})
	table.executeInsert(15, Row{name: "a", email: "a"})
	table.executeInsert(16, Row{name: "a", email: "a"})
	table.executeInsert(17, Row{name: "a", email: "a"})
	table.executeInsert(20, Row{name: "a", email: "a"})
	table.executeInsert(21, Row{name: "a", email: "a"})
	table.executeInsert(22, Row{name: "a", email: "a"})
	table.executeInsert(25, Row{name: "a", email: "a"})
	table.executeInsert(26, Row{name: "a", email: "a"})
	table.executeInsert(27, Row{name: "a", email: "a"})
	table.executeInsert(30, Row{name: "a", email: "a"})
	table.executeInsert(31, Row{name: "a", email: "a"})
	table.executeInsert(32, Row{name: "a", email: "a"})
	table.executeInsert(35, Row{name: "a", email: "a"})
	table.executeInsert(36, Row{name: "a", email: "a"})
	table.executeInsert(37, Row{name: "a", email: "a"})
	table.executeInsert(40, Row{name: "a", email: "a"})
	table.executeInsert(41, Row{name: "a", email: "a"})
	// test
	table.executeInsert(45, Row{name: "a", email: "a"})
	newRootNodeChild := table.rows.children[2]
	rootNode := table.rows
	if !rootNode.Equals(newRootNodeChild.parent) {
		t.Errorf("Root Node Child parent should be %p", &(rootNode))
	}
	for _, newRootNodeChildLeafChild := range newRootNodeChild.children {
		if !newRootNodeChild.Equals(newRootNodeChildLeafChild.parent) {
			t.Errorf("internalChildNode parent should be %p %p %v %v", newRootNodeChild, newRootNodeChildLeafChild.parent, newRootNodeChild, newRootNodeChildLeafChild.parent)
		}
	}
}

func Test3LevelSplit(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	// setup
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(3, Row{name: "a", email: "a"})
	table.executeInsert(4, Row{name: "a", email: "a"})
	table.executeInsert(5, Row{name: "a", email: "a"})
	table.executeInsert(6, Row{name: "a", email: "a"})
	table.executeInsert(7, Row{name: "a", email: "a"})
	table.executeInsert(8, Row{name: "a", email: "a"})
	table.executeInsert(9, Row{name: "a", email: "a"})
	table.executeInsert(10, Row{name: "a", email: "a"})
	table.executeInsert(11, Row{name: "a", email: "a"})
	table.executeInsert(12, Row{name: "a", email: "a"})
	table.executeInsert(13, Row{name: "a", email: "a"})
	table.executeInsert(14, Row{name: "a", email: "a"})
	table.executeInsert(15, Row{name: "a", email: "a"})
	table.executeInsert(16, Row{name: "a", email: "a"})
	table.executeInsert(17, Row{name: "a", email: "a"})
	table.executeInsert(18, Row{name: "a", email: "a"})
	table.executeInsert(19, Row{name: "a", email: "a"})
	table.executeInsert(20, Row{name: "a", email: "a"})
	table.executeInsert(21, Row{name: "a", email: "a"})
	table.executeInsert(22, Row{name: "a", email: "a"})
	table.executeInsert(23, Row{name: "a", email: "a"})
	table.executeInsert(24, Row{name: "a", email: "a"})
	table.executeInsert(25, Row{name: "a", email: "a"})
	table.executeInsert(26, Row{name: "a", email: "a"})
	table.executeInsert(27, Row{name: "a", email: "a"})
	table.executeInsert(28, Row{name: "a", email: "a"})
	table.executeInsert(29, Row{name: "a", email: "a"})
	table.executeInsert(30, Row{name: "a", email: "a"})
	table.executeInsert(31, Row{name: "a", email: "a"})
	table.executeInsert(32, Row{name: "a", email: "a"})
	table.executeInsert(33, Row{name: "a", email: "a"})
	table.executeInsert(34, Row{name: "a", email: "a"})
	table.executeInsert(35, Row{name: "a", email: "a"})
	table.executeInsert(36, Row{name: "a", email: "a"})
	table.executeInsert(37, Row{name: "a", email: "a"})
	table.executeInsert(38, Row{name: "a", email: "a"})
	table.executeInsert(39, Row{name: "a", email: "a"})
	table.executeInsert(40, Row{name: "a", email: "a"})
	table.executeInsert(41, Row{name: "a", email: "a"})
	table.executeInsert(42, Row{name: "a", email: "a"})
	table.executeInsert(43, Row{name: "a", email: "a"})
	table.executeInsert(44, Row{name: "a", email: "a"})
	table.executeInsert(45, Row{name: "a", email: "a"})
	table.executeInsert(46, Row{name: "a", email: "a"})
	table.executeInsert(47, Row{name: "a", email: "a"})
	table.executeInsert(48, Row{name: "a", email: "a"})
	table.executeInsert(49, Row{name: "a", email: "a"})
	table.executeInsert(50, Row{name: "a", email: "a"})
	table.executeInsert(51, Row{name: "a", email: "a"})
	table.executeInsert(52, Row{name: "a", email: "a"})
	// test
	table.executeInsert(53, Row{name: "a", email: "a"})
	for _, rootNodeChild := range table.rows.children {
		if !table.rows.Equals(rootNodeChild.parent) {
			t.Error("Root Node Child parent should be the root node")
		}
		for _, secLevelChild := range rootNodeChild.children {
			if !rootNodeChild.Equals(secLevelChild.parent) {
				t.Error("2nd level node child parent should be the root child node")
			}
			for _, leafChild := range secLevelChild.children {
				if !secLevelChild.Equals(leafChild.parent) {
					t.Error("Leaf node child parent should be the 2nd level child node")
				}
			}
		}
	}
}

func TestMultiLevelTreePersistance(t *testing.T) {
	table := Table{
		length: 0,
		rows: BTreeNode{
			isRoot:   true,
			parent:   nil,
			nodeType: Leaf,
			children: []*BTreeNode{},
			keys:     []int{},
			data:     []Row{},
		},
	}
	table.executeInsert(1, Row{name: "a", email: "a"})
	table.executeInsert(2, Row{name: "a", email: "a"})
	table.executeInsert(3, Row{name: "a", email: "a"})
	table.executeInsert(4, Row{name: "a", email: "a"})
	table.executeInsert(5, Row{name: "a", email: "a"})
	table.executeInsert(6, Row{name: "a", email: "a"})
	table.executeInsert(7, Row{name: "a", email: "a"})
	table.executeInsert(8, Row{name: "a", email: "a"})
	table.executeInsert(9, Row{name: "a", email: "a"})
	table.executeInsert(10, Row{name: "a", email: "a"})
	table.executeInsert(11, Row{name: "a", email: "a"})
	table.executeInsert(12, Row{name: "a", email: "a"})
	table.executeInsert(13, Row{name: "a", email: "a"})
	table.executeInsert(14, Row{name: "a", email: "a"})
	table.executeInsert(15, Row{name: "a", email: "a"})
	table.executeInsert(16, Row{name: "a", email: "a"})
	table.executeInsert(17, Row{name: "a", email: "a"})
	table.executeInsert(18, Row{name: "a", email: "a"})
	table.executeInsert(19, Row{name: "a", email: "a"})
	table.executeInsert(20, Row{name: "a", email: "a"})
	table.executeInsert(21, Row{name: "a", email: "a"})
	table.executeInsert(22, Row{name: "a", email: "a"})
	table.executeInsert(23, Row{name: "a", email: "a"})
	table.executeInsert(24, Row{name: "a", email: "a"})
	table.executeInsert(25, Row{name: "a", email: "a"})
	table.executeInsert(26, Row{name: "a", email: "a"})
	table.executeInsert(27, Row{name: "a", email: "a"})
	table.executeInsert(28, Row{name: "a", email: "a"})
	table.executeInsert(29, Row{name: "a", email: "a"})
	table.executeInsert(30, Row{name: "a", email: "a"})
	table.executeInsert(31, Row{name: "a", email: "a"})
	table.executeInsert(32, Row{name: "a", email: "a"})
	table.executeInsert(33, Row{name: "a", email: "a"})
	table.executeInsert(34, Row{name: "a", email: "a"})
	table.executeInsert(35, Row{name: "a", email: "a"})
	table.executeInsert(36, Row{name: "a", email: "a"})
	table.executeInsert(37, Row{name: "a", email: "a"})
	table.executeInsert(38, Row{name: "a", email: "a"})
	table.executeInsert(39, Row{name: "a", email: "a"})
	table.executeInsert(40, Row{name: "a", email: "a"})
	saveToFile(table, TEST_DB_FILENAME)
	readTable := readFile(TEST_DB_FILENAME)
	if !table.rows.Equals(&readTable.rows) {
		t.Errorf("persistance failed")
	}
}
