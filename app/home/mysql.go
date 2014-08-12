package home

import (
    "errors"
    "fmt"
    "git-go-logictree/app/common"
)

func (t *treeNode) toMysql() (equalityStr, logicStr string, err error) {
    t.attachLeftsAndRights()

    equalityStr, logicStr = t.toMysqlRecursively()

    if len(equalityStr) == 0 {
        return "", "", errors.New(fmt.Sprintf("Error: equality string was empty, which was unexpected. EqualityStr: %s :: LogicStr: %s", equalityStr, logicStr))
    }

    if len(logicStr) == 0 {
        return "", "", errors.New(fmt.Sprintf("Error: logic string was empty, which was unexpected. EqualityStr: %s :: LogicStr: %s", equalityStr, logicStr))
    }

    equalityStr = equalityStr[:(len(equalityStr)-1)]
    logicStr = logicStr[:(len(logicStr)-1)]

    return equalityStr, logicStr, nil
}

func (t *treeNode) toMysqlRecursively() (equalityStr, logicStr string) {
    var equalityTemp, logicTemp string
    for _, child := range t.Children {
        equalityTemp, logicTemp = child.toMysqlRecursively()
        equalityStr += equalityTemp
        logicStr += logicTemp
    }

    switch (t.Node.Type) {
    case "equality":
        equalityStr += fmt.Sprintf("('%s', '%s', '%s', 'equality', %d, %d),", t.Node.Field, t.Node.Operator, t.Node.Value, t.Left, t.Right)
    case "logic":
        logicStr += fmt.Sprintf("('%s', 'logic', %d, %d),", t.Node.Operator, t.Left, t.Right)
    }

    return equalityStr, logicStr
}

func updateDatabase(equalityStr, logicStr string) {
    _, err := common.DB.Query("TRUNCATE TABLE logictree.conditions")
    common.CheckError(err, 2)

    _, err = common.DB.Query("INSERT INTO logictree.conditions (field, operator, value, type, lt, rt) VALUES "+equalityStr)
    common.CheckError(err, 2)
    _, err = common.DB.Query("INSERT INTO logictree.conditions (operator, type, lt, rt) VALUES "+logicStr)
    common.CheckError(err, 2)
}

func getConditions() []Condition {
    conditions := make([]Condition, 0)

    rows, err := common.DB.Query("SELECT field, operator, value FROM logictree.equality")
    common.CheckError(err, 2)

    var field, operator, value string

    i := 0

    for rows.Next() {
        rows.Scan(&field, &operator, &value)
        common.CheckError(err, 2)

        if i != 0 {
            conditions = append(conditions, Condition{
                Text: "AND",
                Operator: "AND",
                Type: "logic",
            })
        }

        conditions = append(conditions, Condition{
            Text: fmt.Sprintf("%s %s %s", field, operator, value),
            Type: "equality",
            Field: field,
            Operator: operator,
            Value: value,
        })

        i++
    }

    return conditions
}