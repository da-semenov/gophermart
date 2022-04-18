package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, status int, message interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, "can't write response", http.StatusBadRequest)
	}
	return nil
}

func getRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	return b, nil
}
