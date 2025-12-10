package documentstore

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value any
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]*Document{}

func Put(doc *Document) {
	if keyField, exists := doc.Fields["key"]; exists && keyField.Type == DocumentFieldTypeString {
		if key, ok := keyField.Value.(string); ok {
			documents[key] = doc
		}
	}
}

func Get(key string) (*Document, bool) {
	doc, exists := documents[key]
	return doc, exists
}

func Delete(key string) bool {
	if _, exists := documents[key]; exists {
		delete(documents, key)
		return true
	}
	return false
}

func List() []*Document {
	docs := make([]*Document, 0, len(documents))

	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return docs
}
