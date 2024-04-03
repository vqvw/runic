package runic

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// token represents a single lexeme returned from the lexer
type token struct {
	Typ    tokenType `json:"type"` // type of token
	Val    string    `json:"val"`  // characters comprising the token
	Line   int       `json:"line"` // line where the token was found
	Pos    int       `json:"pos"`  // position in the input text where the token was found
	indent int       // number of spaces appearing before the token
}

func (t *token) String() string {
	var tokenTypeString string
	switch t.Typ {
	case typeNone:
		tokenTypeString = "typeNone"
	case typeEOF:
		tokenTypeString = "typeEOF"
	case typeText:
		tokenTypeString = "typeText"
	case typeTerminator:
		tokenTypeString = "typeTerminator"
	case typeHeading:
		tokenTypeString = "typeHeading"
	case typeTag:
		tokenTypeString = "typeTag"
	case typeOpeningSquare:
		tokenTypeString = "typeOpeningSquare"
	case typeClosingSquare:
		tokenTypeString = "typeClosingSquare"
	case typeBulletpoint:
		tokenTypeString = "typeBulletpoint"
	}
	var indent string
	if t.indent > 0 {
		indent = ">" + strconv.Itoa(t.indent)
	}
	return fmt.Sprintf("[%d:%d:%s:'%s'%s]", t.Line, t.Pos, tokenTypeString, t.Val, indent)
}

type tokenType int

const (
	typeNone = iota
	typeEOF
	typeText
	typeTerminator
	typeHeading
	typeTag
	typeOpeningSquare
	typeClosingSquare
	typeBulletpoint
)

const (
	eof               = -1
	void              = -2
	charNewline       = '\n'
	charDot           = '.'
	charColon         = ':'
	charOpeningSquare = '['
	charClosingSquare = ']'
	charBackslash     = '\\'
	charHyphen        = '-'
)

// lexer represents the state machine processing the input text
type lexer struct {
	input             string  // input string containing markup
	line              int     // current line number
	pos               int     // current position in the input text
	token             token   // current token
	char              rune    // current character
	lexNext           func()  // next lex function
	skippedNewlines   int     // number of newlines skipped during `skipSpace`
	skippedSpace      int     // number of spaces skipped during `skipSpace`
	tag               string  // current tag accumulated (run of `unicode.isLetter` chars)
	continuousNewline bool    // don't treat a single newline as a terminator
	ctx               ctxType // current context
}

type ctxType int

const (
	ctxList ctxType = iota + 1
)

// lex returns a lexer, initialised to process the given input text
func lex(input string) *lexer {
	l := &lexer{input: input, line: 1}
	l.lexNext = l.lexGlobal
	return l
}

// collectTokens runs the lexer and accumulates the lexed tokens to return,
// used by the test suite
func collectTokens(input string) (tokens []token) {
	lexer := lex(input)
	for lexer.nextToken() {
		tokens = append(tokens, lexer.token)
	}
	return
}

// nextToken progresses the lexer by a single token (or to eof). it returns
// true while the lexer is running. this is called by the parser in order to
// read a new l.token value
func (l *lexer) nextToken() bool {
	if l.lexNext != nil {
		l.lexNext()
		if l.token.Typ == typeNone {
			return l.nextToken()
		}
		return true
	}
	return false
}

// mkToken produces a token of the provided type and value, initialised with
// the current line adnd pos values
func (l *lexer) mkToken(typ tokenType, val string) token {
	return token{Typ: typ, Val: val, Line: l.line, Pos: l.pos}
}

// resetToken sets the current token to an empty token
func (l *lexer) resetToken() {
	l.token = token{}
	l.continuousNewline = false
	l.ctx = 0
}

// next progresses the lexer by a single character
func (l *lexer) next() {
	if l.pos >= len(l.input) || len(l.input) == 0 {
		l.char = eof
		return
	}
	char, byteWidth := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += byteWidth
	if char == charNewline {
		l.line++
		l.skippedSpace = 0
	}
	l.char = char
	l.addToTag()
	l.skipSpace()
}

// backup is used to reverse the action of `l.next`
func (l *lexer) backup() {
	if l.pos == 0 {
		l.char = void
		return
	}
	if l.char == charNewline {
		l.line--
	}
	char, byteWidth := utf8.DecodeLastRuneInString(l.input[:l.pos])
	l.pos -= byteWidth
	l.char = char
}

func (l *lexer) nextN(n int) {
	for range n {
		l.next()
	}
}

func (l *lexer) backupN(n int) {
	for range n {
		l.backup()
	}
}

// skipSpace checks that the current and proceeding characters are both
// whitespace. if they are, the lexer progresses to the next character. if a
// newline has been skipped at any point, the final whitespace character is
// replaced in the input with a newline. the result is that any combination of
// whitespace which includes a newline is reduced to a single newline character
func (l *lexer) skipSpace() {
	if unicode.IsSpace(l.char) && unicode.IsSpace(l.peek()) {
		l.skippedSpace++
		if l.char == charNewline {
			l.skippedNewlines++
		}
		l.next()
		return
	}
	if l.skippedNewlines > 0 && unicode.IsSpace(l.char) {
		if l.char == charNewline {
			l.skippedNewlines++
		}
		l.input = l.input[:l.pos-1] + "\n" + l.input[l.pos:]
		l.char = charNewline
		return
	}
	l.skippedNewlines = 0
}

// addToTag stores the currently lexing tag in `l.tag`. a "tag" is a group of
// a-z characters, typically found before a `charOpeningSquare`
func (l *lexer) addToTag() {
	if l.token.Typ != typeText {
		return
	}
	if unicode.IsLetter(l.char) {
		l.tag += string(l.char)
		return
	}
	if unicode.IsSpace(l.char) {
		l.tag = ""
	}
}

