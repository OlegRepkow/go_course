package main

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name"`
	PrimaryKey    string `json:"primary_key"`
}

type PutDocumentRequest struct {
	CollectionName string                 `json:"collection_name"`
	Document       map[string]interface{} `json:"document"`
}

type GetDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	Key            string `json:"key"`
}

type ListDocumentsRequest struct {
	CollectionName string `json:"collection_name"`
}

type DeleteDocumentRequest struct {
	CollectionName string `json:"collection_name"`
	Key            string `json:"key"`
}

type DeleteCollectionRequest struct {
	CollectionName string `json:"collection_name"`
}

type CreateIndexRequest struct {
	CollectionName string `json:"collection_name"`
	FieldName      string `json:"field_name"`
}

type DeleteIndexRequest struct {
	CollectionName string `json:"collection_name"`
	FieldName      string `json:"field_name"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}

type GetDocumentResponse struct {
	OK        bool                   `json:"ok"`
	Document  map[string]interface{} `json:"document,omitempty"`
	ErrorMsg  string                 `json:"error,omitempty"`
}

type ListDocumentsResponse struct {
	OK        bool                     `json:"ok"`
	Documents []map[string]interface{} `json:"documents,omitempty"`
	ErrorMsg  string                   `json:"error,omitempty"`
}

type ListCollectionsResponse struct {
	OK          bool     `json:"ok"`
	Collections []string `json:"collections,omitempty"`
	ErrorMsg    string   `json:"error,omitempty"`
}
