package runic

import "testing"

type htmlTest struct {
	name         string
	input        string
	expectedHtml string
}

var htmlTests = []htmlTest{
	{
		"empty file",
		"",
		"",
	},
	{
		"paragraph",
		"The quick brown fox jumps over the lazy dog",
		"<p>The quick brown fox jumps over the lazy dog</p>",
	},
	{
		"heading one",
		". This is a level one heading",
		"<h1>This is a level one heading</h1>",
	},
	{
		"heading two",
		": This is a level two heading",
		"<h2>This is a level two heading</h2>",
	},
	{
		"heading three",
		":. This is a level three heading",
		"<h3>This is a level three heading</h3>",
	},
	{
		"heading four",
		":: This is a level four heading",
		"<h4>This is a level four heading</h4>",
	},
	{
		"heading five",
		"::. This is a level five heading",
		"<h5>This is a level five heading</h5>",
	},
	{
		"heading six",
		"::: This is a level six heading",
		"<h6>This is a level six heading</h6>",
	},
	{
		"rich text",
		"italic[quick]",
		"<p><em>quick</em></p>",
	},
	{
		"nested rich text",
		"The quick bold[brown fox italic[jumps] over the] lazy dog",
		"<p>The quick <b>brown fox <em>jumps</em> over the</b> lazy dog</p>",
	},
	{
		"list",
		"- Item one\n- Item two\n- Item three",
		"<ul><li>Item one</li><li>Item two</li><li>Item three</li></ul>",
	},
	{
		"list string literal",
		`
	       - Item one
	       - Item two
	   `,
		"<ul><li>Item one</li><li>Item two</li></ul>",
	},
	{
		"list string literal v2",
		`
	       - Item one
           - Item two
	   `,
		"<ul><li>Item one</li><ul><li>Item two</li></ul></ul>",
	},
	{
		"list with indents v3",
		`
        - Item one
          - Item two
            - Item three
          - Item four
        - Item five
    `,
		"<ul><li>Item one</li><ul><li>Item two</li><ul><li>Item three</li></ul><li>Item four</li></ul><li>Item five</li></ul>",
	},
}

func TestHtml(t *testing.T) {
	for _, test := range htmlTests {
		testParser := New()
		htmlString := testParser.Html(test.input)
		if htmlString != test.expectedHtml {
			t.Errorf("%s ERROR\nexpected: %s\nreceived: %s", test.name, test.expectedHtml, htmlString)
			continue
		}
		t.Log(test.name, "OK")
	}
}

type highlightTextTest struct {
	name                  string
	input                 string
	expectedHighlightText string
}

var highlightTextTests = []highlightTextTest{
	{
		"empty file",
		"",
		"",
	},
	{
		"paragraph",
		"The quick brown fox jumps over the lazy dog",
		`<span class="runic__text">The quick brown fox jumps over the lazy dog</span>`,
	},
	{
		"paragraph with one leading space",
		" The quick brown fox jumps over the lazy dog",
		`<span class="runic__text">&nbsp;The quick brown fox jumps over the lazy dog</span>`,
	},
	{
		"paragraph with multiple leading space",
		"    The quick brown fox jumps over the lazy dog",
		`<span class="runic__text">&nbsp;&nbsp;&nbsp;&nbsp;The quick brown fox jumps over the lazy dog</span>`,
	},
	{
		"paragraph with one trailing space",
		"The quick brown fox jumps over the lazy dog ",
		`<span class="runic__text">The quick brown fox jumps over the lazy dog&nbsp;</span>`,
	},
	{
		"paragraph with multple trailing space",
		"The quick brown fox jumps over the lazy dog    ",
		`<span class="runic__text">The quick brown fox jumps over the lazy dog&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
	},
	{
		"paragraph with multiple inner whitespace",
		"The quick brown                       fox jumps over the lazy dog",
		`<span class="runic__text">The quick brown&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fox jumps over the lazy dog</span>`,
	},
	{
		"paragraph with leading, trailing, and multiple inner whitespace",
		" The quick brown                       fox jumps over the lazy dog ",
		`<span class="runic__text">&nbsp;The quick brown&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;fox jumps over the lazy dog&nbsp;</span>`,
	},
	{
		"heading one",
		". This is a level one heading",
		`<span class="runic__heading">.&nbsp;</span><span class="runic__text">This is a level one heading</span>`,
	},
	{
		"rich text",
		"The quick bold[brown fox jumps over the] lazy dog",
		`<span class="runic__text">The quick&nbsp;</span><span class="runic__tag">bold</span><span class="runic__osq">[</span><span class="runic__text">brown fox jumps over the</span><span class="runic__csq">]&nbsp;</span><span class="runic__text">lazy dog</span>`,
	},
	{
		"paragraph with surrounding newlines",
		`
      The quick brown fox jumps over the lazy dog

      The quick brown fox jumps over the lazy dog
    `,
		`<span class="runic__text"><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;The quick brown fox jumps over the lazy dog<br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>&nbsp;<span class="runic__text">The quick brown fox jumps over the lazy dog<br>&nbsp;&nbsp;&nbsp;&nbsp;</span>`,
	},
	{
		"lists",
		`
      - List item one
        - List item two
          - List item three
        - List item four
      - List item five
    `,
		`<span class="runic__bulletpoint"><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-&nbsp;</span><span class="runic__text">List item one<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="runic__bulletpoint">-&nbsp;</span><span class="runic__text">List item two<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="runic__bulletpoint">-&nbsp;</span><span class="runic__text">List item three<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="runic__bulletpoint">-&nbsp;</span><span class="runic__text">List item four<br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span><span class="runic__bulletpoint">-&nbsp;</span><span class="runic__text">List item five<br>&nbsp;&nbsp;&nbsp;</span>&nbsp;`,
	},
	{
		"heading with paragraph underneath and surround whitespace",
		`
      . This is a level one heading


      The quick brown fox jumps over the lazy dog

    `,
		`<span class="runic__heading"><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;.&nbsp;</span><span class="runic__text">This is a level one heading<br><br><br>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>&nbsp;<span class="runic__text">The quick brown fox jumps over the lazy dog<br><br>&nbsp;&nbsp;&nbsp;</span>&nbsp;`,
	},
}

func TestHighlightText(t *testing.T) {
	for _, test := range highlightTextTests {
		testParser := New()
		highlightText := testParser.HighlightText(test.input)
		if highlightText != test.expectedHighlightText {
			t.Errorf("%s ERROR\nexpected: %s\nreceived: %s", test.name, test.expectedHighlightText, highlightText)
			continue
		}
		t.Log(test.name, "OK")
	}
}
