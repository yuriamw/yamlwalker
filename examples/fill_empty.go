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

	fmt.Printf("--\n%v--\n", string(data))

	walker.SetValue("third.child-1", "Child 1 value changed")

	data, err = yaml.Marshal(walker)
	if err != nil {
		panic("marshal failed")
	}

	fmt.Printf("--\n%v--\n", string(data))
}
