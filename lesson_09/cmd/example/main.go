package main

import (
	"fmt"
	"log"

	"lesson_07/internal/document_store"
)

func main() {
	// Create store and collection
	store := document_store.NewStore()
	config := &document_store.CollectionConfig{PrimaryKey: "id"}
	collection, err := store.CreateCollection("users", config)
	if err != nil {
		log.Fatal(err)
	}

	// Add some documents
	users := []*document_store.Document{
		{Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: "1"},
			"name": {Type: document_store.DocumentFieldTypeString, Value: "Alice"},
			"age":  {Type: document_store.DocumentFieldTypeNumber, Value: 25},
		}},
		{Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: "2"},
			"name": {Type: document_store.DocumentFieldTypeString, Value: "Bob"},
			"age":  {Type: document_store.DocumentFieldTypeNumber, Value: 30},
		}},
		{Fields: map[string]document_store.DocumentField{
			"id":   {Type: document_store.DocumentFieldTypeString, Value: "3"},
			"name": {Type: document_store.DocumentFieldTypeString, Value: "Charlie"},
			"age":  {Type: document_store.DocumentFieldTypeNumber, Value: 35},
		}},
	}

	for _, user := range users {
		if err := collection.Put(user); err != nil {
			log.Fatal(err)
		}
	}

	// Create index on name field
	fmt.Println("Creating index on 'name' field...")
	if err := collection.CreateIndex("name"); err != nil {
		log.Fatal(err)
	}

	// Query all users (sorted by name)
	fmt.Println("\nQuerying all users (sorted by name):")
	results, err := collection.Query("name", document_store.QueryParams{})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		fmt.Printf("ID: %s, Name: %s\n", doc.Fields["id"].Value, doc.Fields["name"].Value)
	}

	// Query users with names starting from "B"
	fmt.Println("\nQuerying users with names >= 'B':")
	minVal := "B"
	results, err = collection.Query("name", document_store.QueryParams{MinValue: &minVal})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		fmt.Printf("ID: %s, Name: %s\n", doc.Fields["id"].Value, doc.Fields["name"].Value)
	}

	// Query users with names in range "B" to "Bob"
	fmt.Println("\nQuerying users with names between 'B' and 'Bob':")
	maxVal := "Bob"
	results, err = collection.Query("name", document_store.QueryParams{MinValue: &minVal, MaxValue: &maxVal})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		fmt.Printf("ID: %s, Name: %s\n", doc.Fields["id"].Value, doc.Fields["name"].Value)
	}

	// Query in descending order
	fmt.Println("\nQuerying all users (descending order):")
	results, err = collection.Query("name", document_store.QueryParams{Desc: true})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		fmt.Printf("ID: %s, Name: %s\n", doc.Fields["id"].Value, doc.Fields["name"].Value)
	}

	// Test dump and restore
	fmt.Println("\nTesting dump and restore...")
	dump, err := store.Dump()
	if err != nil {
		log.Fatal(err)
	}

	restoredStore, err := document_store.NewStoreFromDump(dump)
	if err != nil {
		log.Fatal(err)
	}

	restoredCollection, err := restoredStore.GetCollection("users")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Querying restored collection:")
	results, err = restoredCollection.Query("name", document_store.QueryParams{})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range results {
		fmt.Printf("ID: %s, Name: %s\n", doc.Fields["id"].Value, doc.Fields["name"].Value)
	}

	fmt.Println("\nIndex operations completed successfully!")
}
