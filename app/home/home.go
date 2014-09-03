package home

import (
    "fmt"
    "html/template"
    "net/http"

    "github.com/jadekler/git-go-logictree/app/common"
)

type Condition struct {
    Text, Type, Field, Operator, Value string
}

type treeNode struct {
    Parent *treeNode
    Children []*treeNode
    Node Condition
    Left, Right int
}

type conditionSqlRow struct {
    Field, Operator, Value, Type string
    Left, Right int
}

type userSqlRow struct {
    Name string
    Id, Age, NumPets int
}

func GetHomePage(rw http.ResponseWriter, req *http.Request) {
    type Page struct {
        Title, FrontendJSON string
        Conditions []Condition
        ConditionSqlRows []conditionSqlRow
        UserSqlRows []userSqlRow
    }

    frontendJSON, conditionSqlRows, formattedConditions := getFrontendJSON()
    
    p := Page{
        Title: "home",
        Conditions: formattedConditions,
        FrontendJSON: frontendJSON,
        ConditionSqlRows: conditionSqlRows,
        UserSqlRows: getUserSqlRows(),
    }

    common.Templates = template.Must(template.ParseFiles("templates/home/home.html", common.LayoutPath))
    err := common.Templates.ExecuteTemplate(rw, "base", p)
    common.CheckError(err, 2)
}

func ResetConditions(rw http.ResponseWriter, req *http.Request) {
    resetType := req.FormValue("resetType");

    switch resetType {
    case "simple":
        treeRoot := &treeNode{Parent: nil, Children: nil, Node: Condition{Text: "AND", Type: "logic", Operator: "AND"}}

        child1 := treeNode{Parent: treeRoot, Children: nil, Node: Condition{Text: "age gt 4", Type: "equality", Field: "age", Operator: "gt", Value: "4"}}
        child2 := treeNode{Parent: treeRoot, Children: nil, Node: Condition{Text: "num_pets lt 2", Type: "equality", Field: "num_pets", Operator: "lt", Value: "2"}}
        treeRoot.Children = []*treeNode{&child1, &child2}

        equalitySql, logicSql, err := treeRoot.toMysql()
        common.CheckError(err, 2)

        updateDatabase(equalitySql, logicSql, testingMysqlUsersInput)
    case "advanced":
        fallthrough
    default:
        beforeEach("no")

        updateDatabase(testingMysqlEqualityInput, testingMysqlLogicInput, testingMysqlUsersInput)
    }

    frontendJSON, _, _ := getFrontendJSON()
    rw.Write([]byte(frontendJSON))
}

func UpdateConditions(rw http.ResponseWriter, req *http.Request) {
    conditions := req.FormValue("conditions");

    parsedConditions, _ := parseJSON(conditions);
    treeRoot, err := unserializeFormattedTree(parsedConditions)

    equalityStr, logicStr, err := treeRoot.toMysql()
    if err != nil {
        fmt.Println(err)
        return
    }

    updateDatabase(equalityStr, logicStr, "")

    frontendJSON, _, _ := getFrontendJSON()

    rw.Write([]byte(frontendJSON))
}

func getFrontendJSON() (string, []conditionSqlRow, []Condition) {
    sqlConditions := getConditions()
    conditionsTree := unserializeRawTree(sqlConditions)
    conditionsTreeJSON := conditionsTree.toJSON()
    formattedConditions, err := serializeTree(conditionsTree)
    common.CheckError(err, 2)

    returnedUsers, err := getMatchingUsers()
    common.CheckError(err, 2)
    usersJSON := usersToJSON(returnedUsers)

    combinedJSON := fmt.Sprintf(`{"tree": %s, "matchingUsers": %s}`, conditionsTreeJSON, usersJSON)

    return combinedJSON, sqlConditions, formattedConditions
}


