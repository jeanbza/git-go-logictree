package home

import (
    "encoding/json"
    "fmt"
)

var testingTreeRoot *treeNode
var testingConditions []Condition
var testingJSONFrontend, testingJSONTree, testingMysqlEqualityInput, testingMysqlLogicInput, testingMysqlUsersInput, testingMysqlConditionsInput string
var testingMysqlRows []conditionSqlRow

func usersToJSON(users []userSqlRow) string {
    json := "["

    for key, user := range users {
        if key != 0 {
            json += ","
        }

        json += fmt.Sprintf("{Id: %d, Name: '%s', Age: %d, NumPets: %d}", user.Id, user.Name, user.Age, user.NumPets)
    }

    json += "]"

    return json
}

func (t *treeNode) toJSON() string {
    return "[" + t.toJSONRecursively() + "]"
}

func (t *treeNode) toJSONRecursively() string {
    var childrenJSON string

    for key, child := range t.Children {
        if key != 0 {
            childrenJSON += ","
        }

        childrenJSON += child.toJSONRecursively()
    }

    if childrenJSON != "" {
        return fmt.Sprintf(`
            {"name": "%s", "children": [%s]}
        `, t.Node.textify(), childrenJSON)
    } else {
        return fmt.Sprintf(`
            {"name": "%s"}
        `, t.Node.textify())
    }
}

func (c Condition) textify() string {
    switch c.Type {
    case "equality":
        return c.Field + " " + c.Operator + " " + c.Value
    case "logic":
        return c.Operator
    }

    return ""
}

func parseJSON(conditionsString string) ([]Condition, error) {
    var conditionsSlice []Condition
    
    err := json.Unmarshal([]byte(conditionsString), &conditionsSlice)
    if err != nil {
        return nil, err
    }
    
    return conditionsSlice, nil
}

func (t *treeNode) attachLeftsAndRights() {
    indexStart := 0
    t.attachLeftsAndRightsRecursively(&indexStart)
}

func (t *treeNode) attachLeftsAndRightsRecursively(index *int) {
    *index++
    t.Left = *index

    for _, child := range t.Children {
        child.attachLeftsAndRightsRecursively(index)
    }

    *index++
    t.Right = *index
}

func (t *treeNode) print() string {
    var s string

    if t == nil {
        return ""
    }

    for _, child := range t.Children {
        s += child.print()
    }

    return s + " :: " + fmt.Sprintf("%v, %d, %d", t.Node, t.Left, t.Right)
}

func simplifyConditions(conditions []Condition) string {
    var t string

    for k, c := range conditions {
        if k != 0 {
            t += " "
        }

        t += c.Text
    }

    return t
}

func (t *treeNode) getChildrenConditions() []Condition {
    var children []Condition

    for _, child := range t.Children {
        children = append(children, child.Node)
    }

    return children
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

func (condition conditionSqlRow) conv() Condition {
    switch condition.Type {
    case "logic":
        return Condition{Text: condition.Operator, Type: condition.Type, Field: condition.Field, Operator: condition.Operator, Value: condition.Value}
    case "equality":
        return Condition{Text: condition.Field + " " + condition.Operator + " " + condition.Value, Type: condition.Type, Field: condition.Field, Operator: condition.Operator, Value: condition.Value}
    }

    return Condition{}
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

func conditionSqlMatchesArray(rowsA, rowsB []conditionSqlRow) bool {
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

func usersMatchesArray(rowsA, rowsB []userSqlRow) bool {
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

func (a conditionSqlRow) matches(b conditionSqlRow) bool {
    if a.Field != b.Field {
        return false
    }

    if a.Operator != b.Operator {
        return false
    }

    if a.Value != b.Value {
        return false
    }

    if a.Type != b.Type {
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

func (a userSqlRow) matches(b userSqlRow) bool {
    if a.Name != b.Name {
        return false
    }

    if a.Age != b.Age {
        return false
    }

    if a.NumPets != b.NumPets {
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
    child1 := treeNode{Parent: testingTreeRoot, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
    // Row 2 node 2
    child2 := treeNode{Parent: testingTreeRoot, Children: nil, Node: Condition{Text: "OR", Type: "logic", Operator: "OR"}}
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

    testingJSONTree = `
        [
            {
                "name": "AND",
                "children": [
                    {
                        "name": "OR",
                        "children": [
                            {
                                "name": "AND",
                                "children": [
                                    {"name": "age eq 4"},
                                    {"name": "age eq 5"},
                                    {"name": "age eq 6"},
                                    {"name": "age eq 7"},
                                    {"name": "age eq 8"}
                                ]
                            },
                            {"name": "age eq 1"}
                        ]
                    },
                    {
                        "name": "OR",
                        "children": [
                            {"name": "age eq 2"},
                            {"name": "age eq 3"}
                        ]
                    }
                ]
            }
        ]
    `

    testingJSONFrontend = `
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
    // INSERT INTO logictree.users (name, age, num_pets) VALUES ...
    testingMysqlUsersInput = ""
    for i := 1; i < 500; i++ {
        if i != 1 {
            testingMysqlUsersInput += ","
        }

        testingMysqlUsersInput += fmt.Sprintf("('bob%d', %d, %d)", 5%i, 7%i, 9%i)
    }

    testingMysqlConditionsInput = "(((age = 4 AND age = 5 AND age = 6 AND age = 7 AND age = 8) OR age = 1) AND (age = 2 OR age = 3))"

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



