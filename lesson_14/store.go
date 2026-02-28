package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName             = "documentstore"
	metadataCollection = "_collections"
)

type Store struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewStore(ctx context.Context, uri string) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return &Store{client: client, db: db}, nil
}

func (s *Store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

type collectionMeta struct {
	Name       string `bson:"name"`
	PrimaryKey string `bson:"primary_key"`
}

func (s *Store) getPrimaryKey(ctx context.Context, collectionName string) (string, error) {
	coll := s.db.Collection(metadataCollection)
	var meta collectionMeta
	err := coll.FindOne(ctx, bson.M{"name": collectionName}).Decode(&meta)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("collection not found")
		}
		return "", err
	}
	return meta.PrimaryKey, nil
}

func (s *Store) CreateCollection(ctx context.Context, name, primaryKey string) error {
	coll := s.db.Collection(metadataCollection)
	_, err := coll.InsertOne(ctx, collectionMeta{Name: name, PrimaryKey: primaryKey})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("collection already exists")
		}
		return err
	}
	return nil
}

func (s *Store) ListCollections(ctx context.Context) ([]string, error) {
	coll := s.db.Collection(metadataCollection)
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var names []string
	for cur.Next(ctx) {
		var meta collectionMeta
		if err := cur.Decode(&meta); err != nil {
			return nil, err
		}
		if meta.Name != metadataCollection {
			names = append(names, meta.Name)
		}
	}
	return names, cur.Err()
}

func (s *Store) DeleteCollection(ctx context.Context, name string) error {
	metaColl := s.db.Collection(metadataCollection)
	res, err := metaColl.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("collection not found")
	}
	coll := s.db.Collection(name)
	return coll.Drop(ctx)
}

func (s *Store) PutDocument(ctx context.Context, collectionName string, doc map[string]interface{}) error {
	pk, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return err
	}
	keyVal, ok := doc[pk]
	if !ok {
		return fmt.Errorf("document missing primary key field %q", pk)
	}
	id := keyVal
	doc["_id"] = id
	coll := s.db.Collection(collectionName)
	opts := options.Replace().SetUpsert(true)
	_, err = coll.ReplaceOne(ctx, bson.M{"_id": id}, doc, opts)
	return err
}

func (s *Store) GetDocument(ctx context.Context, collectionName, key string) (map[string]interface{}, error) {
	_, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	coll := s.db.Collection(collectionName)
	var result bson.M
	err = coll.FindOne(ctx, bson.M{"_id": key}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}
	return result, nil
}

func (s *Store) ListDocuments(ctx context.Context, collectionName string) ([]map[string]interface{}, error) {
	_, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	coll := s.db.Collection(collectionName)
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var docs []map[string]interface{}
	for cur.Next(ctx) {
		var d bson.M
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, cur.Err()
}

func (s *Store) DeleteDocument(ctx context.Context, collectionName, key string) error {
	_, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return err
	}
	coll := s.db.Collection(collectionName)
	res, err := coll.DeleteOne(ctx, bson.M{"_id": key})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("document not found")
	}
	return nil
}

func (s *Store) CreateIndex(ctx context.Context, collectionName, fieldName string) error {
	_, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return err
	}
	coll := s.db.Collection(collectionName)
	idx := mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: fieldName, Value: 1}},
		Options: options.Index().SetName(fieldName + "_1"),
	}
	_, err = coll.Indexes().CreateOne(ctx, idx)
	return err
}

func (s *Store) DeleteIndex(ctx context.Context, collectionName, fieldName string) error {
	_, err := s.getPrimaryKey(ctx, collectionName)
	if err != nil {
		return err
	}
	coll := s.db.Collection(collectionName)
	_, err = coll.Indexes().DropOne(ctx, fieldName+"_1")
	return err
}
