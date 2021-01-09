package main

import (
	"strings"
)

type node struct {
	parent *node
	left   *node
	right  *node
	data   string
}

// ParseTree is struct with basic data for parse tree
type ParseTree struct {
	root    *node
	curNode *node
}

// Create tree
func Create(data []string) *ParseTree {
	tree := &ParseTree{root: &node{}, curNode: nil}
	tree.curNode = tree.root

	for _, char := range data {
		// fmt.Println(char)
		// fmt.Println(tree.curNode)
		// If left paren, create new left node
		if char == "(" {
			tree.curNode.insertLeft("")
			tree.curNode = tree.curNode.left
			// If operator, set data in current node and create new right one
		} else if strings.Contains("+-*/d^", char) {
			for tree.curNode.data != "" {
				tree.curNode = tree.curNode.setParent("")
			}
			tree.curNode.setData(char)
			tree.curNode.insertRight("")
			tree.curNode = tree.curNode.right
			// If number, set data in current node and move up a node
		} else if strings.ContainsAny(char, "1234567890") {
			tree.curNode.setData(char)
			tree.curNode = tree.curNode.setParent("")
			// fmt.Println("Up")
			// If right paren, move up a node
		} else if char == ")" {
			tree.curNode = tree.curNode.setParent("")
			// fmt.Println("Up")
		}
	}

	for tree.curNode.parent != nil {
		tree.curNode = tree.curNode.parent
	}
	tree.root = tree.curNode

	return tree
}

func (tree *ParseTree) String() string {
	layer := make([]*node, 1)
	nextLayer := make([]*node, 0)

	layer[0] = tree.root

	output := ""

	for len(layer) != 0 {
		for _, n := range layer {
			// fmt.Println(layer)
			if n != nil {
				output += n.String()
				nextLayer = append(nextLayer, n.left, n.right)
			} else {
				output += " nil "
			}
		}
		output += "\n"
		layer = nextLayer
		nextLayer = nil
	}

	return output
}

func (n *node) String() string {
	var output string
	if n.left != nil {
		output += " /"
	} else {
		output += "  "
	}
	if n.data != "" {
		output += n.data
	} else {
		output += "_"
	}
	if n.right != nil {
		output += "\\ "
	} else {
		output += "  "
	}

	return output
}

// Return the node's parent, or make a new one if needed
func (n *node) setParent(d string) *node {
	// fmt.Println("Up")
	if n.parent == nil {
		n.parent = &node{parent: nil, left: n, right: nil, data: ""}
	}
	return n.parent
}

// Insert a new node to the left
func (n *node) insertLeft(d string) {
	// fmt.Println("Left")
	if n == nil {
		return
	} else if n.left == nil {
		n.left = &node{parent: n, left: nil, right: nil, data: d}
	}
}

// Insert a new node to the right
func (n *node) insertRight(d string) {
	// fmt.Println("Right")
	if n == nil {
		return
	} else if n.right == nil {
		n.right = &node{parent: n, left: nil, right: nil, data: d}
	}
}

// Set the data in the current node
func (n *node) setData(d string) {
	// fmt.Println("Center", d)
	if n == nil {
		return
	} else if n.data == "" {
		n.data = d
	}
}
