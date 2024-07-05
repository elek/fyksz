package k8s

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func TestEnv(t *testing.T) {

	WithInputOutput(t, "testdata/deployment.yaml", "env_test.expected", func(input io.Reader, output io.Writer) error {
		e := Env{
			Key:    "foo",
			Value:  "bar",
			input:  input,
			output: output,
		}
		return e.Run()
	})

}

func WithInputOutput(t *testing.T, input string, expected string, f func(input io.Reader, output io.Writer) error) {
	source, err := os.ReadFile(input)
	require.NoError(t, err)
	result := bytes.NewBuffer([]byte{})
	err = f(bytes.NewReader(source), result)
	require.NoError(t, err)
	exp, err := os.ReadFile(expected)
	require.NoError(t, err)
	fmt.Println(string(result.Bytes()))
	require.Equal(t, strings.TrimSpace(string(exp)), strings.TrimSpace(string(result.Bytes())))
}
