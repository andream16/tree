package tree

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {

	inPkg := new(Node)

	t.Run("should fail because Node is nil", func(t *testing.T) {

		out, err := Get("", nil)

		require.Equal(t, errNode, errors.Cause(err))
		require.Nil(t, out)

	})

	t.Run("should fail because path is empty", func(t *testing.T) {

		out, err := Get("", inPkg)

		require.Equal(t, errPath, errors.Cause(err))
		require.Nil(t, out)

	})

	t.Run("should fail because path is not found", func(t *testing.T) {

		out, err := Get("somePath", inPkg)

		_, ok := err.(*os.PathError)

		require.True(t, ok)
		require.Nil(t, out)

	})

	t.Run("should return expected result", func(t *testing.T) {

		out, err := Get("testdata", inPkg)

		require.NoError(t, err)
		require.Equal(
			t,
			&Node{
				Name: "testdata",
				Nodes: []*Node{
					&Node{
						Name: "somedir",
						Leafs: []Leaf{
							{
								Name: "somefile.go",
								Path: "testdata/somedir/somefile.go",
							}, {
								Path: "testdata/somedir/someotherfile.go",
								Name: "someotherfile.go",
							},
						},
					},
				},
				Leafs: []Leaf{
					{
						Name: "somefile.go",
						Path: "testdata/somefile.go",
					},
				},
			},
			out,
		)

	})

}

func TestLeaf_Tree(t *testing.T) {

	t.Run("should fail because file hasn't been found", func(t *testing.T) {

		l := &Leaf{}

		err := l.Ast()

		require.True(t, os.IsNotExist(err))

	})

	t.Run("should return early because file is empty", func(t *testing.T) {

		l := &Leaf{
			Path: "testdata/somedir/somefile.go",
		}

		err := l.Ast()

		require.NoError(t, err)

	})

	t.Run("should return a valid syntax tree", func(t *testing.T) {

		l := &Leaf{
			Path: "testdata/somefile.go",
		}

		err := l.Ast()

		require.NoError(t, err)

	})

}
