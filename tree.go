package tree

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const goExt = ".go"

var (
	errPath = errors.New("path")
	errNode = errors.New("node")
)

// Treer allows third party mocking
type Treer interface {
	Get(string, *Node) (*Node, error)
	Ast(*Leaf) error
}

// Node represents a package (directory) that can contain other sub-packages.
type Node struct {
	Name  string
	Nodes []*Node
	Leafs []Leaf
}

// Leaf represents a go file.
type Leaf struct {
	Name       string
	Path       string
	SyntaxTree *ast.File
}

type result struct {
	node *Node
	err  error
}

// Get returns all go files in the repository by directory name
func Get(path string, node *Node) (*Node, error) {

	if err := validate(path, node); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir("./" + path)
	if err != nil {
		return nil, err
	}

	node.Name = currentPackage(path)

	if len(files) == 0 {
		return node, nil
	}

	var (
		wg        sync.WaitGroup
		resultsCh = make(chan *result)
		done      = make(chan struct{})
		results   []*result
	)

	go func() {
		for v := range resultsCh {
			results = append(results, v)
		}

		done <- struct{}{}
	}()

	nodes, leafs := filterGoFilesDirs(files)

	wg.Add(len(nodes))

	for _, v := range nodes {
		go func(n *Node) {
			defer wg.Done()

			result := &result{}

			result.node, result.err = Get(path+"/"+n.Name, n)

			resultsCh <- result

		}(v)
	}

	wg.Wait()
	close(resultsCh)

	<-done

	for _, n := range results {
		if n.err != nil {
			return nil, n.err
		}
		node.Nodes = append(node.Nodes, n.node)
	}

	for _, v := range leafs {

		v.Path = path + "/" + v.Name

		node.Leafs = append(node.Leafs, v)
	}

	return node, nil

}

// Ast sets leaf's SyntaxTree
func (l *Leaf) Ast() error {

	b, err := ioutil.ReadFile(l.Path)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, l.Name+goExt, string(b), 0)
	if err != nil {
		return err
	}

	l.SyntaxTree = f

	return nil

}

// Print pretty prints go project structure
func (n *Node) Print() {

	fmt.Print(n.Name + "\n|")

	for _, l := range n.Leafs {
		fmt.Println("\t" + l.Name)
	}

	for _, v := range n.Nodes {
		v.Print()
	}

}

func filterGoFilesDirs(files []os.FileInfo) ([]*Node, []Leaf) {

	var (
		leafs []Leaf
		nodes []*Node
	)

	for _, file := range files {

		fileName := file.Name()

		if filepath.Ext(fileName) == goExt {
			leafs = append(leafs, Leaf{Name: fileName})
		} else if file.IsDir() {
			nodes = append(nodes, &Node{Name: fileName})
		}

	}

	return nodes, leafs

}

func currentPackage(path string) string {
	i := strings.IndexByte(path, '/')

	if i == -1 {
		return path
	}

	return path[i+1:]
}

func validate(path string, node *Node) error {

	if node == nil {
		return errors.Wrap(errNode, "node can't be nil")
	}

	if path == "" {
		return errors.Wrap(errPath, "empty path")
	}

	return nil

}
