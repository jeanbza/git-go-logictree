package home

import (
    "encoding/json"
    "fmt"
)

func (t *treeNode) toJSON() string {
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