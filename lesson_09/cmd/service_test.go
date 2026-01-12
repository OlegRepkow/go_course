package main

import (
	"testing"

	"lesson_07/internal/document_store"
	"lesson_07/internal/document_store/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Service that uses Collection interface
type DocumentService struct {
	collection document_store.Collection
}

func (s *DocumentService) AddUserWithIndex(name, id string) error {
	doc := &document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: id},
			"name": {Type: document_store.DocumentFieldTypeString, Value: name},
		},
	}
	
	if err := s.collection.Put(doc); err != nil {
		return err
	}
	
	return s.collection.CreateIndex("name")
}

func (s *DocumentService) SearchUsers(minName string) ([]document_store.Document, error) {
	return s.collection.Query("name", document_store.QueryParams{MinValue: &minName})
}

func TestDocumentServiceWithMock(t *testing.T) {
	mockCollection := &mocks.MockCollection{}
	service := &DocumentService{collection: mockCollection}

	// Test AddUserWithIndex
	expectedDoc := &document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: "1"},
			"name": {Type: document_store.DocumentFieldTypeString, Value: "Alice"},
		},
	}

	mockCollection.On("Put", expectedDoc).Return(nil)
	mockCollection.On("CreateIndex", "name").Return(nil)

	err := service.AddUserWithIndex("Alice", "1")
	assert.NoError(t, err)

	// Test SearchUsers
	expectedResults := []document_store.Document{*expectedDoc}
	mockCollection.On("Query", "name", mock.MatchedBy(func(params document_store.QueryParams) bool {
		return params.MinValue != nil && *params.MinValue == "A"
	})).Return(expectedResults, nil)

	results, err := service.SearchUsers("A")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Alice", results[0].Fields["name"].Value)

	mockCollection.AssertExpectations(t)
}
