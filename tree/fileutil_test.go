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
							"somefile.go",
							"someotherfile.go",
						},
					},
				},
				Leafs: []Leaf{
					"somefile.go",
				},
			},
			out,
		)

	})

}
