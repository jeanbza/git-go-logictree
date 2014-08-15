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


WRONG



Iterate through conditions
    If current node right is below stored node right OR stored node is not set
        If current node is a branch
            Pop item from conditions
            Store new node including left + right
            Apply following nodes as children
        If node is a leaf
            Pop item from conditions
            Add to children of stored node
    If current node is above stored node right
        Pop all the way back up until we're under the right parent (aka, the parent's left and right encompasses this node)
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
    var node, root *treeNode

    for 0 < len(conditions) {
        // Pop the front item from the slice
        condition = conditions[0]
        conditions = conditions[1:len(conditions)]

        if root == nil {
            // Root hasn't been set - let's do so now and move on immediately
            root = &treeNode{Node: condition.conv(), Left: condition.Left, Right: condition.Right}
            node = root
        } else {
            for condition.Left > node.Right {
                // Current node is above stored node right - keep going up until
                // we find the right parent. The right parent is the one whose left
                // and right encompass this node's left and right
                node = node.Parent
            }

            // Add this node as a child of the parent
            node.Children = append(node.Children, &treeNode{Node: condition.conv(), Left: condition.Left, Right: condition.Right, Parent: node})

            if node.Left != node.Right-1 {
                // This node is a branch - set it as the current node (aka drill down a level)
                node = node.Children[len(node.Children)-1]
            }
        }
    }

    return root, conditions
}





