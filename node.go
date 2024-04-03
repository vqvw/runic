package runic

import (
	"slices"
)

type Node struct {
	Typ      string  `json:"type"`
	Val      string  `json:"value,omitempty"`
	Children []*Node `json:"children,omitempty"`
	parent   *Node
}

const INDENT_WIDTH = 2

const (
	nodeRoot         = "Root"
	nodeError        = "ERROR"
	nodeHeadingOne   = "HeadingOne"
	nodeHeadingTwo   = "HeadingTwo"
	nodeHeadingThree = "HeadingThree"
	nodeHeadingFour  = "HeadingFour"
	nodeHeadingFive  = "HeadingFive"
	nodeHeadingSix   = "HeadingSix"
	nodeParagraph    = "Paragraph"
	nodeText         = "Text"
	nodeBoldTag      = "BoldTag"
	nodeItalicTag    = "ItalicTag"
	nodeList         = "List"
	nodeListItem     = "ListItem"
)

const (
	nodeHeadingOneValue   = "."
	nodeHeadingTwoValue   = ":"
	nodeHeadingThreeValue = ":."
	nodeHeadingFourValue  = "::"
	nodeHeadingFiveValue  = "::."
	nodeHeadingSixValue   = ":::"
)

var (
	errInvalidTag     = "Invalid tag name"
	errInvalidHeading = "Invalid heading value"
)

func isOneOf(nodeType string, nodeTypes ...string) bool {
	return slices.Contains(nodeTypes, nodeType)
}

func getListItemDepth(listItem token) int {
	return listItem.indent / INDENT_WIDTH
}
