package template

import (
	"reflect"
	"testing"

	"github.com/gookit/goutil/dump"
)

func TestParseTextNode(t *testing.T) {
	source := "Hello, world!"
	parser := NewParser()
	tpl, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := &Template{
		Nodes: []*Node{{Type: "text", Text: "Hello, world!"}},
	}

	if !reflect.DeepEqual(tpl, expected) {
		t.Errorf("Expected %v, got %v", expected, tpl)
	}
}

func TestParseTextNodeWithWhitespace(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"OnlySpaces",
			"    ",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "    "}},
			},
		},
		{
			"OnlyLineBreaks",
			"\n\n\n",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "\n\n\n"}},
			},
		},
		{
			"OnlyTabs",
			"\t\t\t",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "\t\t\t"}},
			},
		},
		{
			"SpacesAndText",
			"  Hello, world!  ",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "  Hello, world!  "}},
			},
		},
		{
			"Newlines",
			"\nHello,\nworld!\n",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "\nHello,\nworld!\n"}},
			},
		},
		{
			"TabsAndSpaces",
			"\tHello,  world!\t",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "\tHello,  world!\t"}},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseTextNodeWithMultipleLinesAndVariations(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"MultipleLinesSimple",
			`Hello,
World!
This is a test.`,
			&Template{
				Nodes: []*Node{{Type: "text", Text: "Hello,\nWorld!\nThis is a test."}},
			},
		},
		{
			"MultipleLinesWithEmptyLines",
			`Hello,

World!


This is a test.`,
			&Template{
				Nodes: []*Node{{Type: "text", Text: "Hello,\n\nWorld!\n\n\nThis is a test."}},
			},
		},
		{
			"MultipleLinesWithVariable",
			`User: {{username}}
Welcome back!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Text: "{{username}}"},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndText",
			`User: {{
username
}}
Welcome back!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Text: "{{\nusername\n}}"},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndTextAndSpaces",
			`User: {{
	username
	}}
Welcome back!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Text: `{{
	username
	}}`},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndTextAndTabs",
			"User: {{\t username \n}}\nWelcome back!",
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Text: "{{\t username \n}}"},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndFilters",
			`User: {{username|lower}}
Welcome back!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Filters: []Filter{{Name: "lower"}}, Text: "{{username|lower}}"},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndFiltersAndText",
			`User: {{username|lower}}
Welcome back, {{username}}!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Filters: []Filter{{Name: "lower"}}, Text: "{{username|lower}}"},
					{Type: "text", Text: "\nWelcome back, "},
					{Type: "variable", Variable: "username", Text: "{{username}}"},
					{Type: "text", Text: "!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndFiltersAndTextAndSpaces",
			`User: {{ username | lower }}
Welcome back, {{ username }}!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Filters: []Filter{{Name: "lower"}}, Text: "{{ username | lower }}"},
					{Type: "text", Text: "\nWelcome back, "},
					{Type: "variable", Variable: "username", Text: "{{ username }}"},
					{Type: "text", Text: "!"},
				},
			},
		},
		{
			"MultipleLinesWithVariableAndFiltersAndArgs",
			`User: {{username|replace:"Mr.","Mrs."}}
Welcome back!`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "username", Filters: []Filter{{Name: "replace", Args: []FilterArg{StringArg{val: "Mr."}, StringArg{val: "Mrs."}}}}, Text: `{{username|replace:"Mr.","Mrs."}}`},
					{Type: "text", Text: "\nWelcome back!"},
				},
			},
		},
		{
			"MixedSpacesAndTabs",
			"\tHello,\n  World!  \n\tThis is a test.",
			&Template{
				Nodes: []*Node{{Type: "text", Text: "\tHello,\n  World!  \n\tThis is a test."}},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseVariableNode(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{"NoSpaces", "{{username}}"},
		{"SpacesBeforeVariable", "{{  username}}"},
		{"SpacesAfterVariable", "{{username  }}"},
		{"SpacesBeforeAndAfterVariable", "{{  username  }}"},
		{"LineBreakBeforeVariable", "{{\nusername}}"},
		{"LineBreakAfterVariable", "{{username\n}}"},
		{"LineBreakBeforeAndAfterVariable", "{{\nusername\n}}"},
		{"LineBreaksAndSpacesBeforeVariable", "{{  \nusername}}"},
		{"LineBreaksAndSpacesAfterVariable", "{{username\n  }}"},
		{"LineBreaksAndSpacesBeforeAndAfterVariable", "{{  \nusername\n  }}"},
		{"TabsBeforeVariable", "{{\tusername}}"},
		{"TabsAfterVariable", "{{username\t}}"},
		{"TabsBeforeAndAfterVariable", "{{\tusername\t}}"},
		{"TabsAndSpacesBeforeVariable", "{{\t  username}}"},
		{"TabsAndSpacesAfterVariable", "{{username  \t}}"},
		{"TabsAndSpacesBeforeAndAfterVariable", "{{\t  username  \t}}"},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			expected := &Template{
				Nodes: []*Node{{
					Type:     "variable",
					Variable: "username",
					Text:     tc.source,
				}},
			}

			if !reflect.DeepEqual(tpl, expected) {
				t.Errorf("Expected %v, got %v", expected, tpl)
			}
		})
	}
}

