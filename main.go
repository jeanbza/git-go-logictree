package main

import (
    "net/http"
    "time"
    "flag"

    "git-misc/logic-tree/app/common"
    "git-misc/logic-tree/app/home"

    "github.com/gorilla/mux"
    "github.com/golang/glog"
)

var router *mux.Router

func main() {
    flag.Parse()
    defer glog.Flush()

    router = mux.NewRouter()
    http.HandleFunc("/", httpInterceptor)

    router.HandleFunc("/", home.GetHomePage).Methods("GET")
    router.HandleFunc("/updateConditions", home.UpdateConditions).Methods("POST")

    fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
    http.Handle("/static/", fileServer)

    http.ListenAndServe(":8080", nil)
}

func httpInterceptor(w http.ResponseWriter, req *http.Request) {
    startTime := time.Now()

    router.ServeHTTP(w, req)

    finishTime := time.Now()
    elapsedTime := finishTime.Sub(startTime)

    common.LogAccess(w, req, elapsedTime)
}