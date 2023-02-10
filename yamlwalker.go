package yamlwalker

import (
	"errors"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidType  = errors.New("invalid type conversion")
	ErrKeyMismatch  = errors.New("list of keys does not match map keys")
	ErrDuplicateKey = errors.New("duplicate key name")
	ErrInvalidRange = errors.New("index out of bounds")
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
	// Symbol to separate elements in node path
	Separator string = "."
)

// NewYamlWalker creates new YamlWalker node instance with a specified dataStyle
// Default dataStyle = 0.
func NewYamlWalker(dataStyle ...yaml.Style) *YamlWalker {
	style := yaml.Style(0)
	if len(dataStyle) > 0 {
		style = dataStyle[0]
	}
	return &YamlWalker{
		style: style,
		keys:  make([]yamlKey, 0),
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

// AsMap returns children of the node specified by path
// as map if node is yaml.MappingNode and err set to nil.
// If path does not exists err set to ErrNotFound.
// If the node is not yaml.MappingNode err set to ErrInvalidType.
//
// The map is usefull to iterate over the node children.
// Do not insert, delete or change elements of map directly, use Append(), Delete() or Update() instead.
func (walker *YamlWalker) AsMap(path string) (children map[string]*YamlWalker, err error) {
	return walker.asMap(walker.splitPath(path))
}

// AsSlice returns children of the node specified by path
// as slice if node is yaml.SequenceNode and err set to nil.
// If the node is not yaml.SequenceNode err set to ErrInvalidType.
//
// The slice is usefull to iterate over the node children.
// Do not insert, delete or change elements of slice directly, use Insert(), Remove() or Update() instead.
func (walker *YamlWalker) AsSlice(path string) (children []*YamlWalker, err error) {
	return walker.asSlice(walker.splitPath(path))
}

// Remove removes the item at the index from the slice of children.
// If the node specified by path is yaml.SequenceNode the item at the index is removed
// and err set to nil otherwise err set to ErrInvalidType.
// If index is out of slice bounds err set to ErrInvalidRange.
func (walker *YamlWalker) Remove(path string, index int) error {
	return walker.remove(walker.splitPath(path), index)
}

// Insert inserts the node into the slice of children at the index.
// If the item specified by path is yaml.SequenceNode the node inserted into the slice of children
// at the index and err set to nil otherwise err set to ErrInvalidType.
// If the index == len(children) the node is appnded at the end of slice.
// If index is out of slice bounds err set to ErrInvalidRange.
func (walker *YamlWalker) Insert(path string, index int, node *YamlWalker) error {
	return walker.insert(walker.splitPath(path), index, node)
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

	node, err = walker.findNode(walker.splitPath(path))

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

	style := yaml.Style(0)
	if len(keyStyle) > 0 {
		style = keyStyle[0]
	}

	return walker.appendNode(walker.splitPath(path), node, style)
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

	return walker.deleteNode(walker.splitPath(path))
}
