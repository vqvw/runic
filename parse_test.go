package runic

import (
	"encoding/json"
	"fmt"
	"testing"
)

type parseTest struct {
	name         string
	input        string
	expectedTree *Node
}

var parseTests = []parseTest{
	{
		"empty file",
		"",
		&Node{
			Typ: nodeRoot,
		},
	},
	{
		"plain text",
		"The quick brown fox jumps over the lazy dog",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "The quick brown fox jumps over the lazy dog",
						},
					},
				},
			},
		},
	},
	{
		"nested rich text",
		"The quick bold[brown fox italic[jumps] over the] lazy dog",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "The quick",
						},
						{
							Typ: nodeBoldTag,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "brown fox",
								},
								{
									Typ: nodeItalicTag,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "jumps",
										},
									},
								},
								{
									Typ: nodeText,
									Val: "over the",
								},
							},
						},
						{
							Typ: nodeText,
							Val: "lazy dog",
						},
					},
				},
			},
		},
	},
	{
		"only rich text",
		"bold[fox]",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeBoldTag,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "fox",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"rich text with invalid tag",
		"The quick foo[brown fox jumps over the] lazy dog",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "The quick",
						},
						{
							Typ: nodeError,
							Val: "Invalid tag name: foo",
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "brown fox jumps over the",
								},
							},
						},
						{
							Typ: nodeText,
							Val: "lazy dog",
						},
					},
				},
			},
		},
	},
	{
		"multiple paragraphs",
		"The quick brown fox\n\njumps over the lazy dog",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "The quick brown fox",
						},
					},
				},
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "jumps over the lazy dog",
						},
					},
				},
			},
		},
	},
	{
		"level one heading",
		". This is a level one heading",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeHeadingOne,
					Val: ".",
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "This is a level one heading",
						},
					},
				},
			},
		},
	},
	{
		"level one heading with paragraphs underneath",
		". This is a level one heading\nThe quick brown fox jumps over the lazy dog\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit.",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeHeadingOne,
					Val: ".",
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "This is a level one heading",
						},
					},
				},
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "The quick brown fox jumps over the lazy dog",
						},
					},
				},
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
						},
					},
				},
			},
		},
	},
	{
		"invalid heading with paragraph underneath",
		"..\na",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeError,
					Val: fmt.Sprintf("%s: ..", errInvalidHeading),
				},
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "a",
						},
					},
				},
			},
		},
	},
	{
		"list",
		"- Item one\n-Item two\n-Item three",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item two",
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item three",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"list with indents",
		"         \n        \n  - Item one\n                  - Item two    \n        - Item three",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeList,
							Children: []*Node{
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item two",
										},
									},
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item three",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"list with indents v2",
		"- Item one\n  - Item two\n    - Item three\n      - Item four\n  - Item five",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeList,
							Children: []*Node{
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item two",
										},
									},
								},
								{
									Typ: nodeList,
									Children: []*Node{
										{
											Typ: nodeListItem,
											Children: []*Node{
												{
													Typ: nodeText,
													Val: "Item three",
												},
											},
										},
										{
											Typ: nodeList,
											Children: []*Node{
												{
													Typ: nodeListItem,
													Children: []*Node{
														{
															Typ: nodeText,
															Val: "Item four",
														},
													},
												},
											},
										},
									},
								},
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item five",
										},
									},
								},
							},
						},
					},
				},
			},
		},
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
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeList,
							Children: []*Node{
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item two",
										},
									},
								},
								{
									Typ: nodeList,
									Children: []*Node{
										{
											Typ: nodeListItem,
											Children: []*Node{
												{
													Typ: nodeText,
													Val: "Item three",
												},
											},
										},
									},
								},
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item four",
										},
									},
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item five",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"list with paragraph underneath",
		" - a\nb",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "a",
								},
							},
						},
					},
				},
				{
					Typ: nodeParagraph,
					Children: []*Node{
						{
							Typ: nodeText,
							Val: "b",
						},
					},
				},
			},
		},
	},
	{
		"list items string literal",
		`
	       - Item one
	       - Item two
	  `,
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item two",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"list items string literal v2",
		`
	       - Item one
           - Item two
	  `,
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "Item one",
								},
							},
						},
						{
							Typ: nodeList,
							Children: []*Node{
								{
									Typ: nodeListItem,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "Item two",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"rich text over 2 list items",
		"- The bold[quick\n- brown fox] jumps",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "The",
								},
								{
									Typ: nodeBoldTag,
									Children: []*Node{
										{
											Typ: nodeText,
											Val: "quick",
										},
									},
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "brown fox jumps",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"rich text over 2 list items v2",
		"- The italic[bold[quick\n- brown fox] jumps",
		&Node{
			Typ: nodeRoot,
			Children: []*Node{
				{
					Typ: nodeList,
					Children: []*Node{
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "The",
								},
								{
									Typ: nodeItalicTag,
									Children: []*Node{
										{
											Typ: nodeBoldTag,
											Children: []*Node{
												{
													Typ: nodeText,
													Val: "quick",
												},
											},
										},
									},
								},
							},
						},
						{
							Typ: nodeListItem,
							Children: []*Node{
								{
									Typ: nodeText,
									Val: "brown fox jumps",
								},
							},
						},
					},
				},
			},
		},
	},
}

func checkChildren(parsedChildren, expectedChildren []*Node) bool {
	if len(parsedChildren) != len(expectedChildren) {
		return false
	}

	for i, childNode := range parsedChildren {
		if childNode.Typ != expectedChildren[i].Typ {
			return false
		}
		if childNode.Val != expectedChildren[i].Val {
			return false
		}
		if len(childNode.Children) > 0 {
			if !checkChildren(childNode.Children, expectedChildren[i].Children) {
				return false
			}
		}
	}

	return true
}

func treesAreEqual(parsedTree, expectedTree *Node) bool {
	if parsedTree.Typ != expectedTree.Typ {
		return false
	}
	if parsedTree.Val != expectedTree.Val {
		return false
	}
	return checkChildren(parsedTree.Children, expectedTree.Children)
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		testParser := New()
		parsedTree := testParser.Parse(test.input)
		if !treesAreEqual(parsedTree, test.expectedTree) {
			expectedTreeJSON, _ := json.MarshalIndent(test.expectedTree, "", "  ")
			parsedTreeJSON, _ := json.MarshalIndent(parsedTree, "", "  ")
			t.Errorf("%s ERROR\nexpected: %v\nreceived: %v", test.name, string(expectedTreeJSON), string(parsedTreeJSON))
			continue
		}
		t.Log(test.name, "OK")
	}
}
