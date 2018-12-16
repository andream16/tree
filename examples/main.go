package main

import (
	"fmt"

	"github.com/andream16/adviser/tree"
)

func main() {

	out, err := tree.Get("examples/example", &tree.Package{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Name)                    // example
	fmt.Println(out.Files[0])                // somefile.go
	fmt.Println(out.SubPackages[0].Name)     // subexample
	fmt.Println(out.SubPackages[0].Files[0]) // someotherfile.go

}
