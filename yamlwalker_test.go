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

func (suite *YamlWalkerTestSuite) TestAsMap() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"map": {
				data: map[string]*YamlWalker{
					"first": {
						data: 1,
					},
					"second": {
						data: 2,
					},
				},
				keys: []yamlKey{
					{name: "first"},
					{name: "second"},
				},
			},
		},
		keys: []yamlKey{
			{name: "map"},
		},
	}

	tests := []struct {
		name     string
		key      string
		count    int
		children []string
	}{
		{
			name:     "empty key",
			key:      "",
			count:    1,
			children: []string{"map"},
		},
		{
			name:     "map key",
			key:      "map",
			count:    2,
			children: []string{"first", "second"},
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {

			m, err := y.AsMap(tc.key)
			suite.Assert().Nil(err)
			suite.Assert().Equal(tc.count, len(m))
			for _, k := range tc.children {
				_, found := m[k]
				suite.Assert().True(found)
			}
		})
	}

	errTests := []struct {
		name string
		key  string
		err  error
	}{
		{
			name: "invalid type",
			key:  "map.first",
			err:  ErrInvalidType,
		},
		{
			name: "not found as root",
			key:  "invalid",
			err:  ErrNotFound,
		},
		{
			name: "not found as child",
			key:  "map.invalid",
			err:  ErrNotFound,
		},
	}

	for _, tc := range errTests {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := y.AsMap(tc.key)
			suite.Assert().EqualError(err, tc.err.Error())
		})
	}
}

func (suite *YamlWalkerTestSuite) TestAsSlice() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"parent": {
				data: []*YamlWalker{
					{data: 1},
					{data: 2},
				},
			},
		},
		keys: []yamlKey{
			{name: "parent"},
		},
	}

	tests := []struct {
		name   string
		path   string
		count  int
		values []int
	}{
		{
			name:   "valid",
			path:   "parent",
			count:  2,
			values: []int{1, 2},
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			s, err := y.AsSlice(tc.path)
			suite.Assert().Nil(err)
			suite.Assert().Equal(tc.count, len(s))
			for i := range tc.values {
				v, found := s[i].data.(int)
				suite.Assert().True(found)
				suite.Assert().Equal(tc.values[i], v)
			}
		})
	}

	errTests := []struct {
		name string
		path string
		err  error
	}{
		{
			name: "invalid type",
			path: "",
			err:  ErrInvalidType,
		},
		{
			name: "not found as root",
			path: "invalid",
			err:  ErrNotFound,
		},
	}

	for _, tc := range errTests {
		tc := tc
		suite.Run(tc.name, func() {
			_, err := y.AsSlice(tc.path)
			suite.Assert().EqualError(err, tc.err.Error())
		})
	}
}

func (suite *YamlWalkerTestSuite) TestAsType() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"str": {
				data: "string value",
			},
			"bool": {
				data: true,
			},
			"int": {
				data: 1,
			},
		},
		keys: []yamlKey{
			{name: "str"},
			{name: "bool"},
			{name: "int"},
		},
	}

	s, err := y.AsString("str")
	suite.Assert().Nil(err)
	suite.Assert().Equal("string value", s)

	s, err = y.AsString("bool")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	s, err = y.AsString("int")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	s, err = y.AsString("invalud")
	suite.Assert().EqualError(err, ErrNotFound.Error())

	i, err := y.AsInt("str")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	i, err = y.AsInt("int")
	suite.Assert().Nil(err)
	suite.Assert().Equal(1, i)

	i, err = y.AsInt("bool")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	i, err = y.AsInt("invalud")
	suite.Assert().EqualError(err, ErrNotFound.Error())

	b, err := y.AsBool("str")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	b, err = y.AsBool("int")
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	b, err = y.AsBool("bool")
	suite.Assert().Nil(err)
	suite.Assert().Equal(true, b)

	b, err = y.AsBool("invalud")
	suite.Assert().EqualError(err, ErrNotFound.Error())
}

