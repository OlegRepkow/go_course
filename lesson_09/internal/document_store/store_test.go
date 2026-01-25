package document_store

import (
	"os"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Fatal("NewStore() returned nil")
	}
	if store.collections == nil {
		t.Fatal("NewStore() collections map is nil")
	}
	if len(store.collections) != 0 {
		t.Fatalf("NewStore() should create empty store, got %d collections", len(store.collections))
	}
}

func TestStore_CreateCollection_Success(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	col, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("expected collection to be created, got error: %v", err)
	}
	if col == nil {
		t.Fatal("expected non-nil collection pointer")
	}
	if col.config.PrimaryKey != "id" {
		t.Fatalf("expected collection config PrimaryKey to be 'id', got %q", col.config.PrimaryKey)
	}
	if len(col.documents) != 0 {
		t.Fatalf("expected empty collection, got %d documents", len(col.documents))
	}
}

func TestStore_CreateCollection_DuplicateName(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	_, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("unexpected error creating first collection: %v", err)
	}

	colAgain, err := store.CreateCollection("users", cfg)
	if err == nil {
		t.Fatal("expected error on second CreateCollection with same name")
	}
	if err != ErrCollectionAlreadyExists {
		t.Fatalf("expected ErrCollectionAlreadyExists, got: %v", err)
	}
	if colAgain != nil {
		t.Fatal("expected returned collection to be nil when already exists")
	}
}

func TestStore_CreateCollection_NilConfig(t *testing.T) {
	store := NewStore()

	col, err := store.CreateCollection("users", nil)
	if err == nil {
		t.Fatal("expected error when config is nil")
	}
	if col != nil {
		t.Fatal("expected nil collection when config is nil")
	}
}

func TestStore_CreateCollection_MultipleCollections(t *testing.T) {
	store := NewStore()

	col1, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error creating first collection: %v", err)
	}

	col2, err := store.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})
	if err != nil {
		t.Fatalf("unexpected error creating second collection: %v", err)
	}

	if col1 == col2 {
		t.Fatal("expected different collection instances")
	}
	if col1.config.PrimaryKey != "id" {
		t.Fatalf("expected first collection PrimaryKey 'id', got %q", col1.config.PrimaryKey)
	}
	if col2.config.PrimaryKey != "sku" {
		t.Fatalf("expected second collection PrimaryKey 'sku', got %q", col2.config.PrimaryKey)
	}
}

func TestStore_GetCollection_Success(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	createdCol, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("unexpected error creating collection: %v", err)
	}

	col, err := store.GetCollection("users")
	if err != nil {
		t.Fatalf("expected collection 'users' to exist, got error: %v", err)
	}
	if col == nil {
		t.Fatal("expected non-nil collection from GetCollection")
	}
	if col != createdCol {
		t.Fatal("expected GetCollection to return the same collection instance")
	}
}

func TestStore_GetCollection_NotFound(t *testing.T) {
	store := NewStore()

	col, err := store.GetCollection("nonexistent")
	if err == nil {
		t.Fatal("expected error when collection does not exist")
	}
	if err != ErrCollectionNotFound {
		t.Fatalf("expected ErrCollectionNotFound, got: %v", err)
	}
	if col != nil {
		t.Fatal("expected nil collection when not found")
	}
}

func TestStore_DeleteCollection_Success(t *testing.T) {
	store := NewStore()
	cfg := &CollectionConfig{PrimaryKey: "id"}

	_, err := store.CreateCollection("users", cfg)
	if err != nil {
		t.Fatalf("unexpected error creating collection: %v", err)
	}

	err = store.DeleteCollection("users")
	if err != nil {
		t.Fatalf("expected DeleteCollection to succeed, got: %v", err)
	}

	_, err = store.GetCollection("users")
	if err == nil {
		t.Fatal("expected collection 'users' to be removed")
	}
	if err != ErrCollectionNotFound {
		t.Fatalf("expected ErrCollectionNotFound after deletion, got: %v", err)
	}
}

func TestStore_DeleteCollection_NotFound(t *testing.T) {
	store := NewStore()

	err := store.DeleteCollection("nonexistent")
	if err == nil {
		t.Fatal("expected error when deleting non-existent collection")
	}
	if err != ErrCollectionNotFound {
		t.Fatalf("expected ErrCollectionNotFound, got: %v", err)
	}
}

func TestStore_DeleteCollection_Multiple(t *testing.T) {
	store := NewStore()

	_, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = store.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.DeleteCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// products should still exist
	_, err = store.GetCollection("products")
	if err != nil {
		t.Fatalf("expected 'products' collection to still exist, got error: %v", err)
	}

	// users should be deleted
	_, err = store.GetCollection("users")
	if err == nil {
		t.Fatal("expected 'users' collection to be deleted")
	}
}

func TestStore_Dump_EmptyStore(t *testing.T) {
	store := NewStore()

	dump, err := store.Dump()
	if err != nil {
		t.Fatalf("unexpected error creating dump: %v", err)
	}
	if dump == nil {
		t.Fatal("expected non-nil dump")
	}
	if len(dump) == 0 {
		t.Fatal("expected non-empty dump")
	}
}

func TestStore_Dump_WithCollections(t *testing.T) {
	store := NewStore()

	col, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dump, err := store.Dump()
	if err != nil {
		t.Fatalf("unexpected error creating dump: %v", err)
	}
	if dump == nil {
		t.Fatal("expected non-nil dump")
	}
}

