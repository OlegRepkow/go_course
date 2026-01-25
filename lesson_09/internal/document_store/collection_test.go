package document_store

import (
	"testing"
)

func newTestCollection(t *testing.T) *CollectionImpl {
	t.Helper()
	return &CollectionImpl{
		documents: make(map[string]*Document),
		config:    CollectionConfig{PrimaryKey: "id"},
	}
}

func TestCollectionImpl_Put_Success(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	err := col.Put(doc)
	if err != nil {
		t.Fatalf("unexpected error putting document: %v", err)
	}

	// Verify document was stored
	got, err := col.Get("100")
	if err != nil {
		t.Fatalf("expected document to be stored, got error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil document")
	}
	if got.Fields["name"].Value != "Alice" {
		t.Fatalf("expected name 'Alice', got %v", got.Fields["name"].Value)
	}
}

func TestCollectionImpl_Put_Update(t *testing.T) {
	col := newTestCollection(t)

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	err := col.Put(doc1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Bob"},
		},
	}

	err = col.Put(doc2)
	if err != nil {
		t.Fatalf("unexpected error updating document: %v", err)
	}

	got, err := col.Get("100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Fields["name"].Value != "Bob" {
		t.Fatalf("expected updated name 'Bob', got %v", got.Fields["name"].Value)
	}
}

func TestCollectionImpl_Put_MissingPrimaryKey(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	err := col.Put(doc)
	if err == nil {
		t.Fatal("expected error when primary key is missing")
	}
	if err != ErrUnsupportedDocumentField {
		t.Fatalf("expected ErrUnsupportedDocumentField, got: %v", err)
	}
}

func TestCollectionImpl_Put_WrongPrimaryKeyType(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeNumber, Value: 100},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	err := col.Put(doc)
	if err == nil {
		t.Fatal("expected error when primary key type is not string")
	}
	if err != ErrUnsupportedDocumentField {
		t.Fatalf("expected ErrUnsupportedDocumentField, got: %v", err)
	}
}

func TestCollectionImpl_Put_NilDocument(t *testing.T) {
	col := newTestCollection(t)

	err := col.Put(nil)
	if err == nil {
		t.Fatal("expected error when document is nil")
	}
	if err != ErrUnsupportedDocumentField {
		t.Fatalf("expected ErrUnsupportedDocumentField, got: %v", err)
	}
}

func TestCollectionImpl_Put_MultipleDocuments(t *testing.T) {
	col := newTestCollection(t)

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "2"},
			"name": {Type: DocumentFieldTypeString, Value: "Bob"},
		},
	}

	if err := col.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := col.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(list))
	}
}

func TestCollectionImpl_Get_Success(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := col.Get("100")
	if err != nil {
		t.Fatalf("expected document to be found, got error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil document")
	}
	if got.Fields["id"].Value != "100" {
		t.Fatalf("expected id '100', got %v", got.Fields["id"].Value)
	}
}

func TestCollectionImpl_Get_NotFound(t *testing.T) {
	col := newTestCollection(t)

	doc, err := col.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error when document does not exist")
	}
	if err != ErrDocumentNotFound {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
	if doc != nil {
		t.Fatal("expected nil document when not found")
	}
}

func TestCollectionImpl_Get_EmptyCollection(t *testing.T) {
	col := newTestCollection(t)

	doc, err := col.Get("any")
	if err == nil {
		t.Fatal("expected error when collection is empty")
	}
	if err != ErrDocumentNotFound {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
	if doc != nil {
		t.Fatal("expected nil document")
	}
}

func TestCollectionImpl_Delete_Success(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := col.Delete("100")
	if err != nil {
		t.Fatalf("expected Delete to succeed, got: %v", err)
	}

	_, err = col.Get("100")
	if err == nil {
		t.Fatal("expected document to be removed after Delete")
	}
	if err != ErrDocumentNotFound {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
}

func TestCollectionImpl_Delete_NotFound(t *testing.T) {
	col := newTestCollection(t)

	err := col.Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error when deleting non-existent document")
	}
	if err != ErrDocumentNotFound {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
}

func TestCollectionImpl_Delete_Multiple(t *testing.T) {
	col := newTestCollection(t)

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "1"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "2"},
		},
	}

	if err := col.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := col.Delete("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// doc2 should still exist
	_, err := col.Get("2")
	if err != nil {
		t.Fatalf("expected document '2' to still exist, got error: %v", err)
	}

	// doc1 should be deleted
	_, err = col.Get("1")
	if err == nil {
		t.Fatal("expected document '1' to be deleted")
	}
}

func TestCollectionImpl_List_Empty(t *testing.T) {
	col := newTestCollection(t)

	list := col.List()
	if list == nil {
		t.Fatal("expected non-nil list")
	}
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d documents", len(list))
	}
}

func TestCollectionImpl_List_SingleDocument(t *testing.T) {
	col := newTestCollection(t)

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := col.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 document, got %d", len(list))
	}
	if list[0].Fields["id"].Value != "100" {
		t.Fatalf("expected id '100', got %v", list[0].Fields["id"].Value)
	}
}

func TestCollectionImpl_List_MultipleDocuments(t *testing.T) {
	col := newTestCollection(t)

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "1"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "2"},
		},
	}
	doc3 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "3"},
		},
	}

	if err := col.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col.Put(doc3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := col.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 documents, got %d", len(list))
	}

	// Verify all documents are present
	ids := make(map[string]bool)
	for _, doc := range list {
		if id, ok := doc.Fields["id"].Value.(string); ok {
			ids[id] = true
		}
	}
	if !ids["1"] || !ids["2"] || !ids["3"] {
		t.Fatal("expected all documents to be in the list")
	}
}

func TestCollectionImpl_List_AfterDelete(t *testing.T) {
	col := newTestCollection(t)

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "1"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "2"},
		},
	}

	if err := col.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := col.Delete("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list := col.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 document after deletion, got %d", len(list))
	}
	if list[0].Fields["id"].Value != "2" {
		t.Fatalf("expected remaining document id '2', got %v", list[0].Fields["id"].Value)
	}
}
