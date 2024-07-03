package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Save struct {
	SplitLine string `default:"---"`
	Name      string `usage:"Executable script to determine the name of the fragment (given on the stdin)"`
}

func (s Save) Run() error {
	var input []byte
	var err error

	if os.Getenv("FYKSZ_DEBUG_INPUT") != "" {
		input, err = os.ReadFile(os.Getenv("FYKSZ_DEBUG_INPUT"))
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	for ix, segment := range strings.Split(string(input), "\n"+s.SplitLine+"\n") {
		name := ""
		if s.Name != "" {
			nameCommand := exec.Command("/usr/bin/bash", "-c", resolveAbbrev(s.Name))
			nameCommand.Stdin = bytes.NewBuffer([]byte(segment))
			nameOutput, err := nameCommand.CombinedOutput()
			if err != nil {
				details := ""
				details += fmt.Sprintf("Namer is failed: %s\n", s.Name)
				details += fmt.Sprintln("---stdin---")
				details += fmt.Sprintln(segment)
				details += fmt.Sprintln("---combined output---")
				details += fmt.Sprintln(string(nameOutput))
				details += fmt.Sprintln("---")
				os.Stderr.WriteString(details)
				return errors.WithStack(err)
			}
			name = strings.TrimSpace(string(nameOutput))
		}

		if s.Name == "" && name == "" {
			name = fmt.Sprintf("file-%d", ix)
		}
		err := os.WriteFile(name, []byte(segment), 0644)
		if err != nil {
			return errors.WithMessagef(err, "Failed to save file %s", name)
		}
	}
	return nil
}
