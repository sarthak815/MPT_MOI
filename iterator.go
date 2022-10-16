package smt

import (
	"bytes"
	"errors"
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

//Prove generates Merkle proof of the leaf positioned on
func (it *Iterator) Prove() [][]byte {
	return it.nodeIt.LeafProof()
}

//NodeIterator is an iterator that traverses the trie pre-order
type NodeIterator interface {
	//Next moves the iterator to the next node. The boolean value is false
	//any other child nodes present will be skipped
	Next(bool) bool

	//Error returns the error status of iterator
	Error() error

	//Hash returns the hash of the current node
	Hash() []byte

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

func (n nodeIterator) Error() error {
	//TODO implement me
	panic("implement me")
}

//Hash returns the hash of nodes iterated up until method is called
func (n nodeIterator) Hash() []byte {
	if len(n.stack) == 0 {
		return treeHasher.placeholder()
	}

	return n.trie.th.hasher.Sum(n.stack[len(n.stack)-1].node)
}

func (n nodeIterator) Parent() hash.Hash {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) Path() []byte {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) NodeBlob() MapStore {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) Leaf() bool {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) LeafKey() []byte {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) LeafBlob() []byte {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) LeafProof() [][]byte {
	//TODO implement me
	panic("implement me")
}

func (n nodeIterator) seek(start []byte) error {

}

// errIteratorEnd is stored in nodeIterator.err when iteration is done.
var errIteratorEnd = errors.New("end of iteration")

// seekError is stored in nodeIterator.err if the initial seek has failed.
type seekError struct {
	key []byte
	err error
}

//newNodeIterator creates a new iterator for the sparse merkle tree provided
func newNodeIterator(smt *SparseMerkleTree, start []byte) NodeIterator {
	if bytes.Compare(smt.root, smt.th.placeholder()) == 0 {
		return &nodeIterator{
			trie: smt,
			err:  errIteratorEnd,
		}
	}
	n := &nodeIterator{trie: smt}
	n.err = n.seek(start)
	return n
}

// Next moves the iterator to the next node, returning whether there are any
// further nodes. In case of an internal error this method returns false and
// sets the Error field to the encountered failure. If `descend` is false,
// skips iterating over any subnodes of the current node.
//func (n *nodeIterator) Next(descend bool) bool {
//	if n.err == errIteratorEnd {
//		return false
//	}
//	if seek, ok := n.err.(seekError); ok {
//		if n.err = n.seek(seek.key); n.err != nil {
//			return false
//		}
//	}
//	// Otherwise step forward with the iterator and report any errors.
//	state, parentIndex, path, err := it.peek(descend)
//	n.err = err
//	if n.err != nil {
//		return false
//	}
//	n.push(state, parentIndex, path)
//	return true
//}
//push function pushes the state nodeIteratorState into the stack that is used
//to maintain state data along with nodeIterator
func (n nodeIterator) push(state *nodeIteratorState, parentIndex *int, path []byte) {
	n.path = path
	n.stack = append(n.stack, state)
	if parentIndex != nil {
		*parentIndex++
	}
}
