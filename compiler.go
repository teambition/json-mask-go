package jsonmask

import (
	"github.com/pkg/errors"
)

// Errors exposed.
var (
	ErrEmptyString = errors.New("empty string")
)

var terminalsMap = map[string]int{",": 1, "/": 2, "(": 3, ")": 4}

type valueType int

const (
	typeObject valueType = iota
	typeArray
)

const keyAny = "*"

// Node represents a grammar node.
type Node struct {
	key   string
	typ   valueType
	props []Node
}

type token struct {
	tag   string
	value string
}

// Compile compiles the given mask text to the json mask nodes.
func Compile(text string) ([]Node, error) {
	if text == "" {
		return nil, errors.WithStack(ErrEmptyString)
	}

	return nil, nil
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
