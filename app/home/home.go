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

func GetHomePage(rw http.ResponseWriter, req *http.Request) {
    type Page struct {
        Title string
        Conditions []Condition
    }

    sqlConditions := getConditions()
    conditionsTree := unserializeRawTree(sqlConditions)
    formattedConditions, err := serializeTree(conditionsTree)
    common.CheckError(err, 2)
    
    p := Page{
        Title: "home",
        Conditions: formattedConditions,
    }

    common.Templates = template.Must(template.ParseFiles("templates/home/home.html", common.LayoutPath))
    err = common.Templates.ExecuteTemplate(rw, "base", p)
    common.CheckError(err, 2)
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

    GetHomePage(rw, req)
}

func Truncate(rw http.ResponseWriter, req *http.Request) {
    _, err := common.DB.Query("TRUNCATE TABLE logictree.equality")
    common.CheckError(err, 2)

    _, err = common.DB.Query("TRUNCATE TABLE logictree.logic")
    common.CheckError(err, 2)

    GetHomePage(rw, req)
}




