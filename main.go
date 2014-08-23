package main

import (
    "net/http"
    "time"
    "flag"

    "git-go-logictree/app/common"
    "git-go-logictree/app/home"

    "github.com/gorilla/mux"
    "github.com/golang/glog"
)

var router *mux.Router

func main() {
    flag.Parse()
    defer glog.Flush()
    defer common.DB.Close()

    router = mux.NewRouter()
    http.HandleFunc("/", httpInterceptor)

    router.HandleFunc("/", home.GetHomePage).Methods("GET")
    router.HandleFunc("/conditions", home.UpdateConditions).Methods("PUT")
    router.HandleFunc("/conditions", home.DeleteConditions).Methods("DELETE")

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