func (suite *YamlWalkerTestSuite) TestInsert() {
	getData := func() *YamlWalker {
		return &YamlWalker{
			data: map[string]*YamlWalker{
				"parent": {
					data: []*YamlWalker{
						{data: 1},
						{data: 2},
					},
				},
			},
			keys: []yamlKey{
				{name: "parent"},
			},
		}
	}

	tests := []struct {
		index  int
		values []int
	}{
		{
			index:  0,
			values: []int{3, 1, 2},
		},
		{
			index:  1,
			values: []int{1, 3, 2},
		},
		{
			index:  2,
			values: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(fmt.Sprintf("index:%d", tc.index), func() {
			y := getData()
			node := &YamlWalker{data: 3}
			err := y.Insert("parent", tc.index, node)
			suite.Assert().Nil(err)
			m, found := y.data.(map[string]*YamlWalker)
			suite.Assert().True(found)
			s, found := m["parent"].data.([]*YamlWalker)
			suite.Assert().True(found)
			suite.Assert().Equal(3, len(s))
			v1 := s[0]
			i1, found := v1.data.(int)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.values[0], i1)
			v2 := s[1]
			i2, found := v2.data.(int)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.values[1], i2)
			v3 := s[2]
			i3, found := v3.data.(int)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.values[2], i3)
		})
	}

	y := getData()
	err := y.Insert("invalid", 1, &YamlWalker{})
	suite.Assert().EqualError(err, ErrNotFound.Error())

	y = getData()
	err = y.Insert("", 1, &YamlWalker{})
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	y = getData()
	err = y.Insert("parent", -1, &YamlWalker{})
	suite.Assert().EqualError(err, ErrInvalidRange.Error())

	y = getData()
	err = y.Insert("parent", 10, &YamlWalker{})
	suite.Assert().EqualError(err, ErrInvalidRange.Error())
}

func (suite *YamlWalkerTestSuite) TestRemove() {
	getData := func() *YamlWalker {
		return &YamlWalker{
			data: map[string]*YamlWalker{
				"parent": {
					data: []*YamlWalker{
						{data: 1},
						{data: 2},
						{data: 3},
					},
				},
			},
			keys: []yamlKey{
				{name: "parent"},
			},
		}
	}

	tests := []struct {
		count  int
		index  int
		values []int
	}{
		{
			count:  2,
			index:  0,
			values: []int{2, 3},
		},
		{
			count:  2,
			index:  1,
			values: []int{1, 3},
		},
		{
			count:  2,
			index:  2,
			values: []int{1, 2},
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(fmt.Sprintf("index:%d", tc.index), func() {
			y := getData()
			err := y.Remove("parent", tc.index)
			suite.Assert().Nil(err)
			m, found := y.data.(map[string]*YamlWalker)
			suite.Assert().True(found)
			s, found := m["parent"].data.([]*YamlWalker)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.count, len(s))
			v1 := s[0]
			i1, found := v1.data.(int)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.values[0], i1)
			v2 := s[1]
			i2, found := v2.data.(int)
			suite.Assert().True(found)
			suite.Assert().Equal(tc.values[1], i2)
		})
	}

	y := getData()
	err := y.Remove("invalid", 1)
	suite.Assert().EqualError(err, ErrNotFound.Error())

	y = getData()
	err = y.Remove("", 1)
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	y = getData()
	err = y.Remove("parent", -1)
	suite.Assert().EqualError(err, ErrInvalidRange.Error())

	y = getData()
	err = y.Remove("parent", 10)
	suite.Assert().EqualError(err, ErrInvalidRange.Error())
}

