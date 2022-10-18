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

//nodeIteratorState stores the data pertinent to the iteration state
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

func (n *nodeIterator) iterate([][]byte, []byte) [][]byte {
	// Get tree's root
	root := n.trie.Root()

	if bytes.Equal(root, n.trie.th.placeholder()) {
		// The tree is empty, return the default value.
		return nil
	}

	currentHash := root
	var key [][]byte
	var value [][]byte
	var tempKey [][]byte
	var tempValue [][]byte
	currentData, _ := n.trie.nodes.Get(currentHash)
	tempValue = append(tempValue, currentData)
	tempKey = append(tempKey, currentHash)
	for len(tempKey) > 0 {

		if err != nil {
			return nil, err
		} else if n.trie.th.isLeaf(currentData) {
			// We've reached the end. Is this the actual leaf?
			p, _ := n.trie.th.parseLeaf(currentData)
			if !bytes.Equal(path, p) {
				// Nope. Therefore the key is actually empty.
				return defaultValue, nil
			}
			// Otherwise, yes. Return the value.
			value, err := smt.values.Get(path)
			if err != nil {
				return nil, err
			}
			return value, nil
		}

		leftNode, rightNode := smt.th.parseNode(currentData)
		if getBitAtFromMSB(path, i) == right {
			currentHash = rightNode
		} else {
			currentHash = leftNode
		}

		if bytes.Equal(currentHash, smt.th.placeholder()) {
			// We've hit a placeholder value; this is the end.
			return defaultValue, nil
		}
	}
	//if n.trie.th.isLeaf(currentData) {
	//	// We've reached the end. Is this the actual leaf?
	//	kvpair = append(kvpair, tempkvpair[len(tempkvpair)-1])
	//	return kvpair
	//} else{
	//	leftNode, rightNode := n.trie.th.parseNode(currentData)
	//	n.iterate(tempkvpair,leftNode)
	//	n.iterate(tempkvpair,rightNode)
	//
	//}

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

//Parent returns the hash value of the parent node
func (n *nodeIterator) Parent() []byte {
	if len(n.stack) == 0 {
		return n.trie.th.placeholder()
	}

	return n.trie.th.digest(n.stack[len(n.stack)-1].parent)
}

//Path returns the path to the nodeIterator
func (n *nodeIterator) Path() []byte {
	return n.path
}

func (n nodeIterator) NodeBlob() []byte {
	if bytes.Compare(n.Hash(), n.trie.th.placeholder()) == 0 {
		return nil
	}
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

//Next moves the iterator to the next node, returning whether there are any
//further nodes. In case of an internal error this method returns false and
//sets the Error field to the encountered failure. If `descend` is false,
//skips iterating over any subnodes of the current node.

//TO DO AFTER PEEK
//
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
//	state, parentIndex, path, err := n.peek(descend)
//	n.err = err
//	if n.err != nil {
//		return false
//	}
//	n.push(state, parentIndex, path)
//	return true
//}

// init initializes the iterator.
func (n *nodeIterator) init() *nodeIteratorState {
	root := n.Hash()
	state := &nodeIteratorState{node: n.trie.root, index: -1}
	if bytes.Compare(root, n.trie.th.placeholder()) != 0 {
		state.hash = root
	}
	return state
}

//// peek creates the next state of the iterator.
//func (n *nodeIterator) peek(descend bool) (*nodeIteratorState, *int, []byte) {
//	if len(n.stack) == 0 {
//		baseIt := n.init()
//		return baseIt, nil, nil
//	}
//	//descend is used to skip the other siblings of the current node
//	if !descend {
//		n.pop()
//	}
//	//Continue iteration to next available nodes
//	if len(n.stack) > 0 {
//		parent := n.stack[len(n.stack)-1]
//		grandParent := parent.hash
//		//Checks if parent node has any more parents
//		if bytes.Compare(grandParent, n.trie.th.placeholder()) == 0 {
//			grandParent = parent.parent
//		}
//		//TO DO NEXT CHILD
//	}
//}
//
////nextChild obtains the nodIterator for the next node in the trie
//func (n *nodeIterator) nextChild(parent *nodeIteratorState, ancestor []byte) (*nodeIteratorState, []byte, bool) {
//	sidenodes, pathnodes, leafBool, values, err := n.trie.sideNodesForRoot(, n.trie.root, true)
//	return parent, values, false
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

//pop deletes the most recent value in the stack
func (n nodeIterator) pop() {
	last := n.stack[len(n.stack)-1]
	n.path = n.path[:last.pathlen]
	n.stack[len(n.stack)-1] = nil
	n.stack = n.stack[:len(n.stack)-1]
}
