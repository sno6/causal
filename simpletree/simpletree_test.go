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

func TestTree_Removal(t *testing.T) {
	tree := New[rune](0)
	tree.AddSequence(nil, []rune("Hello World"))

	// Remove the 'H'.
	tree.RemoveNode(ID{
		Timestamp: 2,
		EntityID:  tree.ID.EntityID,
	})

	// Remove ' World'
	for i := 0; i < 6; i++ {
		tree.RemoveNode(ID{
			Timestamp: tree.ID.Timestamp - i,
			EntityID:  tree.ID.EntityID,
		})
	}

	// Add ' there' to the end of the tree.
	tree.AddSequence(&tree.ID, []rune(" there"))

	assert(t, "ello there", treeString(tree))
}

func TestTree_MergeWithNonNullRoots(t *testing.T) {
	t1 := New[rune](0)
	t1.AddSequence(nil, []rune("hi there"))

	t2 := New[rune](1)
	t2.AddSequence(nil, []rune("bye there"))

	t1.Merge(t2)
	t2.Merge(t1)
}

func TestTree_MergeWithDeletes(t *testing.T) {
	t1 := New[rune](0)
	t1.AddSequence(nil, []rune("hi there"))

	t2 := New[rune](1)
	t2.AddSequence(nil, []rune("bye there"))

	// Remove the 'e'
	t2.RemoveNode(ID{
		Timestamp: 10,
		EntityID:  t2.ID.EntityID,
	})

	t1.Merge(t2)
	t2.Merge(t1)

	assert(t, treeString(t1), treeString(t2))
	assert(t, "bye therhi there", treeString(t1))
}

func treeString(tree *Tree[rune]) string {
	var sb strings.Builder
	for _, node := range tree.OrderedNodes(false) {
		if node.Removed {
			panic("wtf")
		}
		sb.WriteRune(node.V)
	}
	return sb.String()
}

func assert[T comparable](t *testing.T, got, wanted T) {
	if got != wanted {
		t.Errorf("Expected %v, Got %v", got, wanted)
	}
}
