package smt

import (
	"bytes"
	"errors"
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
	Parent() []byte

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
	hash    []byte
	node    []byte
	parent  []byte
	index   int
	pathlen int
}

//nodeIterator stores the current position of the iterator in the trie
type nodeIterator struct {
	trie  *SparseMerkleTree
	stack []*nodeIteratorState
	path  []byte
	err   error
}

//func (it *nodeIterator) resolveBlob(hash []byte, path []byte) ([]byte, error) {
//	if it.resolver != nil {
//		if blob, err := it.resolver.Get(hash); err == nil && len(blob) > 0 {
//			return blob, nil
//		}
//	}
//	// Retrieve the specified node from the underlying node reader.
//	// it.trie.resolveAndTrack is not used since in that function the
//	// loaded blob will be tracked, while it's not required here since
//	// all loaded nodes won't be linked to trie at all and track nodes
//	// may lead to out-of-memory issue.
//	return it.trie.reader.nodeBlob(path, common.BytesToHash(hash))
//}
func (n nodeIterator) Error() error {
	//TODO implement me
	panic("implement me")
}

//Hash returns the hash of nodes iterated up until method is called
func (n *nodeIterator) Hash() []byte {
	if len(n.stack) == 0 {
		return n.trie.th.placeholder()
	}

	return n.trie.th.digest(n.stack[len(n.stack)-1].node)
}

func (n *nodeIterator) Parent() []byte {
	if len(n.stack) == 0 {
		return n.trie.th.placeholder()
	}

	return n.trie.th.digest(n.stack[len(n.stack)-1].parent)
}

func (n *nodeIterator) Path() []byte {
	return n.path
}

func (n nodeIterator) NodeBlob() []byte {
	//if bytes.Compare(n.Hash(), n.trie.th.placeholder()) == 0 {
	//	return nil
	//}
	//blob, err := n.resolveBlob(n.Hash(), it.Path())
	//if err != nil {
	//	it.err = err
	//	return nil
	//}
	//return blob
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

// errIteratorEnd is stored in nodeIterator.err when iteration is done.
var errIteratorEnd = errors.New("end of iteration")

// seekError is stored in nodeIterator.err if the initial seek has failed.
type seekError struct {
	key []byte
	err error
}

func (e seekError) Error() string {
	return "seek error: " + e.err.Error()
}

// Next moves the iterator to the next node, returning whether there are any
// further nodes. In case of an internal error this method returns false and
// sets the Error field to the encountered failure. If `descend` is false,
// skips iterating over any subnodes of the current node.
func (n *nodeIterator) Next(descend bool) bool {
	if n.err == errIteratorEnd {
		return false
	}
	if seek, ok := n.err.(seekError); ok {
		if n.err = n.seek(seek.key); n.err != nil {
			return false
		}
	}
	// Otherwise step forward with the iterator and report any errors.
	state, parentIndex, path, err := n.peek(descend)
	n.err = err
	if n.err != nil {
		return false
	}
	n.push(state, parentIndex, path)
	return true
}

// init initializes the iterator.
func (it *nodeIterator) init() *nodeIteratorState {
	root := it.Hash()
	state := &nodeIteratorState{node: it.trie.root, index: -1}
	if bytes.Compare(root, it.trie.th.placeholder()) != 0 {
		state.hash = root
	}
	return state
}

// peek creates the next state of the iterator.
func (it *nodeIterator) peek(descend bool) (*nodeIteratorState, *int, []byte) {
	// Initialize the iterator if we've just started.
	if len(it.stack) == 0 {
		state := it.init()
		return state, nil, nil
	}
	if !descend {
		// If we're skipping children, pop the current node first
		it.pop()
	}

	// Continue iteration to the next child
	for len(it.stack) > 0 {
		parent := it.stack[len(it.stack)-1]
		ancestor := parent.hash
		if bytes.Compare(ancestor, it.trie.th.placeholder()) == 0 {
			ancestor = parent.parent
		}
		state, path, ok := it.nextChild(parent, ancestor)
		if ok {
			if err := state.resolve(it, path); err != nil {
				return parent, &parent.index, path, err
			}
			return state, &parent.index, path, nil
		}
		// No more child nodes, move back up.
		it.pop()
	}
	return nil, nil, nil, errIteratorEnd
}

func (it *nodeIterator) nextChild(parent *nodeIteratorState, ancestor []byte) (*nodeIteratorState, []byte, bool) {
	switch node := parent.node.(type) {
	case *fullNode:
		// Full node, move to the first non-nil child.
		if child, state, path, index := findChild(node, parent.index+1, it.path, ancestor); child != nil {
			parent.index = index - 1
			return state, path, true
		}
	case *shortNode:
		// Short node, return the pointer singleton child
		if parent.index < 0 {
			hash, _ := node.Val.cache()
			state := &nodeIteratorState{
				hash:    common.BytesToHash(hash),
				node:    node.Val,
				parent:  ancestor,
				index:   -1,
				pathlen: len(it.path),
			}
			path := append(it.path, node.Key...)
			return state, path, true
		}
	}
	return parent, it.path, false
}

//push function pushes the state nodeIteratorState into the stack that is used
//to maintain state data along with nodeIterator
func (n nodeIterator) push(state *nodeIteratorState, parentIndex *int, path []byte) {
	n.path = path
	n.stack = append(n.stack, state)
	if parentIndex != nil {
		*parentIndex++
	}
}

func (n nodeIterator) pop() {
	last := n.stack[len(n.stack)-1]
	n.path = n.path[:last.pathlen]
	n.stack[len(n.stack)-1] = nil
	n.stack = n.stack[:len(n.stack)-1]
}
