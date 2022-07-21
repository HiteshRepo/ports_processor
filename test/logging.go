package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func CaptureStandardOut(t *testing.T, f func()) string {
	originalStdOut := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = writer

	f()

	err = writer.Close()
	require.NoError(t, err)

	out, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	os.Stdout = originalStdOut

	return string(out)
}
