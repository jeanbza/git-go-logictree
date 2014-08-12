package home

import (
    "fmt"
    "testing"
)

var testingTreeRoot *treeNode
var testingConditions []Condition
var testingJSON, testingMysqlEqualityInput, testingMysqlLogicInput string
var testingMysqlRows []conditionSqlRow

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

    // Because slices are pointers by default, and unserializeFormatted pops items, we shallow copy a new version for later use
    var originalConditions []Condition
    for _, condition := range testingConditions {
        originalConditions = append(originalConditions, condition)
    }

    // unserializeFormatted into a tree
    treeReturned, errorsReturned := unserializeFormattedTree(conditionsReturned)

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", conditionsReturned, treeReturned.print(), testingTreeRoot.print())
    }

    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", conditionsReturned, errorsReturned, expectedOutErr)
    }

    // Convert tree to mysql input rows
    equalityReturned, logicReturned, _ := treeReturned.toMysql()

    if equalityReturned != testingMysqlEqualityInput {
        t.Errorf("%v.toMysql() equalityReturned - got %v, want %v", treeReturned, equalityReturned, testingMysqlEqualityInput)
    }

    if logicReturned != testingMysqlLogicInput {
        t.Errorf("%v.toMysql() logicReturned - got %v, want %v", treeReturned, logicReturned, testingMysqlLogicInput)
    }
}

// Roundtrip test: if we serialize a tree, then unserializeFormatted the result, we should get the original tree back
func TestSerializationRoundtrip(t *testing.T) {
    beforeEach("home")

    // Because slices are pointers by default, and unserializeFormatted pops items, we shallow copy a new version for later use
    var originalConditions []Condition
    for _, condition := range testingConditions {
        originalConditions = append(originalConditions, condition)
    }

    // unserializeFormatted into a tree
    treeReturned, errorsReturned := unserializeFormattedTree(testingConditions)

    if !treeReturned.matches(testingTreeRoot) {
        t.Errorf("unserializeFormattedTree(%v) - got %v, want %v", testingConditions, treeReturned.print(), testingTreeRoot.print())
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("unserializeFormattedTree(%v) errorsReturned - got %v, want %v", testingConditions, errorsReturned, expectedOutErr)
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
    testingMysqlEqualityInput = "('age', 'eq', '4', 'equality', 4, 5),('age', 'eq', '5', 'equality', 6, 7),('age', 'eq', '6', 'equality', 8, 9),('age', 'eq', '7', 'equality', 10, 11),('age', 'eq', '8', 'equality', 12, 13),('age', 'eq', '1', 'equality', 15, 16),('age', 'eq', '2', 'equality', 19, 20),('age', 'eq', '3', 'equality', 21, 22)"
    // INSERT INTO logictree.logic (operator, lt, rt) VALUES ...
    testingMysqlLogicInput = "('AND', 'logic', 3, 14),('OR', 'logic', 2, 17),('OR', 'logic', 18, 23),('AND', 'logic', 1, 24)"

    testingMysqlRows = []conditionSqlRow{
        conditionSqlRow{Operator: "AND", Type: "logic", Left: 1, Right: 24},
        conditionSqlRow{Operator: "OR", Type: "logic", Left: 2, Right: 17},
        conditionSqlRow{Operator: "AND", Type: "logic", Left: 3, Right: 14},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "4", Type: "equality", Left: 4, Right: 5},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "5", Type: "equality", Left: 6, Right: 7},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "6", Type: "equality", Left: 8, Right: 9},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "7", Type: "equality", Left: 10, Right: 11},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "8", Type: "equality", Left: 12, Right: 13},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "1", Type: "equality", Left: 15, Right: 16},
        conditionSqlRow{Operator: "OR", Type: "logic", Left: 18, Right: 23},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "2", Type: "equality", Left: 19, Right: 20},
        conditionSqlRow{Field: "age", Operator: "eq", Value: "3", Type: "equality", Left: 21, Right: 22},
    }
}