func TestStore_Dump_MultipleCollections(t *testing.T) {
	store := NewStore()

	col1, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	col2, err := store.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "1"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"sku": {Type: DocumentFieldTypeString, Value: "P001"},
		},
	}

	if err := col1.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col2.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dump, err := store.Dump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Restore and verify
	restored, err := NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("unexpected error restoring: %v", err)
	}

	restoredCol1, err := restored.GetCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	restoredCol2, err := restored.GetCollection("products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(restoredCol1.documents) != 1 {
		t.Fatalf("expected 1 document in users, got %d", len(restoredCol1.documents))
	}
	if len(restoredCol2.documents) != 1 {
		t.Fatalf("expected 1 document in products, got %d", len(restoredCol2.documents))
	}
}

func TestNewStoreFromDump_Success(t *testing.T) {
	originalStore := NewStore()

	col, err := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "100"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dump, err := originalStore.Dump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restoredStore, err := NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("unexpected error restoring: %v", err)
	}

	if restoredStore == nil {
		t.Fatal("expected non-nil restored store")
	}

	restoredCol, err := restoredStore.GetCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restoredDoc, err := restoredCol.Get("100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if restoredDoc.Fields["name"].Value != "Alice" {
		t.Fatalf("expected name 'Alice', got %v", restoredDoc.Fields["name"].Value)
	}
}

func TestNewStoreFromDump_InvalidJSON(t *testing.T) {
	invalidDump := []byte("{ invalid json }")

	store, err := NewStoreFromDump(invalidDump)
	if err == nil {
		t.Fatal("expected error with invalid JSON")
	}
	if store != nil {
		t.Fatal("expected nil store with invalid JSON")
	}
}

func TestNewStoreFromDump_EmptyJSON(t *testing.T) {
	emptyDump := []byte("{}")

	store, err := NewStoreFromDump(emptyDump)
	if err != nil {
		t.Fatalf("unexpected error with empty JSON: %v", err)
	}
	if store == nil {
		t.Fatal("expected non-nil store")
	}
	if len(store.collections) != 0 {
		t.Fatalf("expected empty store, got %d collections", len(store.collections))
	}
}

func TestStore_DumpToFile_Success(t *testing.T) {
	store := NewStore()

	col, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Test"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filename := "test_dump.json"
	defer os.Remove(filename)

	err = store.DumpToFile(filename)
	if err != nil {
		t.Fatalf("unexpected error saving to file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}
}

func TestStore_DumpToFile_VerifyContent(t *testing.T) {
	store := NewStore()

	col, err := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filename := "test_dump_verify.json"
	defer os.Remove(filename)

	err = store.DumpToFile(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Restore from file and verify
	restoredStore, err := NewStoreFromFile(filename)
	if err != nil {
		t.Fatalf("unexpected error loading from file: %v", err)
	}

	restoredCol, err := restoredStore.GetCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	restoredDoc, err := restoredCol.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if restoredDoc.Fields["name"].Value != "Alice" {
		t.Fatalf("expected name 'Alice', got %v", restoredDoc.Fields["name"].Value)
	}
}

func TestNewStoreFromFile_Success(t *testing.T) {
	originalStore := NewStore()

	col, err := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Bob"},
		},
	}

	if err := col.Put(doc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	filename := "test_load.json"
	defer os.Remove(filename)

	if err := originalStore.DumpToFile(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadedStore, err := NewStoreFromFile(filename)
	if err != nil {
		t.Fatalf("unexpected error loading from file: %v", err)
	}

	if loadedStore == nil {
		t.Fatal("expected non-nil loaded store")
	}

	loadedCol, err := loadedStore.GetCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadedDoc, err := loadedCol.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loadedDoc.Fields["name"].Value != "Bob" {
		t.Fatalf("expected name 'Bob', got %v", loadedDoc.Fields["name"].Value)
	}
}

func TestNewStoreFromFile_NotFound(t *testing.T) {
	store, err := NewStoreFromFile("nonexistent_file.json")
	if err == nil {
		t.Fatal("expected error when file does not exist")
	}
	if store != nil {
		t.Fatal("expected nil store when file does not exist")
	}
}

func TestStore_DumpAndRestore_RoundTrip(t *testing.T) {
	originalStore := NewStore()

	// Create multiple collections with documents
	col1, err := originalStore.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	col2, err := originalStore.CreateCollection("products", &CollectionConfig{PrimaryKey: "sku"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doc1 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "User1"},
		},
	}
	doc2 := &Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "2"},
			"name": {Type: DocumentFieldTypeString, Value: "User2"},
		},
	}
	doc3 := &Document{
		Fields: map[string]DocumentField{
			"sku":  {Type: DocumentFieldTypeString, Value: "P1"},
			"name": {Type: DocumentFieldTypeString, Value: "Product1"},
		},
	}

	if err := col1.Put(doc1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col1.Put(doc2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := col2.Put(doc3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Dump to memory
	dump, err := originalStore.Dump()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Restore from memory
	restoredStore, err := NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all data
	restoredCol1, err := restoredStore.GetCollection("users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	restoredCol2, err := restoredStore.GetCollection("products")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(restoredCol1.documents) != 2 {
		t.Fatalf("expected 2 documents in users, got %d", len(restoredCol1.documents))
	}
	if len(restoredCol2.documents) != 1 {
		t.Fatalf("expected 1 document in products, got %d", len(restoredCol2.documents))
	}

	// Verify specific documents
	user1, err := restoredCol1.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user1.Fields["name"].Value != "User1" {
		t.Fatalf("expected name 'User1', got %v", user1.Fields["name"].Value)
	}

	user2, err := restoredCol1.Get("2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user2.Fields["name"].Value != "User2" {
		t.Fatalf("expected name 'User2', got %v", user2.Fields["name"].Value)
	}

	product1, err := restoredCol2.Get("P1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product1.Fields["name"].Value != "Product1" {
		t.Fatalf("expected name 'Product1', got %v", product1.Fields["name"].Value)
	}
}
