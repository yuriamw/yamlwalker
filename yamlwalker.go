package yamlwalker

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidType = errors.New("invalid type conversion")
	ErrKeyMismatch = errors.New("list of keys does not match map keys")
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

func (walker *YamlWalker) GetValue(path string) interface{} {
	node, err := walker.Get(path)
	if err != nil {
		return nil
	}
	return node.Value()
}

func (walker *YamlWalker) Get(path string) (node *YamlWalker, err error) {
	if len(path) == 0 {
		node = walker
		return
	}

	parts := strings.Split(path, Separator)
	m, ok := walker.data.(map[string]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}
	for i := 0; i < len(parts); i++ {
		log(fmt.Sprintf("p:%v\n", parts[i]))
		var ok bool
		node, ok = m[parts[i]]
		if !ok {
			err = ErrNotFound
			return
		}
		if i < len(parts)-1 {
			m, ok = node.data.(map[string]*YamlWalker)
			if !ok {
				err = ErrInvalidType
				return
			}
		}
	}

	return
}
