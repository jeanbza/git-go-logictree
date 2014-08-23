package home

import (
    "testing"
)

// Fullstack test: given some conditions in JSON, we should be able to parse to condition slice, serialize
// to a tree, attach lefts and rights, and finally convert to mysql value rows to be inserted
func TestFullstack(t *testing.T) {
    beforeEach("home")

    // Parse from json
    conditionsReturned, errorsReturned := parseJSON(testingJSONFrontend)

    if !conditionsMatchesArray(conditionsReturned, testingConditions) {
        t.Errorf("parseJSON(%v) conditionsReturned - got %v, want %v", testingJSONFrontend, conditionsReturned, testingConditions)
    }

    var expectedOutErr error
    if errorsReturned != expectedOutErr {
        t.Errorf("parseJSON(%v) errorsReturned - got %v, want %v", testingJSONFrontend, errorsReturned, expectedOutErr)
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





