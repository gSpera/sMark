package outTelegraph

import tgraph "github.com/toby3d/telegraph"

//utils.go manages some utilities functions for creating tags.
//They are mostly create for readbility

func createTitle(txt string) tgraph.NodeElement {
	return tgraph.NodeElement{Tag: "h3", Children: []tgraph.Node{txt}}
}

func createSubtitle(txt string) tgraph.NodeElement {
	return tgraph.NodeElement{Tag: "h4", Children: []tgraph.Node{txt}}
}

func createBold(txt string) tgraph.NodeElement {
	return tgraph.NodeElement{Tag: "b", Children: []tgraph.Node{txt}}
}

func createItalic(txt string) tgraph.NodeElement {
	return tgraph.NodeElement{Tag: "i", Children: []tgraph.Node{txt}}
}

func createTag(tag string) tgraph.NodeElement {
	return tgraph.NodeElement{Tag: tag}
}
