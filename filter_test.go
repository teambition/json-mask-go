package jsonmask

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FilterSuite struct {
	suite.Suite
}

type bInner struct {
	D dInner   `json:"d"`
	B []gInner `json:"b"`
	K int      `json:"k"`
}

type dInner struct {
	G gInner `json:"g"`
	B int    `json:"b"`
	C cInner `json:"c"`
}

type gInner struct {
	Z int `json:"z"`
}

type cInner struct {
	A int `json:"a"`
}

type testStruct struct {
	A int      `json:"a"`
	N string   `json:"n"`
	B []bInner `json:"b"`
	C int      `json:"c"`
	G string   `json:"g"`
}

func (s *FilterSuite) TestFilterSimpleObject() {
	mask := nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
		"n": node{typ: typeObject, props: nodeMap{}},
	}

	s.NotNil(mask)

	res, err := filter(testStruct{A: 11, N: "nnn", C: 44, G: "ggg"}, mask)
	s.Nil(err)

	j, err := json.Marshal(res)
	s.Nil(err)

	s.Equal(string(j), `{"a":11,"n":"nnn"}`)
}

func (s *FilterSuite) TestFilterComplexObject() {
	mask := nodeMap{
		"a": node{typ: typeObject, props: nodeMap{}},
		"c": node{typ: typeObject, props: nodeMap{}},
		"b": node{typ: typeArray, props: nodeMap{
			"d": node{typ: typeObject, props: nodeMap{
				keyAny: node{typ: typeObject, props: nodeMap{
					"z": node{typ: typeObject, props: nodeMap{}},
				}},
			}},
			"b": node{typ: typeArray, props: nodeMap{
				"g": node{typ: typeObject, props: nodeMap{}},
			}},
		}},
	}

	obj := testStruct{
		A: 11,
		N: "nn",
		C: 44,
		G: "gg",
		B: []bInner{
			bInner{
				K: 99,
				B: []gInner{gInner{Z: 33}},
				D: dInner{G: gInner{Z: 22}, B: 34, C: cInner{A: 32}},
			},
		},
	}

	res, err := filter(obj, mask)
	s.Nil(err)

	j, err := json.Marshal(res)
	s.Nil(err)

	fmt.Println("json: ", j)
}

func TestFilter(t *testing.T) {
	suite.Run(t, new(FilterSuite))
}
