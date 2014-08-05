package home

import (
    "testing"
)

// SERIALIZE SINGLE NODE: It should be able to serialize a tree with only one node
func TestSerializeTreeOneNodeZeroDepth(t *testing.T) {
    beforeEach("serialize")

    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    expectedOut := []Condition{Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}

    conditionsReturned, errorsReturned := serializeTree(root)

    if !matchesArray(conditionsReturned, expectedOut) {
        t.Errorf("serializeTree(%v) - got %v, want %v", root, conditionsReturned, expectedOut)
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", root, errorsReturned, expectedOutErr)
    }
}

// SERIALIZE SINGLE DEPTH, ENCLOSURE: It should be able to serialize a tree with a node and two children
/**
 * A && B
 *      AND
 *     A   B
 */
func TestSerializeTreeThreeNodeOneDepth(t *testing.T) {
    beforeEach("serialize")

    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child1 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}}
    child2 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    root.Children = []*treeNode{&child1, &child2}

    expectedOut := []Condition{
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"},
        Condition{Text: ")", Type: "scope", Operator: ")"},
    }

    conditionsReturned, errorsReturned := serializeTree(root)

    if !matchesArray(conditionsReturned, expectedOut) {
        t.Errorf("serializeTree(%v) conditionsReturned - got %v, want %v", root, conditionsReturned, expectedOut)
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", root, errorsReturned, expectedOutErr)
    }
}

// SERIALIZE ORDER, ARBITRARY DEPTH: It should be able to serialize a tree with nine nodes and four levels of depth (aka, arbitrary depth) in the correct order
/**
 * ((A && B) || C) && (D || E)
 *             AND
 *       OR           OR
 *   AND     C      D    E
 *  A   B
 */
func TestSerializeTreeArbitraryDepth(t *testing.T) {
    beforeEach("serialize")

    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    root.Children = []*treeNode{&child1, &child2}

    child3 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child4 := treeNode{Parent: &child1, Children: nil, Node: Condition{Text: "age eq 1", Type: "equality", Field: "age", Operator: "eq", Value: "1"}}
    child1.Children = []*treeNode{&child3, &child4}

    child5 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    child6 := treeNode{Parent: &child2, Children: nil, Node: Condition{Text: "age eq 3", Type: "equality", Field: "age", Operator: "eq", Value: "3"}}
    child2.Children = []*treeNode{&child5, &child6}

    child7 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}}
    child8 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    child3.Children = []*treeNode{&child7, &child8}

    expectedOut := []Condition{
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

    conditionsReturned, errorsReturned := serializeTree(root)

    if !matchesArray(conditionsReturned, expectedOut) {
        t.Errorf("serializeTree(%v) conditionsReturned - got %v, want %v", root, simplifyConditions(conditionsReturned), simplifyConditions(expectedOut))
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", root, errorsReturned, expectedOutErr)
    }
}

// SERIALIZE ARBITRARY WIDTH: It should be able to serialize a tree with any amount of children on a branch
func TestSerializeTreeArbitraryWidth(t *testing.T) {
    beforeEach("serialize")

    conditionsReturned, errorsReturned := serializeTree(testingTreeRoot)

    if !matchesArray(conditionsReturned, testingConditions) {
        t.Errorf("serializeTree(%v) conditionsReturned - got %v, want %v", testingTreeRoot, simplifyConditions(conditionsReturned), simplifyConditions(testingConditions))
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", testingTreeRoot, errorsReturned, expectedOutErr)
    }
}

// SERIALIZE ERROR EQUALITY BRANCH: It should return an error when the tree contains an equality condition in a branch
func TestSerializeTreeWithEqualityBranch(t *testing.T) {
    beforeEach("serialize")

    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"}}
    child1 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}}
    child2 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    root.Children = []*treeNode{&child1, &child2}

    expectedOutErr := "ERROR: This tree has an equality condition as a branch. Quitting."

    conditionsReturned, errorsReturned := serializeTree(root) 

    if conditionsReturned != nil {
        t.Errorf("serializeTree(%v) - got %v, want %v", root, conditionsReturned, nil)
    }

    if errorsReturned != nil && errorsReturned.Error() != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", root, errorsReturned, expectedOutErr)
    }
}

// SERIALIZE ERROR LOGIC LEAF: It should return an error when the tree contains a logic condition in a leaf
func TestSerializeTreeWithLogicLeaf(t *testing.T) {
    beforeEach("serialize")

    root := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child1 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}
    child2 := treeNode{Parent: root, Children: nil, Node: Condition{Text: "age eq 2", Type: "equality", Field: "age", Operator: "eq", Value: "2"}}
    root.Children = []*treeNode{&child1, &child2}

    expectedOutErr := "ERROR: This tree has a logic condition as a leaf. Quitting."

    conditionsReturned, errorsReturned := serializeTree(root) 

    if conditionsReturned != nil {
        t.Errorf("serializeTree(%v) - got %v, want %v", root, conditionsReturned, nil)
    }

    if errorsReturned != nil && errorsReturned.Error() != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", root, errorsReturned, expectedOutErr)
    }
}