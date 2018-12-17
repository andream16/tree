package tree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

var (
	errPath = errors.New("path")
	errNode = errors.New("node")
)

// Getter allows third party mocking
type Getter interface {
	Get(string, *Node) (*Node, error)
}

// Node represents a package (directory) that can contain other sub-packages.
type Node struct {
	Name  string
	Nodes []*Node
	Leafs []Leaf
}

// Leaf represents a go file.
type Leaf string

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

	node.Leafs = leafs

	return node, nil

}

func filterGoFilesDirs(files []os.FileInfo) ([]*Node, []Leaf) {

	const goExt = ".go"

	var (
		leafs []Leaf
		nodes []*Node
	)

	for _, file := range files {

		fileName := file.Name()

		if filepath.Ext(fileName) == goExt {
			leafs = append(leafs, Leaf(fileName))
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
