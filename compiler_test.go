package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var tests = map[string]nodeMap{
	"a": nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
	},
	"a,b,c": nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
		"b": node{typ: typeObject, props: nodeMap{}},
		"c": node{typ: typeObject, props: nodeMap{}},
	},
	"a/*/c": nodeMap{
		"a": node{
			typ: typeObject,
			props: nodeMap{
				keyAny: node{
					typ: typeObject, props: nodeMap{"c": node{typ: typeObject, props: nodeMap{}}},
				},
			},
		},
	},
	"a,b(d/*/g,b),c": nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
		"b": node{
			typ: typeArray,
			props: nodeMap{
				"d": node{typ: typeObject, props: nodeMap{
					keyAny: node{
						typ: typeObject,
						props: nodeMap{
							"g": node{typ: typeObject, props: nodeMap{}},
						},
					},
				}},
				"b": node{typ: typeObject, props: nodeMap{}},
			},
		},
		"c": node{typ: typeObject, props: nodeMap{}},
	},
}

type CompilerSuite struct {
	suite.Suite
}

func (s *CompilerSuite) TestCompileEmptyString() {
	_, err := compile("")
	s.Error(err, ErrEmptyString)
}

func (s *CompilerSuite) TestScanOneProp() {
	tokens := scan([]rune("a"))

	s.Equal(len(tokens), 1)
	s.Equal(tokens[0].tag, "_n")
	s.Equal(tokens[0].value, "a")
}

func (s *CompilerSuite) TestScanThreeProps() {
	tokens := scan([]rune("a,b,c"))

	s.Equal(len(tokens), 5)
	s.Equal(tokens[0].tag, "_n")
	s.Equal(tokens[0].value, "a")
	s.Equal(tokens[1].tag, ",")
	s.Equal(tokens[1].value, "")
	s.Equal(tokens[2].tag, "_n")
	s.Equal(tokens[2].value, "b")
	s.Equal(tokens[3].tag, ",")
	s.Equal(tokens[3].value, "")
	s.Equal(tokens[4].tag, "_n")
	s.Equal(tokens[4].value, "c")
}

func (s *CompilerSuite) TestScanMoreProps() {
	tokens := scan([]rune("a,b(d/*/g,b),c"))
	s.Equal(len(tokens), 14)
}

func (s *CompilerSuite) TestCompile() {
	for text, expectedRes := range tests {
		res, err := compile(text)
		s.Nil(err)
		s.Equal(expectedRes, res)
	}
}

func TestCompiler(t *testing.T) {
	suite.Run(t, new(CompilerSuite))
}
