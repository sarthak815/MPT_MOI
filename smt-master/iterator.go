package smt

import (
	"hash"
)

//Iterator is a key-value trie that traverses a Trie
type Iterator struct {
	nodeIt NodeIterator
	node   []byte
	value  []byte
	root   []byte
	Err    error
}

//NewIterator creates a key-value iterator from a node iterator
func NewIterator(it NodeIterator) *Iterator {
	return &Iterator{
		nodeIt: it,
	}
}

//Next moves the iterator forward one key-value entry
func (it *Iterator) Next() bool {
	for it.nodeIt.Next(true) {
		if it.nodeIt.Leaf() {
			it.root = it.nodeIt.LeafKey()
			it.value = it.nodeIt.LeafBlob()

		}
	}
}

//NodeIterator is an iterator that traverses the trie pre-order
type NodeIterator interface {
	//Next moves the iterator to the next node. The boolean value is false
	//any other child nodes present will be skipped
	Next(bool)

	//Error returns the error status of iterator
	Error() error

	//Hash returns the hash of the current node
	Hash() hash.Hash

	//Parent returns the hash of the next available parent node
	Parent() hash.Hash

	//Path returns the hex-encoded path to the current node
	//For leaf elements the last element of the path is-------
	Path() []byte

	//NodeBlob returns the rlp-encoded value of current iterated node
	NodeBlob() MapStore

	//Leaf returns true iff current node is a leaf
	Leaf() bool

	//LeafKey returns the key of the leaf.
	LeafKey() []byte

	//LeafBlob returns the content of the leaf.
	LeafBlob() []byte

	//LeafProof returns the Merkle proof of the leaf
	LeafProof() [][]byte
}

//nodeIteratorState stores the data pertinent to the iteration state at one
//at one particular node in the trie
type nodeIteratorState struct {
	hasher  hash.Hash
	node    []byte
	parent  hash.Hash
	index   int
	pathlen int
}

//NodeIterator stores the current position of the iterator in the trie
type nodeIterator struct {
	trie  *SparseMerkleTree
	stack []*nodeIteratorState
	path  []byte
	err   error
}
