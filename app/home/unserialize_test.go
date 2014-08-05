package home

import (
    "testing"
)

// UNSERIALIZE SINGLE NODE: It should be able to unserialize a tree with only one node
func TestUnserializeTreeOneNodeZeroDepth(t *testing.T) {
    beforeEach("unserialize")

    in := []Condition{Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    expectedOut := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"}}
    var expectedOutErr error

    treeReturned, errorsReturned := unserializeTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", in, treeReturned, expectedOut)
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// UNSERIALIZE SINGLE DEPTH, ENCLOSURE: It should be able to unserialize a tree with a node and two children
/**
 * A && B
 *      AND
 *     A   B
 */
func TestUnserializeTreeThreeNodeOneDepth(t *testing.T) {
    beforeEach("unserialize")

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

    var expectedOutErr error

    treeReturned, errorsReturned := unserializeTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// // UNSERIALIZE ORDER, ARBITRARY DEPTH: It should be able to unserialize a tree with nine nodes and four levels of depth (aka, arbitrary depth) in the correct order
// /**
//  * ((A && B) || C) && (D || E)
//  *             AND
//  *       OR           OR
//  *   AND     C      D    E
//  *  A   B
//  */
func TestUnserializeTreeArbitraryDepth(t *testing.T) {
    beforeEach("unserialize")

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

    var expectedOutErr error

    treeReturned, errorsReturned := unserializeTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}

// UNSERIALIZE ARBITRARY WIDTH: It should be able to serialize a tree with any amount of children on a branch
/**
 * ((A && B) || C) && (D || E)
 *              AND
 *        OR           OR
 *    AND     F      G    H
 * A B C D E
 */
func TestUnserializeTreeArbitraryWidth(t *testing.T) {
    beforeEach("unserialize")

    in := []Condition{
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 4", Type: "equality", Field: "age", Operator: "eq", Value: "4"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 5", Type: "equality", Field: "age", Operator: "eq", Value: "5"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 6", Type: "equality", Field: "age", Operator: "eq", Value: "6"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 7", Type: "equality", Field: "age", Operator: "eq", Value: "7"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
        Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"},
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
    child9 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 6", Type: "equality", Field: "age", Operator: "eq", Value: "6"}}
    child10 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 7", Type: "equality", Field: "age", Operator: "eq", Value: "7"}}
    child11 := treeNode{Parent: &child3, Children: nil, Node: Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"}}
    child3.Children = []*treeNode{&child7, &child8, &child9, &child10, &child11}

    var expectedOutErr error

    treeReturned, errorsReturned := unserializeTree(in)

    if !treeReturned.matches(expectedOut) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", in, treeReturned.print(), expectedOut.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", in, errorsReturned, expectedOutErr)
    }
}