package k8s

import (
	"fyksz/data"
	"fyksz/helper"
	"io"
	"os"
)

type Env struct {
	Key    string `arg:""`
	Value  string `arg:""`
	input  io.Reader
	output io.Writer
}

func (e Env) Run() error {
	if e.input == nil {
		e.input = os.Stdin
	}
	if e.output == nil {
		e.output = os.Stdout
	}
	return helper.ProcessNode(e.input, e.output, func(root data.Node) error {
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
		return nil
	})
}
