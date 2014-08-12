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
    
    p := Page{
        Title: "home",
        Conditions: getConditions(),
    }

    common.Templates = template.Must(template.ParseFiles("templates/home/home.html", common.LayoutPath))
    err := common.Templates.ExecuteTemplate(rw, "base", p)
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




