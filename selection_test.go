package jsonmask

import "testing"

type selectionCase struct {
	fields    string
	shouldErr bool
	res       Selection
}

// a       select a field 'a'
// a,b,c   comma-separated list will select multiple fields
// a/b/c   path will select a field from its parent
// a(b,c)  sub-selection will select many fields from a parent
// a/*/c   the star * wildcard will select all items in a field
// a,b/c(d,e(f,g/h)),i
var selectionCases = []selectionCase{
	{
		fields:    "/",
		shouldErr: true,
	},
	{
		fields:    "(",
		shouldErr: true,
	},
	{
		fields:    ")",
		shouldErr: true,
	},
	{
		fields:    ",",
		shouldErr: true,
	},
	{
		fields:    string([]byte("非utf8")[1:]),
		shouldErr: true,
	},
	{
		fields:    "a/",
		shouldErr: true,
	},
	{
		fields:    "a(",
		shouldErr: true,
	},
	{
		fields:    "a)",
		shouldErr: true,
	},
	{
		fields:    "a,",
		shouldErr: true,
	},
	{
		fields:    "a(b",
		shouldErr: true,
	},
	{
		fields:    "a(b))",
		shouldErr: true,
	},
	{
		fields:    "a(b(c)",
		shouldErr: true,
	},
	{
		fields:    "a(b(c),)",
		shouldErr: true,
	},
	{
		fields:    "a(b/(c))",
		shouldErr: true,
	},
	{
		fields:    "a(b(c)d)",
		shouldErr: true,
	},
	{
		fields:    "a/(b,c)/d",
		shouldErr: true,
	},
	{
		fields:    "a",
		shouldErr: false,
		res:       Selection{"a": Selection{}},
	},
	{
		fields:    "*",
		shouldErr: false,
		res:       Selection{"*": Selection{}},
	},
	{
		fields:    "a,b,c",
		shouldErr: false,
		res: Selection{
			"a": Selection{},
			"b": Selection{},
			"c": Selection{},
		},
	},
	{
		fields:    "a/b/c",
		shouldErr: false,
		res:       Selection{"a": Selection{"b": Selection{"c": Selection{}}}},
	},
	{
		fields:    "a(b,c)",
		shouldErr: false,
		res:       Selection{"a": Selection{"b": Selection{}, "c": Selection{}}},
	},
	{
		fields:    "a(b(c))",
		shouldErr: false,
		res:       Selection{"a": Selection{"b": Selection{"c": Selection{}}}},
	},
	{
		fields:    "a/*/c",
		shouldErr: false,
		res:       Selection{"a": Selection{"*": Selection{"c": Selection{}}}},
	},
	{
		fields:    "a,b/c(d,e(f,g/h)),i",
		shouldErr: false,
		res: Selection{
			"a": Selection{},
			"b": Selection{"c": Selection{
				"d": Selection{},
				"e": Selection{
					"f": Selection{},
					"g": Selection{"h": Selection{}},
				},
			}},
			"i": Selection{},
		},
	},
}

func TestCompile(t *testing.T) {
	for _, c := range selectionCases {
		res, err := Compile(c.fields)
		if c.shouldErr {
			if err == nil {
				t.Errorf("Testing case[%s] failed: should error but got: %#v", c.fields, res)
			}
		} else if err != nil {
			t.Errorf("Testing case[%s] failed: %s", c.fields, err)
		} else if !c.res.equal(res) {
			t.Errorf("Testing case[%s] failed, expected: %#v, got: %#v", c.fields, c.res, res)
		}
	}
}
