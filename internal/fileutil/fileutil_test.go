package fileutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilesByDir(t *testing.T) {

	filesByDir, err := FilesByDir("testdata", map[string]*GoFilesDirs{})

	require.NoError(t, err)
	require.NotEmpty(t, filesByDir)

}
