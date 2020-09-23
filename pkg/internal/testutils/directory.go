package testutils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTempDirectory(t *testing.T) (dir string, cleanup func()) {
	dir, err := ioutil.TempDir("", "noachis")
	require.NoError(t, err, "should not have error creating temp directory")

	return dir, func() {
		require.NoError(t, os.RemoveAll(dir), "should not have error cleaning temp directory")
	}
}
