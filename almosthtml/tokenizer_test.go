package almosthtml

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {

	for _, tc := range []struct {
		name string
		html string
		want []*token
	}{
		{
			name: "Empty html",
			html: "",
			want: []*token{},
		},
		{
			name: "Plain text",
			html: "Hello world",
			want: []*token{
				{
					values:    []string{"Hello world"},
					tokenType: dataState,
					start:     0,
					end:       11,
				},
			},
		},
		{
			name: "Tag without text",
			html: "<html>",
			want: []*token{
				{
					name:      "html",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       5,
				},
			},
		},
		{
			name: "Two consecutive tags",
			html: "<html><html>",
			want: []*token{
				{
					name:      "html",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       5,
				},
				{
					name:      "html",
					values:    []string{},
					tokenType: tagState,
					start:     7,
					end:       11,
				},
			},
		},
		{
			name: "Tag with parameters without values",
			html: "<tag parameter1 >",
			want: []*token{
				{
					name:      "tag",
					values:    []string{"parameter1"},
					tokenType: tagState,
					start:     1,
					end:       16,
				},
			},
		},
		{
			name: "Tag with quoted and unquoted parameters",
			html: "<tag parameter1 \"parameter2 with space\" 'parameter 3' >",
			want: []*token{
				{
					name:      "tag",
					values:    []string{"parameter1", "parameter2 with space", "parameter 3"},
					tokenType: tagState,
					start:     1,
					end:       54,
				},
			},
		},
		{
			name: "Tag with quoted, unquoted values and parameters",
			html: "<tag a b=c, \"d a\"=e f =\"g h\" 'param a'  =  \"param b\">",
			want: []*token{
				{
					name: "tag",
					values: []string{
						"a", "b", "=", "c,", "d a", "=", "e", "f", "=", "g h", "param a", "=", "param b",
					},
					tokenType: tagState,
					start:     1,
					end:       52,
				},
			},
		},
		{
			name: "script with text tags",
			html: "<script>function doit(){\n  console.log(\"<tag>\");\n}</script>",
			want: []*token{
				{
					name:      "script",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       7,
				},
				{
					values:    []string{"function doit(){\n  console.log(\"<tag>\");\n}</script"},
					tokenType: tagScriptState,
					start:     7,
					end:       58,
				},
			},
		},
		{
			name: "script tag without internal tags",
			html: "<script> function somefunc() {\n  console.log(\"HERE\");\n} </script>",
			want: []*token{
				{
					name:      "script",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       7,
				},
				{
					values:    []string{" function somefunc() {\n  console.log(\"HERE\");\n} </script"},
					tokenType: tagScriptState,
					start:     7,
					end:       64,
				},
			},
		},
		{
			name: "Tag comment tag again",
			html: "<tag><!-- a comment --></tag>",
			want: []*token{
				{
					name:      "tag",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       4,
				},
				{
					name:      "!--",
					values:    []string{" a comment --"},
					tokenType: tagCommentState,
					start:     6,
					end:       22,
				},
				{
					name:      "/tag",
					values:    []string{},
					tokenType: tagState,
					start:     24,
					end:       28,
				},
			},
		},
		{
			name: "Only comment",
			html: "<!-- a comment here <andtag><whatever></tag> -->",
			want: []*token{
				{
					name:      "!--",
					values:    []string{" a comment here <andtag><whatever></tag> --"},
					tokenType: tagCommentState,
					start:     1,
					end:       47,
				},
			},
		},
		{
			name: "Comment as value",
			html: "<tag key=\"<!-- not a comment -->\">",
			want: []*token{
				{
					name:      "tag",
					values:    []string{"key", "=", "<!-- not a comment -->"},
					tokenType: tagState,
					start:     1,
					end:       33,
				},
			},
		},
		{
			name: "Text without tags",
			html: "Some text without any tags",
			want: []*token{
				{
					values:    []string{"Some text without any tags"},
					tokenType: dataState,
					start:     0,
					end:       26,
				},
			},
		},
		{
			name: "doc starts with tag ends with tag",
			html: "<html> some text </html>",
			want: []*token{
				{
					name:      "html",
					values:    []string{},
					tokenType: tagState,
					start:     1,
					end:       5,
				},
				{
					values:    []string{" some text "},
					tokenType: dataState,
					start:     5,
					end:       18,
				},
				{
					name:      "/html",
					values:    []string{},
					tokenType: tagState,
					start:     18,
					end:       23,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tokens := tokenize(tc.html)
			if !reflect.DeepEqual(tokens, tc.want) {
				t.Errorf("Result tokens are not equal to expected:\nActual  : %s\nExpected: %s", tokens, tc.want)
			}
		})
	}
}
