package main

import (
	"fmt"

	"github.com/andream16/adviser/internal/fileutil"
)

func main() {

	out, err := fileutil.FilesByDir("examples/example", &fileutil.Package{})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Name)                      // example
	fmt.Println(out.GoFiles[0])                // somefile.go
	fmt.Println(out.SubPackages[0].Name)       // subexample
	fmt.Println(out.SubPackages[0].GoFiles[0]) // someotherfile.go

}
