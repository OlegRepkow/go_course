package mocks

import (
	document_store "lesson_11/internal/document_store"

	"github.com/stretchr/testify/mock"
)

type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) Put(doc *document_store.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockCollection) Get(key string) (*document_store.Document, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document_store.Document), args.Error(1)
}

func (m *MockCollection) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCollection) List() []document_store.Document {
	args := m.Called()
	return args.Get(0).([]document_store.Document)
}

func (m *MockCollection) CreateIndex(fieldName string) error {
	args := m.Called(fieldName)
	return args.Error(0)
}

func (m *MockCollection) DeleteIndex(fieldName string) error {
	args := m.Called(fieldName)
	return args.Error(0)
}

func (m *MockCollection) Query(fieldName string, params document_store.QueryParams) ([]document_store.Document, error) {
	args := m.Called(fieldName, params)
	return args.Get(0).([]document_store.Document), args.Error(1)
}
