package yamlwalker

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed test_data/complex.yaml
	complexFile []byte
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

func (suite *YamlWalkerTestSuite) TestMarshal() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"first": {
				data: map[string]*YamlWalker{
					"one": {
						data: []*YamlWalker{
							{data: "bit", style: yaml.SingleQuotedStyle},
							{data: "byte", style: yaml.SingleQuotedStyle},
							{data: "word", style: yaml.DoubleQuotedStyle},
							{data: "dword", style: yaml.DoubleQuotedStyle},
						},
					},
				},
				keys: []yamlKey{
					{name: "one", style: 0},
				},
			},
			"second": {
				data: map[string]*YamlWalker{
					"two": {
						data: []*YamlWalker{
							{
								data: map[string]*YamlWalker{
									"apple":  {data: "fruit"},
									"walnut": {data: "nut"},
									"pear":   {data: "bean"},
								},
								keys: []yamlKey{
									{name: "apple", style: yaml.DoubleQuotedStyle},
									{name: "walnut", style: yaml.SingleQuotedStyle},
									{name: "pear", style: 0},
								},
							},
							{
								data: map[string]*YamlWalker{
									"husky":     {data: "dog"},
									"main coon": {data: "cat"},
								},
								keys: []yamlKey{
									{name: "husky", style: 0},
									{name: "main coon", style: 0},
								},
							},
						},
					},
				},
				keys: []yamlKey{
					{name: "two", style: 0},
				},
			},
			"3": {
				data: map[string]*YamlWalker{
					"three": {
						data: []*YamlWalker{
							{
								data: map[string]*YamlWalker{
									"1": {
										data: map[string]*YamlWalker{
											"type":    {data: "int"},
											"primary": {data: true},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "primary", style: 0},
										},
									},
									"2": {
										data: map[string]*YamlWalker{
											"type":    {data: "int"},
											"primary": {data: true},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "primary", style: 0},
										},
									},
									"3": {
										data: map[string]*YamlWalker{
											"type":    {data: "int"},
											"primary": {data: true},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "primary", style: 0},
										},
									},
									"4": {
										data: map[string]*YamlWalker{
											"type":    {data: "int"},
											"primary": {data: false},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "primary", style: 0},
										},
									},
								},
								keys: []yamlKey{
									{name: "1", style: 0},
									{name: "2", style: yaml.DoubleQuotedStyle},
									{name: "3", style: yaml.SingleQuotedStyle},
									{name: "4", style: 0},
								},
							},
							{
								data: map[string]*YamlWalker{
									"1.0": {
										data: map[string]*YamlWalker{
											"type":           {data: "float"},
											"representation": {data: "precise"},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "representation", style: 0},
										},
									},
									"3.333": {
										data: map[string]*YamlWalker{
											"type":           {data: "float"},
											"representation": {data: "imprecise"},
										},
										keys: []yamlKey{
											{name: "type", style: 0},
											{name: "representation", style: 0},
										},
									},
								},
								keys: []yamlKey{
									{name: "1.0", style: yaml.SingleQuotedStyle},
									{name: "3.333", style: yaml.SingleQuotedStyle},
								},
							},
						},
					},
				},
				keys: []yamlKey{
					{name: "three", style: 0},
				},
			},
			"flower rating": {
				data: map[string]*YamlWalker{
					"magnolia": {
						data: 3,
					},
					"tulip": {
						data: 1,
					},
					"rose": {
						data: 2,
					},
				},
				keys: []yamlKey{
					{name: "magnolia", style: 0},
					{name: "tulip", style: 0},
					{name: "rose", style: 0},
				},
			},
		},
		keys: []yamlKey{
			{name: "first", style: 0},
			{name: "second", style: 0},
			{name: "3", style: 0},
			{name: "flower rating", style: 0},
		},
	}

	data, err := yaml.Marshal(y)
	suite.Assert().Nil(err)
	// fmt.Printf("---\n%s...\n", string(data))
	suite.Assert().Equal(complexFile, data)
}

func (suite *YamlWalkerTestSuite) TestMarshalError() {
	tests := []struct {
		name string
		yaml *YamlWalker
		err  error
	}{
		{
			name: "more keys than map",
			yaml: &YamlWalker{
				data: map[string]*YamlWalker{
					"first": {
						data: map[string]*YamlWalker{
							"one": {
								data: []*YamlWalker{
									{data: "1", style: 0},
									{data: "2", style: 0},
								},
							},
						},
						keys: []yamlKey{
							{name: "one", style: 0},
							{name: "two", style: 0},
						},
					},
				},
				keys: []yamlKey{
					{name: "first", style: 0},
				},
			},
			err: ErrKeyMismatch,
		},
		{
			name: "map bigger than keys",
			yaml: &YamlWalker{
				data: map[string]*YamlWalker{
					"first": {
						data: []*YamlWalker{
							{
								data: map[string]*YamlWalker{
									"one": {
										data: []*YamlWalker{
											{data: "1", style: 0},
											{data: "2", style: 0},
										},
									},
									"two": {
										data: []*YamlWalker{
											{data: "3", style: 0},
											{data: "4", style: 0},
										},
									},
								},
								keys: []yamlKey{
									{name: "one", style: 0},
								},
							},
						},
					},
				},
				keys: []yamlKey{
					{name: "first", style: 0},
				},
			},
			err: ErrKeyMismatch,
		},
		{
			name: "keys not found",
			yaml: &YamlWalker{
				data: map[string]*YamlWalker{
					"first": {
						data: "something",
					},
				},
				keys: []yamlKey{
					{name: "second", style: 0},
				},
			},
			err: ErrKeyMismatch,
		},
		{
			name: "no keys",
			yaml: &YamlWalker{
				data: map[string]*YamlWalker{
					"first": {
						data: "something",
					},
				},
			},
			err: ErrKeyMismatch,
		},
	}

	for _, tc := range tests {
		tc := tc

		suite.Run(tc.name, func() {
			_, err := yaml.Marshal(tc.yaml)
			suite.Assert().EqualError(err, tc.err.Error())
		})
	}
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
			err:  fmt.Errorf("line 3: unsupported node kind %d(%s)", yaml.AliasNode, decodeKind(yaml.AliasNode)),
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
