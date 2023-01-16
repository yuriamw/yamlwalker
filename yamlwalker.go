package yamlwalker

import (
	"gopkg.in/yaml.v3"
)

type YamlWalker struct {
	data  interface{}
	keys  []yamlKey
	style yaml.Style
}

type yamlKey struct {
	style yaml.Style
	name  string
}

func NewYamlWalker() *YamlWalker {
	return &YamlWalker{
		keys: make([]yamlKey, 0),
	}
}

func (walker *YamlWalker) UnmarshalYAML(value *yaml.Node) error {

	newYW, err := walker.decode(value)
	if err != nil {
		return err
	}

	walker.data = newYW.data
	walker.keys = newYW.keys

	return nil
}

func (walker *YamlWalker) MarshalYAML() (interface{}, error) {
	buffer, err := walker.encode()
	return buffer, err
}

func (walker *YamlWalker) Value() interface{} {
	return walker.data
}
