# json-mask-go
[![Build Status](https://travis-ci.org/DavidCai1993/json-mask-go.svg?branch=master)](https://travis-ci.org/DavidCai1993/json-mask-go)
[![Coverage Status](https://coveralls.io/repos/github/DavidCai1993/json-mask-go/badge.svg?branch=master)](https://coveralls.io/github/DavidCai1993/json-mask-go?branch=master)

JSON mask for Go, selecting specific parts of JSON string. Inspired by:

https://developers.google.com/tasks/performance#partial-response

https://github.com/nemtsov/json-mask

## Installation

```
go get github.com/teambition/json-mask-go
```

## Syntax

The syntax is loosely based on XPath:
```
a       select a field 'a'
a,b,c   comma-separated list will select multiple fields
a/b/c   path will select a field from its parent
a(b,c)  sub-selection will select many fields from a parent
a/*/c   the star * wildcard will select all items in a field
a,b/c(d,e(f,g/h)),i
```

## Examples
```go
doc := `
{
  "kind": "demo",
  "items": [
  {
    "title": "First title",
    "comment": "First comment.",
    "characteristics": {
      "length": "short",
      "accuracy": "high",
      "followers": ["Jo", "Will"]
    },
    "status": "active"
  },
  {
    "title": "Second title",
    "comment": "Second comment.",
    "characteristics": {
      "length": "long",
      "accuracy": "medium",
      "followers": [ ]
    },
    "status": "pending"
  }
  ]
}
`

result, _ := jsonmask.Mask([]byte(doc), "kind,items(title,characteristics/length)")
// OR:
// selection, err := jsonmask.Compile("kind,items(title,characteristics/length)")
// result, err := selection.Mask([]byte(doc))
fmt.Println(string(result))
// Output:
// {"items":[{"characteristics":{"length":"short"},"title":"First title"},{"characteristics":{"length":"long"},"title":"Second title"}],"kind":"demo"}
```
