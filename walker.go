package yamlwalker

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func (walker *YamlWalker) splitPath(path string) []string {
	parts := []string{}
	if len(path) > 0 {
		parts = strings.Split(path, Separator)
	}
	return parts
}

func (walker *YamlWalker) asMap(parts []string) (children map[string]*YamlWalker, err error) {
	w, err := walker.findNode(parts)
	if err != nil {
		return
	}

	m, ok := w.data.(map[string]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}
	children = m
	return
}

func (walker *YamlWalker) asSlice(parts []string) (children []*YamlWalker, err error) {
	w, err := walker.findNode(parts)
	if err != nil {
		return
	}

	s, ok := w.data.([]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}
	children = s
	return
}

func (walker *YamlWalker) asString(parts []string) (value string, err error) {
	w, err := walker.findNode(parts)
	if err != nil {
		return
	}

	s, ok := w.data.(string)
	if !ok {
		err = ErrInvalidType
		return
	}
	value = s
	return
}

func (walker *YamlWalker) asInt(parts []string) (value int, err error) {
	w, err := walker.findNode(parts)
	if err != nil {
		return
	}

	s, ok := w.data.(int)
	if !ok {
		err = ErrInvalidType
		return
	}
	value = s
	return
}

func (walker *YamlWalker) asBool(parts []string) (value bool, err error) {
	w, err := walker.findNode(parts)
	if err != nil {
		return
	}

	s, ok := w.data.(bool)
	if !ok {
		err = ErrInvalidType
		return
	}
	value = s
	return
}

func (walker *YamlWalker) remove(parts []string, index int) error {
	w, err := walker.findNode(parts)
	if err != nil {
		return ErrNotFound
	}

	s, ok := w.data.([]*YamlWalker)
	if !ok {
		return ErrInvalidType
	}

	if index < 0 || index >= len(s) {
		return ErrInvalidRange
	}

	w.data = append(s[:index], s[index+1:]...)

	return nil
}

func (walker *YamlWalker) insert(parts []string, index int, node *YamlWalker) error {
	w, err := walker.findNode(parts)
	if err != nil {
		return err
	}

	s, ok := w.data.([]*YamlWalker)
	if !ok {
		return ErrInvalidType
	}

	if index < 0 || index > len(s) {
		return ErrInvalidRange
	}

	if len(s) == index { // nil or empty slice or after last element
		walker.data = append(s, node)
		return nil
	}
	s = append(s[:index+1], s[index:]...) // index < len(a)
	s[index] = node
	walker.data = s

	return nil
}

func (walker *YamlWalker) findNode(parts []string) (node *YamlWalker, err error) {
	n := walker
	if len(parts) == 0 {
		node = n
		return
	}

	m, ok := walker.data.(map[string]*YamlWalker)
	if !ok {
		err = ErrInvalidType
		return
	}

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
