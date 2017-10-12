package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var tests = map[string][]Node{
	"a": []Node{
		Node{key: "a", typ: typeObject},
	},
	"a,b,c": []Node{
		Node{key: "a", typ: typeObject},
		Node{key: "b", typ: typeObject},
		Node{key: "c", typ: typeObject},
	},
	"a/*/c": []Node{
		Node{key: "a", typ: typeObject, props: []Node{
			Node{
				key:   keyAny,
				typ:   typeObject,
				props: []Node{Node{key: "c", typ: typeObject}},
			},
		}},
	},
	"a,b(d/*/g,b),c": []Node{
		Node{key: "a", typ: typeObject},
		Node{key: "b", typ: typeArray, props: []Node{
			Node{key: "d", typ: typeObject, props: []Node{
				Node{key: keyAny, typ: typeObject, props: []Node{
					Node{key: "g", typ: typeObject},
				}},
			}},
			Node{key: "b", typ: typeObject},
		}},
		Node{key: "c", typ: typeObject},
	},
}

type CompilerSuite struct {
	suite.Suite
}

func (s *CompilerSuite) TestCompileEmptyString() {
	_, err := Compile("")
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

func TestCompiler(t *testing.T) {
	suite.Run(t, new(CompilerSuite))
}
