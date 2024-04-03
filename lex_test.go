package runic

import (
	"strings"
	"testing"
)

type lexTest struct {
	name           string
	input          string
	expectedTokens []token
}

var lexTests = []lexTest{
	{
		"empty file",
		"",
		[]token{
			{Typ: typeEOF, Val: "", Line: 1, Pos: 0},
		},
	},
	{
		"plain text",
		"The quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 43},
		},
	},
	{
		"plain text ending newline",
		"The quick brown fox jumps over the lazy dog\n",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 44},
		},
	},
	{
		"plain text w/ leading space",
		"                                The quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 32},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 75},
		},
	},
	{
		"plain text with a single starting whitespace",
		" The quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 1},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 44},
		},
	},
	{
		"plain text w/ leading space including newlines",
		"\n        \n        \n        \nThe quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 5, Pos: 28},
			{Typ: typeEOF, Val: "", Line: 5, Pos: 71},
		},
	},
	{
		"plain text w/ leading space including newlines v2",
		"          \n        \n          The quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 3, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 73},
		},
	},
	{
		"plain text w/ trailing space",
		"The quick brown fox jumps over the lazy dog                                ",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 75},
		},
	},
	{
		"plain text w/ trailing space including newlines",
		"The quick brown fox jumps over the lazy dog\n        \n        \n        \n",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 4, Pos: 70},
			{Typ: typeEOF, Val: "", Line: 5, Pos: 71},
		},
	},
	{
		"plain text w/ trailing space including newlines v2",
		"The quick brown fox jumps over the lazy dog          \n        \n          ",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 2, Pos: 72},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 73},
		},
	},
	{
		"plain text w/ newline",
		"The quick brown fox\njumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 43},
		},
	},
	{
		"plain text w/ newline v1.2",
		"The quick brown fox \njumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 44},
		},
	},
	{
		"plain text w/ newline v1.3",
		"The quick brown fox\n jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 44},
		},
	},
	{
		"plain text w/ multiple newlines",
		"The quick brown fox\n\njumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 2, Pos: 20},
			{Typ: typeText, Val: "jumps over the lazy dog", Line: 3, Pos: 21},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 44},
		},
	},
	{
		"plain text w/ multiple newlines v2",
		"The quick brown fox\n\n\n\n\njumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 5, Pos: 23},
			{Typ: typeText, Val: "jumps over the lazy dog", Line: 6, Pos: 24},
			{Typ: typeEOF, Val: "", Line: 6, Pos: 47},
		},
	},
	{
		"plain text w/ multiple newlines v3",
		"The quick brown fox \n \n \n \n \n jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 5, Pos: 29},
			{Typ: typeText, Val: "jumps over the lazy dog", Line: 6, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 6, Pos: 53},
		},
	},
	{
		"plain text w/ multiple newlines v4",
		"The quick brown fox\n    \n    \njumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 3, Pos: 29},
			{Typ: typeText, Val: "jumps over the lazy dog", Line: 4, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 4, Pos: 53},
		},
	},
	{
		"plain text containing backslash",
		"The quick brown fox jumps\\\\over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox jumps\\over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 44},
		},
	},
	{
		"heading one",
		". This is a level one heading",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level one heading", Line: 1, Pos: 2},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 29},
		},
	},
	{
		"heading one escaped",
		"\\. This is a level one heading",
		[]token{
			{Typ: typeText, Val: ". This is a level one heading", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 30},
		},
	},
	{
		"heading one escape escaped",
		"\\\\. This is a level one heading",
		[]token{
			{Typ: typeText, Val: "\\. This is a level one heading", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 31},
		},
	},
	{
		"heading two",
		": This is a level two heading",
		[]token{
			{Typ: typeHeading, Val: ":", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level two heading", Line: 1, Pos: 2},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 29},
		},
	},
	{
		"heading three",
		":. This is a level three heading",
		[]token{
			{Typ: typeHeading, Val: ":.", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level three heading", Line: 1, Pos: 3},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 32},
		},
	},
	{
		"heading four",
		":: This is a level four heading",
		[]token{
			{Typ: typeHeading, Val: "::", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level four heading", Line: 1, Pos: 3},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 31},
		},
	},
	{
		"heading five",
		"::. This is a level five heading",
		[]token{
			{Typ: typeHeading, Val: "::.", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level five heading", Line: 1, Pos: 4},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 32},
		},
	},
	{
		"heading six",
		"::: This is a level six heading",
		[]token{
			{Typ: typeHeading, Val: ":::", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level six heading", Line: 1, Pos: 4},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 31},
		},
	},
	{
		"heading w/ excessive characters",
		":::. This is a heading",
		[]token{
			{Typ: typeHeading, Val: ":::", Line: 1, Pos: 0},
			{Typ: typeText, Val: ". This is a heading", Line: 1, Pos: 3},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 22},
		},
	},
	{
		"heading w/ excessive characters v2",
		":::: This is a heading",
		[]token{
			{Typ: typeHeading, Val: ":::", Line: 1, Pos: 0},
			{Typ: typeText, Val: ": This is a heading", Line: 1, Pos: 3},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 22},
		},
	},
	{
		"heading one w/ paragraph underneath",
		". This is a level one heading\nThe quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level one heading", Line: 1, Pos: 2},
			{Typ: typeTerminator, Val: "\n", Line: 1, Pos: 29},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 2, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 73},
		},
	},
	{
		"heading one w/ multiple newlines and paragraph underneath",
		". This is a level one heading\n\n\n\n\n\nThe quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level one heading", Line: 1, Pos: 2},
			{Typ: typeTerminator, Val: "\n", Line: 6, Pos: 34},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 7, Pos: 35},
			{Typ: typeEOF, Val: "", Line: 7, Pos: 78},
		},
	},
	{
		"heading one w/ multiple newlines and spaces and paragraph underneath",
		". This is a level one heading    \n    \n    \n    \n    \n    \n    The quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level one heading", Line: 1, Pos: 2},
			{Typ: typeTerminator, Val: "\n", Line: 6, Pos: 62},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 7, Pos: 63},
			{Typ: typeEOF, Val: "", Line: 7, Pos: 106},
		},
	},
	{
		"heading one w/ multiple newlines and spaces and paragraph underneath v2",
		". This is a level one heading\n     \n      \n      \n      \n     \nThe quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a level one heading", Line: 1, Pos: 2},
			{Typ: typeTerminator, Val: "\n", Line: 6, Pos: 62},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 7, Pos: 63},
			{Typ: typeEOF, Val: "", Line: 7, Pos: 106},
		},
	},
	{
		"heading one with paragraph underneath and excessive whitespace all over",
		"\n    \n    . This      is    a level one   heading    \n    \n    The    quick brown       fox jumps    over the   lazy dog    \n    \n",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 3, Pos: 10},
			{Typ: typeText, Val: "This is a level one heading", Line: 3, Pos: 12},
			{Typ: typeTerminator, Val: "\n", Line: 4, Pos: 62},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 5, Pos: 63},
			{Typ: typeTerminator, Val: "\n", Line: 6, Pos: 129},
			{Typ: typeEOF, Val: "", Line: 7, Pos: 130},
		},
	},
	{
		"rich text",
		"The quick brown fox bold[jumps over] the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 20},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 24},
			{Typ: typeText, Val: "jumps over", Line: 1, Pos: 25},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 35},
			{Typ: typeText, Val: "the lazy dog", Line: 1, Pos: 37},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 49},
		},
	},
	{
		"only rich text with newline",
		"bold[fox]\n",
		[]token{
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 0},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 4},
			{Typ: typeText, Val: "fox", Line: 1, Pos: 5},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 8},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 10},
		},
	},
	{
		"rich text without closing square",
		"bold[The",
		[]token{
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 0},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 4},
			{Typ: typeText, Val: "The", Line: 1, Pos: 5},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 8},
		},
	},
	{
		"rich text with a closing square",
		"The quick brown fox ] jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 20},
			{Typ: typeText, Val: "jumps over the lazy dog", Line: 1, Pos: 22},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 45},
		},
	},
	{
		"rich text escaped",
		"The quick brown fox bold\\[jumps over\\] the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox bold[jumps over] the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 51},
		},
	},
	{
		"rich text nested",
		"The quick bold[brown fox italic[jumps] over the] lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick", Line: 1, Pos: 0},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 10},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 14},
			{Typ: typeText, Val: "brown fox", Line: 1, Pos: 15},
			{Typ: typeTag, Val: "italic", Line: 1, Pos: 25},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 31},
			{Typ: typeText, Val: "jumps", Line: 1, Pos: 32},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 37},
			{Typ: typeText, Val: "over the", Line: 1, Pos: 39},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 47},
			{Typ: typeText, Val: "lazy dog", Line: 1, Pos: 49},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 57},
		},
	},
	{
		"nested rich text compact",
		"bold[italic[quick]]",
		[]token{
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 0},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 4},
			{Typ: typeTag, Val: "italic", Line: 1, Pos: 5},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 11},
			{Typ: typeText, Val: "quick", Line: 1, Pos: 12},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 17},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 18},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 19},
		},
	},
	{
		"square brackets inside plain text",
		"the [quick] fox",
		[]token{
			{Typ: typeText, Val: "the", Line: 1, Pos: 0},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 4},
			{Typ: typeText, Val: "quick", Line: 1, Pos: 5},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 10},
			{Typ: typeText, Val: "fox", Line: 1, Pos: 12},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 15},
		},
	},
	{
		"rich text with no space before tag",
		"The quick brown fox\\bold[jumps over] the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 0},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 20},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 24},
			{Typ: typeText, Val: "jumps over", Line: 1, Pos: 25},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 35},
			{Typ: typeText, Val: "the lazy dog", Line: 1, Pos: 37},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 49},
		},
	},
	{
		"heading with rich text",
		". The quick brown fox bold[jumps over] the lazy dog",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 2},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 22},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 26},
			{Typ: typeText, Val: "jumps over", Line: 1, Pos: 27},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 37},
			{Typ: typeText, Val: "the lazy dog", Line: 1, Pos: 39},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 51},
		},
	},
	{
		"heading with rich text with paragraph underneath",
		". The quick brown fox bold[jumps over] the lazy dog\nLorem ipsum dolor sit amet, consectetur adipiscing elit.",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "The quick brown fox", Line: 1, Pos: 2},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 22},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 26},
			{Typ: typeText, Val: "jumps over", Line: 1, Pos: 27},
			{Typ: typeClosingSquare, Val: "]", Line: 1, Pos: 37},
			{Typ: typeText, Val: "the lazy dog", Line: 1, Pos: 39},
			{Typ: typeTerminator, Val: "\n", Line: 1, Pos: 51},
			{Typ: typeText, Val: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", Line: 2, Pos: 52},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 108},
		},
	},
	{
		"list",
		"- Item one\n- Item two\n- Item three",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "Item one", Line: 1, Pos: 2},
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 11, indent: 0},
			{Typ: typeText, Val: "Item two", Line: 2, Pos: 13},
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 22, indent: 0},
			{Typ: typeText, Val: "Item three", Line: 3, Pos: 24},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 34},
		},
	},
	{
		"list compact",
		"-Item one\n-Item two\n-Item three",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "Item one", Line: 1, Pos: 1},
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 10, indent: 0},
			{Typ: typeText, Val: "Item two", Line: 2, Pos: 11},
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 20, indent: 0},
			{Typ: typeText, Val: "Item three", Line: 3, Pos: 21},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 31},
		},
	},
	{
		"list with indents",
		"         \n        \n         - Item one\n             - Item two    \n- Item three",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 28, indent: 9},
			{Typ: typeText, Val: "Item one", Line: 3, Pos: 30},
			{Typ: typeBulletpoint, Val: "-", Line: 4, Pos: 52, indent: 13},
			{Typ: typeText, Val: "Item two", Line: 4, Pos: 54},
			{Typ: typeBulletpoint, Val: "-", Line: 5, Pos: 67, indent: 0},
			{Typ: typeText, Val: "Item three", Line: 5, Pos: 69},
			{Typ: typeEOF, Val: "", Line: 5, Pos: 79},
		},
	},
	{
		"list with paragraph underneath",
		"- Item one\n- Item two\n- Item three\nThe quick brown fox jumps over the lazy dog",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "Item one", Line: 1, Pos: 2},
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 11, indent: 0},
			{Typ: typeText, Val: "Item two", Line: 2, Pos: 13},
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 22, indent: 0},
			{Typ: typeText, Val: "Item three", Line: 3, Pos: 24},
			{Typ: typeTerminator, Val: "\n", Line: 3, Pos: 34},
			{Typ: typeText, Val: "The quick brown fox jumps over the lazy dog", Line: 4, Pos: 35},
			{Typ: typeEOF, Val: "", Line: 4, Pos: 78},
		},
	},
	{
		"plain text with hyphen",
		"The quick brown fox - jumps over the lazy dog",
		[]token{
			{Typ: typeText, Val: "The quick brown fox - jumps over the lazy dog", Line: 1, Pos: 0},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 45},
		},
	},
	{
		"heading with a hyphen",
		". This is a - level one heading",
		[]token{
			{Typ: typeHeading, Val: ".", Line: 1, Pos: 0},
			{Typ: typeText, Val: "This is a - level one heading", Line: 1, Pos: 2},
			{Typ: typeEOF, Val: "", Line: 1, Pos: 31},
		},
	},
	{
		"invalid heading with paragraph underneath",
		"..\na",
		[]token{
			{Typ: typeHeading, Val: "..", Line: 1, Pos: 0},
			{Typ: typeTerminator, Val: "\n", Line: 1, Pos: 2},
			{Typ: typeText, Val: "a", Line: 2, Pos: 3},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 4},
		},
	},
	{
		"list with rich text",
		"- Item one\n- Item bold[two]\n- Item three",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "Item one", Line: 1, Pos: 2},
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 11, indent: 0},
			{Typ: typeText, Val: "Item", Line: 2, Pos: 13},
			{Typ: typeTag, Val: "bold", Line: 2, Pos: 18},
			{Typ: typeOpeningSquare, Val: "[", Line: 2, Pos: 22},
			{Typ: typeText, Val: "two", Line: 2, Pos: 23},
			{Typ: typeClosingSquare, Val: "]", Line: 2, Pos: 26},
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 28, indent: 0},
			{Typ: typeText, Val: "Item three", Line: 3, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 3, Pos: 40},
		},
	},
	{
		"list with paragraph underneath",
		"- a\nb",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "a", Line: 1, Pos: 2},
			{Typ: typeTerminator, Val: "\n", Line: 1, Pos: 3},
			{Typ: typeText, Val: "b", Line: 2, Pos: 4},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 5},
		},
	},
	{
		"test",
		`
	       - Item one
	       - Item two
	  `,
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 9, indent: 8},
			{Typ: typeText, Val: "Item one", Line: 2, Pos: 11},
			{Typ: typeBulletpoint, Val: "-", Line: 3, Pos: 28, indent: 8},
			{Typ: typeText, Val: "Item two", Line: 3, Pos: 30},
			{Typ: typeTerminator, Val: "\n", Line: 3, Pos: 41},
			{Typ: typeEOF, Val: "", Line: 4, Pos: 42},
		},
	},
	{
		"rich text over 2 list items",
		"- The bold[quick\n- brown fox] jumps",
		[]token{
			{Typ: typeBulletpoint, Val: "-", Line: 1, Pos: 0, indent: 0},
			{Typ: typeText, Val: "The", Line: 1, Pos: 2},
			{Typ: typeTag, Val: "bold", Line: 1, Pos: 6},
			{Typ: typeOpeningSquare, Val: "[", Line: 1, Pos: 10},
			{Typ: typeText, Val: "quick", Line: 1, Pos: 11},
			{Typ: typeBulletpoint, Val: "-", Line: 2, Pos: 17, indent: 0},
			{Typ: typeText, Val: "brown fox", Line: 2, Pos: 19},
			{Typ: typeClosingSquare, Val: "]", Line: 2, Pos: 28},
			{Typ: typeText, Val: "jumps", Line: 2, Pos: 30},
			{Typ: typeEOF, Val: "", Line: 2, Pos: 35},
		},
	},
}

func tokensAreEqual(lexedTokens, expectedTokens []token) bool {
	if len(lexedTokens) != len(expectedTokens) {
		return false
	}

	for i := range lexedTokens {
		if lexedTokens[i].Typ != expectedTokens[i].Typ {
			return false
		}
		if lexedTokens[i].Val != expectedTokens[i].Val {
			return false
		}
		if lexedTokens[i].Line != expectedTokens[i].Line {
			return false
		}
		if lexedTokens[i].Pos != expectedTokens[i].Pos {
			return false
		}
		if lexedTokens[i].indent != expectedTokens[i].indent {
			return false
		}
	}

	return true
}

func stringifyTokens(tokens []token) string {
	var tokensStrings []string
	for _, token := range tokens {
		tokensStrings = append(tokensStrings, token.String())
	}
	return strings.Join(tokensStrings, " ")
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		lexedTokens := collectTokens(test.input)
		if !tokensAreEqual(lexedTokens, test.expectedTokens) {
			t.Errorf("%s ERROR\nexpected: %s\nreceived: %s", test.name, stringifyTokens(test.expectedTokens), stringifyTokens(lexedTokens))
			continue
		}
		t.Log(test.name, "OK")
	}
}
