package document_store

import "errors"

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported focument field")

type Store struct {
	collections map[string]*CollectionImpl
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*CollectionImpl),
	}
}

func (s *Store) CreateCollection(name string, config *CollectionConfig) (*CollectionImpl, error) {
	if _, exists := s.collections[name]; exists {
		return nil, ErrCollectionAlreadyExists
	}
	s.collections[name] = &CollectionImpl{
		documents: make(map[string]*Document),
		config:    *config,
	}
	return s.collections[name], nil
}

func (s *Store) GetCollection(name string) (*CollectionImpl, error) {
	collection, exists := s.collections[name]
	if !exists {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	if _, exists := s.collections[name]; !exists {
		return ErrCollectionNotFound
	}

	delete(s.collections, name)
	return nil
}
