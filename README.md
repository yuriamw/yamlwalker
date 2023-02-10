# yamlwalker

Get, set and walk over generic YAML without prior knowlage of the structure of YAML.

Get/Set/Walk methods take dot separated path to the parameter.
The following document fragment

```yaml
object:
  name:
    param: value
another:
  something:
    interresting: thing
```

should be accessed by path as:

```"object.name.param"```
and
```"another.something.interresting"```

# Usage example

## Build a new yaml from scratch

```golang
package main

import (
	"fmt"

	"github.com/yuriamw/yamlwalker"
	"gopkg.in/yaml.v3"
)

func main() {
	walker := yamlwalker.NewYamlWalker()

	first := yamlwalker.NewYamlWalker()
	first.Update(1)

	second := yamlwalker.NewYamlWalker(yaml.SingleQuotedStyle)
	second.Update("Value for second")

	third := yamlwalker.NewYamlWalker()
	thirdChild1 := yamlwalker.NewYamlWalker()
	thirdChild1.Update("Child 1 value")
	third.Append("child-1", thirdChild1)
	thirdChild2 := yamlwalker.NewYamlWalker(yaml.DoubleQuotedStyle)
	thirdChild2.Update("Child 2 value")
	third.Append("child-2", thirdChild2)

	fourth := yamlwalker.NewYamlWalker()
	fourth.Update(make([]*yamlwalker.YamlWalker, 0))
	for i := 0; i < 10; i++ {
		fourthChild := yamlwalker.NewYamlWalker()
		fourthChild.Update(i)
		err := fourth.Insert("", i, fourthChild)
		if err != nil {
			panic(err)
		}
	}

	walker.Append("first", first)
	walker.Append("second", second)
	walker.Append("third", third)
	walker.Append("fourth", fourth)

	data, err := yaml.Marshal(walker)
	if err != nil {
		panic("marshal failed")
	}

	fmt.Printf("---\n%v...\n", string(data))

	walker.SetValue("third.child-1", "Child 1 value changed")

	data, err = yaml.Marshal(walker)
	if err != nil {
		panic("marshal failed")
	}

	fmt.Printf("---\n%v...\n", string(data))
}
```

Prints:

```
---
first: 1
second: 'Value for second'
third:
    child-1: Child 1 value
    child-2: "Child 2 value"
fourth:
    - 0
    - 1
    - 2
    - 3
    - 4
    - 5
    - 6
    - 7
    - 8
    - 9
...
---
first: 1
second: 'Value for second'
third:
    child-1: Child 1 value changed
    child-2: "Child 2 value"
fourth:
    - 0
    - 1
    - 2
    - 3
    - 4
    - 5
    - 6
    - 7
    - 8
    - 9
...
```

## Read existing file and modify some values

```golang
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yuriamw/yamlwalker"

	"gopkg.in/yaml.v3"
)

var (
	fileName *string
)

func init() {
	fileName = flag.String("f", "test_data/simple.yaml", "File name")

	flag.Parse()
}

func main() {
	data, err := os.ReadFile(*fileName)
	if err != nil {
		panic(err)
	}
	yw := yamlwalker.NewYamlWalker()
	err = yaml.Unmarshal(data, yw)
	if err != nil {
		panic(err)
	}

	next := true
	for next {
		next = false
		sections, ok := yw.Value().(map[string]*yamlwalker.YamlWalker)
		if !ok {
			continue
		}

		srvM, found := sections["servers"]
		if !found {
			continue
		}

		srv, ok := srvM.Value().([]*yamlwalker.YamlWalker)
		if !ok {
			continue
		}
		for _, v := range srv {
			m := v.Value().(map[string]*yamlwalker.YamlWalker)
			url := m["url"].Value().(string)
			fmt.Printf("S:%+v\n", url)
		}
	}

	fmt.Printf("G:%+v\n", yw.GetValue("openapi"))
	fmt.Printf("G:%+v\n", yw.GetValue("info.description"))
	fmt.Printf("G:%+v\n", yw.GetValue("info.contact.name"))
	fmt.Printf("G:%+v\n", yw.GetValue("not exists"))

	url := yamlwalker.NewYamlWalker()
	url.Update("http://example.com")
	err = yw.Insert("servers", 0, url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nUpdate node '%s' value\n", "info.contact.name")
	yw.SetValue("info.contact.name", "My Cool Company")
	fmt.Printf("info.contact.name:%+v\n", yw.GetValue("info.contact.name"))

	out, err := yaml.Marshal(yw)
	if err != nil {
		panic(err)
	}
	fmt.Printf("---\n%+v...\n", string(out))
}
```

Prints:

```
S:{protocol}://localhost:{port}/api/v1.0
G:3.0.2
G:The project is simple example of OpanAPI spec in yaml
G:No real company
G:<nil>

Update node 'info.contact.name' value
info.contact.name:My Cool Company
---
openapi: '3.0.2'
info:
    title: Simple project
    description: The project is simple example of OpanAPI spec in yaml
    version: 3.55.144
    termsOfService: http://some.strange.com/legal-notice/
    contact:
        name: My Cool Company
        url: http://some.strange.com
        email: info@some.strange.com
    license:
        name: License
        url: http://some.strange.com/legal-notice/license.txt
    x-ExtensionBool: true
    x-String: api
servers:
    - http://example.com
    - url: '{protocol}://localhost:{port}/api/v1.0'
      variables:
        protocol:
            enum:
                - http
            default: http
        port:
            default: '8080'
...
```
