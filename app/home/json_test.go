package home

import (
    "testing"
    "strings"
)

func TestParseTreeToJSON(t *testing.T) {
    beforeEach("json")

    jsonReturned := testingTreeRoot.toJSON()
    trimmedReturned := strings.Replace(strings.Replace(jsonReturned, " ", "", -1), "\n", "", -1)
    trimmedExpected := strings.Replace(strings.Replace(testingJSONTree, " ", "", -1), "\n", "", -1)

    if trimmedReturned != trimmedExpected {
        t.Errorf("serializeTree(%v) conditionsReturned - got %v, want %v", testingTreeRoot, trimmedReturned, trimmedExpected)
    }
}

func TestParseJSONFrontend(t *testing.T) {
    beforeEach("json")

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