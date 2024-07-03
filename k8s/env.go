package k8s

import (
	"fmt"
	"fyksz/data"
	"fyksz/yaml"
	"github.com/pkg/errors"
	"io"
	"os"
)

type Env struct {
	Key   string `arg:""`
	Value string `arg:""`
}

func (e Env) Run() error {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return errors.WithStack(err)
	}

	parsed := yaml.MapSlice{}

	err = yaml.Unmarshal(raw, &parsed)
	if err != nil {
		return errors.WithStack(err)
	}

	root, err := data.ConvertToNode(parsed, data.NewPath())
	if err != nil {
		return errors.WithStack(err)
	}
	smartGet := data.SmartGetAll{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "env")}
	root.Accept(&smartGet)
	for _, result := range smartGet.Result {
		found := false
		envs := result.Value.(*data.ListNode)
		for _, env := range envs.Children {
			if env.(*data.MapNode).GetStringValue("name") == e.Key {
				env.(*data.MapNode).Get("value").(*data.KeyNode).Value = e.Value
				found = true
			}
		}
		if !found {
			envEntry := envs.CreateMap()
			envEntry.PutValue("name", e.Key)
			envEntry.PutValue("value", e.Value)
		}
	}

	outYaml := data.ConvertToYaml(root)
	outRaw, err := yaml.Marshal(outYaml)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println(string(outRaw))
	return nil
}
