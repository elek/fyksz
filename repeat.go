package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type Repeat struct {
	Elements string `help:"coma separated list of elements to repeat the original template"`
	Pattern  string `help:"pattern to be replace with the element"`
}

func (s Repeat) Run() error {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}
	res := []string{}
	for _, element := range strings.Split(s.Elements, ",") {
		res = append(res, strings.ReplaceAll(string(raw), s.Pattern, element))
	}
	fmt.Println(strings.Join(res, "\n---\n"))
	return nil
}
