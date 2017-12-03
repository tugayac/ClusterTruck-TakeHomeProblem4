package clustertruck

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// Contains data returned from a call to the GMaps Directions API
// The GMaps Directions API returns a lot more data, but we ignore the unused portions.
type GMapsDirections struct {
	Routes []Route `json:"routes"`
	Status string  `json:"status"`
}

type Route struct {
	// A leg represents a section of the route, between two waypoints
	// For a route with no waypoints (i.e. only start and end points), there will only be 1 leg
	Legs []Leg `json:"legs"`
}

type Leg struct {
	// Display value is in miles, but internal representation is in METERS
	Distance MeasurementValues `json:"distance"`
	// Display value is in hours and minutes. Internal representation is in SECONDS
	Duration MeasurementValues `json:"duration"`
}

// Contains data about a measurement, such as distance or time.
type MeasurementValues struct {
	// Display value of a measurement
	Text string `json:"text"`
	// Internal representation of a measurement
	Value int `json:"value"`
}

func getGoogleMapsDirections(httpClient HttpClient, origin string, destination string) *GMapsDirections {
	apiKey := "AIzaSyB50Zxex5E1MEA_E3F7M4BFYKdKkrFxPkE"
	requestUrl := url.URL{
		Scheme: "https",
		Path:   "maps.googleapis.com/maps/api/directions/json",
	}
	parameters := url.Values{}
	parameters.Add("key", apiKey)
	parameters.Add("origin", origin)
	parameters.Add("destination", destination)
	parameters.Add("alternatives", "true")
	requestUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		panic(err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var directions GMapsDirections
	err = json.Unmarshal(body, &directions)
	if err != nil {
		panic(err)
	}

	return &directions
}
