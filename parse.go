package runic

import (
	"fmt"
	"slices"
)

type parser struct {
	input           string
	tree            *Node
	lexer           *lexer
	error           string
	currentNode     *Node
	ctx             int
	tagDepth        int
	collectedTokens []token
}

func New() *parser {
	return &parser{}
}

func (p *parser) Parse(input string) *Node {
	p.lexer = lex(input)
	p.tree = &Node{Typ: nodeRoot, Val: ""}
	p.currentNode = p.tree
	p.collectedTokens = []token{}
	p.parseGlobal()
	return p.tree
}

func (p *parser) addNewNode(typ, val string) {
	newNode := &Node{
		Typ:    typ,
		Val:    val,
		parent: p.currentNode,
	}
	p.currentNode.Children = append(p.currentNode.Children, newNode)
	p.currentNode = newNode
}

func (p *parser) previousSibling() *Node {
	numOfChildren := len(p.currentNode.Children)
	if numOfChildren > 0 {
		return p.currentNode.Children[numOfChildren-1]
	}
	return &Node{}
}

func (p *parser) nextToken() {
	p.lexer.nextToken()
	p.collectedTokens = append(p.collectedTokens, p.lexer.token)
}

func (p *parser) returnNode() {
	p.currentNode = p.currentNode.parent
}

func (p *parser) isOneOf(t ...tokenType) bool {
	return slices.Contains(t, p.lexer.token.Typ)
}

func (p *parser) parseGlobal() {
	p.tagDepth = 0
	p.nextToken()
	for !p.isOneOf(typeEOF) {
		switch p.lexer.token.Typ {
		case typeHeading:
			p.parseHeading()
		case typeBulletpoint:
			p.parseList(0)
		default:
			p.parseParagraph()
		}
		p.nextToken()
	}
}

func (p *parser) parseParagraph() {
	p.addNewNode(nodeParagraph, "")
	p.parseRichText()
	p.returnNode()
}

func (p *parser) parseText() {
	previousSibling := p.previousSibling()
	// merge consecutive text nodes into one
	if previousSibling.Typ == nodeText {
		previousSibling.Val += " " + p.lexer.token.Val
		return
	}
	p.addNewNode(nodeText, p.lexer.token.Val)
	p.returnNode()
}

func (p *parser) parseRichText() {
	for !p.isOneOf(typeBulletpoint, typeTerminator, typeEOF) {
		switch p.lexer.token.Typ {
		case typeText:
			p.parseText()
		case typeTag:
			p.parseTag()
			if p.tagDepth > 0 && p.isOneOf(typeBulletpoint, typeTerminator) {
				return
			}
		case typeClosingSquare:
			if p.tagDepth > 0 {
				p.tagDepth--
				return
			}
		}
		p.nextToken()
	}
}

func (p *parser) parseHeading() {
	switch p.lexer.token.Val {
	case nodeHeadingOneValue:
		p.addNewNode(nodeHeadingOne, nodeHeadingOneValue)
	case nodeHeadingTwoValue:
		p.addNewNode(nodeHeadingTwo, nodeHeadingTwoValue)
	case nodeHeadingThreeValue:
		p.addNewNode(nodeHeadingThree, nodeHeadingThreeValue)
	case nodeHeadingFourValue:
		p.addNewNode(nodeHeadingFour, nodeHeadingFourValue)
	case nodeHeadingFiveValue:
		p.addNewNode(nodeHeadingFive, nodeHeadingFiveValue)
	case nodeHeadingSixValue:
		p.addNewNode(nodeHeadingSix, nodeHeadingSixValue)
	default:
		p.addNewNode(nodeError, fmt.Sprintf("%s: %s", errInvalidHeading, p.lexer.token.Val))
	}

	p.nextToken()
	p.parseRichText()
	p.returnNode()
}

func (p *parser) parseTag() {
	switch p.lexer.token.Val {
	case "bold":
		p.addNewNode(nodeBoldTag, "")
	case "italic":
		p.addNewNode(nodeItalicTag, "")
	default:
		p.addNewNode(nodeError, fmt.Sprintf("%s: %s", errInvalidTag, p.lexer.token.Val))
	}

	// skip over openSquare token
	p.nextToken()
	p.tagDepth++

	p.parseRichText()
	p.returnNode()
}

func (p *parser) parseList(currentListDepth int) {
	if p.isOneOf(typeBulletpoint) && getListItemDepth(p.lexer.token) < currentListDepth {
		p.returnNode()
		return
	}

	p.addNewNode(nodeList, "")

	for p.isOneOf(typeBulletpoint) {
		currentListDepth = getListItemDepth(p.lexer.token)

		p.nextToken()
		p.parseListItem()
		p.tagDepth = 0

		// bulletpoint is at a lower depth, create nested list
		if p.isOneOf(typeBulletpoint) && getListItemDepth(p.lexer.token) > currentListDepth {
			p.parseList(getListItemDepth(p.lexer.token))
		}

		// bulletpoint is a higher depth, return until no longer shallower
		if p.isOneOf(typeBulletpoint) && getListItemDepth(p.lexer.token) < currentListDepth {
			p.returnNode()
			return
		}
	}

	p.returnNode()
}

func (p *parser) parseListItem() {
	p.addNewNode(nodeListItem, "")
	p.parseRichText()
	p.returnNode()
}
