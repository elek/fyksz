package k8s

import (
	"io"
	"testing"
)

func TestConfigHash(t *testing.T) {
	WithInputOutput(t, "confighash_test.input.yaml", "confighash_test.expected", func(input io.Reader, output io.Writer) error {
		r := ConfigHash{
			input:  input,
			output: output,
		}
		return r.Run()
	})
}
