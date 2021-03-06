package clustertruck

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"sync"
	"math"
	"fmt"
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

func getGoogleMapsDirections(httpClient HttpClient, origin string, destination string,
	kitchenId string, output chan<- *KitchenIDDirectionsPair, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()

	apiKey := os.Getenv("CT_GMAPS_API_KEY")
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
		output <- &KitchenIDDirectionsPair{
			Error: fmt.Sprintf("There was an error creating a request to get direction info: %s", err.Error()),
		}
		return
	}

	res, err := httpClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		output <- &KitchenIDDirectionsPair{
			Error: fmt.Sprintf("There was an error performing a request to the GMaps Directions API: %s",
				err.Error()),
		}
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		output <- &KitchenIDDirectionsPair{
			Error: fmt.Sprintf("There was an error reading the response from the GMaps Directions API: %s",
				err.Error()),
		}
		return
	}

	var directions GMapsDirections
	err = json.Unmarshal(body, &directions)
	if err != nil {
		output <- &KitchenIDDirectionsPair{
			Error: fmt.Sprintf("There was an error deserializing the response from the GMaps Directions API: %s",
				err.Error()),
		}
		return
	}

	if directions.Status != "OK" {
		output <- &KitchenIDDirectionsPair{
			Error: fmt.Sprintf("Status of GMaps Directions API response was %s", directions.Status),
		}
		return
	}

	output <- &KitchenIDDirectionsPair{
		ID:         kitchenId,
		Directions: &directions,
	}
}

func findShortestRouteByDriveTime(routes []Route) *Route {
	shortestDriveTimeRoute := Route{
		Legs: []Leg{
			{
				Duration: MeasurementValues{
					Value: math.MaxInt32,
				},
			},
		},
	}
	for _, route := range routes {
		driveTime := route.Legs[0].Duration.Value
		if driveTime < shortestDriveTimeRoute.Legs[0].Duration.Value {
			shortestDriveTimeRoute = route
		}
	}

	return &shortestDriveTimeRoute
}
