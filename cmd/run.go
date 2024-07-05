package cmd

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Run struct {
}

func (s Run) Run() error {
	raw, err := os.ReadFile("build.fyksz")
	if err != nil {
		return errors.WithStack(err)
	}
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}
	output, err := RunString(string(raw), in)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(output)
	return nil
}

func RunString(def string, in []byte) (output string, err error) {
	input := in
	for _, line := range strings.Split(def, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = resolveAbbrev(line)
		cmd := exec.Command("/usr/bin/bash", "-c", line)
		cmd.Stdin = bytes.NewReader(input)
		cmd.Stderr = os.Stderr
		result := bytes.NewBuffer([]byte{})
		cmd.Stdout = result
		_, _ = os.Stderr.WriteString("Executing " + line + "\n")
		err = cmd.Run()
		if err != nil {
			return string(input), errors.WithStack(err)
		}
		input = result.Bytes()

	}
	return string(input), nil
}

func resolveAbbrev(line string) string {
	if strings.HasPrefix(line, "@") {
		return "fyksz " + line[1:]
	}
	if strings.Contains(line, "@@") {
		parts := strings.SplitN(line, "@@", 2)
		return fmt.Sprintf("fyksz apply --filter \"%s\" %s", strings.ReplaceAll(parts[0], "\"", "\\\""), parts[1])
	}
	return line
}
