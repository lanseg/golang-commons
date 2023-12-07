package almosthtml

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lanseg/golang-commons/collections"
	"github.com/lanseg/golang-commons/optional"
)

var (
	selfClosingTags = collections.NewSet([]string{
		"area", "base", "br", "col", "embed",
		"hr", "img", "input", "link", "meta",
		"param", "source", "track", "wbr",
		"!DOCTYPE", "#text",
	})
	emptyParam = &optional.Nothing[string]{}
)

func getChildren(n *Node) []*Node {
	if n.Children == nil {
		return []*Node{}
	}
	return n.Children
}

// Node is something like a html element: tag or text data.
type Node struct {
	Name     string
	Raw      string
	Params   map[string]string
	Children []*Node
}

func (n *Node) iterateChildren() collections.Iterator[*Node] {
	return collections.IterateTree(n, collections.DepthFirst, func(node *Node) []*Node {
		if node.Children == nil {
			return []*Node{}
		}
		return node.Children
	})
}

// InnerHTML generates a html representation of the node and it's contents.
func (n *Node) InnerHTML() string {
	result := strings.Builder{}
	n.iterateChildren().ForEachRemaining(func(node *Node) bool {
		result.WriteString(node.Raw)
		return false
	})
	return result.String()
}

// GetElementsByTagAndClass finds all tags with given name, which have all the given classes.
func (n *Node) GetElementsByTagAndClass(tag string, classes ...string) []*Node {
	return n.iterateChildren().Filter(func(node *Node) bool {
		data, ok := node.Params["class"]
		if tag != node.Name || (!ok && len(classes) > 0) {
			return false
		}
		return collections.NewSet(strings.Split(data, " ")).ContainsAll(classes)
	}).Collect()
}

// GetElementsByTags finds all elements with given tag names.
func (n *Node) GetElementsByTags(tags ...string) []*Node {
	tagSet := collections.NewSet(tags)
	return n.iterateChildren().Filter(func(node *Node) bool {
		return tagSet.Contains(node.Name)
	}).Collect()
}

// GetAttribute returns attribute value or Nothing optional if there is no such attribute.
func (n *Node) GetAttribute(attr string) optional.Optional[string] {
	if value, ok := n.Params[attr]; ok {
		return optional.Of(value)
	}
	return emptyParam
}

func (n *Node) String() string {
	params := []string{}
	for k, v := range n.Params {
		params = append(params, fmt.Sprintf("%q=%q", k, v))
	}
	sort.Strings(params)
	return fmt.Sprintf("Node { %q [%s] %q}", n.Name, strings.Join(params, ", "), n.Raw)
}

func newNode(name string) *Node {
	return &Node{
		Name:     name,
		Params:   map[string]string{},
		Children: []*Node{},
	}
}

func newDataNode(data string) *Node {
	return &Node{
		Name:     "#text",
		Params:   map[string]string{},
		Raw:      data,
		Children: []*Node{},
	}
}

// ParseHTML converts a html text into a Node tree.
func ParseHTML(doc string) (*Node, error) {
	nodes := []*Node{}
	runes := []rune(doc)
	for _, t := range tokenize(doc) {
		raw := string(runes[t.start:t.end])
		switch t.tokenType {
		case dataState:
			nodes = append(nodes, newDataNode(raw))
		case tagScriptState:
			nodes = append(nodes, newDataNode(raw))
		case tagState:
			node := newNode(t.name)
			node.Raw = raw
			for i, value := range t.values {
				if value == "=" {
					continue
				}
				if i > 1 && t.values[i-1] == "=" {
					node.Params[t.values[i-2]] = value
				} else {
					node.Params[value] = ""
				}
			}
			nodes = append(nodes, node)
		}
	}

	root := newNode("#root")
	stack := []*Node{root}
	for _, n := range nodes {
		if n.Name[0] == '/' {
			stack = stack[:len(stack)-1]
			continue
		}
		parent := stack[len(stack)-1]
		parent.Children = append(parent.Children, n)
		if !selfClosingTags.Contains(n.Name) {
			stack = append(stack, n)
		}
	}
	return root, nil
}

func dump(n *Node, prefix string) {
	fmt.Printf("%s %s\n", prefix, n)
	for _, nn := range n.Children {
		dump(nn, prefix+"  ")
	}
}

// GetTitle returns value of the title tag of this document
func GetTitle(s string) string {
	nodes, _ := ParseHTML(s)
	titles := nodes.GetElementsByTagAndClass("title")
	if len(titles) == 0 || len(titles[0].Children) == 0 {
		return ""
	}
	return titles[0].Children[0].Raw
}

// StripTags removes all tags and keeps only text values
func StripTags(s string) string {
	nodes, _ := ParseHTML(s)
	result := strings.Builder{}
	nodes.iterateChildren().ForEachRemaining(func(n *Node) bool {
		if n.Name != "#text" || n.Raw == "" {
			return false
		}
		result.WriteString(n.Raw)
		return false
	})
	return result.String()
}
