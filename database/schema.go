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

	LikeMapping = `{
		"mappings": {
		  "properties": {
			"id": {
			  "type": "keyword"
			},
			"userId": {
			  "type": "keyword"
			},
			"postId": {
			  "type": "keyword"
			},
			"createdAt": {
			  "type": "date"
			},
			"updatedAt": {
			  "type": "date"
			}
		  }
		}
	  }`

	CommentMapping = `{
		"mappings": {
		  "properties": {
			"id": {
			  "type": "keyword"
			},
			"userId": {
			  "type": "keyword"
			},
			"text": {
			  "type": "text"
			},
			"postId": {
			  "type": "keyword"
			},
			"createdAt": {
			  "type": "date"
			},
			"updatedAt": {
			  "type": "date"
			}
		  }
		}
	  }`

	ReplyMapping = `{
		"mappings": {
		  "properties": {
			"id": {
			  "type": "keyword"
			},
			"userId": {
			  "type": "keyword"
			},
			"text": {
			  "type": "text"
			},
			"commentId": {
			  "type": "keyword"
			},
			"createdAt": {
			  "type": "date"
			},
			"updatedAt": {
			  "type": "date"
			}
		  }
		}
	  }`

	ShareMapping = `{
		"mappings": {
		  "properties": {
			"id": {
			  "type": "keyword"
			},
			"userId": {
			  "type": "keyword"
			},
			"postId": {
			  "type": "keyword"
			},
			"text": {
			  "type": "text"
			},
			"createdAt": {
			  "type": "date"
			},
			"updatedAt": {
			  "type": "date"
			}
		  }
		}
	  }`
)
