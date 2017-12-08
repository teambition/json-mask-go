package jsonmask_test

import (
	"encoding/json"
	"fmt"

	"github.com/DavidCai1993/json-mask-go"
)

func ExampleMask() {
	type ExampleStruct struct {
		A string             `json:"a"`
		B int                `json:"b"`
		C map[string]float32 `json:"c"`
	}

	result, _ := jsonmask.Mask(ExampleStruct{
		A: "aaa",
		B: 234,
		C: map[string]float32{"c1": 12, "c2": 33},
	}, "a,c")

	j, _ := json.Marshal(result)

	fmt.Println("json output: ", string(j))
	// json output:  {"a":"aaa","c":{"c1":12,"c2":33}}
}
