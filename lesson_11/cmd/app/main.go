package main

import (
	"fmt"
	"lesson_11/internal/document_store"

	"sync"
)

func main() {
	store := document_store.NewStore()
	config := &document_store.CollectionConfig{PrimaryKey: "id"}
	col, err := store.CreateCollection("test", config)
	if err != nil {
		panic(err)
	}
	if err := col.CreateIndex("name"); err != nil {
		panic(err)
	}

	const numGoroutines = 1000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("doc-%d", id)
			doc := &document_store.Document{
				Fields: map[string]document_store.DocumentField{
					"id":   {Type: document_store.DocumentFieldTypeString, Value: key},
					"name": {Type: document_store.DocumentFieldTypeString, Value: fmt.Sprintf("name-%d", id)},
				},
			}
			_ = col.Put(doc)
			_, _ = col.Get(key)
			minVal := "name-0"
			_, _ = col.Query("name", document_store.QueryParams{MinValue: &minVal})
			_ = col.Delete(key)
		}(i)
	}

	wg.Wait()
	fmt.Println("All 1000 goroutines finished")
}
