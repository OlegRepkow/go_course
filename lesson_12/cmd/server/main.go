package main

import (
	"bufio"
	"lesson_12/internal/conv"
	"lesson_12/internal/document_store"
	"lesson_12/internal/protocol"
	"log"
	"net"
	"strings"
)

func main() {
	store := document_store.NewStore()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleConn(conn, store)
	}
}

func handleConn(conn net.Conn, store *document_store.Store) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		req, err := protocol.ReadRequest(r)
		if err != nil {
			return
		}

		resp := handleRequest(store, req)
		if err := protocol.WriteResponse(w, resp); err != nil {
			return
		}
		if err := w.Flush(); err != nil {
			return
		}
	}
}

func handleRequest(store *document_store.Store, req *protocol.Request) *protocol.Response {
	switch strings.TrimSpace(req.Cmd) {
	case protocol.CmdCreateCollection:
		return handleCreateCollection(store, req)
	case protocol.CmdGetCollection:
		return handleGetCollection(store, req)
	case protocol.CmdDeleteCollection:
		return handleDeleteCollection(store, req)
	case protocol.CmdListCollections:
		return handleListCollections(store)
	case protocol.CmdPut:
		return handlePut(store, req)
	case protocol.CmdGet:
		return handleGet(store, req)
	case protocol.CmdDelete:
		return handleDelete(store, req)
	case protocol.CmdList:
		return handleList(store, req)
	case protocol.CmdCreateIndex:
		return handleCreateIndex(store, req)
	case protocol.CmdDeleteIndex:
		return handleDeleteIndex(store, req)
	case protocol.CmdQuery:
		return handleQuery(store, req)
	default:
		return &protocol.Response{OK: false, Err: "unknown command: " + req.Cmd}
	}
}

func handleCreateCollection(store *document_store.Store, req *protocol.Request) *protocol.Response {
	if req.Config == nil {
		return &protocol.Response{OK: false, Err: "config required"}
	}
	config := &document_store.CollectionConfig{PrimaryKey: req.Config.PrimaryKey}
	_, err := store.CreateCollection(req.Name, config)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleGetCollection(store *document_store.Store, req *protocol.Request) *protocol.Response {
	_, err := store.GetCollection(req.Name)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleDeleteCollection(store *document_store.Store, req *protocol.Request) *protocol.Response {
	err := store.DeleteCollection(req.Name)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleListCollections(store *document_store.Store) *protocol.Response {
	names := store.ListCollections()
	return &protocol.Response{OK: true, Names: names}
}

func handlePut(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	doc := conv.WireToDocument(req.Doc)
	if doc == nil {
		return &protocol.Response{OK: false, Err: "doc required"}
	}
	err = col.Put(doc)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleGet(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	doc, err := col.Get(req.Key)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true, Doc: conv.DocumentToWire(doc)}
}

func handleDelete(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	err = col.Delete(req.Key)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleList(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	docs := col.List()
	wires := make([]protocol.DocWire, 0, len(docs))
	for i := range docs {
		wires = append(wires, *conv.DocumentToWire(&docs[i]))
	}
	return &protocol.Response{OK: true, Docs: wires}
}

func handleCreateIndex(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	err = col.CreateIndex(req.FieldName)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleDeleteIndex(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	err = col.DeleteIndex(req.FieldName)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	return &protocol.Response{OK: true}
}

func handleQuery(store *document_store.Store, req *protocol.Request) *protocol.Response {
	col, err := store.GetCollection(req.Collection)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	params := conv.WireQueryParams(req.Params)
	docs, err := col.Query(req.FieldName, params)
	if err != nil {
		return &protocol.Response{OK: false, Err: err.Error()}
	}
	wires := make([]protocol.DocWire, 0, len(docs))
	for i := range docs {
		wires = append(wires, *conv.DocumentToWire(&docs[i]))
	}
	return &protocol.Response{OK: true, Docs: wires}
}