func TestParseNestedContextVariable(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{"DirectNestedVariable", "{{user.details.name}}"},
		{"SpacesInsideBraces", "{{ user.details.name }}"},
		{"SpacesBeforeVariable", "{{  user.details.name}}"},
		{"SpacesAfterVariable", "{{user.details.name  }}"},
		{"SpacesBeforeAndAfterVariable", "{{  user.details.name  }}"},
		{"TabsBeforeNestedVariable", "{{\tuser.details.name}}"},
		{"TabsAfterNestedVariable", "{{user.details.name\t}}"},
		{"TabsAroundNestedVariable", "{{\tuser.details.name\t}}"},
		{"LineBreakBeforeNestedVariable", "{{\nuser.details.name}}"},
		{"LineBreakAfterNestedVariable", "{{user.details.name\n}}"},
		{"LineBreaksAroundNestedVariable", "{{\nuser.details.name\n}}"},
		{"MixedWhitespaceAroundNestedVariable", "{{ \t\nuser.details.name\t\n }}"},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			expected := &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "user.details.name",
						Text:     tc.source,
					},
				},
			}

			if !reflect.DeepEqual(tpl, expected) {
				t.Errorf("Case %s: Expected %+v, got %+v", tc.name, expected, tpl)
			}
		})
	}
}

