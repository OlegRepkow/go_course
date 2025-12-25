package main

import (
	"fmt"
	"lesson_07/internal/document_store"
)

func main() {
	fmt.Println("Document Store Application")

	store := document_store.NewStore()
	coll, err := store.CreateCollection("users", &document_store.CollectionConfig{
		PrimaryKey: "id",
	})
	if err != nil {
		fmt.Printf("Error creating collection: %v\n", err)
		return
	}

	doc := &document_store.Document{
		Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: "1"},
			"name": {Type: document_store.DocumentFieldTypeString, Value: "John Doe"},
		},
	}

	if err := coll.Put(doc); err != nil {
		fmt.Printf("Error putting document: %v\n", err)
		return
	}

	retrieved, err := coll.Get("1")
	if err != nil {
		fmt.Printf("Error getting document: %v\n", err)
		return
	}

	fmt.Printf("Retrieved document: %+v\n", retrieved)
}
