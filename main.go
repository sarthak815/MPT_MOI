package main

import (
	"MPT_MOI/smt-master"
	"crypto/sha256"
)

func main() {
	// Initialise two new key-value store to store the nodes and values of the tree
	nodeStore := smt.NewSimpleMap()
	valueStore := smt.NewSimpleMap()
	// Initialise the tree
	tree := smt.NewSparseMerkleTree(nodeStore, valueStore, sha256.New())
	//fmt.Println(tree.Root())
	// Update the key "foo" with the value "bar"
	_, _ = tree.Update([]byte("foo"), []byte("bar"))
	_, _ = tree.Update([]byte("tom"), []byte("lin"))
	_, _ = tree.Update([]byte("sam"), []byte("fam"))

	// Generate a Merkle proof for foo=bar
	//proof, _ := tree.Prove([]byte("foo"))
	//root := tree.Root() // We also need the current tree root for the proof\
	//Generating a NodeIterator
	iterator := smt.NodeIteratorSMT{tree}
	//iterator traverses the tree and obtains all the hash values of the keys in pre order
	keys := iterator.Iterate()
	//iterator prints the hashed values of the keys
	iterator.PrintKeys(keys)
	// Verify the Merkle proof for foo=bar
	//if smt.VerifyProof(proof, root, []byte("foo"), []byte("bar"), sha256.New()) {
	//	fmt.Println("Proof verification succeeded.")
	//} else {
	//	fmt.Println("Proof verification failed.")
	//}
}