func TestParseMixedTextAndVariableNodes(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"BasicMixedContent",
			"Hello, {{username}}! Welcome to the site.",
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello, "},
					{Type: "variable", Variable: "username", Text: "{{username}}"},
					{Type: "text", Text: "! Welcome to the site."},
				},
			},
		},
		{
			"SpacesInsideVariableBraces",
			"Hello, {{ username }}! Welcome to our world.",
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello, "},
					{Type: "variable", Variable: "username", Text: "{{ username }}"},
					{Type: "text", Text: "! Welcome to our world."},
				},
			},
		},
		{
			"MultipleVariables",
			"User: {{ firstName }} {{ lastName }} - Welcome back!",
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "User: "},
					{Type: "variable", Variable: "firstName", Text: "{{ firstName }}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "lastName", Text: "{{ lastName }}"},
					{Type: "text", Text: " - Welcome back!"},
				},
			},
		},
		{
			"VariableStartOfLine",
			"{{ greeting }} John, have a great day!",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "greeting", Text: "{{ greeting }}"},
					{Type: "text", Text: " John, have a great day!"},
				},
			},
		},
		{
			"VariableEndOfLine",
			"Goodbye, {{ username }}",
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Goodbye, "},
					{Type: "variable", Variable: "username", Text: "{{ username }}"},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseVariableNodeWithFilterNoParams(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			name:   "SingleFilterNoSpace",
			source: "{{username|upper}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "upper"}}, Text: "{{username|upper}}"}},
			},
		},
		{
			name:   "SpaceBeforePipe",
			source: "{{username |upper}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "upper"}}, Text: "{{username |upper}}"}},
			},
		},
		{
			name:   "SpaceAfterPipe",
			source: "{{username| upper}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "upper"}}, Text: "{{username| upper}}"}},
			},
		},
		{
			name:   "SpacesAroundPipe",
			source: "{{username | upper}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "upper"}}, Text: "{{username | upper}}"}},
			},
		},
		{
			name:   "MultipleFiltersNoSpaces",
			source: "{{username|lower|capitalize}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "lower"}, {Name: "capitalize"}}, Text: "{{username|lower|capitalize}}"}},
			},
		},
		{
			name:   "SpacesAroundMultipleFilters",
			source: "{{username | lower | capitalize}}",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "lower"}, {Name: "capitalize"}}, Text: "{{username | lower | capitalize}}"}},
			},
		},
		{
			name:   "TextNodesAroundVariableWithFilter",
			source: "Hello {{name|capitalize}}, welcome!",
			expected: &Template{
				Nodes: []*Node{{Type: "text", Text: "Hello "}, {Type: "variable", Variable: "name", Filters: []Filter{{Name: "capitalize"}}, Text: "{{name|capitalize}}"}, {Type: "text", Text: ", welcome!"}},
			},
		},
		{
			name:   "TextNodeBeforeVariableMultipleFilters",
			source: "User: {{username|trim|lower}}",
			expected: &Template{
				Nodes: []*Node{{Type: "text", Text: "User: "}, {Type: "variable", Variable: "username", Filters: []Filter{{Name: "trim"}, {Name: "lower"}}, Text: "{{username|trim|lower}}"}},
			},
		},
		{
			name:   "TextNodeAfterVariableMultipleFilters",
			source: "{{username|trim|capitalize}} logged in",
			expected: &Template{
				Nodes: []*Node{{Type: "variable", Variable: "username", Filters: []Filter{{Name: "trim"}, {Name: "capitalize"}}, Text: "{{username|trim|capitalize}}"}, {Type: "text", Text: " logged in"}},
			},
		},
		{
			name:   "ComplexMixedTextAndVariables",
			source: "Dear {{name|capitalize}}, your score is {{score|round}}.",
			expected: &Template{
				Nodes: []*Node{{Type: "text", Text: "Dear "}, {Type: "variable", Variable: "name", Filters: []Filter{{Name: "capitalize"}}, Text: "{{name|capitalize}}"}, {Type: "text", Text: ", your score is "}, {Type: "variable", Variable: "score", Filters: []Filter{{Name: "round"}}, Text: "{{score|round}}"}, {Type: "text", Text: "."}},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

// TestFilterWithDoubleQuotesStringLiteralArguments covers various scenarios for filters with string literal arguments.
func TestFilterWithDoubleQuotesStringLiteralArguments(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			name:   "SingleFilterSingleStringArgument",
			source: `{{ name|append:"!" }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "name",
						Filters: []Filter{
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
						},
						Text: `{{ name|append:"!" }}`,
					},
				},
			},
		},
		{
			name:   "SingleFilterMultipleStringArguments",
			source: `{{ greeting|replace:"Hello","Hi" }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "greeting",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "Hello"}, StringArg{val: "Hi"}}},
						},
						Text: `{{ greeting|replace:"Hello","Hi" }}`,
					},
				},
			},
		},
		{
			name:   "MultipleFiltersSingleStringArgument",
			source: `{{ name|append:"!"|uppercase }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "name",
						Filters: []Filter{
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
							{Name: "uppercase"},
						},
						Text: `{{ name|append:"!"|uppercase }}`,
					},
				},
			},
		},
		{
			name:   "MultipleFiltersMultipleStringArguments",
			source: `{{ greeting|replace:"Hello","Hi"|append:" everyone" }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "greeting",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "Hello"}, StringArg{val: "Hi"}}},
							{Name: "append", Args: []FilterArg{StringArg{val: " everyone"}}},
						},
						Text: `{{ greeting|replace:"Hello","Hi"|append:" everyone" }}`,
					},
				},
			},
		},
		{
			name:   "MultipleVariablesWithFilters",
			source: `Hello {{name|capitalize|append:""}}, you have {{count|pluralize:"item","items"}}.`,
			expected: &Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello "},
					{
						Type:     "variable",
						Variable: "name",
						Filters:  []Filter{{Name: "capitalize"}, {Name: "append", Args: []FilterArg{StringArg{val: ""}}}},
						Text:     `{{name|capitalize|append:""}}`,
					},
					{Type: "text", Text: ", you have "},
					{
						Type:     "variable",
						Variable: "count",
						Filters:  []Filter{{Name: "pluralize", Args: []FilterArg{StringArg{val: "item"}, StringArg{val: "items"}}}},
						Text:     `{{count|pluralize:"item","items"}}`,
					},
					{Type: "text", Text: "."},
				},
			},
		},
		{
			name:   "MultipleVariablesNoDelimiter",
			source: `{{firstName}}{{lastName}}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "firstName",
						Text:     "{{firstName}}",
					},
					{
						Type:     "variable",
						Variable: "lastName",
						Text:     "{{lastName}}",
					},
				},
			},
		},
		{
			name:   "MultipleVariablesSpaceDelimiter",
			source: `{{ firstName }} {{ lastName }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "firstName",
						Text:     "{{ firstName }}",
					},
					{
						Type: "text",
						Text: " ",
					},
					{
						Type:     "variable",
						Variable: "lastName",
						Text:     "{{ lastName }}",
					},
				},
			},
		},
		{
			name:   "MultipleVariablesOtherDelimiters",
			source: `{{firstName}},{{lastName}}!`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "firstName",
						Text:     "{{firstName}}",
					},
					{
						Type: "text",
						Text: ",",
					},
					{
						Type:     "variable",
						Variable: "lastName",
						Text:     "{{lastName}}",
					},
					{
						Type: "text",
						Text: "!",
					},
				},
			},
		},
		{
			name:   "MultipleVariablesWithFiltersAndDelimiters",
			source: `{{firstName|replace:"Mr.",""|replace:"Mrs.",""}}, {{lastName|lower}}!`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "firstName",
						Filters:  []Filter{{Name: "replace", Args: []FilterArg{StringArg{val: "Mr."}, StringArg{val: ""}}}, {Name: "replace", Args: []FilterArg{StringArg{val: "Mrs."}, StringArg{val: ""}}}},
						Text:     `{{firstName|replace:"Mr.",""|replace:"Mrs.",""}}`,
					},
					{
						Type: "text",
						Text: ", ",
					},
					{
						Type:     "variable",
						Variable: "lastName",
						Filters:  []Filter{{Name: "lower"}},
						Text:     `{{lastName|lower}}`,
					},
					{
						Type: "text",
						Text: "!",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser()
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error for case '%s': %v", tc.name, err)
			}
			if !reflect.DeepEqual(tpl, tc.expected) {
				dump.P(tc.expected)
				t.Errorf("For case '%s', expected %+v, got %+v", tc.name, tc.expected, tpl)
			}
		})
	}
}

