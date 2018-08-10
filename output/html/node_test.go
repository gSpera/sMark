package htmlout

import (
	"eNote/utils"
	"fmt"
	"testing"
)

func TestAddChildren(t *testing.T) {
	tm := []struct {
		name     string
		n        *HTMLNode
		children *HTMLNode
		output   *HTMLNode
	}{
		{
			name: "simple",
			n: &HTMLNode{
				tag: "h1",
			},
			children: &HTMLNode{
				tag: "h2",
			},

			output: &HTMLNode{
				tag: "h1",
				children: []Node{
					&HTMLNode{
						tag: "h2",
					},
				},
			},
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.n
			in.AddChildren(tt.children)
			if err := checkNode(in, tt.output); err != nil {
				t.Error(err)
			}
		})
	}
}

func checkNode(n1, n2 *HTMLNode) error {
	if n1.tag != n2.tag {
		return fmt.Errorf("Tags are different: %v != %v", n1.tag, n2.tag)
	}
	if len(n1.children) != len(n2.children) {
		return fmt.Errorf("Len of childrens are different: %v != %v", len(n1.children), len(n2.children))
	}

	return nil
}

func TestTextNodeHTML(t *testing.T) {
	tm := []struct {
		name        string
		n           TextNode
		indentation int
		output      string
	}{
		{
			name:        "simple",
			n:           TextNode("test"),
			indentation: 0,
			output:      "test",
		},
		{
			name:        "indentation",
			n:           TextNode("test"),
			indentation: 1,
			output:      "\ttest",
		},
		{
			name:        "text multi-line indentation",
			n:           TextNode("test\ntest"),
			indentation: 1,
			output:      "\ttest\n\ttest",
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.n.HTML(tt.indentation)
			if tt.output != out {
				t.Errorf("Expected:\n%v\ngot:\n%v", tt.output, out)
			}
		})
	}
}

func TestNodeHTML(t *testing.T) {
	tm := []struct {
		name   string
		n      Node
		output string
	}{
		{
			name: "simple",
			n: &HTMLNode{
				tag: "h1",
			},
			output: "<h1>\n</h1>",
		},
		{
			name: "children",
			n: func() *HTMLNode {
				root := &HTMLNode{
					tag: "h1",
				}
				root.AddChildren(&HTMLNode{
					tag: "h2",
				})

				return root
			}(),
			output: `<h1>
	<h2>
	</h2>
</h1>`,
		},
		{
			name: "text",
			n: &HTMLNode{
				tag: "h1",
				children: []Node{
					TextNode("test"),
				},
			},
			output: `<h1>
	test
</h1>`,
		},
		{
			name: "attrs",
			n: &HTMLNode{
				tag: "h1",
				attrs: map[string]string{
					"test": "value",
				},
			},
			output: `<h1 test="value">
</h1>`,
		},
	}

	for _, tt := range tm {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.n.HTML(0)
			if tt.output != out {
				t.Errorf("Expected:\n%s;got:\n%s", tt.output, out)
			}
		})
	}
}

func TestToString(t *testing.T) {
	ToString(nil, eNote.Options{})
}
