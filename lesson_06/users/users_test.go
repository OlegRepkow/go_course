package users

import (
	"errors"
	"lesson_06/document_store"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	store := document_store.NewStore()
	coll, err := store.CreateCollection("users", &document_store.CollectionConfig{
		PrimaryKey: "id",
	})
	if err != nil {
		t.Fatalf("failed to create collection: %v", err)
	}

	return &Service{Coll: coll}
}

func TestService_CreateUser(t *testing.T) {
	svc := newTestService(t)

	user := User{
		ID:   "1",
		Name: "Alice",
	}

	created, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if created == nil {
		t.Fatalf("expected non-nil user from CreateUser")
	}
	if created.ID != user.ID || created.Name != user.Name {
		t.Fatalf("unexpected user returned from CreateUser: %+v", created)
	}

	doc, err := svc.Coll.Get("1")
	if err != nil {
		t.Fatalf("expected document to be stored in collection, got error: %v", err)
	}
	if doc == nil {
		t.Fatalf("expected non-nil document in collection")
	}
}

func TestService_GetUser_Success(t *testing.T) {
	svc := newTestService(t)

	user := User{
		ID:   "100",
		Name: "Bob",
	}

	_, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	got, err := svc.GetUser("100")
	if err != nil {
		t.Fatalf("GetUser returned error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil user from GetUser")
	}
	if got.ID != user.ID || got.Name != user.Name {
		t.Fatalf("unexpected user from GetUser: %+v", got)
	}
}

func TestService_GetUser_NotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.GetUser("does-not-exist")
	if err == nil {
		t.Fatalf("expected error for non-existing user")
	}
	if !errors.Is(err, document_store.ErrDocumentNotFound) {
		t.Fatalf("expected ErrDocumentNotFound, got: %v", err)
	}
}

func TestService_DeleteUser_Success(t *testing.T) {
	svc := newTestService(t)

	user := User{
		ID:   "42",
		Name: "Charlie",
	}

	_, err := svc.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}

	if err := svc.DeleteUser("42"); err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}

	if _, err := svc.GetUser("42"); err == nil {
		t.Fatalf("expected error when getting deleted user")
	}
}

func TestService_DeleteUser_NotFound(t *testing.T) {
	svc := newTestService(t)

	err := svc.DeleteUser("missing")
	if err == nil {
		t.Fatalf("expected error for non-existing user")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestService_ListUsers(t *testing.T) {
	svc := newTestService(t)

	usersToCreate := []User{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}

	for _, u := range usersToCreate {
		if _, err := svc.CreateUser(u); err != nil {
			t.Fatalf("CreateUser(%q) returned error: %v", u.ID, err)
		}
	}

	list, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("ListUsers returned error: %v", err)
	}

	if len(list) != len(usersToCreate) {
		t.Fatalf("expected %d users, got %d", len(usersToCreate), len(list))
	}
}
