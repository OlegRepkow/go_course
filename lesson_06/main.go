package main

import (
	"fmt"
	"lesson_06/document_store"
	"lesson_06/users"
)

func main() {
	usersToStore := []users.User{
		{ID: "1234567890", Name: "John Doe"},
		{ID: "1234567891", Name: "Jane Test"},
		{ID: "1234567892", Name: "Jane Test2"},
		{ID: "1234567893", Name: "Jane Test3"},
	}

	store := document_store.NewStore()
	coll, _ := store.CreateCollection("users", &document_store.CollectionConfig{
		PrimaryKey: "id",
	})

	service := users.Service{Coll: coll}

	for i := range usersToStore {
		user, err := service.CreateUser(usersToStore[i])
		if err != nil {
			fmt.Println("Failed to create user:", err)
		}
		fmt.Println("User created:", i, user.ID, user.Name)
	}

	users, _ := service.ListUsers()

	fmt.Println("User list:", users)

	user, _ := service.GetUser("1234567890")

	fmt.Println("User:", user.ID, user.Name)

	//Receives custom error because the user does not exist in the database
	err := service.DeleteUser("123456789111")
	if err != nil {
		fmt.Println("Failed to delete user:", err)
	}

	//Deletes the user successfully
	err = service.DeleteUser("1234567890")
	if err != nil {
		fmt.Println("Failed to delete user:", err)
	}

	//Receives custom error because the user does not exist in the database
	_, err = service.GetUser("1234567890")
	if err != nil {
		fmt.Println("Failed to get user:", err)
	}
	users, _ = service.ListUsers()

	// List without deleted user
	fmt.Println("User list:", users)

	// Testing new functionality: Dump, DumpToFile, NewStoreFromDump, NewStoreFromFile
	testDumpFunctionality(store)
}

func testDumpFunctionality(originalStore *document_store.Store) {
	fmt.Println("\n=== Testing Dump Functionality ===")

	// 1. Testing Dump() - creating dump
	fmt.Println("\n1. Creating dump from store...")
	dump, err := originalStore.Dump()
	if err != nil {
		fmt.Printf("Error creating dump: %v\n", err)
		return
	}
	fmt.Printf("Dump created successfully, size: %d bytes\n", len(dump))

	// 2. Testing NewStoreFromDump() - restoring from dump
	fmt.Println("\n2. Restoring store from dump...")
	restoredStore, err := document_store.NewStoreFromDump(dump)
	if err != nil {
		fmt.Printf("Error restoring from dump: %v\n", err)
		return
	}
	fmt.Println("Store restored successfully from dump")

	// Verifying that data was restored
	restoredColl, err := restoredStore.GetCollection("users")
	if err != nil {
		fmt.Printf("Error getting collection from restored store: %v\n", err)
		return
	}
	restoredDocs := restoredColl.List()
	fmt.Printf("Restored store contains %d documents in 'users' collection\n", len(restoredDocs))

	// 3. Testing DumpToFile() - saving to file
	fmt.Println("\n3. Saving store to file...")
	filename := "store_backup.json"
	err = originalStore.DumpToFile(filename)
	if err != nil {
		fmt.Printf("Error saving to file: %v\n", err)
		return
	}
	fmt.Printf("Store saved successfully to file: %s\n", filename)

	// 4. Testing NewStoreFromFile() - loading from file
	fmt.Println("\n4. Loading store from file...")
	fileStore, err := document_store.NewStoreFromFile(filename)
	if err != nil {
		fmt.Printf("Error loading from file: %v\n", err)
		return
	}
	fmt.Println("Store loaded successfully from file")

	// Verifying that data was loaded
	fileColl, err := fileStore.GetCollection("users")
	if err != nil {
		fmt.Printf("Error getting collection from file store: %v\n", err)
		return
	}
	fileDocs := fileColl.List()
	fmt.Printf("File store contains %d documents in 'users' collection\n", len(fileDocs))

	// 5. Comparing data
	fmt.Println("\n5. Comparing data...")
	originalColl, _ := originalStore.GetCollection("users")
	originalDocs := originalColl.List()

	if len(originalDocs) == len(restoredDocs) && len(restoredDocs) == len(fileDocs) {
		fmt.Printf("✓ All stores have the same number of documents: %d\n", len(originalDocs))
	} else {
		fmt.Printf("✗ Document count mismatch: original=%d, restored=%d, file=%d\n",
			len(originalDocs), len(restoredDocs), len(fileDocs))
	}

	fmt.Println("\n=== Dump Functionality Test Completed ===")
}
