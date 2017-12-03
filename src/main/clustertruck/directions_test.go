package clustertruck

import (
	"testing"
	"net/http"
	"bytes"
	"fmt"
	"os"
	"io/ioutil"
)

func TestGetGoogleMapsDirectionsWithSingleRoute(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_single_route.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
			}, nil
		},
	}

	directions := GetGoogleMapsDirections(client, "origin", "destination")

	expected := "54.2 mi"
	actual := directions.Routes[0].Legs[0].Distance.Text
	if actual != expected {
		t.Fatal(fmt.Sprintf("Expected %s, but got %s", expected, actual))
	}
}

func TestGetGoogleMapsDirectionsWithMultipleRoutes(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_multiple_routes.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
			}, nil
		},
	}

	directions := GetGoogleMapsDirections(client, "origin", "destination")

	expected := 3
	actual := len(directions.Routes)
	if actual != expected {
		t.Fatal(fmt.Sprintf("Expected %d, but got %d", expected, actual))
	}

	expected = 4854
	actual = directions.Routes[1].Legs[0].Duration.Value
	if actual != expected {
		t.Fatal(fmt.Sprintf("Expected %d, but got %d", expected, actual))
	}
}

func readMockFile(filename string) []byte {
	raw, err := ioutil.ReadFile(fmt.Sprintf("resources/mock-data/%s", filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return raw
}
