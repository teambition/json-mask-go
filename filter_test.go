package jsonmask

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FilterSuite struct {
	suite.Suite
}

type testStruct struct {
	A int    `json:"a"`
	N string `json:"n"`
	B []struct {
		D struct {
			G struct {
				Z int `json:"z"`
			} `json:"g"`
			B int `json:"b"`
			C struct {
				A int `json:"a"`
			} `json:"c"`
		} `json:"d"`
		B []struct {
			Z int `json:"z"`
		} `json:"b"`
		K int `json:"k"`
	} `json:"b"`
	C int    `json:"c"`
	G string `json:"g"`
}

func (s *FilterSuite) TestFilter() {
	// mask := nodeMap{
	// 	"a": node{typ: typeObject, props: nodeMap{}},
	// 	"b": node{
	// 		typ: typeArray,
	// 		props: nodeMap{
	// 			"d": node{typ: typeObject, props: nodeMap{
	// 				keyAny: node{
	// 					typ: typeObject,
	// 					props: nodeMap{
	// 						"z": node{typ: typeObject, props: nodeMap{}},
	// 					},
	// 				},
	// 			}},
	// 			"b": node{typ: typeObject, props: nodeMap{}},
	// 		},
	// 	},
	// 	"c": node{typ: typeObject, props: nodeMap{}},
	// }

	// obj := testStruct{
	// 	A: 11,
	// 	N: "nnn",
	// 	C: 44,
	// 	G: "ggg",
	// }

	// obj.B[0].D.G.Z = 22
	// obj.B[0].D.B = 34
	// obj.B[0].D.C.A = 32
	// obj.B[0].B[0].Z = 33
	// obj.B[0].K = 99
}

func TestFilter(t *testing.T) {
	suite.Run(t, new(FilterSuite))
}
