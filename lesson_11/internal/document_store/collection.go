package document_store

import (
	"errors"
	"sort"
	"sync"
)

var ErrIndexAlreadyExists = errors.New("index already exists")
var ErrIndexNotFound = errors.New("index not found")

type Collection interface {
	Put(doc *Document) error
	Get(key string) (*Document, error)
	Delete(key string) error
	List() []Document
	CreateIndex(fieldName string) error
	DeleteIndex(fieldName string) error
	Query(fieldName string, params QueryParams) ([]Document, error)
}

type index struct {
	data map[string][]*Document
	keys []string
}

type CollectionImpl struct {
	mu        sync.RWMutex
	documents map[string]*Document
	config    CollectionConfig
	indexes   map[string]*index
}

type CollectionConfig struct {
	PrimaryKey string
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

var _ Collection = (*CollectionImpl)(nil)

func (c *CollectionImpl) Put(doc *Document) error {
	if doc == nil {
		return ErrUnsupportedDocumentField
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if keyField, exists := doc.Fields[c.config.PrimaryKey]; exists && keyField.Type == DocumentFieldTypeString {
		if key, ok := keyField.Value.(string); ok {
			oldDoc := c.documents[key]
			c.documents[key] = doc
			
			for fieldName, idx := range c.indexes {
				if oldDoc != nil {
					c.removeFromIndex(idx, fieldName, oldDoc)
				}
				c.addToIndex(idx, fieldName, doc)
			}
			return nil
		}
	}
	return ErrUnsupportedDocumentField
}

func (c *CollectionImpl) Get(key string) (*Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, exists := c.documents[key]
	if !exists {
		return nil, ErrDocumentNotFound
	}
	return doc, nil
}

func (c *CollectionImpl) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	doc, exists := c.documents[key]
	if !exists {
		return ErrDocumentNotFound
	}

	for fieldName, idx := range c.indexes {
		c.removeFromIndex(idx, fieldName, doc)
	}
	
	delete(c.documents, key)
	return nil
}

func (c *CollectionImpl) List() []Document {
	c.mu.RLock()
	defer c.mu.RUnlock()
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docs = append(docs, *doc)
	}
	return docs
}

func (c *CollectionImpl) CreateIndex(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.indexes[fieldName]; exists {
		return ErrIndexAlreadyExists
	}

	idx := &index{
		data: make(map[string][]*Document),
		keys: make([]string, 0),
	}
	for _, doc := range c.documents {
		c.addToIndex(idx, fieldName, doc)
	}
	
	sort.Strings(idx.keys)
	
	c.indexes[fieldName] = idx
	return nil
}

func (c *CollectionImpl) DeleteIndex(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}
	delete(c.indexes, fieldName)
	return nil
}

func (c *CollectionImpl) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	idx, exists := c.indexes[fieldName]
	if !exists {
		return nil, ErrIndexNotFound
	}

	startIdx := 0
	endIdx := len(idx.keys)
	
	if params.MinValue != nil {
		startIdx = sort.Search(len(idx.keys), func(i int) bool {
			return idx.keys[i] >= *params.MinValue
		})
	}
	
	if params.MaxValue != nil {
		endIdx = sort.Search(len(idx.keys), func(i int) bool {
			return idx.keys[i] > *params.MaxValue
		})
	}
	
	keys := make([]string, endIdx-startIdx)
	copy(keys, idx.keys[startIdx:endIdx])
	
	if params.Desc {
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}
	
	var result []Document
	for _, key := range keys {
		for _, doc := range idx.data[key] {
			result = append(result, *doc)
		}
	}
	
	return result, nil
}

func (c *CollectionImpl) addToIndex(idx *index, fieldName string, doc *Document) {
	if field, exists := doc.Fields[fieldName]; exists && field.Type == DocumentFieldTypeString {
		if value, ok := field.Value.(string); ok {
			if _, exists := idx.data[value]; !exists {
				insertPos := sort.SearchStrings(idx.keys, value)
				idx.keys = append(idx.keys, "")
				copy(idx.keys[insertPos+1:], idx.keys[insertPos:])
				idx.keys[insertPos] = value
			}
			idx.data[value] = append(idx.data[value], doc)
		}
	}
}

func (c *CollectionImpl) removeFromIndex(idx *index, fieldName string, doc *Document) {
	if field, exists := doc.Fields[fieldName]; exists && field.Type == DocumentFieldTypeString {
		if value, ok := field.Value.(string); ok {
			docs := idx.data[value]
			for i, d := range docs {
				if d == doc {
					idx.data[value] = append(docs[:i], docs[i+1:]...)
					if len(idx.data[value]) == 0 {
						delete(idx.data, value)
						keyPos := sort.SearchStrings(idx.keys, value)
						if keyPos < len(idx.keys) && idx.keys[keyPos] == value {
							copy(idx.keys[keyPos:], idx.keys[keyPos+1:])
							idx.keys = idx.keys[:len(idx.keys)-1]
						}
					}
					break
				}
			}
		}
	}
}
