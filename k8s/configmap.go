package k8s

import (
	"fmt"
	"fyksz/yaml"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

type AsConfigMap struct {
	Input string `arg:""`
	Name  string `default:"config"`
	Key   string `default:""`
}

func (a AsConfigMap) Run() error {
	raw, err := os.ReadFile(a.Input)
	if err != nil {
		return errors.WithStack(err)
	}
	result := map[string]interface{}{}
	result["apiVersion"] = "v1"
	result["kind"] = "ConfigMap"
	result["metadata"] = map[string]interface{}{
		"name": a.Name,
	}
	result["data"] = map[string]interface{}{
		filepath.Base(a.Input): string(raw),
	}
	serialized, err := yaml.Marshal(result)

	var input []byte
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
	fmt.Println(string(input))

	_, err = os.Stdout.WriteString("---\n")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = os.Stdout.Write(serialized)
	return err
}
