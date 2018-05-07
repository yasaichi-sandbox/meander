package main

import (
	"encoding/json"
	"github.com/yasaichi-sandbox/meander"
	"net/http"
)

func main() {
	// meander.APIKey = "TODO"

	http.HandleFunc("/journeys", func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	})

	http.ListenAndServe(":8080", nil)
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))

	for i, value := range data {
		publicData[i] = meander.Public(value)
	}

	return json.NewEncoder(w).Encode(publicData)
}
