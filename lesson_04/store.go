package documentstore

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, config *CollectionConfig) (bool, *Collection) {
	if _, exists := s.collections[name]; exists {
		return false, nil
	}

	if config == nil {
		return false, nil
	}

	s.collections[name] = &Collection{
		documents: make(map[string]*Document),
		config:    *config,
	}
	return true, s.collections[name]
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	collection, exists := s.collections[name]
	return collection, exists
}

func (s *Store) DeleteCollection(name string) bool {
	if _, exists := s.collections[name]; exists {
		delete(s.collections, name)
		return true
	}
	return false
}
