# tree [![CircleCI](https://circleci.com/gh/AndreaM16/tree/tree/master.svg?style=svg)](https://circleci.com/gh/AndreaM16/tree/tree/master) [![GoDoc](https://godoc.org/github.com/AndreaM16/tree?status.svg)](https://godoc.org/github.com/AndreaM16/tree) [![Go Report Card](https://goreportcard.com/badge/github.com/AndreaM16/tree)](https://goreportcard.com/report/github.com/AndreaM16/tree) [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/andream16/tree/master/LICENSE)

Simple go project to tree structure. It supports only go files.

# ???

![alt text](https://raw.githubusercontent.com/AndreaM16/tree/master/assets/structure.png)

```
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

	// example
	// |       somefile.go
	// example/subexample
	// |       someotherfile.go
	out.Print()

}
```
