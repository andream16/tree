package main

import (
	"fmt"

	"github.com/andream16/tree"
)

func main() {

	out, err := tree.Get("examples/example", &tree.Node{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Name)              // example
	fmt.Println(out.Leafs[0])          // somefile.go
	fmt.Println(out.Nodes[0].Name)     // example/subexample
	fmt.Println(out.Nodes[0].Leafs[0]) // someotherfile.go

}
