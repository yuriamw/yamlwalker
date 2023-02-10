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

	fmt.Printf("\nUpdate node '%s' value\n", "info.contact.name")
	yw.SetValue("info.contact.name", "My Cool Company")
	fmt.Printf("info.contact.name:%+v\n", yw.GetValue("info.contact.name"))

	out, err := yaml.Marshal(yw)
	if err != nil {
		panic(err)
	}
	fmt.Printf("---\n%+v...\n", string(out))
}
