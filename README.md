# json-mask-go
[![Build Status](https://travis-ci.org/DavidCai1993/json-mask-go.svg?branch=master)](https://travis-ci.org/DavidCai1993/json-mask-go)
[![Coverage Status](https://coveralls.io/repos/github/DavidCai1993/json-mask-go/badge.svg?branch=master)](https://coveralls.io/github/DavidCai1993/json-mask-go?branch=master)

JSON mask for Go

> Warning: for now only support filtering top level properties.

## Installation

```
go get -u github.com/DavidCai1993/json-mask-go
```

## Documentation

API documentation can be found here: https://godoc.org/github.com/DavidCai1993/json-mask-go


## Usage

```go
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
```
