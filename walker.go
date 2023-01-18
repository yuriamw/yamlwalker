package yamlwalker

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (walker *YamlWalker) findNode(parts []string) (node *YamlWalker, err error) {
	m, ok := walker.data.(map[string]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}

	n := walker

	for i := 0; i < len(parts); i++ {
		log(fmt.Sprintf("p:%v\n", parts[i]))
		var ok bool
		n, ok = m[parts[i]]
		if !ok {
			err = ErrNotFound
			return
		}
		if i < len(parts)-1 {
			m, ok = n.data.(map[string]*YamlWalker)
			if !ok {
				err = ErrInvalidType
				return
			}
		}
	}

	node = n

	return
}

func (walker *YamlWalker) findParent(parts []string) (parent *YamlWalker, err error) {
	if len(parts) == 0 {
		err = ErrKeyMismatch
		return
	}

	parent = walker

	if len(parts) > 1 {
		parentParts := parts[:len(parts)-1]
		var e error
		parent, e = walker.findNode(parentParts)
		if e != nil {
			err = e
			return
		}
	}

	return
}

// func (walker *YamlWalker) setNode(parts []string, node *YamlWalker) (err error) {
// 	n, err := walker.findNode(parts)
// 	if err != nil {
// 		return
// 	}

// 	n.data = node.data
// 	n.keys = node.keys
// 	n.style = node.style

// 	return
// }

func (walker *YamlWalker) appendNode(parts []string, node *YamlWalker, keyStyle yaml.Style) (err error) {
	parent, err := walker.findParent(parts)
	if err != nil {
		return
	}

	childName := parts[len(parts)-1]

	if parent.keyExists(childName) {
		return ErrDuplicateKey
	}

	var m map[string]*YamlWalker
	if parent.data == nil {
		m = make(map[string]*YamlWalker)
	} else {
		var ok bool
		m, ok = parent.data.(map[string]*YamlWalker)
		if !ok {
			return ErrInvalidType
		}
	}

	m[childName] = node
	parent.keys = append(parent.keys, yamlKey{name: childName, style: keyStyle})
	parent.data = m

	return
}

func (walker *YamlWalker) deleteNode(parts []string) (err error) {
	parent, err := walker.findParent(parts)
	if err != nil {
		return
	}
	childName := parts[len(parts)-1]

	m, ok := parent.data.(map[string]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}

	found := false
	for i := range parent.keys {
		if parent.keys[i].name == childName {
			found = true
			parent.keys = append(parent.keys[:i], parent.keys[i+1:]...)
			break
		}
	}
	if !found {
		err = ErrNotFound
		return
	}

	delete(m, childName)

	return
}

func (walker *YamlWalker) keyExists(keyName string) bool {
	for _, k := range walker.keys {
		if k.name == keyName {
			return true
		}
	}
	return false
}
