package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func registerCallbacks() {
	js.Global().Set("handleSend", js.FuncOf(handleSend))
}

func handleSend(v js.Value, inputs []js.Value) any {
	println("Button clicked!")
	b, _ := json.Marshal(inputs)
	fmt.Println(string(b))

	for _, v := range inputs {
		fmt.Println("got input", v) // works
	}
	return nil
}

func main() {
	registerCallbacks()
	select {}
}
