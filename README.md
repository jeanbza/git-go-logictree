logic-tree
===============

[![Build Status](https://travis-ci.org/jadekler/git-go-logictree.svg?branch=master)](https://travis-ci.org/jadekler/git-go-logictree)

### Installation

1. `go get github.com/jadekler/git-go-logictree`
1. Navigate to the project directory
1. `go test ./...` (make sure mysql is running and root/no pwd is enabled)
1. `go run $GOPATH/src/github.com/jadekler/git-go-logictree/main.go`
1. Navigate to `localhost:8080` in your browser

### Testing

To test locally, simply run `go test -v ./...` from the project root.

Note: Useful testing command: `printf "$(go test $GOPATH/src/github.com/jadekler/git-go-logictree/app/home | sed 's/::/\\n/g')" && echo;`

### What's it about?

This program is a pretty simple [lexical analyzer](http://en.wikipedia.org/wiki/Lexical_analysis) (which only parses) that accepts human-readable conditions, such as 'age greater than 5', and translates it to mysql-queryable language. The basic case is fairly simple, but the more advanced cases, such as '(age is greater than 5 or number of pets is less than 2) and age is less than 9', require a bit more work. The intermediary step between human-readable conditions and mysql conditions is implemented as a tree. Additionally, the human-readable conditions are stored in mysql as a left-right hierarchy tree (which represents nested sets). See more detail on the tree and nested sets below.

### Three-stage representation

The setup for this app is: human-readable on the front, golang tree in the middle, mysql nested sets in the back. Here is a basic overview of how the stages interact:

![](https://raw.github.com/jadekler/git-go-logictree/master/logictree-static/images/flowchart.png)

Below are the explanation for each stage.

#### Frontend representation of conditions - Human-readable language

The front end is represented as a set of parenthesis, equality conditions, and logical conditions. Combined, they look like this:

![](https://raw.github.com/jadekler/git-go-logictree/master/logictree-static/images/conditions.png)

Each block is draggable - by clicking the 'Save Re-ordering' button on the page, the re-ordered conditions are saved in mysql. The conditions will also be converted to a mysql query and executed against a dummy users table, with results displayed at the bottom of the page.

#### Server representation of conditions - Tree

We use a tree as an intermediary between the human-readable conditions on the front and interactions with mysql, including saving the conditions, executing the conditions, and converting the conditions (via tree) into human-readable conditions for the front. The tree is a simple n-child tree that is generally traversed post-order. The branches are logical conditions, and the leaves are equality conditions. For instance, (A OR B) would be represented as OR being the branch and A and B being two children of OR, and so on. See below a larger example (taken from the app):

![](https://raw.github.com/jadekler/git-go-logictree/master/logictree-static/images/tree.png)

#### Database representation of conditions - Nested sets

The general idea behind storing a tree in mysql as nested sets is well explained [here](http://mikehillyer.com/articles/managing-hierarchical-data-in-mysql/). The basic idea, however, is that each condition is a set, and each child (in the aforementioned tree) is a set within the set that is its parent condition. Once we have so ordered our sets, we can assign lefts and rights to them as below:

![](https://raw.github.com/jadekler/git-go-logictree/master/logictree-static/images/nested_sets.png)

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

### Time and Space Complexity

The space complexity is O(n) in both the server language and mysql - that is, we store each node and no additional information both in mysql and in our tree (we use pointers in the tree for parents and children).

The time complexity is as follows:

##### Unserialize From Mysql To Go

The function that deals with this is called [unserializeRawTree](https://github.com/jadekler/git-go-logictree/blob/master/app/home/unserialize.go#L71). This function assumes a set of conditions pulled from mysql ordered by the LEFT column ascending. It iterates through the conditions, popping conditions off and creating the tree until the conditions are empty. So, we touch each condition only once - this function runs in O(n).

##### Unserialize From Frontend * To Go

The function that deals with this is called [unserializeFormattedTree](https://github.com/jadekler/git-go-logictree/blob/master/app/home/unserialize.go#L11). I realize my naming scheme is awful. Anyways, this function iterates over leaves and recurses through branches, building the tree as it goes and relying on the recursive call to glue it all together. So once again, we touch each condition only once - again, this function runs in O(n).

##### Serialize From Go To Frontend *

The function that deals with this is called [serializeTree](https://github.com/jadekler/git-go-logictree/blob/master/app/home/serialize.go#L8). It traverses the tree post-order, building a linear array of conditions as it goes. This function runs in O(n).

##### Serialize From Go To Mysql

The function that deals with this is called [toJSON](https://github.com/jadekler/git-go-logictree/blob/master/app/home/helpers.go#L30) (actually, the helper function toJSONRecursively does the real work) and simply traverses the tree post-order in O(n) time.

\* By frontend, I mean a linear array of conditions, including parenthesis that surround children.

### Questions and feedback

Please shoot any questions or feedback over to jadekler@gmail.com.
