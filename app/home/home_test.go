package home

import (
    "fmt"
    "testing"
)

var testingTreeRoot *treeNode
var testingConditions []Condition
var testingJSON, testingMysqlEqualityInput, testingMysqlLogicInput string
var testingMysqlEqualityOutput []equalitySqlRow
var testingMysqlLogicOutput []logicSqlRow

// Fullstack test: given some conditions in JSON, we should be able to parse to condition slice, serialize
// to a tree, attach lefts and rights, and finally convert to mysql value rows to be inserted
func TestFullstack(t *testing.T) {
    beforeEach("home")

    // Parse from json
    conditionsReturned, errorsReturned := parseJSON(testingJSON)

    if !conditionsMatchesArray(conditionsReturned, testingConditions) {
        t.Errorf("parseJSON(%v) conditionsReturned - got %v, want %v", testingJSON, conditionsReturned, testingConditions)
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("parseJSON(%v) errorsReturned - got %v, want %v", testingJSON, errorsReturned, expectedOutErr)
    }

    // Because slices are pointers by default, and unserialize pops items, we shallow copy a new version for later use
    var originalConditions []Condition
    for _, condition := range testingConditions {
        originalConditions = append(originalConditions, condition)
    }

    // Unserialize into a tree
    treeReturned, errorsReturned := unserializeTree(conditionsReturned)

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", conditionsReturned, treeReturned.print(), testingTreeRoot.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", conditionsReturned, errorsReturned, expectedOutErr)
    }

    // Convert tree to mysql input rows
    equalityReturned, logicReturned := treeReturned.toMysql()

    if equalityReturned != testingMysqlEqualityInput {
        t.Errorf("%v.toMysql() equalityReturned - got %v, want %v", treeReturned, equalityReturned, testingMysqlEqualityInput)
    }

    if logicReturned != testingMysqlLogicInput {
        t.Errorf("%v.toMysql() logicReturned - got %v, want %v", treeReturned, logicReturned, testingMysqlLogicInput)
    }
}

// Roundtrip test: if we serialize a tree, then unserialize the result, we should get the original tree back
func TestSerializationRoundtrip(t *testing.T) {
    beforeEach("home")

    // Because slices are pointers by default, and unserialize pops items, we shallow copy a new version for later use
    var originalConditions []Condition
    for _, condition := range testingConditions {
        originalConditions = append(originalConditions, condition)
    }

    // Unserialize into a tree
    treeReturned, errorsReturned := unserializeTree(testingConditions)

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeTree(%v) - got %v, want %v", testingConditions, treeReturned.print(), testingTreeRoot.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeTree(%v) errorsReturned - got %v, want %v", testingConditions, errorsReturned, expectedOutErr)
    }

    // Serialize back into conditions array
    conditionsReturned, errorsReturned := serializeTree(treeReturned)

    if !conditionsMatchesArray(conditionsReturned, originalConditions) {
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

    if !conditionsMatchesArray(conditionsReturned, expectedOut) {
        t.Errorf("parseJSON(%v) conditionsReturned - got %v, want %v", expectedOut, conditionsReturned, expectedOut)
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("parseJSON(%v) errorsReturned - got %v, want %v", expectedOut, errorsReturned, expectedOutErr)
    }
}

