package main

import (
	"bytes"
	"fmt"
	"github.com/google/shlex"
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
			namer, err := shlex.Split(s.Name)
			if err != nil {
				return errors.WithStack(err)
			}
			nameCommand := exec.Command(namer[0], namer[1:]...)
			nameCommand.Stdin = bytes.NewBuffer([]byte(segment))
			nameOutput, err := nameCommand.CombinedOutput()
			if err != nil {
				fmt.Printf("Namer is failed: %s", s.Name)
				fmt.Println("---stdin---")
				fmt.Println(segment)
				fmt.Println("---combined output---")
				fmt.Println(nameOutput)
				fmt.Println("---")
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
