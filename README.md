# tree [![CircleCI](https://circleci.com/gh/AndreaM16/tree/tree/master.svg?style=svg)](https://circleci.com/gh/AndreaM16/tree/tree/master)
Simple go project to tree structure. It supports only go files.

# ???

![alt text](https://raw.githubusercontent.com/AndreaM16/tree/master/assets/structure.png)

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

**Todo get AST per each file.**
