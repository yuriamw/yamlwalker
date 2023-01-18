package yamlwalker

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (walker *YamlWalker) decode(node *yaml.Node) (*YamlWalker, error) {
	log(fmt.Sprintf("+++++++++++++++++++++++\n%s", printNodeContent(node)))

	newYW := NewYamlWalker()
	newYW.style = node.Style

	switch node.Kind {
	case yaml.MappingNode:
		var err error
		newYW.keys, newYW.data, err = walker.decodeMap(node)
		if err != nil {
			return nil, err
		}
	case yaml.SequenceNode:
		var err error
		newYW.data, err = walker.decodeSeq(node)
		if err != nil {
			return nil, err
		}
	case yaml.ScalarNode:
		newYW.data = walker.decodeScalar(node)
	default:
		err := fmt.Errorf("line %d: unsupported node kind %d(%s)", node.Line, node.Kind, decodeKind(node.Kind))
		return nil, err
	}

	log(fmt.Sprintf("-----------------------\n"))

	return newYW, nil
}

func (walker *YamlWalker) decodeMap(node *yaml.Node) (keys []yamlKey, data map[string]*YamlWalker, err error) {
	count := len(node.Content) / 2

	keys = make([]yamlKey, count)
	data = make(map[string]*YamlWalker)

	for i := 0; i < count; i++ {
		contentIdx := i * 2
		contentKey := node.Content[contentIdx]
		keyName := contentKey.Value
		keyStyle := contentKey.Style
		contentValue := node.Content[contentIdx+1]
		value, e := NewYamlWalker().decode(contentValue)
		if e != nil {
			err = e
			return
		}
		keys[i] = yamlKey{
			style: keyStyle,
			name:  keyName,
		}
		data[keyName] = value

		log(fmt.Sprintf("=    map[%s]=%v\n", keyName, value.Value()))
	}

	return
}

func (walker *YamlWalker) decodeSeq(node *yaml.Node) ([]*YamlWalker, error) {
	slice := make([]*YamlWalker, len(node.Content))
	for i, v := range node.Content {
		sibling, err := NewYamlWalker().decode(v)
		if err != nil {
			return nil, err
		}
		slice[i] = sibling

		log(fmt.Sprintf("=    [%d]=%v\n", i, sibling.Value()))

	}
	return slice, nil
}

func (walker *YamlWalker) decodeScalar(node *yaml.Node) interface{} {
	return node.Value
}

// func convert(i interface{}) interface{} {
// 	switch x := i.(type) {
// 	case map[interface{}]interface{}:
// 		m2 := map[string]interface{}{}
// 		for k, v := range x {
// 			m2[k.(string)] = convert(v)
// 		}
// 		return m2
// 	case []interface{}:
// 		for i, v := range x {
// 			x[i] = convert(v)
// 		}
// 	}
// 	return i
// }
