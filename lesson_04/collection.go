package documentstore

type Collection struct {
	documents map[string]*Document
	config    CollectionConfig
}

type CollectionConfig struct {
	PrimaryKey string
}

func (c *Collection) Put(doc *Document) {
	if keyField, exists := doc.Fields[c.config.PrimaryKey]; exists && keyField.Type == DocumentFieldTypeString {
		if key, ok := keyField.Value.(string); ok {
			c.documents[key] = doc
		}
	}
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, exists := c.documents[key]
	return doc, exists
}

func (c *Collection) Delete(key string) bool {
	if _, exists := c.documents[key]; exists {
		delete(c.documents, key)
		return true
	}
	return false
}

func (c *Collection) List() []Document {
	docs := make([]Document, 0, len(c.documents))

	for _, doc := range c.documents {
		docs = append(docs, *doc)
	}
	return docs
}
