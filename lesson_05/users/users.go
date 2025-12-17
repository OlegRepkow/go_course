package users

import (
	"errors"
	"lesson_05/convert"
	"lesson_05/document_store"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Service struct {
	Coll document_store.Collection
}

func (s *Service) CreateUser(user User) (*User, error) {
	doc, err := convert.MarshalDocument(user)
	if err != nil {
		return nil, err
	}
	s.Coll.Put(doc)

	return &user, nil
}

func (s *Service) GetUser(id string) (*User, error) {
	doc, err := s.Coll.Get(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	var user User
	err = convert.UnmarshalDocument(doc, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) DeleteUser(id string) error {
	err := s.Coll.Delete(id)

	if err == document_store.ErrDocumentNotFound {

		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListUsers() ([]User, error) {
	docs := s.Coll.List()
	users := make([]User, 0, len(docs))

	for _, doc := range docs {
		var user User

		err := convert.UnmarshalDocument(&doc, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
