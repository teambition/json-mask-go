package main

import (
	"fmt"

	jsonmask "github.com/teambition/json-mask-go"
)

func main() {
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

	result, err := jsonmask.Mask([]byte(doc), "kind,items(title,characteristics/length)")

	fmt.Println("json output: ", err, string(result))
	// json output:  {"a":"aaa","c":{"c1":12,"c2":33}}
}
