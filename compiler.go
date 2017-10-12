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

// Node represents a grammar node.
type Node struct {
	key        string
	typ        valueType
	properties *Node
}

// Compile compiles the given mask text to the json mask nodes.
func Compile(text string) (*Node, error) {
	if text == "" {
		return nil, errors.WithStack(ErrEmptyString)
	}

	return nil, nil
}
