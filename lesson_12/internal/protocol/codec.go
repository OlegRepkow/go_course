package protocol

import (
	"bufio"
	"encoding/json"
)

func ReadMessage(r *bufio.Reader) ([]byte, error) {
	return r.ReadBytes('\n')
}

func WriteMessage(w *bufio.Writer, msg []byte) error {
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.WriteByte('\n')
}

func DecodeRequest(data []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func EncodeResponse(resp *Response) ([]byte, error) {
	return json.Marshal(resp)
}

func EncodeRequest(req *Request) ([]byte, error) {
	return json.Marshal(req)
}

func DecodeResponse(data []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func ReadRequest(r *bufio.Reader) (*Request, error) {
	line, err := ReadMessage(r)
	if err != nil {
		return nil, err
	}
	return DecodeRequest(line)
}

func WriteResponse(w *bufio.Writer, resp *Response) error {
	data, err := EncodeResponse(resp)
	if err != nil {
		return err
	}
	return WriteMessage(w, data)
}
