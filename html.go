package runic

import (
	"fmt"
	"regexp"
	"strings"
)

type htmlCtxType int

const (
	htmlCtxNone htmlCtxType = iota
	htmlCtxParagraph
)

func (p *parser) Html(input string) string {
	tree := p.Parse(input)
	s := ""
	return p.toHtml(tree, &s, htmlCtxNone)
}

func (p *parser) toHtml(currentNode *Node, htmlString *string, htmlCtx htmlCtxType) string {
	for _, child := range currentNode.Children {
		switch child.Typ {
		case nodeError:
			*htmlString += "<span class='error'>"
		case nodeHeadingOne:
			*htmlString += "<h1>"
		case nodeHeadingTwo:
			*htmlString += "<h2>"
		case nodeHeadingThree:
			*htmlString += "<h3>"
		case nodeHeadingFour:
			*htmlString += "<h4>"
		case nodeHeadingFive:
			*htmlString += "<h5>"
		case nodeHeadingSix:
			*htmlString += "<h6>"
		case nodeParagraph:
			*htmlString += "<p>"
			htmlCtx = htmlCtxParagraph
		case nodeBoldTag:
			*htmlString += "<b>"
		case nodeItalicTag:
			*htmlString += "<em>"
		case nodeList:
			*htmlString += "<ul>"
		case nodeListItem:
			*htmlString += "<li>"
		}

		if child.Typ == nodeText {
			*htmlString += child.Val + " "
		}

		if len(child.Children) > 0 {
			p.toHtml(child, htmlString, htmlCtx)
			*htmlString = strings.TrimSpace(*htmlString)
		}

		switch child.Typ {
		case nodeError:
			*htmlString += "</span>"
		case nodeHeadingOne:
			*htmlString += "</h1>"
		case nodeHeadingTwo:
			*htmlString += "</h2>"
		case nodeHeadingThree:
			*htmlString += "</h3>"
		case nodeHeadingFour:
			*htmlString += "</h4>"
		case nodeHeadingFive:
			*htmlString += "</h5>"
		case nodeHeadingSix:
			*htmlString += "</h6>"
		case nodeParagraph:
			*htmlString += "</p>"
			htmlCtx = htmlCtxNone
		case nodeBoldTag:
			*htmlString += "</b> "
		case nodeItalicTag:
			*htmlString += "</em> "
		case nodeList:
			*htmlString += "</ul>"
		case nodeListItem:
			*htmlString += "</li>"
		}
	}

	return *htmlString
}

func (p *parser) htmlSanitiseSlice(start, end int) (s string) {
	re := regexp.MustCompile("^\\s|\\s\\s+|\\s$")
	s = strings.ReplaceAll(p.input[start:end], "\n", "<br>")
	s = re.ReplaceAllStringFunc(s, func(s string) string {
		return strings.Repeat("&nbsp;", len(s))
	})
	return
}

func (p *parser) HighlightText(input string) (highlightedText string) {
	p.input = input
	lexer := lex(input)

	var prevToken token

	for lexer.nextToken() {
		if prevToken.Typ == typeNone {
			prevToken = lexer.token
			prevToken.Pos = 0
			continue
		}

		start := prevToken.Pos
		end := lexer.token.Pos

		switch prevToken.Typ {
		case typeText:
			highlightedText += fmt.Sprintf(`<span class="runic__text">%s</span>`, p.htmlSanitiseSlice(start, end))
		case typeHeading:
			highlightedText += fmt.Sprintf(`<span class="runic__heading">%s</span>`, p.htmlSanitiseSlice(start, end))
		case typeTag:
			highlightedText += fmt.Sprintf(`<span class="runic__tag">%s</span>`, p.htmlSanitiseSlice(start, end))
		case typeOpeningSquare:
			highlightedText += fmt.Sprintf(`<span class="runic__osq">%s</span>`, p.htmlSanitiseSlice(start, end))
		case typeClosingSquare:
			highlightedText += fmt.Sprintf(`<span class="runic__csq">%s</span>`, p.htmlSanitiseSlice(start, end))
		case typeBulletpoint:
			highlightedText += fmt.Sprintf(`<span class="runic__bulletpoint">%s</span>`, p.htmlSanitiseSlice(start, end))
		default:
			highlightedText += p.htmlSanitiseSlice(start, end)
		}

		prevToken = lexer.token
	}

	return
}

type editorData struct {
	Html          string `json:"html"`
	HighlightText string `json:"highlightText"`
}

func (p *parser) EditorData(input string) (*editorData, error) {
	return &editorData{
		Html:          p.Html(input),
		HighlightText: p.HighlightText(input),
	}, nil
}
