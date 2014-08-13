package home

import "fmt"

/** Treat conditions like a queue. Rules:
 * If you reach a (, pop the condition, drop down a depth and assign results to root's children
 * If you reach a ), pop the condition, pop back up a depth with the root
 * If you reach a logical condition, pop the condition, assign it as the root's node
 * If you reach an equality condition, pop the condition, assign it as one of the children of the root
 * At the end of the loop, return the root's first child (since we have parans around all conditions we're going to be one level too deep)
**/
func unserializeFormattedTree(conditions []Condition) (*treeNode, error) {
    var root treeNode
    var emptyNode, condition Condition

    key := 0

    for key < len(conditions) {
        // Pop the front item from the slice
        condition = conditions[0]
        conditions = append(conditions[:key], conditions[key+1:]...)

        switch condition.Type {
        case "scope":
            if condition.Operator == "(" {
                children, _ := unserializeFormattedTree(conditions)

                if len(root.Children) == 0 {
                    root.Children = []*treeNode{children}
                } else {
                    root.Children = append(root.Children, children)
                }
            }

            if condition.Operator == ")" {
                if root.Node == emptyNode {
                    return root.Children[0], nil
                } else {
                    return &root, nil
                }
            }
        case "logic":
            root.Node = condition
        case "equality":
            root.Children = append(root.Children, &treeNode{Parent: &root, Node: condition})
        }
    }

    if len(root.Children) > 0 && !root.Node.matches(emptyNode) {
        return &root, nil
    }

    return root.Children[0], nil
}

func unserializeRawTree(conditions []conditionSqlRow) *treeNode {
    conditionsToReturn, _ := unserializeRawTreeRecursively(conditions)

    return conditionsToReturn
}

/**
 * Assumption: data should be ordered by LEFT
 * Recursive steps:
 * If current node is a branch:
 *      Create node with it
 *      Pop item from conditions
 *      Add children recursively
 * If node is a leaf:
 *      Create node with it
 *      Pop item from conditions
 * Return node


WRONG



Iterate through conditions
    If current node right is below stored node right
        If current node is a branch
            Store new node including left + right
            Pop item from conditions
            Apply following nodes as children
        If node is a leaf
            Pop item from conditions
            Return to be added to children of stored node
    If current node is above stored node right
        Stored node becomes previously stored node
        Apply conditions of part 1
        


RECURSIVELY:
     
     ?????
     
Iterate through conditions
    If current node right is below stored node right
        If current node is a branch
            Store new node including left + right
            Pop item from conditions
            Apply following nodes as children
        If node is a leaf
            Pop item from conditions
            Return to be added to children of stored node
    If current node is above stored node right
        Pop up a level
        Do NOT pop item off conditions

**/
func unserializeRawTreeRecursively(conditions []conditionSqlRow) (*treeNode, []conditionSqlRow) {
    var condition conditionSqlRow
    var node *treeNode
    fmt.Println("asdd")
    root := &treeNode{}

    // for 0 < len(conditions) {
    if len(conditions) > 0 {
        // Pop the front item from the slice
        condition = conditions[0]
        conditions = conditions[1:len(conditions)]

        root.Node = condition.conv()
        root.Left = condition.Left
        root.Right = condition.Right

        if condition.Left != condition.Right-1 {
            // Has children
            once := false

            // (the once is basically a do while)
            for !once || node.Left == node.Right-1 {
                once = true
                node, conditions = unserializeRawTreeRecursively(conditions)
                
                if node.Left != 0 {
                    root.Children = append(root.Children, node)
                }
            }
        } else {
            // Has no children
            return root, conditions
        }
    }

    return root, conditions
}





