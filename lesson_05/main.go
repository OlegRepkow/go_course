package main

import (
	"fmt"
	"lesson_05/document_store"
	"lesson_05/users"
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

	//list without deleted user
	fmt.Println("User list:", users)
}