// TestFilterWithSingleQuotesStringLiteralArguments covers various scenarios for filters with string literal arguments.
func TestFilterWithSingleQuotesStringLiteralArguments(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			name:   "SingleFilterSingleStringArgument",
			source: `{{ name|append:'!' }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "name",
						Filters: []Filter{
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
						},
						Text: `{{ name|append:'!' }}`,
					},
				},
			},
		},
		{
			name:   "SingleFilterMultipleStringArguments",
			source: `{{ greeting|replace:'Hello','Hi' }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "greeting",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "Hello"}, StringArg{val: "Hi"}}},
						},
						Text: `{{ greeting|replace:'Hello','Hi' }}`,
					},
				},
			},
		},
		{
			name:   "MultipleFiltersSingleStringArgument",
			source: `{{ name|append:'!'|uppercase }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "name",
						Filters: []Filter{
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
							{Name: "uppercase"},
						},
						Text: `{{ name|append:'!'|uppercase }}`,
					},
				},
			},
		},
		{
			name:   "MultipleFiltersMultipleStringArguments",
			source: `{{ greeting|replace:'Hello','Hi'|append:' everyone' }}`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "greeting",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "Hello"}, StringArg{val: "Hi"}}},
							{Name: "append", Args: []FilterArg{StringArg{val: " everyone"}}},
						},
						Text: `{{ greeting|replace:'Hello','Hi'|append:' everyone' }}`,
					},
				},
			},
		},
		{
			name:   "MultipleVariablesWithFilters",
			source: `Hello {{name|capitalize|append:''}}, you have {{count|pluralize:'item','items'}}.`,
			expected: &Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello "},
					{
						Type:     "variable",
						Variable: "name",
						Filters:  []Filter{{Name: "capitalize"}, {Name: "append", Args: []FilterArg{StringArg{val: ""}}}},
						Text:     `{{name|capitalize|append:''}}`,
					},
					{Type: "text", Text: ", you have "},
					{
						Type:     "variable",
						Variable: "count",
						Filters:  []Filter{{Name: "pluralize", Args: []FilterArg{StringArg{val: "item"}, StringArg{val: "items"}}}},
						Text:     `{{count|pluralize:'item','items'}}`,
					},
					{Type: "text", Text: "."},
				},
			},
		},
		{
			name:   "MultipleVariablesWithFiltersAndDelimiters",
			source: `{{firstName|replace:'Mr.',''|replace:'Mrs.',''}}, {{lastName|lower}}!`,
			expected: &Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "firstName",
						Filters:  []Filter{{Name: "replace", Args: []FilterArg{StringArg{val: "Mr."}, StringArg{val: ""}}}, {Name: "replace", Args: []FilterArg{StringArg{val: "Mrs."}, StringArg{val: ""}}}},
						Text:     `{{firstName|replace:'Mr.',''|replace:'Mrs.',''}}`,
					},
					{
						Type: "text",
						Text: ", ",
					},
					{
						Type:     "variable",
						Variable: "lastName",
						Filters:  []Filter{{Name: "lower"}},
						Text:     `{{lastName|lower}}`,
					},
					{
						Type: "text",
						Text: "!",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser()
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error for case '%s': %v", tc.name, err)
			}
			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("For case '%s', expected %+v, got %+v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseFilterWithMultipleParameters(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		// Basic case with a single filter and multiple string literal arguments
		{
			"SingleFilterWithMultipleArgs",
			"{{ value|replace:'hello','world' }}",
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "value",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "hello"}, StringArg{val: "world"}}},
						},
						Text: `{{ value|replace:'hello','world' }}`,
					},
				},
			},
		},
		// Spaces around arguments
		{
			"SpacesAroundArguments",
			"{{ value|replace: 'hello', 'world' }}",
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "value",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "hello"}, StringArg{val: "world"}}},
						},
						Text: `{{ value|replace: 'hello', 'world' }}`,
					},
				},
			},
		},
		// Multiple filters with multiple arguments
		{
			"MultipleFiltersWithMultipleArgs",
			"{{ greeting|replace:'Hello','Hi'|append: '!', ' Have a great day' }}",
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "greeting",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "Hello"}, StringArg{val: "Hi"}}},
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}, StringArg{val: " Have a great day"}}},
						},
						Text: `{{ greeting|replace:'Hello','Hi'|append: '!', ' Have a great day' }}`,
					},
				},
			},
		},
		// Complex scenario with mixed text and multiple variables with filters and multiple arguments
		{
			"ComplexMixedTextAndVariables",
			`Hello {{ name|capitalize }}, you have {{ unread|pluralize:"message","messages" }}.`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello "},
					{
						Type:     "variable",
						Variable: "name",
						Filters:  []Filter{{Name: "capitalize"}},
						Text:     "{{ name|capitalize }}",
					},
					{Type: "text", Text: ", you have "},
					{
						Type:     "variable",
						Variable: "unread",
						Filters:  []Filter{{Name: "pluralize", Args: []FilterArg{StringArg{val: "message"}, StringArg{val: "messages"}}}},
						Text:     `{{ unread|pluralize:"message","messages" }}`,
					},
					{Type: "text", Text: "."},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseVariableWithMultiplePipelineFiltersWithMultipleParameters(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		// Original case
		{
			"ReplaceAndAppend",
			`{{ username|replace:"hello","world"|append:"!" }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "username",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "hello"}, StringArg{val: "world"}}},
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
						},
						Text: `{{ username|replace:"hello","world"|append:"!" }}`,
					},
				},
			},
		},
		// Additional case with space around pipe symbols
		{
			"SpacesAroundPipes",
			`{{ username | replace:"hello","world" | append:"!" }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "username",
						Filters: []Filter{
							{Name: "replace", Args: []FilterArg{StringArg{val: "hello"}, StringArg{val: "world"}}},
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
						},
						Text: `{{ username | replace:"hello","world" | append:"!" }}`,
					},
				},
			},
		},
		// Multiple filters with varied arguments
		{
			"MultipleFiltersVariedArgs",
			`{{ date|date:"YYYY-MM-DD"|prepend:"Date: " }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "date",
						Filters: []Filter{
							{Name: "date", Args: []FilterArg{StringArg{val: "YYYY-MM-DD"}}},
							{Name: "prepend", Args: []FilterArg{StringArg{val: "Date: "}}},
						},
						Text: `{{ date|date:"YYYY-MM-DD"|prepend:"Date: " }}`,
					},
				},
			},
		},
		// Complex scenario mixing text and multiple variables with multiple filters
		{
			"MixedTextMultipleVarsFilters",
			`Hello {{ name|capitalize|append:"!" }}, you have {{ unread|pluralize:"1 message","%d messages"|replace:"%d","many" }}.`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Hello "},
					{
						Type:     "variable",
						Variable: "name",
						Filters: []Filter{
							{Name: "capitalize"},
							{Name: "append", Args: []FilterArg{StringArg{val: "!"}}},
						},
						Text: `{{ name|capitalize|append:"!" }}`,
					},
					{Type: "text", Text: ", you have "},
					{
						Type:     "variable",
						Variable: "unread",
						Filters: []Filter{
							{Name: "pluralize", Args: []FilterArg{StringArg{val: "1 message"}, StringArg{val: "%d messages"}}},
							{Name: "replace", Args: []FilterArg{StringArg{val: "%d"}, StringArg{val: "many"}}},
						},
						Text: `{{ unread|pluralize:"1 message","%d messages"|replace:"%d","many" }}`,
					},
					{Type: "text", Text: "."},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseMultipleAdjacentVariables(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"TwoVariablesNoSpace",
			"{{firstName}}{{lastName}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Text: "{{firstName}}"},
					{Type: "variable", Variable: "lastName", Text: "{{lastName}}"},
				},
			},
		},
		{
			"TwoVariablesSpaceBetween",
			"{{firstName}} {{lastName}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Text: "{{firstName}}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "lastName", Text: "{{lastName}}"},
				},
			},
		},
		{
			"ThreeVariablesMixedWhitespace",
			"{{firstName}}\t{{lastName}}\n{{email}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Text: "{{firstName}}"},
					{Type: "text", Text: "\t"},
					{Type: "variable", Variable: "lastName", Text: "{{lastName}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "email", Text: "{{email}}"},
				},
			},
		},
		{
			"FourVariablesLineBreaks",
			"{{firstName}}\n{{lastName}}\n{{email}}\n{{username}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Text: "{{firstName}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "lastName", Text: "{{lastName}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "email", Text: "{{email}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "username", Text: "{{username}}"},
				},
			},
		},
		{
			"TwoVariablesOneFilter",
			"{{firstName|upper}}{{lastName}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "upper"}}, Text: "{{firstName|upper}}"},
					{Type: "variable", Variable: "lastName", Text: "{{lastName}}"},
				},
			},
		},
		{
			"VariablesWithFilterAndArgument",
			"{{user|default:'Anonymous'}}{{age|default:18}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "user", Filters: []Filter{{Name: "default", Args: []FilterArg{StringArg{val: "Anonymous"}}}}, Text: "{{user|default:'Anonymous'}}"},
					{Type: "variable", Variable: "age", Filters: []Filter{{Name: "default", Args: []FilterArg{NumberArg{val: 18}}}}, Text: "{{age|default:18}}"},
				},
			},
		},
		{
			"VariablesFilterPipeline",
			"{{firstName|trim|capitalize}}{{lastName|lower}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "trim"}, {Name: "capitalize"}}, Text: "{{firstName|trim|capitalize}}"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "lower"}}, Text: "{{lastName|lower}}"},
				},
			},
		},
		{
			"ThreeVariablesWithMultipleFilters",
			"{{firstName|trim}}{{lastName|lower|capitalize}}{{age|default:30}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "trim"}}, Text: "{{firstName|trim}}"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "lower"}, {Name: "capitalize"}}, Text: "{{lastName|lower|capitalize}}"},
					{Type: "variable", Variable: "age", Filters: []Filter{{Name: "default", Args: []FilterArg{NumberArg{val: 30}}}}, Text: "{{age|default:30}}"},
				},
			},
		},
		{
			"TwoVariablesTabBetweenWithFilter",
			"{{firstName|capitalize}}\t{{lastName|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "capitalize"}}, Text: "{{firstName|capitalize}}"},
					{Type: "text", Text: "\t"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "upper"}}, Text: "{{lastName|upper}}"},
				},
			},
		},
		{
			"VariablesWithLineBreakAndFilter",
			"{{firstName|lower}}\n{{lastName|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "lower"}}, Text: "{{firstName|lower}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "upper"}}, Text: "{{lastName|upper}}"},
				},
			},
		},
		{
			"ThreeVariablesSpaceAndFilter",
			"{{firstName|upper}} {{lastName|upper}} {{email|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "upper"}}, Text: "{{firstName|upper}}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "upper"}}, Text: "{{lastName|upper}}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "email", Filters: []Filter{{Name: "upper"}}, Text: "{{email|upper}}"},
				},
			},
		},
		{
			"NestedVariablesWithFilters",
			"{{user.firstName|capitalize}}{{user.lastName|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "user.firstName", Filters: []Filter{{Name: "capitalize"}}, Text: "{{user.firstName|capitalize}}"},
					{Type: "variable", Variable: "user.lastName", Filters: []Filter{{Name: "upper"}}, Text: "{{user.lastName|upper}}"},
				},
			},
		},
		{
			"VariablesWithMultipleFiltersAndWhitespace",
			"{{ firstName | trim | capitalize }}\n{{ lastName | lower | trim }}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "trim"}, {Name: "capitalize"}}, Text: "{{ firstName | trim | capitalize }}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "lower"}, {Name: "trim"}}, Text: "{{ lastName | lower | trim }}"},
				},
			},
		},
		{
			"FourVariablesWithMixedFiltersAndWhitespace",
			"{{firstName|capitalize}} {{lastName|lower}}\t{{email|upper}}\n{{username|reverse}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "capitalize"}}, Text: "{{firstName|capitalize}}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "lower"}}, Text: "{{lastName|lower}}"},
					{Type: "text", Text: "\t"},
					{Type: "variable", Variable: "email", Filters: []Filter{{Name: "upper"}}, Text: "{{email|upper}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "username", Filters: []Filter{{Name: "reverse"}}, Text: "{{username|reverse}}"},
				},
			},
		},
		{
			"VariablesWithSpecialCharactersAndFilters",
			"{{firstName|capitalize|replace:'John','Jonathan'}}{{lastName|append:' Smith'}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "capitalize"}, {Name: "replace", Args: []FilterArg{StringArg{val: "John"}, StringArg{val: "Jonathan"}}}}, Text: "{{firstName|capitalize|replace:'John','Jonathan'}}"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "append", Args: []FilterArg{StringArg{val: " Smith"}}}}, Text: "{{lastName|append:' Smith'}}"},
				},
			},
		},
		{
			"NestedAndComplexFilters",
			"{{user.details.address.city|capitalize}}\n{{user.details.phoneNumber|default:'N/A'}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "user.details.address.city", Filters: []Filter{{Name: "capitalize"}}, Text: "{{user.details.address.city|capitalize}}"},
					{Type: "text", Text: "\n"},
					{Type: "variable", Variable: "user.details.phoneNumber", Filters: []Filter{{Name: "default", Args: []FilterArg{StringArg{val: "N/A"}}}}, Text: "{{user.details.phoneNumber|default:'N/A'}}"},
				},
			},
		},
		{
			"ComplexNestedVariablesWithMultipleFilters",
			"{{user.address|trim}} {{user.phone|default:'Unknown'|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "user.address", Filters: []Filter{{Name: "trim"}}, Text: "{{user.address|trim}}"},
					{Type: "text", Text: " "},
					{Type: "variable", Variable: "user.phone", Filters: []Filter{{Name: "default", Args: []FilterArg{StringArg{val: "Unknown"}}}, {Name: "upper"}}, Text: "{{user.phone|default:'Unknown'|upper}}"},
				},
			},
		},
		{
			"MultipleAdjacentVariablesWithMixedFilters",
			"{{firstName|lower}}{{middleName|capitalize}}{{lastName|upper}}",
			&Template{
				Nodes: []*Node{
					{Type: "variable", Variable: "firstName", Filters: []Filter{{Name: "lower"}}, Text: "{{firstName|lower}}"},
					{Type: "variable", Variable: "middleName", Filters: []Filter{{Name: "capitalize"}}, Text: "{{middleName|capitalize}}"},
					{Type: "variable", Variable: "lastName", Filters: []Filter{{Name: "upper"}}, Text: "{{lastName|upper}}"},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %+v, got %+v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseVariableWithFilterHavingCommaInArguments(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"DateFilterWithComma",
			`{{ current | date:"F j, Y" }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "current",
						Filters: []Filter{
							{Name: "date", Args: []FilterArg{StringArg{val: "F j, Y"}}},
						},
						Text: `{{ current | date:"F j, Y" }}`,
					},
				},
			},
		},
		{
			"DateFilterWithQuotedComma",
			`{{ current | date:'F j, Y' }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "current",
						Filters: []Filter{
							{Name: "date", Args: []FilterArg{StringArg{val: "F j, Y"}}},
						},
						Text: `{{ current | date:'F j, Y' }}`,
					},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %+v, got %+v", tc.name, tc.expected, tpl)
			}
		})
	}
}

