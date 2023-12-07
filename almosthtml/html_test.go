package almosthtml

import (
	_ "os"
	"reflect"
	"testing"
)

func makeNode(name string, raw string, children ...*Node) *Node {
	if children == nil {
		children = []*Node{}
	}
	return &Node{
		Name:     name,
		Raw:      raw,
		Params:   map[string]string{},
		Children: children,
	}
}

func makeParamNode(name string, raw string, k string, v string) *Node {
	return &Node{
		Name: name,
		Raw:  raw,
		Params: map[string]string{
			k: v,
		},
		Children: []*Node{},
	}
}

func makeText(raw string) *Node {
	return &Node{
		Name:     "#text",
		Raw:      raw,
		Params:   map[string]string{},
		Children: []*Node{},
	}
}

func TestRegressions(t *testing.T) {
}
func TestHTML(t *testing.T) {

	for _, tc := range []struct {
		name string
		html string
		want *Node
	}{
		{
			name: "Plain text only",
			html: "Hello world",
			want: makeNode("#root", "", makeText("Hello world")),
		},
		{
			name: "Simple tag no nested",
			html: "<node>Whatever</node>",
			want: makeNode("#root", "",
				makeNode("node", "node",
					makeText("Whatever"),
				),
			),
		},
		{
			name: "Simple tag with parameters",
			html: "<tag key=value></tag><tag key=\"Value quoted\"></tag>",
			want: makeNode("#root", "",
				makeParamNode("tag", "tag key=value", "key", "value"),
				makeParamNode("tag", "tag key=\"Value quoted\"", "key", "Value quoted")),
		},
		{
			name: "Comments are ignored",
			html: "<html><!-- Hello <world>intag</world>",
			want: makeNode("#root", "",
				makeNode("html", "html"),
			),
		},
		{
			name: "Script contents are set as data",
			html: "<html><script>function () {\n  console.log('<a tag></tag>');\n} </script>",
			want: makeNode("#root", "",
				makeNode("html", "html",
					makeNode("script", "script",
						makeText("function () {\n  console.log('<a tag></tag>');\n} </script")))),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node, err := ParseHTML(string(tc.html))
			if err != nil {
				t.Errorf("Error while parsing html: %s", err)
			}
			if !reflect.DeepEqual(tc.want, node) {
				t.Errorf("Result nodes are not equal to expected:\nActual  : %s\nExpected: %s", node, tc.want)
				dump(tc.want, "")
				dump(node, "")
			}
		})
	}
}
