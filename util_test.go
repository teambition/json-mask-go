package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilSuite struct {
	suite.Suite
}

func (s *UtilSuite) TestgetFiledNamesByJSONKeys() {
	fieldNames, err := getFiledNamesByJSONKeys(testStruct{A: 11, N: "nnn", C: 44, G: "ggg"}, []string{"a", "n", "fff"})

	s.Nil(err)
	s.Len(fieldNames, 2)
}

func (s *UtilSuite) TestgetFiledNamesByJSONKeysAnonymousStruct() {
	fieldNames, err := getFiledNamesByJSONKeys(anonymousStruct{BInner: &BInner{}}, []string{"a", "n", "d", "g"})

	s.Nil(err)
	s.Len(fieldNames, 2)
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilSuite))
}