func conditionsMatchesArray(conditionsA, conditionsB []Condition) bool {
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

func equalitySqlMatchesArray(rowsA, rowsB []equalitySqlRow) bool {
    var truth bool

    if rowsA == nil || len(rowsA) != len(rowsB) {
        return false
    }

    for _, valA := range rowsA {
        truth = false

        for _, valB := range rowsB {
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

func logicSqlMatchesArray(rowsA, rowsB []logicSqlRow) bool {
    var truth bool

    if rowsA == nil || len(rowsA) != len(rowsB) {
        return false
    }

    for _, valA := range rowsA {
        truth = false

        for _, valB := range rowsB {
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

func (a logicSqlRow) matches(b logicSqlRow) bool {
    if a.Operator != b.Operator {
        return false
    }

    if a.Left != b.Left {
        return false
    }

    if a.Right != b.Right {
        return false
    }

    return true
}

func (a equalitySqlRow) matches(b equalitySqlRow) bool {
    if a.Field != b.Field {
        return false
    }

    if a.Operator != b.Operator {
        return false
    }

    if a.Value != b.Value {
        return false
    }

    if a.Left != b.Left {
        return false
    }

    if a.Right != b.Right {
        return false
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

func beforeEach(testName string) {
    fmt.Printf("Starting %s tests..\n", testName)

    /**
     * lt-node-rt
     *                                     1-AND-24
     *                          2-OR-17                     18-OR-23
     *              3-AND-14                15-F-16     19-G-20  21-H-22
     * 4-A-5 6-B-7 8-C-9 10-D-11 12-E-13
     */
    testingTreeRoot = nil

    // Row 1 node 1
    testingTreeRoot = &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

    // Row 2 node 1
    child1 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    // Row 2 node 2
    child2 := treeNode{Parent: nil, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    testingTreeRoot.Children = []*treeNode{&child1, &child2}

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

    testingJSON = `
        [
            {"Text": "(", "Type": "scope", "Operator": "("},
            {"Text": "(", "Type": "scope", "Operator": "("},
            {"Text": "(", "Type": "scope", "Operator": "("},
            {"Text": "age eq 4", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "4"},
            {"Text": "AND", "Type": "logic", "Operator": "AND"},
            {"Text": "age eq 5", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "5"},
            {"Text": "AND", "Type": "logic", "Operator": "AND"},
            {"Text": "age eq 6", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "6"},
            {"Text": "AND", "Type": "logic", "Operator": "AND"},
            {"Text": "age eq 7", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "7"},
            {"Text": "AND", "Type": "logic", "Operator": "AND"},
            {"Text": "age eq 8", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "8"},
            {"Text": ")", "Type": "scope", "Operator": ")"},
            {"Text": "OR", "Type": "logic", "Operator": "OR"},
            {"Text": "age eq 1", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "1"},
            {"Text": ")", "Type": "scope", "Operator": ")"},
            {"Text": "AND", "Type": "logic", "Operator": "AND"},
            {"Text": "(", "Type": "scope", "Operator": "("},
            {"Text": "age eq 2", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "2"},
            {"Text": "OR", "Type": "logic", "Operator": "OR"},
            {"Text": "age eq 3", "Type": "equality", "Field": "age", "Operator": "eq", "Value": "3"},
            {"Text": ")", "Type": "scope", "Operator": ")"},
            {"Text": ")", "Type": "scope", "Operator": ")"}
        ]
    `

    testingConditions = []Condition{
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

    // INSERT INTO logictree.equality (field, operator, value, lt, rt) VALUES ...
    testingMysqlEqualityInput = "('age', 'eq', '4', 4, 5),('age', 'eq', '5', 6, 7),('age', 'eq', '6', 8, 9),('age', 'eq', '7', 10, 11),('age', 'eq', '8', 12, 13),('age', 'eq', '1', 15, 16),('age', 'eq', '2', 19, 20),('age', 'eq', '3', 21, 22)"
    // INSERT INTO logictree.logic (operator, lt, rt) VALUES ...
    testingMysqlLogicInput = "('AND', 3, 14),('OR', 2, 17),('OR', 18, 23),('AND', 1, 24)"

    testingMysqlEqualityOutput = []equalitySqlRow{
        equalitySqlRow{Field: "age", Operator: "eq", Value: "4", Left: 4, Right: 5},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "5", Left: 6, Right: 7},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "6", Left: 8, Right: 9},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "7", Left: 10, Right: 11},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "8", Left: 12, Right: 13},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "1", Left: 15, Right: 16},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "2", Left: 19, Right: 20},
        equalitySqlRow{Field: "age", Operator: "eq", Value: "3", Left: 21, Right: 22},
    }

    testingMysqlLogicOutput = []logicSqlRow{
        logicSqlRow{Operator: "AND", Left: 3, Right: 14},
        logicSqlRow{Operator: "OR", Left: 2, Right: 17},
        logicSqlRow{Operator: "OR", Left: 18, Right: 23},
        logicSqlRow{Operator: "AND", Left: 1, Right: 24},
    }
}





