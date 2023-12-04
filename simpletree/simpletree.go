// Package simpletree is an educational implementation of a 'causal tree'
// based on the idea from Victor Grishchenko. It is purposely verbose
// with comments, and way too slow for production use.
package simpletree

import (
	"slices"
	"sort"
)

// An ID is a Lamport timestamp and a unique identifier.
//
// The EntityID should be unique across trees, think of it
// like a node (in the networking sense) or user ID. It is
// used to order nodes on tie-breaks between Timestamps.
type ID struct {
	Timestamp int
	EntityID  int
}

// A Node (or Atom in the literature) represents a single operation.
// Nodes are recursive, and belong to a Tree, starting at a root.
// They hold a single generic value V.
type Node[T any] struct {
	ID       ID
	V        T
	Removed  bool
	Parent   *Node[T]
	Children []*Node[T]
}

// A Tree holds Nodes recursively from the root.
//
// The Tree.ID.Timestamp increases monotonically, when a Node
// is appended to a Tree the Timestamp increases by 1. There
// should never be more than one node with the same ID in a tree.
type Tree[T any] struct {
	ID   ID
	Root *Node[T]
}

// New creates a new Tree and initiates it with a root id.
func New[T any](entityID int) *Tree[T] {
	return &Tree[T]{
		ID: ID{
			Timestamp: 1,
			EntityID:  entityID,
		},
		Root: &Node[T]{
			ID: ID{
				Timestamp: 0,
				EntityID:  entityID,
			},
		},
	}
}

func (t *Tree[T]) AddSequence(parentID *ID, seq []T) *ID {
	for _, v := range seq {
		parentID = t.AddNode(parentID, v)
	}
	return parentID
}

// AddNode finds the parent/causing node by parentID
// and then adds a new node as a child of that parent.
func (t *Tree[T]) AddNode(parentID *ID, v T) *ID {
	parent := t.Root
	if parentID != nil {
		parent = t.Find(*parentID)
	}

	// We just add Nodes in the order they come in.
	// To get the correct ordering of Nodes call OrderedNodes.
	parent.Children = append(parent.Children, &Node[T]{
		ID:     t.IncrTimestamp(),
		V:      v,
		Parent: parent,
	})

	return &t.ID
}

func (t *Tree[T]) IncrTimestamp() ID {
	t.ID.Timestamp += 1
	return t.ID
}

func (t *Tree[T]) RemoveNode(id ID) {
	node := t.Find(id)
	if node == nil {
		return
	}
	node.Removed = true
}

// OrderedNodes returns all the tree's nodes in depth-first pre-order.
// A partial ordering is turned into a total order by ordering sibling branches
// by node ID and then EntityID (secondary).
func (t *Tree[T]) OrderedNodes(includeRemoved bool) []*Node[T] {
	var nodes []*Node[T]
	t.traverseFunc(t.Root, func(n *Node[T]) {
		if n.Removed && includeRemoved == false {
			return
		}
		nodes = append(nodes, n)
	})
	return nodes[1:] // Don't return "root" placeholder node.
}

func (t *Tree[T]) traverseFunc(current *Node[T], f func(*Node[T])) {
	f(current)

	// This is where we do the actual ordering of Nodes by timestamps (and tree/node order).
	children := slices.Clone(current.Children)
	sort.Slice(children, func(i, j int) bool {
		return children[i].ID.Less(children[j].ID)
	})

	for _, node := range children {
		t.traverseFunc(node, f)
	}
}

// Merge two trees together by essentially doing a diff and patch.
// After running this function the dst tree (the caller) should have
// all nodes that src has.
func (t *Tree[T]) Merge(src *Tree[T]) {
	// If we've added a Node with a higher Timestamp than we have
	// (or have seen) then we need to keep track of that.
	t.ID.Timestamp = max(src.ID.Timestamp, t.ID.Timestamp)

	for _, n := range src.OrderedNodes(true) {
		if t.Exists(n.ID) {
			// If we have a common node make sure that they both
			// share the same attributes, such as being removed.
			ours := t.Find(n.ID)
			if n.Removed && !ours.Removed {
				t.RemoveNode(n.ID)
			}

			continue
		}

		var parentID ID
		if n.Parent != nil {
			parentID = n.Parent.ID
		}

		// src root has a nil parent.
		if !t.Exists(parentID) {
			panic("parent should exist for Node to be added")
		}

		parent := t.Find(parentID)
		if parent == nil {
			panic("parent should exist for Node to be added")
		}
		parent.Children = append(parent.Children, &Node[T]{
			ID:       n.ID,
			V:        n.V,
			Parent:   parent,
			Removed:  n.Removed,
			Children: nil,
		})
	}
}

func (id ID) IsRoot() bool {
	return id.Timestamp == 0
}

func (t *Tree[T]) Exists(id ID) bool {
	if id.IsRoot() {
		return true
	}
	var exists bool
	t.traverseFunc(t.Root, func(n *Node[T]) {
		if !exists && n.ID.Equals(id) {
			exists = true
		}
	})
	return exists
}

func (t *Tree[T]) Find(id ID) *Node[T] {
	if id.IsRoot() {
		return t.Root
	}
	var node *Node[T]
	t.traverseFunc(t.Root, func(n *Node[T]) {
		if node == nil && n.ID.Equals(id) {
			node = n
		}
	})
	return node
}

func (id ID) Equals(id2 ID) bool {
	return id.Timestamp == id2.Timestamp && id.EntityID == id2.EntityID
}

// Less is used to order sibling trees and thus needs to do the following:
//
//   - If node A has a larger timestamp than node B it returns true, because
//     events that happen later a fuller picture of the world. In other words,
//     subtrees with larger timestamps should be traversed first.
//
//   - If nodes A and B have the same timestamp, we want to order Entity's with
//     higher id's first. Why? because that's the arbitrary rule I made up, that's why.
func (id ID) Less(id2 ID) bool {
	if id.Timestamp == id2.Timestamp {
		return id.EntityID > id2.EntityID
	}
	return id.Timestamp > id2.Timestamp
}
