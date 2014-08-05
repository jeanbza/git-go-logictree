package home

import (
    "testing"
)

// ATTACH LEFTS AND RIGHTS TO TREE: It should be able to assign lefts and rights to a tree
/**
 * lt-node-rt
 *                                     1-AND-24
 *                          2-OR-17                     18-OR-23
 *              3-AND-14                15-F-16     19-G-20  21-H-22
 * 4-A-5 6-B-7 8-C-9 10-D-11 12-E-13
 */
func TestAttachLeftsAndRights(t *testing.T) {
    beforeEach("mysql")

    // Row 1 node 1
    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    // Row 2 node 1
    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    // Row 2 node 2
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    root.Children = []*treeNode{&child1, &child2}

    // Row 3 node 1
    child3 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    // Row 3 node 2
    child4 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"}}
    child1.Children = []*treeNode{&child3, &child4}

    // Row 3 node 3
    child5 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    // Row 3 node 4
    child6 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"}}
    child2.Children = []*treeNode{&child5, &child6}

    // Row 4 nodes 1-5
    child7 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}}
    child8 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    child9 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 6", Type: "equality", Field: "age", Operator: "eq", Value: "6"}}
    child10 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 7", Type: "equality", Field: "age", Operator: "eq", Value: "7"}}
    child11 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}}
    child3.Children = []*treeNode{&child7, &child8, &child9, &child10, &child11}

    // Row 1 node 1
    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}, Left: 1, Right: 24}

    // Row 2 node 1
    outChild1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}, Left: 2, Right: 17}
    // Row 2 node 2
    outChild2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}, Left: 18, Right: 23}
    expectedOut.Children = []*treeNode{&outChild1, &outChild2}

    // Row 3 node 1
    outChild3 := treeNode{Parent: &outChild1, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}, Left: 3, Right: 14}
    // Row 3 node 2
    outChild4 := treeNode{Parent: &outChild1, Children: nil, Node: Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"}, Left: 15, Right: 16}
    outChild1.Children = []*treeNode{&outChild3, &outChild4}

    // Row 3 node 3
    outChild5 := treeNode{Parent: &outChild2, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}, Left: 19, Right: 20}
    // Row 3 node 4
    outChild6 := treeNode{Parent: &outChild2, Children: nil, Node: Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"}, Left: 21, Right: 22}
    outChild2.Children = []*treeNode{&outChild5, &outChild6}

    // Row 4 nodes 1-5
    outChild7 := treeNode{Parent: &outChild3, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}, Left: 4, Right: 5}
    outChild8 := treeNode{Parent: &outChild3, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}, Left: 6, Right: 7}
    outChild9 := treeNode{Parent: &outChild3, Children: nil, Node: Condition{Text: "age eq 6", Type: "equality", Field: "age", Operator: "eq", Value: "6"}, Left: 8, Right: 9}
    outChild10 := treeNode{Parent: &outChild3, Children: nil, Node: Condition{Text: "age eq 7", Type: "equality", Field: "age", Operator: "eq", Value: "7"}, Left: 10, Right: 11}
    outChild11 := treeNode{Parent: &outChild3, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}, Left: 12, Right: 13}
    outChild3.Children = []*treeNode{&outChild7, &outChild8, &outChild9, &outChild10, &outChild11}

    root.attachLeftsAndRights()

    if !root.matches(expectedOut) {
        t.Errorf("%v.attachLeftsAndRights - got %v, want %v", root, root.print(), expectedOut.print())
    }
}

// CONVERT TREE TO MYSQL: It should be able to convert a tree to mysql input values
/**
 * lt-node-rt
 *                                     1-AND-24
 *                          2-OR-17                     18-OR-23
 *              3-AND-14                15-F-16     19-G-20  21-H-22
 * 4-A-5 6-B-7 8-C-9 10-D-11 12-E-13
 */
func TestToMysql(t *testing.T) {
    beforeEach("mysql")

    // Row 1 node 1
    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    // Row 2 node 1
    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    // Row 2 node 2
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    root.Children = []*treeNode{&child1, &child2}

    // Row 3 node 1
    child3 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    // Row 3 node 2
    child4 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"}}
    child1.Children = []*treeNode{&child3, &child4}

    // Row 3 node 3
    child5 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    // Row 3 node 4
    child6 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"}}
    child2.Children = []*treeNode{&child5, &child6}

    // Row 4 nodes 1-5
    child7 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}}
    child8 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    child9 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 6", Type: "equality", Field: "age", Operator: "eq", Value: "6"}}
    child10 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 7", Type: "equality", Field: "age", Operator: "eq", Value: "7"}}
    child11 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}}
    child3.Children = []*treeNode{&child7, &child8, &child9, &child10, &child11}

    // INSERT INTO logictree.equality (field, operator, value, lt, rt) VALUES ...
    expectedOutEquality := "('age', 'eq', '4', 4, 5),('age', 'eq', '5', 6, 7),('age', 'eq', '6', 8, 9),('age', 'eq', '7', 10, 11),('age', 'eq', '8', 12, 13),('age', 'eq', '1', 15, 16),('age', 'eq', '2', 19, 20),('age', 'eq', '3', 21, 22)"

    // INSERT INTO logictree.logic (operator, lt, rt) VALUES ...
    expectedOutLogic := "('AND', 3, 14),('OR', 2, 17),('OR', 18, 23),('AND', 1, 24)"

    equalityReturned, logicReturned := root.toMysql()

    if equalityReturned != expectedOutEquality {
        t.Errorf("%v.toMysql() equalityReturned - got %v, want %v", root, equalityReturned, expectedOutEquality)
    }

    if logicReturned != expectedOutLogic {
        t.Errorf("%v.toMysql() logicReturned - got %v, want %v", root, logicReturned, expectedOutLogic)
    }
}
