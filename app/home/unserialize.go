package home

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
**/
func unserializeRawTreeRecursively(conditions []conditionSqlRow) (*treeNode, []conditionSqlRow) {
    var condition conditionSqlRow
    var node *treeNode
    root := &treeNode{}

    // PROBLEM: How to add children recursively

    for 0 < len(conditions) {
        // Pop the front item from the slice
        condition = conditions[0]
        conditions = conditions[1:len(conditions)]

        root.Node = condition.conv()
        root.Left = condition.Left
        root.Right = condition.Right

        if condition.Left != condition.Right-1 {
            // Has children
            node, conditions = unserializeRawTreeRecursively(conditions)
            for node.Left == node.Right-1 {
                root.Children = append(root.Children, node)
                node, conditions = unserializeRawTreeRecursively(conditions)
            }
        } else {
            // Has no children
            return root, conditions
        }
    }

    return root, conditions
}