// peek reads the following rune in the input but does not progress the lexer
func (l *lexer) peek() rune {
	if l.pos == len(l.input) {
		return eof
	}
	char, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return char
}

// peekBehind reads the previous rune in the input but does not backtrack
func (l *lexer) peekBehind() rune {
	if l.pos == 0 {
		return void
	}
	if l.pos == len(l.input) {
		char, _ := utf8.DecodeLastRuneInString(l.input)
		return char
	}
	char, _ := utf8.DecodeLastRuneInString(l.input[:l.pos-1])
	return char
}

func (l *lexer) peekNextNonSpace() rune {
	if l.pos == len(l.input) {
		return eof
	}
	char, _ := utf8.DecodeRuneInString(strings.TrimSpace(l.input[l.pos:]))
	return char
}

func (l *lexer) addToToken(c rune) {
	l.token.Val += string(c)
}

func (l *lexer) truncateToken(n int) {
	l.token.Val = l.token.Val[:len(l.token.Val)-n]
}

func (l *lexer) trimTrailingSpace() {
	if len(l.token.Val) == 0 {
		return
	}
	if unicode.IsSpace(l.peekBehind()) {
		l.token.Val = l.token.Val[:len(l.token.Val)-1]
	}
}

func (l *lexer) isHeadingChar() bool {
	if l.char == charDot || l.char == charColon {
		return true
	}
	return false
}

// lexGlobal is the starting state of the lexer
func (l *lexer) lexGlobal() {
	l.resetToken()
	l.next()
	if l.char == eof {
		l.token = l.mkToken(typeEOF, "")
		l.lexNext = nil
		return
	}
	if unicode.IsSpace(l.char) {
		l.lexNext = l.lexGlobal
		return
	}
	if l.char == charBackslash {
		l.backup()
		l.lexNext = l.lexText
		return
	}
	if l.isHeadingChar() {
		l.backup()
		l.lexNext = l.lexHeading
		return
	}
	if l.char == charHyphen {
		l.backup()
		l.lexNext = l.lexHypen
		return
	}
	l.backup()
	l.continuousNewline = true
	l.lexNext = l.lexText
}

// lexText is used to parse standard text
func (l *lexer) lexText() {
	l.token = l.mkToken(typeText, "")
	defer func() {
		if l.token.Val == "" {
			l.token.Typ = typeNone
		}
	}()
	for {
		l.next()
		if l.char == eof {
			l.trimTrailingSpace()
			l.lexNext = l.lexGlobal
			return
		}
		if l.char == charBackslash && l.peek() != charBackslash {
			l.tag = ""
			continue
		}
		if l.peekBehind() == charBackslash {
			l.addToToken(l.char)
			continue
		}
		if l.char == charNewline && l.ctx == ctxList {
			if l.peekNextNonSpace() == charHyphen {
				l.lexNext = l.lexHypen
				return
			}
			l.lexNext = l.lexTerminator
			return
		}
		if l.char == charNewline && (!l.continuousNewline || l.skippedNewlines >= 2) {
			l.lexNext = l.lexTerminator
			return
		}
		if l.char == charNewline && l.continuousNewline {
			l.addToToken(' ')
			continue
		}
		if l.char == charOpeningSquare {
			if len(l.tag) == 0 {
				l.trimTrailingSpace()
				l.backup()
				l.lexNext = l.lexOpeningSquare
				return
			}
			l.backupN(len(l.tag))
			l.truncateToken(len(l.tag))
			l.trimTrailingSpace()
			l.backup()
			l.lexNext = l.lexTag
			return
		}
		if l.char == charClosingSquare && l.peekBehind() != charBackslash {
			l.trimTrailingSpace()
			l.backup()
			l.lexNext = l.lexClosingSquare
			return
		}
		if l.token.Val == "" && unicode.IsSpace(l.char) {
			l.token = l.mkToken(typeText, "")
			continue
		}
		l.addToToken(l.char)
	}
}

// lexHeading lexes the heading symbols and returns to `lexText`
func (l *lexer) lexHeading() {
	l.token = l.mkToken(typeHeading, "")
	for {
		l.next()
		if l.char == eof {
			l.lexNext = l.lexGlobal
			return
		}
		if !l.isHeadingChar() || len(l.token.Val) == 3 {
			l.backup()
			l.lexNext = l.lexText
			return
		}
		l.token.Val += string(l.char)
	}
}

// lexTerminator is intended to produce a terminator token and return to `lexGlobal`
func (l *lexer) lexTerminator() {
	l.backup()
	l.token = l.mkToken(typeTerminator, string(charNewline))
	l.next()
	l.lexNext = l.lexGlobal
}

func (l *lexer) lexTag() {
	l.token = l.mkToken(typeTag, l.tag)
	l.nextN(len(l.tag))
	l.tag = ""
	l.lexNext = l.lexOpeningSquare
}

func (l *lexer) lexOpeningSquare() {
	l.token = l.mkToken(typeOpeningSquare, string(charOpeningSquare))
	l.next()
	l.lexNext = l.lexText
}

func (l *lexer) lexClosingSquare() {
	l.token = l.mkToken(typeClosingSquare, string(charClosingSquare))
	l.next()
	l.lexNext = l.lexText
}

func (l *lexer) lexHypen() {
	l.token = l.mkToken(typeBulletpoint, "-")
	l.token.indent = l.skippedSpace
	l.next()
	if unicode.IsSpace(l.char) {
		l.next()
	}
	l.ctx = ctxList
	l.lexNext = l.lexText
}