func (suite *YamlWalkerTestSuite) TestFindNode() {
	y := &YamlWalker{}
	n, err := y.findNode([]string{})
	suite.Assert().Nil(err)
	suite.Assert().Equal(y, n)

	y = &YamlWalker{data: 0}
	_, err = y.findNode([]string{"item"})
	suite.Assert().EqualError(err, ErrInvalidType.Error())

	y = &YamlWalker{
		data: map[string]*YamlWalker{
			"first-0": {
				data: map[string]*YamlWalker{
					"first-1": {
						data: map[string]*YamlWalker{
							"first-2": {
								data: 1,
							},
						},
						keys: []yamlKey{
							{name: "first-2"},
						},
					},
					"second-1": {
						data: 4,
					},
				},
				keys: []yamlKey{
					{name: "first-1"},
					{name: "second-1"},
				},
			},
			"second-0": {
				data: 2,
			},
			"third-0": {
				data: 3,
			},
		},
		keys: []yamlKey{
			{name: "first-0"},
			{name: "second-0"},
			{name: "third-0"},
		},
	}

	w := y.data.(map[string]*YamlWalker)["first-0"]
	n, err = w.findNode([]string{})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(2, len(n.keys))
	suite.Assert().Equal("first-1", n.keys[0].name)
	suite.Assert().Equal("second-1", n.keys[1].name)

	n, err = y.findNode([]string{"first-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(2, len(n.keys))
	suite.Assert().Equal("first-1", n.keys[0].name)
	suite.Assert().Equal("second-1", n.keys[1].name)

	n, err = y.findNode([]string{"second-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(0, len(n.keys))
	i, ok := n.data.(int)
	suite.Assert().True(ok)
	suite.Assert().Equal(2, i)

	n, err = y.findNode([]string{"third-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(0, len(n.keys))
	i, ok = n.data.(int)
	suite.Assert().True(ok)
	suite.Assert().Equal(3, i)

	n, err = y.findNode([]string{"first-0", "first-1"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(1, len(n.keys))
	suite.Assert().Equal("first-2", n.keys[0].name)

	n, err = y.findNode([]string{"first-0", "first-1", "first-2"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(0, len(n.keys))
	i, ok = n.data.(int)
	suite.Assert().True(ok)
	suite.Assert().Equal(1, i)

	_, err = y.findNode([]string{"first-0", "invalid"})
	suite.Assert().EqualError(err, ErrNotFound.Error())
	_, err = y.findNode([]string{"second-0", "second-1", "second-2"})
	suite.Assert().EqualError(err, ErrInvalidType.Error())
}

func (suite *YamlWalkerTestSuite) TestFindParent() {
	y := &YamlWalker{}
	_, err := y.findParent([]string{})
	suite.Assert().EqualError(err, ErrKeyMismatch.Error())

	y = &YamlWalker{
		data: map[string]*YamlWalker{
			"first-0": {
				data: map[string]*YamlWalker{
					"first-1": {
						data: map[string]*YamlWalker{
							"first-2": {
								data: 1,
							},
						},
						keys: []yamlKey{
							{name: "first-2"},
						},
					},
					"second-1": {
						data: 4,
					},
				},
				keys: []yamlKey{
					{name: "first-1"},
					{name: "second-1"},
				},
			},
			"second-0": {
				data: 2,
			},
			"third-0": {
				data: 3,
			},
		},
		keys: []yamlKey{
			{name: "first-0"},
			{name: "second-0"},
			{name: "third-0"},
		},
	}

	w := y.data.(map[string]*YamlWalker)["first-0"]
	_, err = w.findParent([]string{})
	suite.Assert().EqualError(err, ErrKeyMismatch.Error())

	n, err := y.findParent([]string{"first-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(3, len(n.keys))
	suite.Assert().Equal("first-0", n.keys[0].name)
	suite.Assert().Equal("second-0", n.keys[1].name)
	suite.Assert().Equal("third-0", n.keys[2].name)

	n, err = y.findParent([]string{"second-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(3, len(n.keys))
	suite.Assert().Equal("first-0", n.keys[0].name)
	suite.Assert().Equal("second-0", n.keys[1].name)
	suite.Assert().Equal("third-0", n.keys[2].name)

	n, err = y.findParent([]string{"third-0"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(3, len(n.keys))
	suite.Assert().Equal("first-0", n.keys[0].name)
	suite.Assert().Equal("second-0", n.keys[1].name)
	suite.Assert().Equal("third-0", n.keys[2].name)

	n, err = y.findParent([]string{"first-0", "first-1"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(2, len(n.keys))
	suite.Assert().Equal("first-1", n.keys[0].name)
	suite.Assert().Equal("second-1", n.keys[1].name)

	n, err = y.findParent([]string{"first-0", "not-exists"})
	suite.Assert().Nil(err)
	suite.Assert().NotNil(n)
	suite.Assert().Equal(2, len(n.keys))
	suite.Assert().Equal("first-1", n.keys[0].name)
	suite.Assert().Equal("second-1", n.keys[1].name)

	_, err = y.findParent([]string{"second-0", "second-1", "second-2"})
	suite.Assert().EqualError(err, ErrInvalidType.Error())
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

func (suite *YamlWalkerTestSuite) TestSet() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"first": {data: 1},
			"second": {
				data: map[string]*YamlWalker{
					"child": {
						data: 2,
					},
				},
				keys: []yamlKey{
					{name: "child"},
				},
			},
		},
		keys: []yamlKey{
			{name: "first"},
			{name: "second"},
		},
	}

	n := NewYamlWalker()
	n.Update(3)
	err := y.Set("second.child", n)
	suite.Assert().Nil(err)
	m := y.Value().(map[string]*YamlWalker)
	d, found := m["second"]
	suite.Assert().True(found)
	d, found = d.Value().(map[string]*YamlWalker)["child"]
	suite.Assert().True(found)
	suite.Assert().Equal(3, d.Value())
	y.SetValue("first", 5)
	m = y.Value().(map[string]*YamlWalker)
	d, found = m["first"]
	suite.Assert().True(found)
	i, found := d.Value().(int)
	suite.Assert().True(found)
	suite.Assert().Equal(5, i)
	y.SetValue("invalid", 5)
	m = y.Value().(map[string]*YamlWalker)
	_, found = m["invalid"]
	suite.Assert().False(found)
}

func (suite *YamlWalkerTestSuite) TestSetError() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"first":  {data: 1},
			"second": {data: 2},
		},
		keys: []yamlKey{
			{name: "first"},
			{name: "second"},
		},
	}

	n := NewYamlWalker()
	n.Update(3)
	err := y.Set("second.child", n)
	suite.Assert().EqualError(err, ErrNotFound.Error())
}

func (suite *YamlWalkerTestSuite) TestSetInArray() {
	y := &YamlWalker{
		data: map[string]*YamlWalker{
			"first": {
				data: []*YamlWalker{
					{
						data: map[string]*YamlWalker{
							"level": {
								data: 1,
							},
							"value": {
								data: "abc",
							},
						},
						keys: []yamlKey{
							{name: "level"},
							{name: "value"},
						},
					},
					{
						data: map[string]*YamlWalker{
							"level": {
								data: 2,
							},
							"value": {
								data: "def",
							},
						},
						keys: []yamlKey{
							{name: "level"},
							{name: "value"},
						},
					},
				},
			},
			"second": {data: 2},
		},
		keys: []yamlKey{
			{name: "first"},
			{name: "second"},
		},
	}

	data := []struct {
		level int
		value string
	}{
		{
			level: 10,
			value: "ABC",
		},
		{
			level: 20,
			value: "DEF",
		},
	}

	iface, err := y.Get("first")
	suite.Assert().Nil(err)
	array, ok := iface.Value().([]*YamlWalker)
	suite.Assert().True(ok)
	suite.Assert().Equal(2, len(array))

	for i := range data {
		nl := NewYamlWalker()
		nl.Update(data[i].level)
		nv := NewYamlWalker()
		nv.Update(data[i].value)

		err := array[i].Set("level", nl)
		suite.Assert().Nil(err)
		err = array[i].Set("value", nv)
		suite.Assert().Nil(err)
	}

	m, ok := y.data.(map[string]*YamlWalker)
	suite.Assert().True(ok)
	a, ok := m["first"].data.([]*YamlWalker)
	suite.Assert().True(ok)
	suite.Assert().Equal(2, len(a))

	for i := range a {
		m, ok := a[i].data.(map[string]*YamlWalker)
		suite.Assert().True(ok)
		l, found := m["level"]
		suite.Assert().True(found)
		lvl, ok := l.data.(int)
		suite.Assert().True(ok)
		suite.Assert().Equal(data[i].level, lvl)
		v, found := m["value"]
		suite.Assert().True(found)
		val, ok := v.data.(string)
		suite.Assert().True(ok)
		suite.Assert().Equal(data[i].value, val)
	}
}

func (suite *YamlWalkerTestSuite) TestDelete() {
	getData := func() *YamlWalker {
		return &YamlWalker{
			data: map[string]*YamlWalker{
				"first": {
					data: map[string]*YamlWalker{
						"first-subitem": {
							data: 1,
						},
					},
					keys: []yamlKey{{name: "first-subitem"}},
				},
				"second": {
					data: map[string]*YamlWalker{
						"second-submap-1": {
							data: map[string]*YamlWalker{
								"second-subitem-1": {
									data: 1,
								},
							},
							keys: []yamlKey{{name: "second-subitem-1"}},
						},
						"second-submap-2": {
							data: map[string]*YamlWalker{
								"second-subitem-2-1": {
									data: 21,
								},
								"second-subitem-2-2": {
									data: 22,
								},
								"second-subitem-2-3": {
									data: 23,
								},
							},
							keys: []yamlKey{
								{name: "second-subitem-2-1"},
								{name: "second-subitem-2-2"},
								{name: "second-subitem-2-3"},
							},
						},
						"second-submap-3": {
							data: map[string]*YamlWalker{
								"second-subitem-3": {
									data: 2,
								},
							},
							keys: []yamlKey{{name: "second-subitem-3"}},
						},
					},
					keys: []yamlKey{
						{name: "second-submap-1"},
						{name: "second-submap-2"},
						{name: "second-submap-3"},
					},
				},
			},
			keys: []yamlKey{
				{name: "first"},
				{name: "second"},
			},
		}
	}

	y := getData()
	err := y.Delete("")
	suite.Assert().EqualError(err, ErrKeyMismatch.Error())

	y = getData()
	err = y.Delete("first")
	suite.Assert().Nil(err)
	m, found := y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().False(found)
	_, found = m["second"]
	suite.Assert().True(found)

	y = getData()
	err = y.Delete("second")
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["second"]
	suite.Assert().False(found)
	_, found = m["first"]
	suite.Assert().True(found)

	y = getData()
	err = y.Delete("second.second-submap-2")
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().True(found)
	d, found := m["second"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = d["second-submap-1"]
	suite.Assert().True(found)
	_, found = d["second-submap-2"]
	suite.Assert().False(found)
	_, found = d["second-submap-3"]
	suite.Assert().True(found)

	y = getData()
	err = y.Delete("second.second-submap-2.second-subitem-2-2")
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().True(found)
	d, found = m["second"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = d["second-submap-1"]
	suite.Assert().True(found)
	_, found = d["second-submap-2"]
	suite.Assert().True(found)
	_, found = d["second-submap-3"]
	suite.Assert().True(found)
	d, found = d["second-submap-2"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = d["second-subitem-2-1"]
	suite.Assert().True(found)
	_, found = d["second-subitem-2-2"]
	suite.Assert().False(found)
	_, found = d["second-subitem-2-3"]
	suite.Assert().True(found)

	y = getData()
	err = y.Delete("third.something")
	suite.Assert().EqualError(err, ErrNotFound.Error())
	err = y.Delete("first.something")
	suite.Assert().EqualError(err, ErrNotFound.Error())
	err = y.Delete("second.second-submap-2.second-subitem-2-2.something")
	suite.Assert().EqualError(err, ErrInvalidType.Error())
}

func (suite *YamlWalkerTestSuite) TestAppend() {
	getData := func() *YamlWalker {
		return &YamlWalker{
			data: map[string]*YamlWalker{
				"first": {
					data: map[string]*YamlWalker{
						"first-subitem": {
							data: 1,
						},
					},
					keys: []yamlKey{{name: "first-subitem"}},
				},
				"second": {
					data: map[string]*YamlWalker{
						"second-submap-1": {
							data: map[string]*YamlWalker{
								"second-subitem-1": {
									data: 1,
								},
							},
							keys: []yamlKey{{name: "second-subitem-1"}},
						},
						"second-submap-2": {
							data: map[string]*YamlWalker{
								"second-subitem-2-1": {
									data: 21,
								},
								"second-subitem-2-2": {
									data: 22,
								},
								"second-subitem-2-3": {
									data: 23,
								},
							},
							keys: []yamlKey{
								{name: "second-subitem-2-1"},
								{name: "second-subitem-2-2"},
								{name: "second-subitem-2-3"},
							},
						},
						"second-submap-3": {
							data: map[string]*YamlWalker{
								"second-subitem-3": {},
							},
							keys: []yamlKey{{name: "second-subitem-3"}},
						},
					},
					keys: []yamlKey{
						{name: "second-submap-1"},
						{name: "second-submap-2"},
						{name: "second-submap-3"},
					},
				},
			},
			keys: []yamlKey{
				{name: "first"},
				{name: "second"},
			},
		}
	}

	y := getData()
	n := &YamlWalker{
		data: 1,
	}
	err := y.Append("", n)
	suite.Assert().EqualError(err, ErrKeyMismatch.Error())

	y = getData()
	err = y.Append("third", n)
	suite.Assert().Nil(err)
	m, found := y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	suite.Assert().Equal(3, len(m))
	_, found = m["first"]
	suite.Assert().True(found)
	_, found = m["second"]
	suite.Assert().True(found)
	v, found := m["third"]
	suite.Assert().True(found)
	i, found := v.data.(int)
	suite.Assert().True(found)
	suite.Assert().Equal(1, i)

	y = getData()
	err = y.Append("first.appended", n)
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().True(found)
	_, found = m["second"]
	suite.Assert().True(found)
	d, found := m["first"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	v, found = d["appended"]
	suite.Assert().True(found)
	i, found = v.data.(int)
	suite.Assert().True(found)
	suite.Assert().Equal(1, i)

	y = getData()
	err = y.Append("second.second-submap-2.second-subitem-2-4", n)
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().True(found)
	_, found = m["second"]
	suite.Assert().True(found)
	d, found = m["second"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	d, found = d["second-submap-2"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	v, found = d["second-subitem-2-4"]
	suite.Assert().True(found)
	i, found = v.data.(int)
	suite.Assert().True(found)
	suite.Assert().Equal(1, i)

	err = y.Append("second.second-submap-3.second-subitem-3.appended", n, yaml.SingleQuotedStyle)
	suite.Assert().Nil(err)
	m, found = y.data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	_, found = m["first"]
	suite.Assert().True(found)
	_, found = m["second"]
	suite.Assert().True(found)
	d, found = m["second"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	d, found = d["second-submap-3"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	d, found = d["second-subitem-3"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	v, found = d["appended"]
	suite.Assert().True(found)
	i, found = v.data.(int)
	suite.Assert().True(found)
	suite.Assert().Equal(1, i)
	d, found = m["second"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	d, found = d["second-submap-3"].data.(map[string]*YamlWalker)
	suite.Assert().True(found)
	v, found = d["second-subitem-3"]
	suite.Assert().True(found)
	suite.Assert().Equal(1, len(v.keys))
	suite.Assert().Equal("appended", v.keys[0].name)
	suite.Assert().Equal(yaml.SingleQuotedStyle, v.keys[0].style)

	y = getData()
	err = y.Append("something.missing", n)
	suite.Assert().EqualError(err, ErrNotFound.Error())
	err = y.Append("second.second-submap-2", n)
	suite.Assert().EqualError(err, ErrDuplicateKey.Error())
	err = y.Append("second.second-submap-2.second-subitem-2-2.something", n)
	suite.Assert().EqualError(err, ErrInvalidType.Error())
}
