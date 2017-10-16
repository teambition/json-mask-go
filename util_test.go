package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/oleiade/reflections.v1"
)

type UtilSuite struct {
	suite.Suite
}

type testUtilStruct struct {
	A int    `json:"a"`
	N string `json:"n"`
	C int    `json:"c"`
	G string `json:"g"`
}

func (s *UtilSuite) TestClearJSONTag() {
	t := testUtilStruct{A: 1, N: "2", C: 3, G: "4"}

	newT, err := clearJSONTag(t, "A")
	s.Nil(err)

	tagValueA, err := reflections.GetFieldTag(newT, "A", "json")
	s.Nil(err)
	s.Equal("-", tagValueA)

	tagValueG, err := reflections.GetFieldTag(newT, "G", "json")
	s.Nil(err)
	s.Equal("g", tagValueG)
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
