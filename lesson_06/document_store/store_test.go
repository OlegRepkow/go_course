package document_store

import (
	"fmt"
	"testing"
)

func TestStore_CreateCollection(t *testing.T) {
	store := NewStore()

	cfg := &CollectionConfig{PrimaryKey: "id"}

	col, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("expected collection to be created on first call, got error: %v", err)
	}
	if col == nil {
		t.Fatalf("expected non-nil collection pointer")
	}
	if col.config.PrimaryKey != "id" {
		t.Fatalf("expected collection config PrimaryKey to be 'id', got %q", col.config.PrimaryKey)
	}

	colAgain, err := store.CreateCollection("users", cfg)
	if err == nil {
		t.Fatalf("expected error on second CreateCollection with same name")
	}
	if colAgain != nil {
		t.Fatalf("expected returned collection to be nil when already exists")
	}
}

func TestStore_GetAndDeleteCollection(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	_, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("unexpected error creating collection: %v", err)
	}

	col, err := store.GetCollection("users")
	if err != nil {
		t.Fatalf("expected collection 'users' to exist, got error: %v", err)
	}
	if col == nil {
		t.Fatalf("expected non-nil collection from GetCollection")
	}

	if err := store.DeleteCollection("users"); err != nil {
		t.Fatalf("expected DeleteCollection to succeed for existing collection, got: %v", err)
	}

	if _, err := store.GetCollection("users"); err == nil {
		t.Fatalf("expected collection 'users' to be removed")
	}

	if err := store.DeleteCollection("users"); err == nil {
		t.Fatalf("expected DeleteCollection to return error for non-existing collection")
	}
}

func TestStore_Collection_DocumentFlow(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	col, err := store.CreateCollection("users", cfg)
	if err != nil || col == nil {
		t.Fatalf("failed to create collection 'users': %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error putting document: %v", err)
	}

	got, err := col.Get("100")
	fmt.Println("got:", got)
	if err != nil {
		t.Fatalf("expected document with key '100' to be stored in collection, got error: %v", err)
	}
	if got == nil || got.Fields["name"].Value != "Alice" {
		t.Fatalf("unexpected document data in collection")
	}

	list := col.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 document in collection list, got %d", len(list))
	}

	if err := col.Delete("100"); err != nil {
		t.Fatalf("expected Delete to succeed for existing document, got: %v", err)
	}
	if _, err := col.Get("100"); err == nil {
		t.Fatalf("expected document to be removed after Delete")
	}
}
