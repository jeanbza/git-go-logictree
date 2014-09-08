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

### What's it about?

This program is a pretty simple [lexical analyzer](http://en.wikipedia.org/wiki/Lexical_analysis) that accepts human-readable conditions, such as 'age greater than 5', and translates it to mysql-queryable language. The basic case is fairly simple, but the more advanced cases, such as '(age is greater than 5 or number of pets is less than 2) and age is less than 9', require a bit more work. The intermediary step between human-readable conditions and mysql conditions is implemented as a tree. Additionally, the human-readable conditions are stored in mysql as a left-right hierarchy tree (which represents nested sets). See more detail on the tree and nested sets below.

### Three-stage representation

The setup for this app is: human-readable on the front, golang tree in the middle, mysql nested sets in the back. Here is a basic overview of how the stages interact:

![alt tag](https://raw.github.com/jadekler/git-go-logictree/master/static/images/flowchart.png)

Below are the explanation for each stage.

#### Frontend representation of conditions - Human-readable language

The front end is represented as a set of parenthesis, equality conditions, and logical conditions. Combined, they look like this:

![alt tag](https://raw.github.com/jadekler/git-go-logictree/master/static/images/conditions.png)

Each block is draggable - by clicking the 'Save Re-ordering' button on the page, the re-ordered conditions are saved in mysql. The conditions will also be converted to a mysql query and executed against a dummy users table, with results displayed at the bottom of the page.

#### Server representation of conditions - Tree

We use a tree as an intermediary between the human-readable conditions on the front and interactions with mysql, including saving the conditions, executing the conditions, and converting the conditions (via tree) into human-readable conditions for the front. The tree is a simple n-child tree that is generally traversed post-order. The branches are logical conditions, and the leaves are equality conditions. For instance, (A OR B) would be represented as OR being the branch and A and B being two children of OR, and so on. See below a larger example (taken from the app):

![alt tag](https://raw.github.com/jadekler/git-go-logictree/master/static/images/tree.png)

#### Database representation of conditions - Nested sets

The general idea behind storing a tree in mysql as nested sets is well explained [here](http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/). The basic idea, however, is that each condition is a set, and each child (in the aforementioned tree) is a set within the set that is its parent condition. Once we have so ordered our sets, we can assign lefts and rights to them as below:

![alt tag](https://raw.github.com/jadekler/git-go-logictree/master/static/images/nested_sets.png)

In mysql, this becomes:

```
condition, left, right
----------------------
AND, 1, 24
OR, 2, 17
AND, 3, 14
age = 4, 4, 5
age = 5, 6, 7
age = 6, 8, 9
age = 7, 10, 11
age = 8, 12, 13
age = 1, 15, 16
OR, 18, 23
age = 2, 19, 20
age = 3, 21, 22
```
