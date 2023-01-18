package yamlwalker

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidType  = errors.New("invalid type conversion")
	ErrKeyMismatch  = errors.New("list of keys does not match map keys")
	ErrDuplicateKey = errors.New("duplicate key name")
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

// NewYamlWalker creates new YamlWalker node instance
func NewYamlWalker() *YamlWalker {
	return &YamlWalker{
		keys: make([]yamlKey, 0),
	}
}

// UnmarshalYAML decode YAML into internal representation
func (walker *YamlWalker) UnmarshalYAML(value *yaml.Node) error {

	newYW, err := walker.decode(value)
	if err != nil {
		return err
	}

	walker.data = newYW.data
	walker.keys = newYW.keys

	return nil
}

// MarshalYAML encode internal representation to YAML
func (walker *YamlWalker) MarshalYAML() (interface{}, error) {
	buffer, err := walker.encode()
	return buffer, err
}

// Value returns the value of the node
func (walker *YamlWalker) Value() interface{} {
	return walker.data
}

// Update updates the value of the node.
// All previouse data is lost.
// To assign mapped tree of new nodes or sequence of nodes use Set() instead.
func (walker *YamlWalker) Update(value interface{}) {
	walker.data = value
	walker.keys = make([]yamlKey, 0)
}

// GetValue returns the value of the node specified by path or <nil> if node does not exists
// It is equivalent to call Get(path).Value()
// If the node's underlying type is mapping node (Kind == yaml.MappingNode)
// it returns map[string]*YamlWalker.
// If the node's underlying type is sequence node (Kind == yaml.SequenceNode)
// it returns []*YamlWalker.
func (walker *YamlWalker) GetValue(path string) interface{} {
	node, err := walker.Get(path)
	if err != nil {
		return nil
	}
	return node.Value()
}

// Get returns the node specified by path or ErrNotFound if node does not exists
// It searches through the tree of mapping nodes (Kind == yaml.MappingNode).
// Empty path returns the top node.
// If node Kind other than yaml.MappingNode occurs in the middle of the tree it returns ErrInvalidType.
func (walker *YamlWalker) Get(path string) (node *YamlWalker, err error) {
	if len(path) == 0 {
		node = walker
		return
	}

	parts := strings.Split(path, Separator)

	node, err = walker.findNode(parts)

	return
}

// SetValue sets the value of the node at the specified path.
// All previouse data is lost.
// If path does not exists it silently does nothing.
func (walker *YamlWalker) SetValue(path string, value interface{}) {
	node, err := walker.Get(path)
	if err != nil {
		return
	}
	node.Update(value)
}

// Set sets the node at the specified path by making a copy of node properties.
// All previouse data is lost.
// It returns ErrNotFound if node does not exists.
// It searches through the tree of mapping nodes (Kind == yaml.MappingNode).
// Empty calling with empty path is equivalent to call SetValue(node.Value()).
// If node Kind other than yaml.MappingNode occurs in the middle of the tree it returns ErrInvalidType.
func (walker *YamlWalker) Set(path string, node *YamlWalker) error {
	existing, err := walker.Get(path)
	if err != nil {
		return ErrNotFound
	}

	existing.Update(node.Value())
	return nil
}

// Append appends the node to the map at the path.
// Default keyStyle = 0.
// It searches through the tree of mapping nodes (Kind == yaml.MappingNode).
// If node Kind other than yaml.MappingNode occurs in the middle of the tree it returns ErrInvalidType.
// If key name (last part of the path) already exists it returns ErrDuplicateKey.
func (walker *YamlWalker) Append(path string, node *YamlWalker, keyStyle ...yaml.Style) error {
	if len(path) == 0 {
		return ErrKeyMismatch
	}

	parts := strings.Split(path, Separator)

	style := yaml.Style(0)
	if len(keyStyle) > 0 {
		style = keyStyle[0]
	}

	return walker.appendNode(parts, node, style)
}

// Delete deletes the node from the map at the path.
// It searches through the tree of mapping nodes (Kind == yaml.MappingNode).
// If node Kind other than yaml.MappingNode occurs in the middle of the tree it returns ErrInvalidType.
// If key name (last part of the path) does not exists it returns ErrNotFound.
// All node's children are lost.
func (walker *YamlWalker) Delete(path string) error {
	if len(path) == 0 {
		return ErrKeyMismatch
	}

	parts := strings.Split(path, Separator)

	return walker.deleteNode(parts)
}
