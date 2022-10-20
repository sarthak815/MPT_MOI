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
	//fmt.Println("Root is ", root)
	if bytes.Equal(root, n.Trie.th.placeholder()) {
		// The tree is empty, return the default value.
		fmt.Println("Tree is empty")
		return nil
	}

	currentHash := root
	var key [][]byte
	var tempKey [][]byte

	tempKey = append(tempKey, currentHash)
	for len(tempKey) > 0 {
		//fmt.Println(len(tempKey))
		var poppedKey []byte
		poppedKey, tempKey = pop(tempKey)
		//fmt.Println("Popped Key: ", poppedKey)
		//fmt.Println("Len after popping: ", len(tempKey))
		key = append(key, poppedKey)
		currentData, _ := n.Trie.nodes.Get(poppedKey)
		if n.Trie.th.isLeaf(currentData) {
			//fmt.Println("This is a leaf node")
			continue
		} else {
			//fmt.Println("Not a leaf")
			leftNode, rightNode := n.Trie.th.parseNode(currentData)
			if !bytes.Equal(leftNode, n.Trie.th.placeholder()) && !bytes.Equal(leftNode, n.Trie.th.nullLeaf()) {
				tempKey = append(tempKey, leftNode)
			}
			if !bytes.Equal(rightNode, n.Trie.th.placeholder()) && !bytes.Equal(rightNode, n.Trie.th.nullLeaf()) {
				tempKey = append(tempKey, rightNode)
			}

			//fmt.Println("Left Key: ", leftNode, "\nRight Key: ", rightNode)
		}

	}
	return key
}

//PrintKeys() prints all the values stored of each key
func (n *NodeIteratorSMT) PrintKeys(key [][]byte) {
	fmt.Println("Number of keys: ", len(key))
	for i := 0; i < len(key); i++ {
		currentData, _ := n.Trie.nodes.Get(key[i])
		fmt.Println(currentData)
	}
}
