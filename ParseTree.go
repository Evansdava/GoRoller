package main

import (
	"fmt"
	"strings"
)

type node struct {
	parent *node
	left   *node
	right  *node
	data   string
}

type ParseTree struct {
	root    *node
	curNode *node
}

// Create tree
func Create(data []string) *ParseTree {
	tree := &ParseTree{root: &node{}, curNode: nil}
	tree.curNode = tree.root

	for _, char := range data {
		fmt.Println(char)

		// If left paren, create new left node
		if char == "(" {
			tree.curNode.insertLeft("")
			tree.curNode = tree.curNode.left
			// If operator, set data in current node and create new right one
		} else if strings.Contains("+-*/d^", char) {
			tree.curNode.setData(char)
			tree.curNode.insertRight("")
			tree.curNode = tree.curNode.right
			// If number, set data in current node and move up a node
		} else if strings.ContainsAny(char, "1234567890") {
			tree.curNode.setData(char)
			if tree.curNode.parent == nil {
				tree.curNode.parent = &node{parent: nil, left: tree.curNode.parent, right: nil, data: ""}
			}
			tree.curNode = tree.curNode.parent
			fmt.Println("Up")
			// If right paren, move up a node
		} else if char == ")" {
			if tree.curNode.parent == nil {
				tree.curNode.parent = &node{parent: nil, left: tree.curNode.parent, right: nil, data: ""}
			}
			tree.curNode = tree.curNode.parent
			fmt.Println("Up")
		}
	}

	return tree
}

// Insert a new node to the left
func (n *node) insertLeft(d string) {
	fmt.Println("Left")
	if n == nil {
		return
	} else if n.left == nil {
		n.left = &node{parent: n, left: nil, right: nil, data: d}
	}
}

// Insert a new node to the right
func (n *node) insertRight(d string) {
	fmt.Println("Right")
	if n == nil {
		return
	} else if n.right == nil {
		n.right = &node{parent: n, left: nil, right: nil, data: d}
	}
}

// Set the data in the current node
func (n *node) setData(d string) {
	fmt.Println("Center", d)
	if n == nil {
		return
	} else if n.data == "" {
		n.data = d
	}
}
