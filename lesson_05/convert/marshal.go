package convert

import (
	"encoding/json"
	"lesson_05/document_store"
)

func MarshalDocument(v any) (*document_store.Document, error) {

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}

	doc := &document_store.Document{
		Fields: make(map[string]document_store.DocumentField),
	}

	for key, value := range data {
		field := document_store.DocumentField{
			Value: value,
		}

		switch value.(type) {
		case string:
			field.Type = document_store.DocumentFieldTypeString
		case float64:
			field.Type = document_store.DocumentFieldTypeNumber
		case bool:
			field.Type = document_store.DocumentFieldTypeBool
		case []any:
			field.Type = document_store.DocumentFieldTypeArray
		case map[string]any:
			field.Type = document_store.DocumentFieldTypeObject
		default:
			return nil, document_store.ErrUnsupportedDocumentField
		}

		doc.Fields[key] = field
	}

	return doc, nil
}

func UnmarshalDocument(doc *document_store.Document, v any) error {
	data := make(map[string]any)
	for key, field := range doc.Fields {
		data[key] = field.Value
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, v)
}