func TestParseMalformedVariableNodeAsText(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{
			"MissingClosingBracket",
			"Welcome back, {{username",
		},
		{
			"MissingOpeningBracket",
			"Hello, username}}!",
		},
		{
			"UnfinishedFilter",
			"Your account balance is {{balance|}} today.",
		},
		{
			"PipeWithoutFilterName",
			"Good morning, {{name| . Have a nice day!",
		},
		{
			"MissingFilterNameWithArguments",
			"Record: {{record||upper}}",
		},
		{
			"NestedBracesMalformed",
			"Error: {{user.details.{name}}",
		},
		{
			"MissingVariableName",
			"New message: {{|capitalize}}",
		},
		{
			"MalformedWithTextAround",
			"Hello, {{user|trim} in the system.",
		},
		{
			"MultipleMalformedInText",
			"Start {{of something |middle|end}} incomplete.",
		},
		{
			"SpaceBeforeClosingBracket",
			"Attempt: {{attempt | }}",
		},
		{
			"RandomCharactersInBraces",
			"Code: {{1234abcd!@#$}}",
		},
		{
			"MalformedFilterSyntax",
			"Discount: {{price|*0.85}}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser()
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			// Expecting the problematic part, or in some cases the entire input, to be treated as a text node
			expected := &Template{
				Nodes: []*Node{
					{Type: "text", Text: tc.source},
				},
			}

			if !reflect.DeepEqual(tpl, expected) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, expected, tpl)
			}
		})
	}
}

