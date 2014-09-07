logic-tree
===============

[![Build Status](https://travis-ci.org/jadekler/git-go-logictree.svg?branch=master)](https://travis-ci.org/jadekler/git-go-logictree)

### What's it about?

This program is a pretty simple [lexical analyzer](http://en.wikipedia.org/wiki/Lexical_analysis) that accepts human-readable conditions, such as 'age greater than 5', and translates it to mysql-queryable language. The basic case is fairly simple, but the more advanced cases, such as '(age is greater than 5 or number of pets is less than 2) and age is less than 9', require a bit more work. The intermediary step between human-readable conditions and mysql conditions is implemented as a tree. Additionally, the human-readable conditions are stored in mysql as a left-right hierarchy tree (which represents nested sets). See more detail on the tree and nested sets below.

##### Frontend representation of conditions - Human-readable language

The front end is represented as a set of parenthesis, equality conditions, and logical conditions. Combined, they look like this:

![alt tag](https://raw.github.com/jadekler/git-go-logictree/master/static/images/conditions.png)

Each block is draggable - by clicking the 'Save Re-ordering' button on the page, the re-ordered conditions are saved in mysql. The conditions will also be converted to a mysql query and executed against a dummy users table, with results displayed at the bottom of the page.

##### Server representation of conditions - Tree

##### Database representation of conditions - Nested sets

### Installation

1. `go get github.com/jadekler/git-go-logictree`
1. `go run $GOPATH/src/github.com/jadekler/git-go-logictree/main.go`
1. Navigate to `localhost:8080` in your browser

### Testing

To test locally, simply run `go test -v ./...` from the project root.

Note: Useful testing command: `printf "$(go test $GOPATH/src/github.com/jadekler/git-go-logictree/app/home | sed 's/::/\\n/g')" && echo;`
