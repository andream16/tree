package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var errReadingFilesInDir = errors.New("read_files_in_dir")

// GoFilesDirs contains useful files
type GoFilesDirs struct {
	goFiles []string
	dirs    []string
}

// FilesByDir returns all go files in the repository by directory name
func FilesByDir(dir string, dirGoFiles map[string]*GoFilesDirs) (map[string]*GoFilesDirs, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(errReadingFilesInDir, "unable to read files in %s directory", dir)
	}

	if len(files) > 0 {
		goFilesDirs := filterGoFilesDirs(files)
		if len(goFilesDirs.goFiles) > 0 || len(goFilesDirs.dirs) > 0 {
			dirGoFiles[dir] = goFilesDirs
		}
		if len(goFilesDirs.dirs) > 0 {
			for _, dirName := range goFilesDirs.dirs {
				return FilesByDir(dir+"/"+dirName, dirGoFiles)
			}
		}
	}

	return dirGoFiles, nil

}

func filterGoFilesDirs(files []os.FileInfo) *GoFilesDirs {

	const goExt = ".go"

	var (
		goFiles []string
		dirs    []string
	)

	for _, file := range files {

		fileName := file.Name()

		if filepath.Ext(fileName) == goExt {
			goFiles = append(goFiles, fileName)
		} else if file.IsDir() {
			dirs = append(dirs, fileName)
		}

	}

	return &GoFilesDirs{
		goFiles: goFiles,
		dirs:    dirs,
	}

}
