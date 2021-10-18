package jsonmask

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Selection map[string]Selection

// used for testing purposes
func (s Selection) equal(other Selection) bool {
	if other == nil {
		return false
	}
	if len(s) != len(other) {
		return false
	}
	for k, v := range s {
		if !v.equal(other[k]) {
			return false
		}
	}
	return true
}

// The syntax is loosely based on XPath:
//
// a       select a field 'a'
// a,b,c   comma-separated list will select multiple fields
// a/b/c   path will select a field from its parent
// a(b,c)  sub-selection will select many fields from a parent
// a/*/c   the star * wildcard will select all items in a field
// a,b/c(d,e(f,g/h)),i
//
func Compile(str string) (Selection, error) {
	if !utf8.ValidString(str) {
		return nil, fmt.Errorf("invalid fields")
	}
	tokens := make([]string, 0, len(str)/2)
	var state uint8
	var b strings.Builder
	b.Grow(len(str) / 2)
	for _, r := range str {
		switch r {
		case ',':
			switch state {
			case 4:
				tokens = append(tokens, b.String())
			case 3:
				// nothing to do
			default:
				return nil, fmt.Errorf("invalid char before ','")
			}
			tokens = append(tokens, ",")
			b.Reset()
			state = 1
		case '/':
			if state != 4 {
				return nil, fmt.Errorf("invalid char before '/'")
			}
			tokens = append(tokens, b.String())
			tokens = append(tokens, "/")
			b.Reset()
			state = 2
		case '(':
			if state != 4 {
				return nil, fmt.Errorf("invalid char before '('")
			}
			tokens = append(tokens, b.String())
			tokens = append(tokens, "(")
			b.Reset()
			state = 2
		case ')':
			switch state {
			case 4:
				tokens = append(tokens, b.String())
			case 3:
				// nothing to do
			default:
				return nil, fmt.Errorf("invalid char before ')'")
			}
			tokens = append(tokens, ")")
			b.Reset()
			state = 3
		default:
			if state == 3 {
				return nil, fmt.Errorf("invalid ')' before a char")
			}
			b.WriteRune(r)
			state = 4
		}
	}

	switch state {
	case 4:
		tokens = append(tokens, b.String())
	case 3:
		// nothing to do
	default:
		return nil, fmt.Errorf("invalid end")
	}

	node := make(Selection)
	err := buildSelection(tokens, node)
	return node, err
}

func buildSelection(tokens []string, root Selection) error {
	if len(tokens) == 0 {
		return nil
	}

	var child Selection
	node := root
	for i := 0; i < len(tokens); i++ {
		switch tokens[i] {
		case ",":
			node = root
		case "/":
			node = child
		case "(":
			end := findCloseIndex(tokens, i+1)
			if end == -1 {
				return fmt.Errorf("sub-selector not close")
			}
			if err := buildSelection(tokens[i+1:end], child); err != nil {
				return err
			}
			i = end
		case ")":
			return fmt.Errorf("invalid field char: ')'")
		default:
			child = make(Selection)
			node[tokens[i]] = child
		}
	}
	return nil
}

func findCloseIndex(tokens []string, start int) int {
	for i := len(tokens) - 1; i >= start; i-- {
		if tokens[i] == ")" {
			return i
		}
	}
	return -1
}
