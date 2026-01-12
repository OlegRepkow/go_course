package document_store

import (
	"testing"
)

func TestDocumentFieldType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		fieldType DocumentFieldType
		expected string
	}{
		{"String", DocumentFieldTypeString, "string"},
		{"Number", DocumentFieldTypeNumber, "number"},
		{"Bool", DocumentFieldTypeBool, "bool"},
		{"Array", DocumentFieldTypeArray, "array"},
		{"Object", DocumentFieldTypeObject, "object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.fieldType) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(tt.fieldType))
			}
		})
	}
}

func TestDocumentField(t *testing.T) {
	field := DocumentField{
		Type:  DocumentFieldTypeString,
		Value: "test",
	}

	if field.Type != DocumentFieldTypeString {
		t.Errorf("expected type String, got %v", field.Type)
	}
	if field.Value != "test" {
		t.Errorf("expected value 'test', got %v", field.Value)
	}
}

func TestDocument_EmptyFields(t *testing.T) {
	doc := Document{
		Fields: make(map[string]DocumentField),
	}

	if doc.Fields == nil {
		t.Fatal("expected non-nil Fields map")
	}
	if len(doc.Fields) != 0 {
		t.Fatalf("expected empty Fields, got %d fields", len(doc.Fields))
	}
}

func TestDocument_WithFields(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
			"age":  {Type: DocumentFieldTypeNumber, Value: 30.0},
			"active": {Type: DocumentFieldTypeBool, Value: true},
		},
	}

	if len(doc.Fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(doc.Fields))
	}

	if doc.Fields["id"].Value != "100" {
		t.Errorf("expected id '100', got %v", doc.Fields["id"].Value)
	}
	if doc.Fields["name"].Value != "Alice" {
		t.Errorf("expected name 'Alice', got %v", doc.Fields["name"].Value)
	}
	if doc.Fields["age"].Value != 30.0 {
		t.Errorf("expected age 30.0, got %v", doc.Fields["age"].Value)
	}
	if doc.Fields["active"].Value != true {
		t.Errorf("expected active true, got %v", doc.Fields["active"].Value)
	}
}

func TestDocument_FieldTypes(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"string": {Type: DocumentFieldTypeString, Value: "text"},
			"number": {Type: DocumentFieldTypeNumber, Value: 42.0},
			"bool":   {Type: DocumentFieldTypeBool, Value: true},
			"array":  {Type: DocumentFieldTypeArray, Value: []interface{}{1, 2, 3}},
			"object": {Type: DocumentFieldTypeObject, Value: map[string]interface{}{"key": "value"}},
		},
	}

	if doc.Fields["string"].Type != DocumentFieldTypeString {
		t.Error("expected string type")
	}
	if doc.Fields["number"].Type != DocumentFieldTypeNumber {
		t.Error("expected number type")
	}
	if doc.Fields["bool"].Type != DocumentFieldTypeBool {
		t.Error("expected bool type")
	}
	if doc.Fields["array"].Type != DocumentFieldTypeArray {
		t.Error("expected array type")
	}
	if doc.Fields["object"].Type != DocumentFieldTypeObject {
		t.Error("expected object type")
	}
}
