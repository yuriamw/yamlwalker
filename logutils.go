package yamlwalker

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	enableLogs = false // used for development
)

func log(str string) {
	if enableLogs {
		fmt.Printf("%s", str)
	}
}

func decodeKind(kind yaml.Kind) string {
	switch kind {
	case yaml.DocumentNode:
		return "DocumentNode"
	case yaml.SequenceNode:
		return "SequenceNode"
	case yaml.MappingNode:
		return "MappingNode"
	case yaml.ScalarNode:
		return "ScalarNode"
	case yaml.AliasNode:
		return "AliasNode"
	}
	return fmt.Sprintf("%d", kind)
}

func decodeStyle(style yaml.Style) string {
	switch style {
	case yaml.TaggedStyle:
		return "TaggedStyle"
	case yaml.DoubleQuotedStyle:
		return "DoubleQuotedStyle"
	case yaml.SingleQuotedStyle:
		return "SingleQuotedStyle"
	case yaml.LiteralStyle:
		return "LiteralStyle"
	case yaml.FoldedStyle:
		return "FoldedStyle"
	case yaml.FlowStyle:
		return "FlowStyle"
	}
	return fmt.Sprintf("%d", style)
}

func printNodeContent(node *yaml.Node) string {
	str := fmt.Sprintf("%s\n", printNode(node))
	for _, c := range node.Content {
		str += fmt.Sprintf(">  %s\n", printNode(c))
	}
	return str
}

func printNode(n *yaml.Node) string {
	str := fmt.Sprintf("&{Kind:%d=%s Style:%d=%s Tag:%s Value:%v Anchor:%v Alias:%v Content:[", n.Kind, decodeKind(n.Kind), n.Style, decodeStyle(n.Style), n.Tag, n.Value, n.Anchor, n.Alias)
	for i, c := range n.Content {
		s := ""
		if i > 0 {
			s += " "
		}
		s += fmt.Sprintf("%p", c)
		str += s
	}
	str += "]}"

	return str
}
