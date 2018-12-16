# tree
Simple go project to tree structure. It supports only go files.

# ???

Given:

![alt text](https://raw.githubusercontent.com/AndreaM16/tree/master/assets/structure.png)

By Running:

```
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
```
