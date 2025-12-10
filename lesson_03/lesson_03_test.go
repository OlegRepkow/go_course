package documentstore

import "testing"

func TestPut(t *testing.T) {
	documents = map[string]*Document{}
	
	doc := &Document{
		Fields: map[string]DocumentField{
			"key": {Type: DocumentFieldTypeString, Value: "test1"},
			"name": {Type: DocumentFieldTypeString, Value: "Test Doc"},
		},
	}
	
	Put(doc)
	
	if len(documents) != 1 {
		t.Errorf("Expected 1 document, got %d", len(documents))
	}
	
	if documents["test1"] != doc {
		t.Error("Document not stored correctly")
	}
}

func TestPutWithoutKey(t *testing.T) {
	documents = map[string]*Document{}
	
	doc := &Document{
		Fields: map[string]DocumentField{
			"name": {Type: DocumentFieldTypeString, Value: "Test Doc"},
		},
	}
	
	Put(doc)
	
	if len(documents) != 0 {
		t.Error("Document without key should not be stored")
	}
}

func TestGet(t *testing.T) {
	documents = map[string]*Document{}
	
	doc := &Document{Fields: map[string]DocumentField{}}
	documents["test1"] = doc
	
	result, exists := Get("test1")
	if !exists || result != doc {
		t.Error("Get should return existing document")
	}
	
	result, exists = Get("nonexistent")
	if exists || result != nil {
		t.Error("Get should return nil for nonexistent document")
	}
}

func TestDelete(t *testing.T) {
	documents = map[string]*Document{}
	
	doc := &Document{Fields: map[string]DocumentField{}}
	documents["test1"] = doc
	
	if !Delete("test1") {
		t.Error("Delete should return true for existing document")
	}
	
	if len(documents) != 0 {
		t.Error("Document should be deleted")
	}
	
	if Delete("nonexistent") {
		t.Error("Delete should return false for nonexistent document")
	}
}

func TestList(t *testing.T) {
	documents = map[string]*Document{}
	
	doc1 := &Document{Fields: map[string]DocumentField{}}
	doc2 := &Document{Fields: map[string]DocumentField{}}
	documents["test1"] = doc1
	documents["test2"] = doc2
	
	result := List()
	
	if len(result) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(result))
	}
}
