package http

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, statusCode int, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

type ErrorRespMsg struct {
	Message string `example:"error message" json:"message"`
	Code    string `example:"IR001"         json:"status,omitempty"`
} // @name Error

func WriteError(rw http.ResponseWriter, statusCode int, err error) {
	_ = WriteJSON(rw, statusCode, ErrorRespMsg{
		Message: err.Error(),
	})
}

func DecodeJSON(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return json.NewDecoder(req.Body).Decode(obj)
}

func ParseQuery(req *http.Request, obj interface{}) error {
	params := req.URL.Query()

	m := map[string]string{}
	for k, v := range params {
		m[k] = v[0]
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}
