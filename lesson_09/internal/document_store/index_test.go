package document_store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexOperations(t *testing.T) {
	store := NewStore()
	config := &CollectionConfig{PrimaryKey: "id"}
	collection, err := store.CreateCollection("test", config)
	require.NoError(t, err)

	// Add test documents
	docs := []*Document{
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Alice"},
		}},
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "2"},
			"name": {Type: DocumentFieldTypeString, Value: "Bob"},
		}},
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "3"},
			"name": {Type: DocumentFieldTypeString, Value: "Charlie"},
		}},
	}

	for _, doc := range docs {
		require.NoError(t, collection.Put(doc))
	}

	// Test CreateIndex
	err = collection.CreateIndex("name")
	assert.NoError(t, err)

	// Test duplicate index creation
	err = collection.CreateIndex("name")
	assert.Equal(t, ErrIndexAlreadyExists, err)

	// Test Query
	minVal := "B"
	results, err := collection.Query("name", QueryParams{MinValue: &minVal})
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Test Query with range
	maxVal := "Bob"
	results, err = collection.Query("name", QueryParams{MinValue: &minVal, MaxValue: &maxVal})
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Bob", results[0].Fields["name"].Value)

	// Test DeleteIndex
	err = collection.DeleteIndex("name")
	assert.NoError(t, err)

	// Test query on deleted index
	_, err = collection.Query("name", QueryParams{})
	assert.Equal(t, ErrIndexNotFound, err)

	// Test delete non-existent index
	err = collection.DeleteIndex("nonexistent")
	assert.Equal(t, ErrIndexNotFound, err)
}
