package jsonmask

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FilterSuite struct {
	suite.Suite
}

type testStruct struct {
	A int    `json:"a"`
	N string `json:"n"`
	C int    `json:"c"`
	G string `json:"g"`
}

func (s *FilterSuite) TestFilter() {
	mask := nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
		"c": node{typ: typeObject, props: nodeMap{}},
	}

	s.NotNil(mask)

	res, err := filter(testStruct{A: 11, N: "nnn", C: 44, G: "ggg"}, mask)
	s.Nil(err)

	j, err := json.Marshal(res)
	s.Nil(err)

	s.Equal(string(j), `{"a":11,"c":44}`)
}

func TestFilter(t *testing.T) {
	suite.Run(t, new(FilterSuite))
}
