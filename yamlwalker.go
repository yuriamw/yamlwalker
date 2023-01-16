package yamlwalker

import (
	"fmt"
	"strings"

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

const (
	Separator string = "."
)

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

func (walker *YamlWalker) Get(path string) interface{} {
	if len(path) == 0 {
		return walker.Value()
	}

	parts := strings.Split(path, Separator)
	m := walker.data.(map[string]*YamlWalker)
	var node *YamlWalker
	for i := 0; i < len(parts); i++ {
		log(fmt.Sprintf("p:%v\n", parts[i]))
		var ok bool
		node, ok = m[parts[i]]
		if !ok {
			return nil
		}
		if i < len(parts)-1 {
			m = node.data.(map[string]*YamlWalker)
		}
	}
	return node.Value()
}
