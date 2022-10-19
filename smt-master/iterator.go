package smt

import (
	"bytes"
	"fmt"
)

//NodeIteratorSMT stores the trie to be iterated
type NodeIteratorSMT struct {
	Trie *SparseMerkleTree
}

//pop() is uses to make the byte array act like a stack
func pop(stack [][]byte) ([]byte, [][]byte) {
	return stack[len(stack)-1], stack[:len(stack)-1]
}

//Iterate() iterates over all the key-value pairs present
func (n *NodeIteratorSMT) Iterate() [][]byte {
	// Get tree's root
	root := n.Trie.Root()
	fmt.Println("Root is ", root)
	if bytes.Equal(root, n.Trie.th.placeholder()) {
		// The tree is empty, return the default value.
		fmt.Println("Tree is empty")
		return nil
	}

	currentHash := root
	var key [][]byte
	var tempKey [][]byte
	currentData, _ := n.Trie.nodes.Get(currentHash)
	fmt.Println("Data is:", string(currentData))
	tempKey = append(tempKey, currentHash)
	for len(tempKey) > 0 {
		fmt.Println(len(tempKey))
		var poppedKey []byte
		poppedKey, tempKey = pop(tempKey)
		key = append(key, poppedKey)
		leftNode, rightNode := n.Trie.th.parseNode(currentData)
		if n.Trie.th.isLeaf(currentData) {
			continue
		} else {
			tempKey = append(tempKey, rightNode)
			tempKey = append(tempKey, leftNode)
		}

	}
	return key
}

//PrintKeys() prints all the values stored of each key
func (n *NodeIteratorSMT) PrintKeys(key [][]byte) {
	for i := 0; i < len(key); i++ {
		currentData, _ := n.Trie.nodes.Get(key[i])
		fmt.Println(string(currentData))
	}
}
