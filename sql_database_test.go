package main

import (
	"slices"
	"testing"
)

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

func TestBTreeSplit(t *testing.T) {
	table := Table{length: 5, rows: BTreeNode{isRoot: true, nodeType: Leaf, parent: nil, children: []*BTreeNode{}, keys: []int{1, 2, 3, 4, 5}, data: []Row{{"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}}}}
	table.rows.split([]int{0})
	if table.rows.keys[0] != 3 {
		t.Errorf("didnt propogate the middle key up to parent")
	}
	if table.rows.children[0].keys[0] != 1 {
		t.Errorf("didn't save the left keys correctly")
	}
	if table.rows.children[1].keys[0] != 4 {
		t.Errorf("didn't save the right keys correctly")
	}
}

func TestBTreeMultipleSplit(t *testing.T) {
	table := Table{
		rows: BTreeNode{
			isRoot:   true,
			nodeType: Leaf,
			parent:   nil,
			keys:     []int{1, 2, 10, 11, 12},
			data:     []Row{{"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}},
			children: []*BTreeNode{},
		},
	}
	table.rows.split([]int{0})
	keysToInsert := []int{3, 4}
	dataToInsert := []Row{{"a", "a"}, {"a", "a"}}
	for i, key := range keysToInsert {
		table.rows.children[0].keys = slices.Insert(table.rows.children[0].keys, 2+i, key)
		table.rows.children[0].data = slices.Insert(table.rows.children[0].data, 2+i, dataToInsert[i])
	}
	table.rows.children[0].split([]int{0})
	if len(table.rows.children) != 3 {
		t.Errorf("didn't create 2 new children")
	}
	if table.rows.children[0].keys[0] != 1 {
		t.Errorf("didn't split left keys correctly")
	}
	if table.rows.children[1].keys[0] != 4 {
		t.Errorf("didn't split right keys correctly")
	}
	if table.rows.keys[0] != 3 {
		t.Errorf("did't propogate the middle key correctly")
	}
}

func TestInsert(t *testing.T) {
	table := Table{rows: BTreeNode{isRoot: true, nodeType: Leaf, parent: nil, children: []*BTreeNode{}, keys: []int{}, data: []Row{}}}
	table.rows.insertKey(1, Row{"a", "a"}, []int{0})
	if table.rows.keys[0] != 1 {
		t.Errorf("didn't insert key")
	}
}

func TestPersistance(t *testing.T) {
	table := Table{length: 1, rows: BTreeNode{isRoot: true, nodeType: Leaf, parent: nil, children: []*BTreeNode{}, keys: []int{1}, data: []Row{{}}}}
	saveToFile(table, TEST_DB_FILENAME)
	readTable := readFile(TEST_DB_FILENAME)
	if !table.rows.Equals(&readTable.rows) {
		t.Errorf("persistance failed")
	}
}
