package jsonmask_test

import (
	"fmt"

	jsonmask "github.com/teambition/json-mask-go"
)

func Example_mask() {
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

	fmt.Println(string(result))
	// Output:
	// {"items":[{"characteristics":{"length":"short"},"title":"First title"},{"characteristics":{"length":"long"},"title":"Second title"}],"kind":"demo"}
}
