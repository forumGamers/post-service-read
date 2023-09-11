package database

const (
	PostMapping = `{
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"userId": {"type": "keyword"},
				"text": {"type": "text"},
				"media": {
					"properties": {
						"url": {"type": "text"},
						"type": {"type": "keyword"},
						"id": {"type": "keyword"}
					}
				},
				"allowComment": {"type": "boolean"},
				"createdAt": {"type": "date"},
				"updatedAt": {"type": "date"},
				"tags": {"type": "keyword"},
				"privacy": {"type": "keyword"}
			}
		}
	}`
)
