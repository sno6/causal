package simpletree

import (
	"strings"
	"testing"
)

// Test that OrderedNodes returns nodes in a depth-first pre-order.
// Example from https://en.wikipedia.org/wiki/Tree_traversal#Depth-first_search
func TestTree_OrderedNodes(t *testing.T) {
	tree := New[rune](0)
	f := tree.AddNode(nil, 'F') // Root -> F
	b := tree.AddNode(f, 'B')   // F -> B
	tree.AddNode(b, 'A')        // B -> A
	d := tree.AddNode(b, 'D')   // B -> D
	tree.AddNode(d, 'C')        // D -> C
	tree.AddNode(d, 'E')        // D -> E
	g := tree.AddNode(f, 'G')   // F -> G
	i := tree.AddNode(g, 'I')   // G -> I
	tree.AddNode(i, 'H')        // I -> H

	assert(t, treeString(tree), "FBADCEGIH")
}

func TestTree_Kleppmann(t *testing.T) {
	tree1 := New[rune](0)
	tree2 := New[rune](1)

	tree1.AddSequence(nil, []rune("Hello!"))
	tree2.Merge(tree1)

	assert(t, treeString(tree1), "Hello!")
	assert(t, treeString(tree2), "Hello!")

	// The id of the '!' character.
	id := ID{
		Timestamp: tree1.ID.Timestamp - 1,
		EntityID:  tree1.ID.EntityID,
	}

	// The first tree adds 'World' before '!'.
	tree1.AddSequence(&id, []rune(" World"))

	// The second tree adds ':-)' after '!'.
	tree2.AddSequence(&ID{
		Timestamp: id.Timestamp + 1,
		EntityID:  id.EntityID,
	}, []rune(":-)"))

	assert(t, treeString(tree1), "Hello World!")
	assert(t, treeString(tree2), "Hello!:-)")

	tree1.Merge(tree2)
	tree2.Merge(tree1)

	assert(t, treeString(tree1), "Hello World!:-)")
	assert(t, treeString(tree2), "Hello World!:-)")
}

func TestTree_SimpleConflict(t *testing.T) {
	// Create a tree and append 'Hi' to it.
	tree1 := New[rune](0)
	tree1.AddSequence(nil, []rune("Hi "))

	// Create another tree and merge.
	tree2 := New[rune](1)
	tree2.Merge(tree1)

	// Assume that at this point both trees go offline.
	// This is the ID of the last write to the tree before going offline.
	discrepID := tree1.ID

	tree1.AddSequence(&tree1.ID, []rune("World"))
	tree2.AddSequence(&discrepID, []rune("Coders"))

	assert(t, treeString(tree1), "Hi World")
	assert(t, treeString(tree2), "Hi Coders")

	tree1.Merge(tree2)
	tree2.Merge(tree1)

	// Hey, we didn't say it would make sense... it just needs to be consistent.
	assert(t, treeString(tree1), "Hi CodersWorld")
	assert(t, treeString(tree2), "Hi CodersWorld")
}

func treeString(tree *Tree[rune]) string {
	var sb strings.Builder
	for _, node := range tree.OrderedNodes() {
		sb.WriteRune(node.V)
	}
	return sb.String()
}

func assert[T comparable](t *testing.T, got, wanted T) {
	if got != wanted {
		t.Errorf("Expected %v, Got %v", got, wanted)
	}
}
