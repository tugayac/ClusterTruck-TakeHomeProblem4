package clustertruck

import (
	"testing"
	"net/http"
	"bytes"
	"sync"
)

func TestGetGoogleMapsDirectionsWithSingleRoute(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_single_route.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockGmapsResponseData)), nil
		},
	}

	kitchenDirectionsPair := make(chan *KitchenIDDirectionsPair, 1)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	getGoogleMapsDirections(client, "origin", "destination", "kitchenId",
		kitchenDirectionsPair, &waitGroup)
	close(kitchenDirectionsPair)

	expected := "54.2 mi"
	actual := (<-kitchenDirectionsPair).Directions.Routes[0].Legs[0].Distance.Text
	assertResult(t, expected, actual)
}

func TestGetGoogleMapsDirectionsWithMultipleRoutes(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_multiple_routes.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockGmapsResponseData)), nil

		},
	}

	kitchenDirectionsPair := make(chan *KitchenIDDirectionsPair, 1)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	getGoogleMapsDirections(client, "origin", "destination", "kitchenId",
		kitchenDirectionsPair, &waitGroup)
	close(kitchenDirectionsPair)

	directions := (<-kitchenDirectionsPair).Directions
	expected := 3
	actual := len(directions.Routes)
	assertResult(t, expected, actual)

	expected = 4854
	actual = directions.Routes[1].Legs[0].Duration.Value
	assertResult(t, expected, actual)
}

func TestGetGoogleMapsDirectionsError(t *testing.T) {
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return createHttpResponseForTest(http.StatusOK, bytes.NewBufferString("invalidJson")), nil
		},
	}

	kitchenDirectionsPair := make(chan *KitchenIDDirectionsPair, 1)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	getGoogleMapsDirections(client, "origin", "destination", "kitchenId",
		kitchenDirectionsPair, &waitGroup)
	close(kitchenDirectionsPair)

	expected := "There was an error deserializing the response from the GMaps Directions API: " +
		"invalid character 'i' looking for beginning of value"
	actual := (<-kitchenDirectionsPair).Error
	assertResult(t, expected, actual)
}
