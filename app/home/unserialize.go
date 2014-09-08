package home

// Recursively creates a tree from an array of conditions using the logic below
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

// Iteratively creates a tree from an array of mysql conditions using the logic below
/*
Iterative approach, since we're handed a flat array:
Iterate through conditions
    Pop a condition from conditions
    If root is unset
        Set it with popped condition
        Set root as the current node
    Else
        If popped condition's left is greater than popped item's right
            You've gone past the encompassed (left and right within parent's left and right) children. This node may be a sibling, uncle, etc. Pop 
            back up the parent chain until you find a parent that once again encompasses this child
        Add the popped condition to the current node's children
        If the current node is a branch
            Go down a level - set it as the current branch
*/
func unserializeRawTree(conditions []conditionSqlRow) *treeNode {
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

    return root
}