func TestParserWithMultipleFiltersAndNumericArguments(t *testing.T) {
	cases := []struct {
		name     string
		source   string
		expected *Template
	}{
		{
			"MultipleFiltersOnSingleVariable",
			`{{ price|plus:10|minus:5|times:2|divide:3|round }}`,
			&Template{
				Nodes: []*Node{
					{
						Type:     "variable",
						Variable: "price",
						Filters: []Filter{
							{Name: "plus", Args: []FilterArg{NumberArg{val: 10}}},
							{Name: "minus", Args: []FilterArg{NumberArg{val: 5}}},
							{Name: "times", Args: []FilterArg{NumberArg{val: 2}}},
							{Name: "divide", Args: []FilterArg{NumberArg{val: 3}}},
							{Name: "round"},
						},
						Text: `{{ price|plus:10|minus:5|times:2|divide:3|round }}`,
					},
				},
			},
		},
		{
			"MultipleVariablesAndFilters",
			`Total: {{ price|plus:taxes|minus:discount }} and {{ shipping|plus:5 }}`,
			&Template{
				Nodes: []*Node{
					{Type: "text", Text: "Total: "},
					{
						Type:     "variable",
						Variable: "price",
						Filters: []Filter{
							{Name: "plus", Args: []FilterArg{VariableArg{name: "taxes"}}},
							{Name: "minus", Args: []FilterArg{VariableArg{name: "discount"}}},
						},
						Text: `{{ price|plus:taxes|minus:discount }}`,
					},
					{Type: "text", Text: " and "},
					{
						Type:     "variable",
						Variable: "shipping",
						Filters: []Filter{
							{Name: "plus", Args: []FilterArg{NumberArg{val: 5}}},
						},
						Text: `{{ shipping|plus:5 }}`,
					},
				},
			},
		},
	}

	parser := NewParser()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := parser.Parse(tc.source)
			if err != nil {
				t.Fatalf("Unexpected error in %s: %v", tc.name, err)
			}

			if !reflect.DeepEqual(tpl, tc.expected) {
				t.Errorf("Case %s: Expected %+v, got %+v", tc.name, tc.expected, tpl)
			}
		})
	}
}
