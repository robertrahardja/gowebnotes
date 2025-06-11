package main

import (
	"encoding/json"
	"net/http"
)

// writing JSON responses in HTTP handlers
//
// w http.ResponseWriter - the HTTP response writer to send data back to the client
// status int - the HTTP status code (like 200, 404, 500)
// data any - any Go value that can be converted to JSON (using the any type which is an alias for interface{})
func writeJSON(w http.ResponseWriter, status int, data any) error {
	// w.header().set("content-type", "application/json") - sets the response header to indicate json content
	w.Header().Set("content-type", "application/json")

	// w.WriteHeader(status) - sends the HTTP status code
	w.WriteHeader(status)

	// creates a JSON encoder that writes directly to the response writer
	// then encodes the data as JSON
	return json.NewEncoder(w).Encode(data)

	// returns an error if json encoding fails, otherwise nil
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	// return an error if the JSON contains fields that don't exist in the target struct.
	// This helps catch typos and prevents clients from sending unexpected data.
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}
