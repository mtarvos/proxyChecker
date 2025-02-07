package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Error  string `json:"error"`
	Status string `json:"status"`
}

const (
	statusOK    = "OK"
	statusError = "ERROR"
)

func OK() Response                 { return Response{Status: statusOK} }
func Error(errMsg string) Response { return Response{Status: statusError, Error: errMsg} }

func Text(w http.ResponseWriter, v interface{}, status int) {
	var res []byte

	if str, ok := v.(string); ok {
		res = []byte(str)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write(res)
}

func JSON(w http.ResponseWriter, v interface{}, status int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes()) //nolint:errcheck
}
