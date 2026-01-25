package document_store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryWithDifferentParams(t *testing.T) {
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
		{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "4"},
			"name": {Type: DocumentFieldTypeString, Value: "David"},
		}},
	}

	for _, doc := range docs {
		require.NoError(t, collection.Put(doc))
	}

	require.NoError(t, collection.CreateIndex("name"))

	tests := []struct {
		name           string
		params         QueryParams
		expectedCount  int
		expectedFirst  string
		expectedLast   string
	}{
		{
			name:          "All documents ascending",
			params:        QueryParams{},
			expectedCount: 4,
			expectedFirst: "Alice",
			expectedLast:  "David",
		},
		{
			name:          "All documents descending",
			params:        QueryParams{Desc: true},
			expectedCount: 4,
			expectedFirst: "David",
			expectedLast:  "Alice",
		},
		{
			name:          "From B onwards",
			params:        QueryParams{MinValue: stringPtr("B")},
			expectedCount: 3,
			expectedFirst: "Bob",
			expectedLast:  "David",
		},
		{
			name:          "Up to Charlie",
			params:        QueryParams{MaxValue: stringPtr("Charlie")},
			expectedCount: 3,
			expectedFirst: "Alice",
			expectedLast:  "Charlie",
		},
		{
			name:          "Between Bob and Charlie",
			params:        QueryParams{MinValue: stringPtr("Bob"), MaxValue: stringPtr("Charlie")},
			expectedCount: 2,
			expectedFirst: "Bob",
			expectedLast:  "Charlie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := collection.Query("name", tt.params)
			require.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)
			
			if tt.expectedCount > 0 {
				assert.Equal(t, tt.expectedFirst, results[0].Fields["name"].Value)
				assert.Equal(t, tt.expectedLast, results[len(results)-1].Fields["name"].Value)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
