package jsonmask

import (
	"github.com/pkg/errors"
)

// Errors exposed.
var (
	ErrEmptyString = errors.New("empty string")
)

type valueType int

const (
	typeObject valueType = iota
	typeArray
)

const keyAny = "*"

type node struct {
	typ   valueType
	props nodeMap
}

type token struct {
	tag   string
	value string
	node
}

type nodeMap map[string]node

// compile compiles the fieldmask text (example: 'a,b.c') to
// the nodeMap sturct.
// For more information:
// https://developers.google.com/discovery/v1/performance#partial-response
func compile(text string) (nodeMap, error) {
	if text == "" {
		return nil, errors.WithStack(ErrEmptyString)
	}

	return parse(scan([]rune(text))), nil
}

func scan(text []rune) []token {
	var (
		tokens = []token{}
		name   = ""
	)

	var maybePush = func() {
		if name == "" {
			return
		}
		tokens = append(tokens, token{tag: "_n", value: name})
		name = ""
	}

	for i := 0; i < len(text); i++ {
		ch := string(text[i])

		switch ch {
		case ",", "/", "(", ")":
			maybePush()
			tokens = append(tokens, token{tag: ch})
		default:
			name += ch
		}
	}

	maybePush()

	return tokens
}

func parse(tokens []token) nodeMap {
	return buildTree(&tokens, &token{}, &[]token{})
}

func buildTree(tokens *[]token, parent *token, stack *[]token) nodeMap {
	var (
		t     = token{}
		props = nodeMap{}
	)

	for {
		if len(*tokens) == 0 {
			break
		}

		t = (*tokens)[0]
		*tokens = (*tokens)[1:]

		if t.tag == "_n" {
			t.typ = typeObject
			t.props = buildTree(tokens, &t, stack)

			if len(*stack) != 0 {
				peek := (*stack)[len(*stack)-1]
				if peek.tag == "/" {
					*stack = (*stack)[:len(*stack)-1]
					addToken(t, props)
					return props
				}
			}
		} else if t.tag == "," {
			return props
		} else if t.tag == "(" {
			*stack = append(*stack, t)
			parent.typ = typeArray
			continue
		} else if t.tag == ")" {
			*stack = (*stack)[:len(*stack)-1]
			return props
		} else if t.tag == "/" {
			*stack = append(*stack, t)
			continue
		}

		addToken(t, props)
	}

	return props
}

func addToken(t token, props nodeMap) {
	n := node{typ: t.typ}

	if t.props != nil {
		n.props = t.props
	}

	props[t.value] = n
}
