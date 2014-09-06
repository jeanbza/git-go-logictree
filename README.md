logic-tree
===============

[![Build Status](https://travis-ci.org/jadekler/git-go-logictree.svg?branch=master)](https://travis-ci.org/jadekler/git-go-logictree)

### Installation

1. `go get github.com/jadekler/git-go-logictree`
1. `go run $GOPATH/src/github.com/jadekler/git-go-logictree/main.go`
1. Navigate to `localhost:8080` in your browser

### Testing

To test locally, simply run `go test -v ./...` from the project root.

Note: Useful testing command: `printf "$(go test $GOPATH/src/github.com/jadekler/git-go-logictree/app/home | sed 's/::/\\n/g')" && echo;`
