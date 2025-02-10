package main

import (
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
	_, _, res := parseStatement("select")
	if res != nil {
		t.Errorf("%v", res)
	}
}

func TestBTreeSplit(t *testing.T) {
	table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []Row{{1, "a", "a"}, {2, "a", "a"}, {3, "a", "a"}, {4, "a", "a"}, {5, "a", "a"}}, Children: []*BTreeNode{}}
	split(&table, 0)
	if table.Keys[0].Id != 3 {
		t.Errorf("didnt propogate the middle key up to parent")
	}
	if table.Children[0].Keys[0].Id != 1 {
		t.Errorf("didn't save the left keys correctly")
	}
	if table.Children[1].Keys[0].Id != 4 {
		t.Errorf("didn't save the right keys correctly")
	}
}

func TestBTreeMultipleSplit(t *testing.T) {
	table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []Row{{1, "a", "a"}, {2, "a", "a"}, {10, "a", "a"}, {11, "a", "a"}, {12, "a", "a"}}, Children: []*BTreeNode{}}
	split(&table, 0)
	keysToInsert := []Row{{3, "a", "a"}, {4, "a", "a"}}
	for i, key := range keysToInsert {
		table.Children[0].Keys = insert(table.Children[0].Keys, 2+i, key)
	}
	split(table.Children[0], 0)
	if len(table.Children) != 3 {
		t.Errorf("didn't create 2 new children")
	}
	if table.Children[0].Keys[0].Id != 1 {
		t.Errorf("didn't split left keys correctly")
	}
	if table.Children[1].Keys[0].Id != 4 {
		t.Errorf("didn't split right keys correctly")
	}
	if table.Keys[0].Id != 3 {
		t.Errorf("did't propogate the middle key correctly")
	}
}

func TestInsert(t *testing.T) {
	table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []Row{}, Children: []*BTreeNode{}}
	executeInsert(&table, Row{1, "a", "a"}, 0)
	if table.Keys[0].Id != 1 {
		t.Errorf("didn't insert key")
	}
}

func TestPersistance(t *testing.T) {
	table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []Row{}, Children: []*BTreeNode{}}
	saveToFile(table, TEST_DB_FILENAME)
	readTable := readFile(TEST_DB_FILENAME)
	if !table.Equals(&readTable) {
		t.Errorf("persistance failed")
	}
}
