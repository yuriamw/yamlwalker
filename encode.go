package yamlwalker

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (walker *YamlWalker) encode() (node *yaml.Node, err error) {
	switch walker.data.(type) {
	case map[string]*YamlWalker:
		node, err = walker.encodeMap()
		return
	case []*YamlWalker:
		node, err = walker.encodeSeq()
		return
	default:
		node, err = walker.encodeScalar()
		return
	}
}

func (walker *YamlWalker) encodeMap() (node *yaml.Node, err error) {
	x := walker.data.(map[string]*YamlWalker)
	count := len(x) * 2

	node = &yaml.Node{
		Kind:    yaml.MappingNode,
		Content: make([]*yaml.Node, count),
		Style:   walker.style,
	}

	for i := 0; i < count; i += 2 {
		keyIdx := i / 2
		key := walker.keys[keyIdx]

		keyNode := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: key.name,
			Style: key.style,
		}
		node.Content[i] = keyNode

		value := x[key.name]
		content, e := value.encode()
		if err != nil {
			err = e
			return
		}
		node.Content[i+1] = content
	}

	return
}

func (walker *YamlWalker) encodeSeq() (node *yaml.Node, err error) {
	x := walker.data.([]*YamlWalker)

	node = &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: make([]*yaml.Node, len(x)),
		Style:   walker.style,
	}

	for i, value := range x {
		n, e := value.encode()
		if e != nil {
			err = e
			return
		}
		node.Content[i] = n
	}

	return
}

func (walker *YamlWalker) encodeScalar() (node *yaml.Node, err error) {
	node = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: walker.style,
	}

	var ok bool
	node.Value, ok = walker.data.(string)
	if !ok {
		err = fmt.Errorf("conversion from type %T failed", walker.data)
		return
	}

	return
}
