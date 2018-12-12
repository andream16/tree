package main

import (
	"github.com/andream16/adviser/internal/fileutil"
)

func main() {

	out, err := fileutil.FilesByDir("example", new(fileutil.Package))
	if err != nil {
		panic(err)
	}

	fmt.Println(out.name) // example
	fmt.Println(out.goFiles[0]) // somefile.go
	fmt.Println(out.subPackages[0].name) // subexample
	fmt.Println(out.subPackages[0].goFiles[0]) // someotherfile.go

}