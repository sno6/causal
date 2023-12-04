//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"github.com/sno6/causal/simpletree"
	"sync"
	"syscall/js"
	"time"
)

var mu sync.Mutex
var clientTreeMap = make(map[string]*simpletree.Tree[rune])

func registerGlobalFunctions() {
	js.Global().Set("addClient", js.FuncOf(addClient))
	js.Global().Set("onAdd", js.FuncOf(onAdd))
	js.Global().Set("onRemove", js.FuncOf(onRemove))
	js.Global().Set("getNodes", js.FuncOf(getNodes))

}

func getNodes(_ js.Value, inputs []js.Value) any {
	clientID := inputs[0].Int()
	tree := clientTreeMap[mapKey(clientID)]
	if tree == nil {
		return "[]"
	}
	return orderedNodeList(tree)
}

func addClient(_ js.Value, inputs []js.Value) any {
	clientID := inputs[0].Int()

	mu.Lock()
	clientTreeMap[mapKey(clientID)] = simpletree.New[rune](clientID)
	mu.Unlock()

	return nil
}

func onAdd(_ js.Value, inputs []js.Value) any {
	var parentID *simpletree.ID
	if !inputs[0].IsNull() && !inputs[0].IsUndefined() {
		parentID = &simpletree.ID{
			Timestamp: inputs[0].Get("Timestamp").Int(),
			EntityID:  inputs[0].Get("EntityID").Int(),
		}
	}

	newVal := inputs[1].String()
	clientID := inputs[2].Int()

	tree := clientTreeMap[mapKey(clientID)]
	tree.AddSequence(parentID, []rune(newVal))
	return nil
}

func onRemove(_ js.Value, inputs []js.Value) any {
	nodeID := simpletree.ID{
		Timestamp: inputs[0].Get("Timestamp").Int(),
		EntityID:  inputs[0].Get("EntityID").Int(),
	}

	clientID := inputs[1].Int()

	tree := clientTreeMap[mapKey(clientID)]
	tree.RemoveNode(nodeID)
	return nil
}

func mapKey(id int) string {
	return fmt.Sprintf("tree-%d", id)
}

func orderedNodeList(t *simpletree.Tree[rune]) string {
	type nodeItem struct {
		ID       simpletree.ID `json:"id"`
		ParentID simpletree.ID `json:"parent_id"`
		Value    string        `json:"value"`
		Removed  bool          `json:"removed"`
	}

	var nodes []nodeItem
	for _, n := range t.OrderedNodes(true) {
		item := nodeItem{
			ID:      n.ID,
			Value:   string(n.V),
			Removed: n.Removed,
		}
		if n.Parent != nil {
			item.ParentID = n.Parent.ID
		}
		nodes = append(nodes, item)
	}

	b, _ := json.Marshal(nodes)
	return string(b)
}

func main() {
	registerGlobalFunctions()

	// Merge all trees 5 times a second.
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if clientTreeMap != nil {
				for k := range clientTreeMap {
					for k2 := range clientTreeMap {
						if k != k2 {
							clientTreeMap[k].Merge(clientTreeMap[k2])
						}
					}
				}
			}
		}
	}()

	// Block so the WASM module doesn't exit early.
	select {}
}
