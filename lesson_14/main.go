package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultMongoURI = "mongodb://localhost:27017"

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = defaultMongoURI
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	store, err := NewStore(ctx, uri)
	if err != nil {
		log.Fatalf("connect to MongoDB: %v", err)
	}
	defer store.Close(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/create_collection", postOnly(handleCreateCollection(store)))
	mux.HandleFunc("/list_collections", postOnly(handleListCollections(store)))
	mux.HandleFunc("/delete_collection", postOnly(handleDeleteCollection(store)))
	mux.HandleFunc("/put_document", postOnly(handlePutDocument(store)))
	mux.HandleFunc("/get_document", postOnly(handleGetDocument(store)))
	mux.HandleFunc("/list_documents", postOnly(handleListDocuments(store)))
	mux.HandleFunc("/delete_document", postOnly(handleDeleteDocument(store)))
	mux.HandleFunc("/create_index", postOnly(handleCreateIndex(store)))
	mux.HandleFunc("/delete_index", postOnly(handleDeleteIndex(store)))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func postOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, OkResponse{OK: false})
			return
		}
		h(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleCreateCollection(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCollectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" || req.PrimaryKey == "" {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.CreateCollection(ctx, req.CollectionName, req.PrimaryKey)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}

func handleListCollections(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		names, err := store.ListCollections(ctx)
		if err != nil {
			writeJSON(w, http.StatusOK, ListCollectionsResponse{OK: false, ErrorMsg: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, ListCollectionsResponse{OK: true, Collections: names})
	}
}

func handleDeleteCollection(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeleteCollectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.DeleteCollection(ctx, req.CollectionName)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}

func handlePutDocument(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" || req.Document == nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.PutDocument(ctx, req.CollectionName, req.Document)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}

func handleGetDocument(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GetDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, GetDocumentResponse{OK: false})
			return
		}
		if req.CollectionName == "" {
			writeJSON(w, http.StatusBadRequest, GetDocumentResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		doc, err := store.GetDocument(ctx, req.CollectionName, req.Key)
		if err != nil {
			writeJSON(w, http.StatusOK, GetDocumentResponse{OK: false, ErrorMsg: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, GetDocumentResponse{OK: true, Document: doc})
	}
}

func handleListDocuments(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ListDocumentsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ListDocumentsResponse{OK: false})
			return
		}
		if req.CollectionName == "" {
			writeJSON(w, http.StatusBadRequest, ListDocumentsResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		docs, err := store.ListDocuments(ctx, req.CollectionName)
		if err != nil {
			writeJSON(w, http.StatusOK, ListDocumentsResponse{OK: false, ErrorMsg: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, ListDocumentsResponse{OK: true, Documents: docs})
	}
}

func handleDeleteDocument(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeleteDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.DeleteDocument(ctx, req.CollectionName, req.Key)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}

func handleCreateIndex(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateIndexRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" || req.FieldName == "" {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.CreateIndex(ctx, req.CollectionName, req.FieldName)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}

func handleDeleteIndex(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeleteIndexRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		if req.CollectionName == "" || req.FieldName == "" {
			writeJSON(w, http.StatusBadRequest, OkResponse{OK: false})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := store.DeleteIndex(ctx, req.CollectionName, req.FieldName)
		if err != nil {
			writeJSON(w, http.StatusOK, OkResponse{OK: false})
			return
		}
		writeJSON(w, http.StatusOK, OkResponse{OK: true})
	}
}
