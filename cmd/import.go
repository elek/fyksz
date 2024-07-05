package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type Import struct {
	Path string `arg:""`
}

func (i Import) Run() error {
	stat, err := os.Stat(i.Path)
	if err != nil {
		return errors.WithStack(err)
	}
	if stat.IsDir() {
		entries, err := os.ReadDir(i.Path)
		if err != nil {
			return errors.WithStack(err)
		}
		for _, entry := range entries {
			content, err := os.ReadFile(filepath.Join(i.Path, entry.Name()))
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Println(string(content))
			fmt.Println("---")
		}
	} else {
		content, err := os.ReadFile(i.Path)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println(string(content))
	}
	return nil
}
