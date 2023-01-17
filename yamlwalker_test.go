package yamlwalker

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed test_data/not-key-value.yaml
	notKeyValFile []byte
	//go:embed test_data/unsupported.yaml
	aliasFile []byte
	//go:embed test_data/x.yaml
	xFile []byte
	//go:embed test_data/y.yaml
	yFile []byte
	//go:embed test_data/slice.yaml
	sliceFile []byte
)

type YamlWalkerTestSuite struct {
	suite.Suite
}

func NewYamlWalkerTestSuite() *YamlWalkerTestSuite {
	return &YamlWalkerTestSuite{}
}

func TestYamlWalkerTestSuite(t *testing.T) {
	suite.Run(t, NewYamlWalkerTestSuite())
}

func (suite *YamlWalkerTestSuite) TestUnmarshal() {
	valid := NewYamlWalker()
	err := yaml.Unmarshal(xFile, valid)
	suite.Assert().Nil(err)

	errTests := []struct {
		name string
		body []byte
		err  error
	}{
		{
			name: "not a key-value",
			body: notKeyValFile,
			err:  fmt.Errorf("yaml: line 2: could not find expected ':'"),
		},
		{
			name: "alias not supported",
			body: aliasFile,
			err:  fmt.Errorf("line 2: unsupported node kind %d(%s)", yaml.AliasNode, decodeKind(yaml.AliasNode)),
		},
	}

	for _, tc := range errTests {
		tc := tc

		suite.Run(tc.name, func() {
			invalid := NewYamlWalker()
			err := yaml.Unmarshal(tc.body, invalid)
			suite.Assert().EqualError(err, tc.err.Error())
		})
	}
}

func (suite *YamlWalkerTestSuite) TestGet() {
	y := NewYamlWalker()
	err := yaml.Unmarshal(xFile, y)
	suite.Assert().Nil(err)

	tests := []struct {
		name     string
		path     string
		verifier func(value interface{}) bool
	}{
		{
			name: "path empty",
			path: "",
			verifier: func(value interface{}) bool {
				var y *YamlWalker
				var found bool
				mm := value.(map[string]*YamlWalker)
				y, found = mm["object"]
				if !found {
					return false
				}
				if y.keys[0].name != "name" {
					return false
				}
				y, found = y.data.(map[string]*YamlWalker)["name"]
				if !found {
					return false
				}
				if y.keys[0].name != "param" {
					return false
				}
				y, found = y.data.(map[string]*YamlWalker)["param"]
				if !found {
					return false
				}

				y, found = mm["another"]
				if !found {
					return false
				}
				if y.keys[0].name != "something" {
					return false
				}
				y, found = y.data.(map[string]*YamlWalker)["something"]
				if !found {
					return false
				}
				if y.keys[0].name != "interresting" {
					return false
				}
				y, found = y.data.(map[string]*YamlWalker)["interresting"]
				if !found {
					return false
				}

				return true
			},
		},
		{
			name: "path in a middle",
			path: "object",
			verifier: func(value interface{}) bool {
				var y *YamlWalker
				var found bool
				mm := value.(map[string]*YamlWalker)
				y, found = mm["name"]
				if !found {
					return false
				}
				if y.keys[0].name != "param" {
					return false
				}
				return true
			},
		},
		{
			name: "value at path 1",
			path: "object.name.param",
			verifier: func(value interface{}) bool {
				x := interface{}("value")
				return value == x
			},
		},
		{
			name: "value at path 2",
			path: "another.something.interresting",
			verifier: func(value interface{}) bool {
				x := interface{}("thing")
				return value == x
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		suite.Run(tc.name, func() {
			node, err := y.Get(tc.path)
			suite.Assert().Nil(err)
			suite.Assert().True(tc.verifier(node.Value()))
			value := y.GetValue(tc.path)
			suite.Assert().NotNil(value)
			suite.Assert().True(tc.verifier(value))
		})
	}
}

func (suite *YamlWalkerTestSuite) TestGetError() {
	mixedKeys := NewYamlWalker()
	err := yaml.Unmarshal(yFile, mixedKeys)
	suite.Assert().Nil(err)

	slice := NewYamlWalker()
	err = yaml.Unmarshal(sliceFile, slice)
	suite.Assert().Nil(err)

	tests := []struct {
		name   string
		walker *YamlWalker
		path   string
		err    error
	}{
		{
			name:   "not found",
			walker: mixedKeys,
			path:   "missing",
			err:    ErrNotFound,
		},
		{
			name:   "invalid type at path",
			walker: mixedKeys,
			path:   "object.array.name",
			err:    ErrInvalidType,
		},
		{
			name:   "invalid type at start",
			walker: slice,
			path:   "object",
			err:    ErrInvalidType,
		},
	}

	for _, tc := range tests {
		tc := tc

		suite.Run(tc.name, func() {
			_, err := tc.walker.Get(tc.path)
			suite.Assert().EqualError(err, tc.err.Error())
			value := tc.walker.GetValue(tc.path)
			suite.Assert().Nil(value)
		})
	}
}
