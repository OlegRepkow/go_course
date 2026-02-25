package conv

import (
	"lesson_13/internal/document_store"
	"lesson_13/internal/protocol"
)

func WireToDocument(w *protocol.DocWire) *document_store.Document {
	if w == nil {
		return nil
	}
	fields := make(map[string]document_store.DocumentField)
	for k, f := range w.Fields {
		fields[k] = document_store.DocumentField{
			Type:  document_store.DocumentFieldType(f.Type),
			Value: f.Value,
		}
	}
	return &document_store.Document{Fields: fields}
}

func DocumentToWire(d *document_store.Document) *protocol.DocWire {
	if d == nil {
		return nil
	}
	fields := make(map[string]protocol.DocFieldWire)
	for k, f := range d.Fields {
		fields[k] = protocol.DocFieldWire{
			Type:  string(f.Type),
			Value: f.Value,
		}
	}
	return &protocol.DocWire{Fields: fields}
}

func WireQueryParams(p *protocol.QueryParamsWire) document_store.QueryParams {
	if p == nil {
		return document_store.QueryParams{}
	}
	return document_store.QueryParams{
		Desc:     p.Desc,
		MinValue: p.MinValue,
		MaxValue: p.MaxValue,
	}
}
