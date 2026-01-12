package document_store

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

var ErrDocumentNotFound = errors.New("document not found")
var ErrCollectionAlreadyExists = errors.New("collection already exists")
var ErrCollectionNotFound = errors.New("collection not found")
var ErrUnsupportedDocumentField = errors.New("unsupported focument field")

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

func NewStoreFromDump(dump []byte) (*Store, error) {
	var dumpData storeDump

	if err := json.Unmarshal(dump, &dumpData); err != nil {
		return nil, err
	}

	store := &Store{
		collections: make(map[string]*CollectionImpl),
	}

	collectionCount := 0
	totalDocs := 0
	for name, collData := range dumpData.Collections {
		docCount := len(collData.Documents)
		store.collections[name] = &CollectionImpl{
			documents: collData.Documents,
			config:    collData.Config,
		}
		collectionCount++
		totalDocs += docCount
	}

	log.Printf("[STORE] Store restored from dump: collections=%d, totalDocuments=%d", collectionCount, totalDocs)
	return store, nil
}

func (s *Store) Dump() ([]byte, error) {
	dumpData := storeDump{
		Collections: make(map[string]collectionDump),
	}

	collectionCount := 0
	totalDocs := 0
	for name, collection := range s.collections {
		docCount := len(collection.documents)
		dumpData.Collections[name] = collectionDump{
			Config:    collection.config,
			Documents: collection.documents,
		}
		collectionCount++
		totalDocs += docCount
	}

	dump, err := json.Marshal(dumpData)
	if err != nil {
		return nil, err
	}

	log.Printf("[STORE] Store dumped: collections=%d, totalDocuments=%d, size=%d bytes", collectionCount, totalDocs, len(dump))
	return dump, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	log.Printf("[STORE] Loading store from file: filename=%s", filename)
	dump, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("[STORE] Failed to read file: filename=%s, error=%v", filename, err)
		return nil, err
	}

	store, err := NewStoreFromDump(dump)
	if err != nil {
		log.Printf("[STORE] Failed to restore from file: filename=%s, error=%v", filename, err)
		return nil, err
	}

	log.Printf("[STORE] Store successfully loaded from file: filename=%s", filename)
	return store, nil
}

func (s *Store) DumpToFile(filename string) error {
	log.Printf("[STORE] Saving store to file: filename=%s", filename)
	dump, err := s.Dump()
	if err != nil {
		log.Printf("[STORE] Failed to create dump: filename=%s, error=%v", filename, err)
		return err
	}

	if err := os.WriteFile(filename, dump, 0644); err != nil {
		log.Printf("[STORE] Failed to write file: filename=%s, error=%v", filename, err)
		return err
	}

	log.Printf("[STORE] Store successfully saved to file: filename=%s", filename)
	return nil
}
func NewStore() *Store {
	log.Printf("[STORE] New store created")
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
	log.Printf("[STORE] Collection created: name=%s, primaryKey=%s", name, config.PrimaryKey)
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
	log.Printf("[STORE] Collection deleted: name=%s", name)
	return nil
}
