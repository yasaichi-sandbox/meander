package main

import (
	"encoding/json"
	"github.com/yasaichi-sandbox/meander"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	meander.APIKey = os.Getenv("GOOGLE_PLACES_API_KEY")

	http.HandleFunc("/journeys", cors(func(w http.ResponseWriter, r *http.Request) {
		respond(w, r, meander.Journeys)
	}))
	http.HandleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query()

		q := &meander.Query{
			Journey:      strings.Split(urlQuery.Get("journey"), "|"),
			CostRangeStr: urlQuery.Get("cost"),
		}
		q.Lat, _ = strconv.ParseFloat(urlQuery.Get("lat"), 64)
		q.Lng, _ = strconv.ParseFloat(urlQuery.Get("lng"), 64)
		q.Radius, _ = strconv.Atoi(urlQuery.Get("radius"))

		respond(w, r, q.Run())
	}))

	http.ListenAndServe(":8080", nil)
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		f(w, r)
	}
}

func respond(w http.ResponseWriter, r *http.Request, data []interface{}) error {
	publicData := make([]interface{}, len(data))

	for i, value := range data {
		publicData[i] = meander.Public(value)
	}

	return json.NewEncoder(w).Encode(publicData)
}
