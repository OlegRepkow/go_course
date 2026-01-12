package document_store

import (
	"encoding/json"
	"errors"
	"os"
)

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")

type Store struct {
	collections map[string]*CollectionImpl
}

type collectionDump struct {
	Config    CollectionConfig     `json:"config"`
	Documents map[string]*Document `json:"documents"`
}

type storeDump struct {
	Collections map[string]collectionDump `json:"collections"`
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*CollectionImpl),
	}
}

func (s *Store) CreateCollection(name string, config *CollectionConfig) (*CollectionImpl, error) {
	if config == nil {
		return nil, ErrUnsupportedDocumentField
	}
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

func NewStoreFromDump(dump []byte) (*Store, error) {
	var dumpData storeDump

	if err := json.Unmarshal(dump, &dumpData); err != nil {
		return nil, err
	}

	store := &Store{
		collections: make(map[string]*CollectionImpl),
	}

	for name, collData := range dumpData.Collections {
		store.collections[name] = &CollectionImpl{
			documents: collData.Documents,
			config:    collData.Config,
		}
	}

	return store, nil
}

func (s *Store) Dump() ([]byte, error) {
	dumpData := storeDump{
		Collections: make(map[string]collectionDump),
	}

	for name, collection := range s.collections {
		dumpData.Collections[name] = collectionDump{
			Config:    collection.config,
			Documents: collection.documents,
		}
	}

	return json.Marshal(dumpData)
}

func NewStoreFromFile(filename string) (*Store, error) {
	dump, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewStoreFromDump(dump)
}

func (s *Store) DumpToFile(filename string) error {
	dump, err := s.Dump()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, dump, 0644)
}
