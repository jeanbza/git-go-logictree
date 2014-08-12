package home

import (
    "errors"
)

func serializeTree(node *treeNode) ([]Condition, error) {
    if node.Children == nil || len(node.Children) == 0 {
        // Has no children - should be equality

        if node.Node.Type != "equality" {
            return nil, errors.New("ERROR: This tree has a logic condition as a leaf. Quitting.")
        }

        return []Condition{node.Node}, nil
    } else {
        // Has children - should be logic

        if node.Node.Type != "logic" {
            return nil, errors.New("ERROR: This tree has an equality condition as a branch. Quitting.")
        }
    }

    linearConditions := []Condition{Condition{Text: "(", Type: "scope", Operator: "("}}

    for key, child := range node.Children {
        if key != 0 {
            linearConditions = append(linearConditions, node.Node)
        }

        serializedChild, err := serializeTree(child)

        if err != nil {
            return nil, err
        }

        linearConditions = append(linearConditions, serializedChild...)
    }

    linearConditions = append(linearConditions, Condition{Text: ")", Type: "scope", Operator: ")"})

    return linearConditions, nil
}