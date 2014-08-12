package home

import (
    "testing"
)

// UnserializeRaw SINGLE NODE: It should be able to unserialize a tree with only one node
func TestUnserializeRawTreeOneNodeZeroDepth(t *testing.T) {
    beforeEach("unserializeRaw")

    in := []conditionSqlRow{conditionSqlRow{Type: "equality", Field: "age", Operator: "eq", Value: "5", Left: 1, Right: 2}}
    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    treeReturned := unserializeRawTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned, expectedOut)
    }
}

// UnserializeRaw SINGLE DEPTH, NON-ENCLOSURE: It should be able to unserialize a tree with multiple nodes on one branch
func TestUnserializeRawTreeThreeNodeOneDepth(t *testing.T) {
    beforeEach("unserializeRaw")

    in := []conditionSqlRow{
        conditionSqlRow{Type: "logic", Operator: "AND", Left: 1, Right: 6},
        conditionSqlRow{Type: "equality", Field: "age", Operator: "eq", Value: "81", Left: 2, Right: 3},
        conditionSqlRow{Type: "equality", Field: "age", Operator: "eq", Value: "27", Left: 4, Right: 5},
    }

    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}, Left: 1, Right: 6}
    child1 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 81", Type: "equality", Field: "age", Operator: "eq", Value: "81"}, Left: 2, Right: 3}
    child2 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 27", Type: "equality", Field: "age", Operator: "eq", Value: "27"}, Left: 4, Right: 5}
    expectedOut.Children = []*treeNode{&child1, &child2}

    treeReturned := unserializeRawTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }
}

// UnserializeRaw ARBITRARY WIDTH: It should be able to unserialize a tree with any amount of children on a branch
func TestUnserializeRawTreeArbitraryWidth(t *testing.T) {
    beforeEach("unserializeFormatted")

    treeReturned := unserializeRawTree(testingMysqlRows)

    testingTreeRoot.attachLeftsAndRights()

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", testingMysqlRows, treeReturned.print(), testingTreeRoot.print())
    }
}

// UnserializeFormatted SINGLE NODE: It should be able to unserialize a tree with only one node
func TestUnserializeFormattedTreeOneNodeZeroDepth(t *testing.T) {
    beforeEach("unserializeFormatted")

    in := []Condition{Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    var expectedOutErr error

    treeReturned, errorsReturned := unserializeFormattedTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned, expectedOut)
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// unserializeFormatted SINGLE DEPTH, NON-ENCLOSURE: It should be able to unserialize a tree with multiple nodes on one branch
func TestUnserializeFormattedTreeThreeNodeOneDepth(t *testing.T) {
    beforeEach("unserializeFormatted")

    in := []Condition{
        Condition{Text: "age eq 81", Type: "equality", Field: "age", Operator: "eq", Value: "81"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 27", Type: "equality", Field: "age", Operator: "eq", Value: "27"},
    }

    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child1 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 81", Type: "equality", Field: "age", Operator: "eq", Value: "81"}}
    child2 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 27", Type: "equality", Field: "age", Operator: "eq", Value: "27"}}
    expectedOut.Children = []*treeNode{&child1, &child2}

    treeReturned, errorsReturned := unserializeFormattedTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// UnserializeFormatted SINGLE DEPTH, ENCLOSURE: It should be able to unserialize a tree with a node and two children
/**
 * A && B
 *      AND
 *     A   B
 */
func TestUnserializeFormattedTreeThreeNodeOneDepthEnclosure(t *testing.T) {
    beforeEach("unserializeFormatted")

    in := []Condition{
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 81", Type: "equality", Field: "age", Operator: "eq", Value: "81"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 27", Type: "equality", Field: "age", Operator: "eq", Value: "27"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
    }

    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child1 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 81", Type: "equality", Field: "age", Operator: "eq", Value: "81"}}
    child2 := treeNode{Parent: expectedOut, Children: nil, Node: Condition{Text: "age eq 27", Type: "equality", Field: "age", Operator: "eq", Value: "27"}}
    expectedOut.Children = []*treeNode{&child1, &child2}

    treeReturned, errorsReturned := unserializeFormattedTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// UnserializeFormatted ORDER, ARBITRARY DEPTH: It should be able to unserialize a tree with nine nodes and four levels of depth (aka, arbitrary depth) in the correct order
// /**
//  * ((A && B) || C) && (D || E)
//  *             AND
//  *       OR           OR
//  *   AND     C      D    E
//  *  A   B
//  */
func TestUnserializeFormattedTreeArbitraryDepth(t *testing.T) {
    beforeEach("unserializeFormatted")

    in := []Condition{
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
        Condition{Text: "OR", Type: "logic", Operator: "OR"},
        Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"},
        Condition{Text: "OR", Type: "logic", Operator: "OR"},
        Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
    }

    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    expectedOut.Children = []*treeNode{&child1, &child2}

    child3 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child4 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"}}
    child1.Children = []*treeNode{&child3, &child4}

    child5 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    child6 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"}}
    child2.Children = []*treeNode{&child5, &child6}

    child7 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}}
    child8 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    child3.Children = []*treeNode{&child7, &child8}

    treeReturned, errorsReturned := unserializeFormattedTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// UnserializeFormatted ARBITRARY WIDTH: It should be able to unserialize a tree with any amount of children on a branch
func TestUnserializeFormattedTreeArbitraryWidth(t *testing.T) {
    beforeEach("unserializeFormatted")

    treeReturned, errorsReturned := unserializeFormattedTree(testingConditions)

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", testingConditions, treeReturned.print(), testingTreeRoot.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", testingConditions, errorsReturned, expectedOutErr)
    }
}