package document_store

import "log"

// CollectionInterface визначає контракт для роботи з колекцією документів
type Collection interface {
	Put(doc *Document) error
	Get(key string) (*Document, error)
	Delete(key string) error
	List() []Document
}

type CollectionImpl struct {
	documents map[string]*Document
	config    CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

var _ Collection = (*CollectionImpl)(nil)

func (c *CollectionImpl) Put(doc *Document) error {
	if keyField, exists := doc.Fields[c.config.PrimaryKey]; exists && keyField.Type == DocumentFieldTypeString {
		if key, ok := keyField.Value.(string); ok {
			_, isUpdate := c.documents[key]
			c.documents[key] = doc
			if isUpdate {
				log.Printf("[COLLECTION] Document updated: key=%s", key)
			} else {
				log.Printf("[COLLECTION] Document created: key=%s", key)
			}
			return nil
		}
	}

	return ErrUnsupportedDocumentField
}

func (c *CollectionImpl) Get(key string) (*Document, error) {
	doc, exists := c.documents[key]
	if !exists {
		return nil, ErrDocumentNotFound
	}

	return doc, nil
}

func (c *CollectionImpl) Delete(key string) error {
	if _, exists := c.documents[key]; !exists {
		return ErrDocumentNotFound
	}

	delete(c.documents, key)
	log.Printf("[COLLECTION] Document deleted: key=%s", key)
	return nil
}

func (c *CollectionImpl) List() []Document {
	docs := make([]Document, 0, len(c.documents))

	for _, doc := range c.documents {
		docs = append(docs, *doc)
	}
	return docs
}
