package document_store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexDumpRestore(t *testing.T) {
	// Create store with indexed collection
	store := NewStore()
	config := &CollectionConfig{PrimaryKey: "id"}
	collection, err := store.CreateCollection("test", config)
	require.NoError(t, err)

	// Add documents
	docs := []*Document{
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		}},
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "2"},
			"name": {Type: DocumentFieldTypeString, Value: "Bob"},
		}},
	}

	for _, doc := range docs {
		require.NoError(t, collection.Put(doc))
	}

	// Create index
	require.NoError(t, collection.CreateIndex("name"))

	// Dump store
	dump, err := store.Dump()
	require.NoError(t, err)

	// Restore from dump
	restoredStore, err := NewStoreFromDump(dump)
	require.NoError(t, err)

	// Get restored collection
	restoredCollection, err := restoredStore.GetCollection("test")
	require.NoError(t, err)

	// Test that index works after restore
	results, err := restoredCollection.Query("name", QueryParams{})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify documents are sorted by name (Alice, Bob)
	assert.Equal(t, "Alice", results[0].Fields["name"].Value)
	assert.Equal(t, "Bob", results[1].Fields["name"].Value)
}
