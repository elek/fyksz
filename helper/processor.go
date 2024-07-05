package helper

import (
	"fmt"
	"fyksz/data"
	"fyksz/yaml"
	"github.com/pkg/errors"
	"io"
	"strings"
)

func ProcessAll(in io.Reader, out io.Writer, f func(input string) (string, error)) error {
	raw, err := io.ReadAll(in)
	if err != nil {
		return errors.WithStack(err)
	}
	var output []string
	for _, part := range strings.Split(string(raw), "---\n") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		result, err := f(part)
		if err != nil {
			return errors.WithStack(err)
		}
		output = append(output, result)
	}
	_, err = fmt.Fprintln(out, strings.Join(output, "\n---\n"))
	return err
}

func ProcessNode(in io.Reader, out io.Writer, f func(node data.Node) error) error {
	return ProcessAll(in, out, func(input string) (string, error) {
		parsed := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(input), &parsed)
		if err != nil {
			return "", errors.WithStack(err)
		}

		root, err := data.ConvertToNode(parsed, data.NewPath())
		if err != nil {
			return "", errors.WithStack(err)
		}
		err = f(root)
		if err != nil {
			return "", errors.WithStack(err)
		}
		outYaml := data.ConvertToYaml(root)
		outRaw, err := yaml.Marshal(outYaml)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return string(outRaw), nil
	})
}
