package home

import (
    "fmt"
    "testing"
)

func beforeEach(testName string) {
    fmt.Printf("Starting %s tests..\n", testName)
}

// Fullstack test
/**
 * ((A && B) || C) && (D || E)
 *              AND
 *        OR           OR
 *    AND     F      G    H
 * A B C D E
 */
func TestFullstack(t *testing.T) {
    beforeEach("home")

    conditionsIn := []Condition{
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

    // Because slices are pointers by default, and unserialize pops items, we shallow copy a new version for later use
    var originalConditions []Condition
    for _, condition := range conditionsIn {
        originalConditions = append(originalConditions, condition)
    }

    treeRootExpectation := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    treeRootExpectation.Children = []*treeNode{&child1, &child2}

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

    // Unserialize into a tree
    var expectedOutErr error

    treeReturned, errorsReturned := unserializeTree(conditionsIn)

    if !treeReturned.matches(treeRootExpectation) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", conditionsIn, treeReturned.print(), treeRootExpectation.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", conditionsIn, errorsReturned, expectedOutErr)
    }

    // Serialize back into conditions array
    conditionsReturned, errorsReturned := serializeTree(treeReturned)

    if !matchesArray(conditionsReturned, originalConditions) {
        t.Errorf("serializeTree(%v) conditionsReturned - got %v, want %v", treeReturned, simplifyConditions(conditionsReturned), simplifyConditions(originalConditions))
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("serializeTree(%v) errorsReturned - got %v, want %v", treeReturned, errorsReturned, expectedOutErr)
    }
}

func TestParseJSON(t *testing.T) {
    beforeEach("home")

    in := `
        [
            {
                "Text": "age eq 8",
                "Type": "equality",
                "Field": "age",
                "Operator": "eq",
                "Value": "8"
            },
            {
                "Text": "(",
                "Type": "scope",
                "Operator": "("
            },
            {
                "Text": "AND",
                "Type": "logic",
                "Operator": "AND"
            }
        ]
    `
    expectedOut := []Condition{
        Condition{Text: "(", Type: "scope", Operator: "("},
        Condition{Text: "age eq 8", Type: "equality", Field: "age", Operator: "eq", Value: "8"},
        Condition{Text: "AND", Type: "logic", Operator: "AND"},
    }
    var expectedOutErr error

    conditionsReturned, errorsReturned := parseJSON(in)

    if !matchesArray(conditionsReturned, expectedOut) {
        t.Errorf("parseJSON(%v) conditionsReturned - got %v, want %v", expectedOut, conditionsReturned, expectedOut)
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("parseJSON(%v) errorsReturned - got %v, want %v", expectedOut, errorsReturned, expectedOutErr)
    }
}

func matchesArray(conditionsA []Condition, conditionsB []Condition) bool {
    var truth bool

    if conditionsA == nil || len(conditionsA) != len(conditionsB) {
        return false
    }

    for _, valA := range conditionsA {
        truth = false

        for _, valB := range conditionsB {
            if valA.matches(valB) {
                truth = true
            }
        }

        if !truth {
            return false
        }
    }

    return true
}

// Only matches DOWNWARDS - not up the parent chain
func (treeNodeA *treeNode) matches(treeNodeB *treeNode) bool {
    if treeNodeA == nil || treeNodeB == nil {
        return false
    }

    if len(treeNodeA.Children) != len(treeNodeB.Children) {
        return false
    }

    for key, child := range treeNodeA.Children {
        if !child.matches(treeNodeB.Children[key]) || child.Left != treeNodeB.Children[key].Left || child.Right != treeNodeB.Children[key].Right {
            return false
        }
    }
    
    if !treeNodeA.Node.matches(treeNodeB.Node) {
        return false
    }

    return true
}

func (conditionA Condition) matches(conditionB Condition) bool {
    if conditionA.Text != conditionB.Text {
        return false
    }

    if conditionA.Type != conditionB.Type {
        return false
    }

    if conditionA.Field != conditionB.Field {
        return false
    }

    if conditionA.Operator != conditionB.Operator {
        return false
    }

    if conditionA.Value != conditionB.Value {
        return false
    }

    return true
}

