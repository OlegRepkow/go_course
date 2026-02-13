package document_store

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported document field")

type Store struct {
	mu          sync.RWMutex
	collections map[string]*CollectionImpl
}

type collectionDump struct {
	Config     CollectionConfig     `json:"config"`
	Documents  map[string]*Document `json:"documents"`
	IndexNames []string             `json:"index_names"`
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
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.collections[name]; exists {
		return nil, ErrCollectionAlreadyExists
	}
	s.collections[name] = &CollectionImpl{
		documents: make(map[string]*Document),
		config:    *config,
		indexes:   make(map[string]*index),
	}
	return s.collections[name], nil
}

func (s *Store) GetCollection(name string) (*CollectionImpl, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	collection, exists := s.collections[name]
	if !exists {
		return nil, ErrCollectionNotFound
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.collections[name]; !exists {
		return ErrCollectionNotFound
	}

	delete(s.collections, name)
	return nil
}

func (s *Store) ListCollections() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.collections))
	for name := range s.collections {
		names = append(names, name)
	}
	return names
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
		collection := &CollectionImpl{
			documents: collData.Documents,
			config:    collData.Config,
			indexes:   make(map[string]*index),
		}

		for _, indexName := range collData.IndexNames {
			collection.CreateIndex(indexName)
		}

		store.collections[name] = collection
	}

	return store, nil
}

func (s *Store) Dump() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dumpData := storeDump{
		Collections: make(map[string]collectionDump),
	}

	for name, collection := range s.collections {
		collection.mu.RLock()
		indexNames := make([]string, 0, len(collection.indexes))
		for indexName := range collection.indexes {
			indexNames = append(indexNames, indexName)
		}

		dumpData.Collections[name] = collectionDump{
			Config:     collection.config,
			Documents:  collection.documents,
			IndexNames: indexNames,
		}
		collection.mu.RUnlock()
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
