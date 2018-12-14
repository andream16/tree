package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	errPath    = errors.New("path")
	errPackage = errors.New("package")
)

// GoFilesDirs contains useful files
type GoFilesDirs struct {
	goFiles []string
	dirs    []string
}

// Package represents a package (directory) that can contain other sub-packages.
// Can be intended as a node.
type Package struct {
	Name        string
	SubPackages []*Package
	GoFiles     []GoFile
}

// GoFile represents a go file.GoFile
// Can be intended as a leaf.
type GoFile string

// FilesByDir returns all go files in the repository by directory name
func FilesByDir(path string, pkg *Package) (*Package, error) {

	if err := validate(path, pkg); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir("./" + path)
	if err != nil {
		return nil, err
	}

	pkg.Name = pathPostfix(path)

	if len(files) > 0 {

		subPackages, goFiles := filterGoFilesDirs(files)
		for _, p := range subPackages {
			sPkg, err := FilesByDir(path+"/"+p.Name, p)
			if err != nil {
				return nil, err
			}
			pkg.SubPackages = append(pkg.SubPackages, sPkg)
		}

		pkg.GoFiles = goFiles

	}

	return pkg, nil

}

func filterGoFilesDirs(files []os.FileInfo) ([]*Package, []GoFile) {

	const goExt = ".go"

	var (
		goFiles []GoFile
		pkgs    []*Package
	)

	for _, file := range files {

		fileName := file.Name()

		if filepath.Ext(fileName) == goExt {
			goFiles = append(goFiles, GoFile(fileName))
		} else if file.IsDir() {
			pkgs = append(pkgs, &Package{Name: fileName})
		}

	}

	return pkgs, goFiles

}

func pathPostfix(path string) string {
	s := strings.Split(path, "/")
	return s[len(s)-1]
}

func validate(path string, pkg *Package) error {

	if pkg == nil {
		return errors.Wrap(errPackage, "initial Package can't be nil")
	}

	if path == "" {
		return errors.Wrap(errPath, "empty path")
	}

	return nil

}
