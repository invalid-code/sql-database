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
	_, _, _, res := parseStatement("select")
	if res != nil {
		t.Errorf("%v", res)
	}
}

func TestBTreeSplit(t *testing.T) {
	// table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []int{1, 2, 3, 4, 5}, Data: []Row{{"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}}, Children: []*BTreeNode{}}
	// table.split(0)
	// if table.Keys[0] != 3 {
	// 	t.Errorf("didnt propogate the middle key up to parent")
	// }
	// if table.Children[0].Keys[0] != 1 {
	// 	t.Errorf("didn't save the left keys correctly")
	// }
	// if table.Children[1].Keys[0] != 4 {
	// 	t.Errorf("didn't save the right keys correctly")
	// }
}

func TestBTreeMultipleSplit(t *testing.T) {
	// table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []int{1, 2, 10, 11, 12}, Data: []Row{{"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}, {"a", "a"}}, Children: []*BTreeNode{}}
	// table.split(0)
	// keysToInsert := []int{3, 4}
	// dataToInsert := []Row{{"a", "a"}, {"a", "a"}}
	// for i, key := range keysToInsert {
	// 	table.Children[0].Keys = insert(table.Children[0].Keys, 2+i, key)
	// 	table.Children[0].Data = insert(table.Children[0].Data, 2+i, dataToInsert[i])
	// }
	// table.Children[0].split(0)
	// if len(table.Children) != 3 {
	// 	t.Errorf("didn't create 2 new children")
	// }
	// if table.Children[0].Keys[0] != 1 {
	// 	t.Errorf("didn't split left keys correctly")
	// }
	// if table.Children[1].Keys[0] != 4 {
	// 	t.Errorf("didn't split right keys correctly")
	// }
	// if table.Keys[0] != 3 {
	// 	t.Errorf("did't propogate the middle key correctly")
	// }
}

func TestInsert(t *testing.T) {
	// table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []int{}, Data: []Row{}, Children: []*BTreeNode{}}
	// table.executeInsert(1, Row{"a", "a"}, 0)
	// if table.Keys[0] != 1 {
	// 	t.Errorf("didn't insert key")
	// }
}

func TestPersistance(t *testing.T) {
	// table := BTreeNode{IsRoot: true, NodeType: Leaf, Parent: nil, Keys: []int{}, Data: []Row{}, Children: []*BTreeNode{}}
	// saveToFile(table, TEST_DB_FILENAME)
	// readTable := readFile(TEST_DB_FILENAME)
	// if !table.Equals(&readTable) {
	// 	t.Errorf("persistance failed")
	// }
}
