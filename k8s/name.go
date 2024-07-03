package k8s

import (
	"fmt"
	"fyksz/yaml"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type Name struct {
}

func (n Name) Run() error {
	var structured map[string]interface{}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}
	err = yaml.Unmarshal(raw, &structured)
	if err != nil {
		return errors.WithStack(err)
	}
	name, _ := structured["metadata"].(yaml.MapSlice).Get("name")
	kind := strings.ToLower(structured["kind"].(string))
	fmt.Printf("%s-%s.yaml", name, kind)
	return nil
}
