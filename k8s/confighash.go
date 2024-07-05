package k8s

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fyksz/data"
	"fyksz/helper"
	"fyksz/yaml"
	"github.com/pkg/errors"
	"io"
	"os"
)

type ConfigHash struct {
	input  io.Reader
	output io.Writer
}

func (c ConfigHash) Run() error {
	if c.input == nil {
		c.input = os.Stdin
	}
	if c.output == nil {
		c.output = os.Stdout
	}
	nameToHash := map[string]string{}

	cachedInput, err := io.ReadAll(c.input)
	if err != nil {
		return errors.WithStack(err)
	}
	err = helper.ProcessAll(bytes.NewReader(cachedInput), bytes.NewBuffer([]byte{}), func(input string) (string, error) {
		parsed := yaml.MapSlice{}
		err := yaml.Unmarshal([]byte(input), &parsed)
		if err != nil {
			return "", errors.WithStack(err)
		}

		kind, ok := parsed.Get("kind")
		if !ok {
			return "", nil
		}
		metadata, ok := parsed.Get("metadata")
		if !ok {
			return "", nil
		}
		name, ok := metadata.(yaml.MapSlice).Get("name")
		if !ok {
			return "", nil
		}

		if kind == "ConfigMap" {
			hash := md5.Sum([]byte(input))
			nameToHash[name.(string)] = hex.EncodeToString(hash[:md5.Size])
		}

		return "", nil
	})
	if err != nil {
		return err
	}
	err = helper.ProcessNode(bytes.NewReader(cachedInput), c.output, func(node data.Node) error {
		getAll := data.GetAll{
			//Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "envFrom", ".*", "configMapRef", "name"),
			Path: data.NewPath("spec", "template", "spec", "volumes", ".*", "configMap", "name"),
		}
		node.Accept(&getAll)
		for _, match := range getAll.Result {
			configName := match.Value.(*data.KeyNode).Value.(string)
			if val, ok := nameToHash[configName]; ok {
				annotations := data.SmartGetAll{
					Path: data.NewPath("metadata", "annotations"),
				}
				node.Accept(&annotations)
				for _, annotation := range annotations.Result {
					annotation.Value.(*data.MapNode).PutValue("fyksz-config-hash/"+configName, val)
				}
				return nil
			}
		}
		return nil
	})
	return err
}
