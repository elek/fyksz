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

type Apply struct {
	Target    []string `arg:""`
	SplitLine string   `default:"---"`
	Filter    string   `usage:"Executable command. If returns with 0 but not empty stdout, the target will be executed on the input"`
}

func (a Apply) Run() error {
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

	for ix, segment := range strings.Split(string(input), "\n"+a.SplitLine+"\n") {
		if ix > 0 {
			fmt.Println(a.SplitLine)
		}
		if strings.HasSuffix(a.Target[0], ".sh") {
			a.Target = append([]string{"bash", "-c"}, strings.Join(a.Target, " "))
		}

		include := true
		if a.Filter != "" {
			filter, err := shlex.Split(a.Filter)
			if err != nil {
				return errors.WithStack(err)
			}
			filterCommand := exec.Command(filter[0], filter[1:]...)
			filterCommand.Stdin = bytes.NewBuffer([]byte(segment))
			filterOutput, err := filterCommand.CombinedOutput()
			if err != nil {
				fmt.Printf("Filter command is failed: %s", a.Filter)
				return errors.WithStack(err)
			}
			if strings.TrimSpace(string(filterOutput)) == "" {
				include = false
			}
		}
		if include {
			if strings.HasPrefix(a.Target[0], "@") {
				a.Target[0] = a.Target[0][1:]
				a.Target = append([]string{"fyksz"}, a.Target...)
			}
			command := exec.Command(a.Target[0], a.Target[1:]...)
			command.Stdin = bytes.NewBuffer([]byte(segment))
			output, err := command.CombinedOutput()

			if err != nil {
				fmt.Println(a.Target)
				fmt.Println(string(output))
				return errors.WithStack(err)
			}
			out := strings.TrimSpace(string(output))
			if out == "" {
				fmt.Println(segment)
			}
			fmt.Println(string(output))
		} else {
			fmt.Println(segment)
		}
	}
	return nil
}
