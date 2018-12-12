package fileutil

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestFilesByDir(t *testing.T) {

	inPkg := new(Package)

	t.Run("should fail because Package is nil", func(t *testing.T) {

		out, err := FilesByDir("", nil)

		require.Equal(t, errPackage, errors.Cause(err))
		require.Nil(t, out)

	})

	t.Run("should fail because path is empty", func(t *testing.T) {

		out, err := FilesByDir("", inPkg)

		require.Equal(t, errPath, errors.Cause(err))
		require.Nil(t, out)

	})

	t.Run("should fail because path is not found", func(t *testing.T) {

		out, err := FilesByDir("somePath", inPkg)

		_, ok := err.(*os.PathError)

		require.True(t, ok)
		require.Nil(t, out)

	})

	t.Run("should return expected result", func(t *testing.T) {

		out, err := FilesByDir("testdata", inPkg)

		require.NoError(t, err)
		require.Equal(
			t,
			&Package{
				name: "testdata",
				subPackages: []*Package{
					&Package{
						name: "somedir",
						goFiles: []GoFile{
							"somefile.go",
							"someotherfile.go",
						},
					},
				},
				goFiles: []GoFile{
					"somefile.go",
				},
			},
			out,
		)

	})

}
