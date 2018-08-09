package htmlout

import (
	"fmt"
	"strings"
)

//Node is a generic node
type Node interface {
	//HTML produces an HTML element from the Node
	//the argument is the indentation
	HTML(indentation int) string
}

//HtmlNode is a node rapresenting an HTML element
type HtmlNode struct {
	//tag is the HTML tag
	tag      string
	parent   *HtmlNode
	children []Node
	attrs    map[string]string
	single   bool
}

//AddChildren adds a children to an HtmlNode
func (n *HtmlNode) AddChildren(children *HtmlNode) {
	n.AddChildrenNode(children)
	children.parent = n
}

//AddChildrenNode adds a generic children to the node
func (n *HtmlNode) AddChildrenNode(children Node) {
	if n.children == nil {
		n.children = []Node{}
	}

	n.children = append(n.children, children)
}

//HTML generate HTML from the node
func (n *HtmlNode) HTML(indentation int) string {
	var tags string
	for k, v := range n.attrs {
		tags += fmt.Sprintf(" %s=\"%s\"", k, v)
	}

	if n.single {
		if len(n.children) != 0 {
			panic("Single Html Element cannot have children")
		}
		return indent(indentation) + fmt.Sprintf("<%s%s/>", n.tag, tags)
	}

	var text string
	if n.tag != "" {
		text = indent(indentation) + fmt.Sprintf("<%s%s>\n", n.tag, tags)
	}
	for _, c := range n.children {
		text += c.HTML(indentation+1) + "\n"
	}
	if n.tag != "" {
		text += indent(indentation) + fmt.Sprintf("</%s>", n.tag)
	}
	return text
}

//TextNode is a node containing string
type TextNode string

//HTML return the content of the node
//Correctly indented
func (n TextNode) HTML(indentation int) string {
	return indent(indentation) + strings.Replace(string(n), "\n", "\n"+indent(indentation), -1)
}

func indent(indentation int) string {
	return strings.Repeat("\t", indentation)
}

//CodeNode is a TextNode with doen't adds indentation
type CodeNode TextNode

//HTML return the content of the CodeNode but without indentation
func (n CodeNode) HTML(_ int) string {
	nn := TextNode(n)
	return nn.HTML(0)
}
