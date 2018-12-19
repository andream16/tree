package main

import (
	"fmt"
	"go/ast"

	"github.com/andream16/tree"
)

func main() {

	out, err := tree.Get("examples/example", &tree.Node{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Name)                   // example
	fmt.Println(out.Leafs[0].Name)          // somefile.go
	fmt.Println(out.Leafs[0].Path)          // examples/example/somefile.go
	fmt.Println(out.Nodes[0].Name)          // example/subexample
	fmt.Println(out.Nodes[0].Leafs[0].Name) // someotherfile.go
	fmt.Println(out.Nodes[0].Leafs[0].Path) // examples/example/subexample/someotherfile.go

	err = out.Leafs[0].Ast()
	if err != nil {
		panic(err)
	}

	ast.Inspect(out.Leafs[0].SyntaxTree, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident:
			s = x.Name
		}
		if s != "" {
			fmt.Printf("%s\n", s) // example
		}
		return true
	})

}
