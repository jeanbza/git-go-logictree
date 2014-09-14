package common

import (
    "html/template"
    "runtime"
    "path"
)

var Templates *template.Template
var AppDir string
var LayoutPath string

func init() {
    _, filename, _, _ := runtime.Caller(1)
    AppDir = path.Dir(filename) + "/../../"
    LayoutPath = AppDir + "templates/layout.html"
}