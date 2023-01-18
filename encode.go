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
		node = walker.encodeScalar()
		return
	}
}

func (walker *YamlWalker) encodeMap() (node *yaml.Node, err error) {
	x := walker.data.(map[string]*YamlWalker)

	numKeys := len(walker.keys)
	numKeysInMap := len(x)
	if numKeys != numKeysInMap {
		err = ErrKeyMismatch
		return
	}

	count := len(x) * 2

	n := &yaml.Node{
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
		n.Content[i] = keyNode

		value, found := x[key.name]
		if !found {
			err = ErrKeyMismatch
			fmt.Printf("%v\n", err)
			return
		}
		content, e := value.encode()
		if e != nil {
			err = e
			return
		}
		n.Content[i+1] = content
	}

	node = n

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

func (walker *YamlWalker) encodeScalar() (node *yaml.Node) {
	node = &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: walker.style,
	}

	node.Value = fmt.Sprintf("%v", walker.data)

	return
}
