package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CompilerSuite struct {
	suite.Suite
}

func (s *CompilerSuite) TestCompileEmptyString() {
	_, err := Compile("")
	s.Error(err, ErrEmptyString)
}

func TestCompiler(t *testing.T) {
	suite.Run(t, new(CompilerSuite))
}
