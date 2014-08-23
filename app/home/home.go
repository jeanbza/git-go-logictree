package home

import (
    "fmt"
    "html/template"
    "net/http"

    "git-go-logictree/app/common"
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

func getFrontendJSON() (string, []Condition) {
    sqlConditions := getConditions()
    conditionsTree := unserializeRawTree(sqlConditions)
    formattedConditions, err := serializeTree(conditionsTree)
    common.CheckError(err, 2)

    return conditionsTree.toJSON(), formattedConditions
}

func GetHomePage(rw http.ResponseWriter, req *http.Request) {
    type Page struct {
        Title string
        Conditions []Condition
        TreeJSON string
    }

    frontendJSON, formattedConditions := getFrontendJSON()
    
    p := Page{
        Title: "home",
        Conditions: formattedConditions,
        TreeJSON: frontendJSON,
    }

    common.Templates = template.Must(template.ParseFiles("templates/home/home.html", common.LayoutPath))
    err := common.Templates.ExecuteTemplate(rw, "base", p)
    common.CheckError(err, 2)
}

func ResetConditions(rw http.ResponseWriter, req *http.Request) {
    beforeEach("no")

    updateDatabase(testingMysqlEqualityInput, testingMysqlLogicInput)

    frontendJSON, _ := getFrontendJSON()

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

    updateDatabase(equalityStr, logicStr)

    frontendJSON, _ := getFrontendJSON()

    rw.Write([]byte(frontendJSON))
}




