package meander

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var APIKey string

type Place struct {
	// NOTE: You can access `Lat` by `place.Lat` because the field name for
	// `googleGeometry` isn't specified and `Lat` is unique among the attribute
	// names of `Place` struct
	*googleGeometry `json:"geometry"`
	Name            string         `json:"name"`
	Icon            string         `json:"icon"`
	Photos          []*googlePhoto `json:"photos"`
	Vicinity        string         `json:"vicinity"`
}

type googleResponse struct {
	Results []*Place `json:"results"`
}

type googleGeometry struct {
	*googleLocation `json:"location"`
}

type googleLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type googlePhoto struct {
	PhotoRef string `json:"photo_reference"`
	URL      string `json:"url"`
}

func (p *Place) Public() interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"icon":     p.Icon,
		"photos":   p.Photos,
		"vicinity": p.Vicinity,
		"lat":      p.Lat,
		"lng":      p.Lng,
	}
}

type Query struct {
	Lat          float64
	Lng          float64
	Journey      []string
	Radius       int
	CostRangeStr string
}

func (q *Query) find(types string) (*googleResponse, error) {
	u, _ := url.Parse("https://maps.googleapis.com/maps/api/place/nearbysearch/json")
	vals := u.Query()

	vals.Set("location", fmt.Sprintf("%g,%g", q.Lat, q.Lng))
	vals.Set("radius", fmt.Sprintf("%d", q.Radius))
	vals.Set("types", types)
	vals.Set("key", APIKey)
	if len(q.CostRangeStr) > 0 {
		r := ParseCostRange(q.CostRangeStr)
		vals.Set("minprice", fmt.Sprintf("%d", int(r.From)-1))
		vals.Set("maxprice", fmt.Sprintf("%d", int(r.To)-1))
	}

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response googleResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (q *Query) Run() []interface{} {
	rand.Seed(time.Now().UnixNano())

	var w sync.WaitGroup
	places := make([]interface{}, len(q.Journey))

	for i, r := range q.Journey {
		w.Add(1)

		go func(types string, i int) {
			defer w.Done() // NOTE: `Done` decrements the WaitGroup counter by one.

			response, err := q.find(types)
			if err != nil {
				log.Println("施設の検索に失敗しました:", err)
				return
			}
			if len(response.Results) == 0 {
				log.Println("施設が見つかりませんでした:", types)
				return
			}

			for _, result := range response.Results {
				for _, photo := range result.Photos {
					// TODO: Investigate URL values returned by the API
					photo.URL = "https://maps.googleapis.com/maps/api/place/photo?" +
						"maxwidth=1000&photoreference=" + photo.PhotoRef +
						"&key=" + APIKey
				}
			}

			// NOTE: We don't have to use `sync.Mutex` because we assign the value to
			// a different element of `places` in each goroutine (`i` is an argument of
			// each goroutine, so it is constant).
			randI := rand.Intn(len(response.Results))
			places[i] = response.Results[randI]
		}(r, i)
	}
	// NOTE: All goroutines blocked on Wait are released when the counter becomes zero.
	w.Wait()

	return places
}
