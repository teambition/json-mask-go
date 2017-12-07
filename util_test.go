package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) TestGetFiledByJSONKey() {
	field, ok := getFiledByJSONKey(testStruct{A: 11, N: "nnn", C: 44, G: "ggg"}, "a")

	s.True(ok)
	s.Equal(field.Name, "A")
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
