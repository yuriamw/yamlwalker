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

	fmt.Printf("\n=========================\n\n")

	fmt.Printf("W:%T:%+v\n", yw, yw)
	fmt.Printf("V:%T:%+v\n", yw.Value(), yw.Value())

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
			fmt.Printf("S:%T:%+v\n", url, url)
		}
	}

	fmt.Printf("G:%+v\n", yw.GetValue("another.something.interresting"))

	out, err := yaml.Marshal(yw)
	if err != nil {
		panic(err)
	}
	fmt.Printf("---\n%+v...\n", string(out))
}
