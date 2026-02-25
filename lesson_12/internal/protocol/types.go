package protocol

type DocFieldWire struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type DocWire struct {
	Fields map[string]DocFieldWire `json:"fields"`
}

type QueryParamsWire struct {
	Desc     bool    `json:"desc"`
	MinValue *string `json:"min_value,omitempty"`
	MaxValue *string `json:"max_value,omitempty"`
}

type Request struct {
	Cmd string `json:"cmd"`

	Name   string `json:"name,omitempty"`
	Config *struct {
		PrimaryKey string `json:"primary_key"`
	} `json:"config,omitempty"`

	Collection string        `json:"collection,omitempty"`
	Key        string        `json:"key,omitempty"`
	Doc        *DocWire      `json:"doc,omitempty"`
	FieldName  string        `json:"field_name,omitempty"`
	Params     *QueryParamsWire `json:"params,omitempty"`
}

type Response struct {
	OK   bool     `json:"ok"`
	Err  string   `json:"err,omitempty"`
	Doc  *DocWire `json:"doc,omitempty"`
	Docs []DocWire `json:"docs,omitempty"`
	Names []string `json:"names,omitempty"`
}
