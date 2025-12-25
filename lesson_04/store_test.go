package documentstore

import (
	"fmt"
	"testing"
)

func TestStore_CreateCollection(t *testing.T) {
	store := NewStore()

	cfg := CollectionConfig{PrimaryKey: "id"}

	created, col := store.CreateCollection("users", &cfg)
	if !created {
		t.Fatalf("expected collection to be created on first call")
	}
	if col == nil {
		t.Fatalf("expected non-nil collection pointer")
	}
	if col.config.PrimaryKey != "id" {
		t.Fatalf("expected collection config PrimaryKey to be 'id', got %q", col.config.PrimaryKey)
	}

	createdAgain, colAgain := store.CreateCollection("users", &cfg)
	if createdAgain {
		t.Fatalf("expected second CreateCollection with same name to return false")
	}
	if colAgain != nil {
		t.Fatalf("expected returned collection to be nil when already exists")
	}
}

func TestStore_GetAndDeleteCollection(t *testing.T) {
	store := NewStore()
	cfg := CollectionConfig{PrimaryKey: "id"}

	store.CreateCollection("users", &cfg)

	col, ok := store.GetCollection("users")
	if !ok {
		t.Fatalf("expected collection 'users' to exist")
	}
	if col == nil {
		t.Fatalf("expected non-nil collection from GetCollection")
	}

	if deleted := store.DeleteCollection("users"); !deleted {
		t.Fatalf("expected DeleteCollection to return true for existing collection")
	}

	if _, ok := store.GetCollection("users"); ok {
		t.Fatalf("expected collection 'users' to be removed")
	}

	if deletedAgain := store.DeleteCollection("users"); deletedAgain {
		t.Fatalf("expected DeleteCollection to return false for non-existing collection")
	}
}

func TestStore_Collection_DocumentFlow(t *testing.T) {
	store := NewStore()
	cfg := CollectionConfig{PrimaryKey: "id"}

	created, col := store.CreateCollection("users", &cfg)
	if !created || col == nil {
		t.Fatalf("failed to create collection 'users'")
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	col.Put(doc)

	got, ok := col.Get("100")
	fmt.Println("got:", got)
	if !ok {
		t.Fatalf("expected document with key '100' to be stored in collection")
	}
	if got == nil || got.Fields["name"].Value != "Alice" {
		t.Fatalf("unexpected document data in collection")
	}

	list := col.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 document in collection list, got %d", len(list))
	}

	if deleted := col.Delete("100"); !deleted {
		t.Fatalf("expected Delete to return true for existing document")
	}
	if _, ok := col.Get("100"); ok {
		t.Fatalf("expected document to be removed after Delete")
	}
}
