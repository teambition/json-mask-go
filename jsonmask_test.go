package jsonmask

import "testing"

type jsonmaskCase struct {
	doc       string
	fields    string
	shouldErr bool
	res       string
}

var jsonmaskCases = []jsonmaskCase{
	{
		doc:       "",
		fields:    "a",
		shouldErr: true,
	},
	{
		doc:       "null",
		fields:    "a",
		shouldErr: true,
	},
	{
		doc:       "0",
		fields:    "a",
		shouldErr: true,
	},
	{
		doc:       string([]byte("非utf8")[1:]),
		fields:    "a",
		shouldErr: true,
	},
	{
		doc: `
		{
			"a": "a",
			"b": "b"
		}
		`,
		fields:    "a",
		shouldErr: false,
		res:       `{"a": "a"}`,
	},
	{
		doc: `
		[{
			"a": 1,
			"b": "b"
		}, {
			"a": 2,
			"b": "b"
		}]
		`,
		fields:    "a",
		shouldErr: false,
		res:       `[{"a": 1}, {"a": 2}]`,
	},
	{
		doc: `
		{
			"nextToken": "",
			"result": [
				{
					"name": "name1",
					"data": null
				}, {
					"name": "name2",
					"data": []
				}
			]
		}
		`,
		fields:    "nextToken,result(name)",
		shouldErr: false,
		res: `
		{
			"nextToken": "",
			"result": [
				{
					"name": "name1"
				}, {
					"name": "name2"
				}
			]
		}
		`,
	},
	{
		doc: `
		{
			"nextToken": "",
			"result": [
				{
					"name": "name1",
					"data": {
						"tasks": 1,
						"events": 2
					}
				}, {
					"name": "name2",
					"data": {
						"tasks": 3,
						"events": 4
					}
				}
			]
		}
		`,
		fields:    "result(data/tasks,name),nextToken",
		shouldErr: false,
		res: `
		{
			"nextToken": "",
			"result": [
				{
					"name": "name1",
					"data": {
						"tasks": 1
					}
				}, {
					"name": "name2",
					"data": {
						"tasks": 3
					}
				}
			]
		}
		`,
	},
	{
		doc: `
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
		`,
		fields:    "kind,items(title,characteristics(length,followers))",
		shouldErr: false,
		res:       `{"items":[{"characteristics":{"length":"short","followers":["Jo", "Will"]},"title":"First title"},{"characteristics":{"length":"long","followers": []},"title":"Second title"}],"kind":"demo"}`,
	},
	{
		doc: `
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
		`,
		fields:    "*/title",
		shouldErr: true,
		res:       `{"items":[{"title":"First title"},{"title":"Second title"}]}`,
	},
	{
		doc: `
		{
			"result": {
				"name": "name",
				"title": "title"
			},
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
		`,
		fields:    "*/title",
		shouldErr: false,
		res:       `{"items":[{"title":"First title"},{"title":"Second title"}],"result": {"title": "title"}}`,
	},
}

func TestJSONMask(t *testing.T) {
	for _, c := range jsonmaskCases {
		res, err := Mask([]byte(c.doc), c.fields)
		if c.shouldErr {
			if err == nil {
				t.Errorf("Testing case[%s] failed: should error but got: %#v", c.doc, string(res))
			}
		} else if err != nil {
			t.Errorf("Testing case[%s] failed: %s", c.fields, err)
		} else if !jsonDeepEqual([]byte(c.res), []byte(res)) {
			t.Errorf("Testing case[%s] failed, expected: %#v, got: %#v", c.doc, c.res, string(res))
		}
	}
}
