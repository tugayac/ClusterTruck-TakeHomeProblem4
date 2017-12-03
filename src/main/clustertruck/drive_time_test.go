package clustertruck

import (
	"testing"
	"net/http"
	"bytes"
	"strings"
)

func TestFindDriveTimeToClosestClusterTruckKitchen(t *testing.T) {
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "Indianapolis") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_1.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else if strings.Contains(req.URL.String(), "Bloomington") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_2.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else if strings.Contains(req.URL.String(), "Columbus") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_3.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else if strings.Contains(req.URL.String(), "Kansas+City") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_4.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else if strings.Contains(req.URL.String(), "Denver") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_5.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else if strings.Contains(req.URL.String(), "Cleveland") {
				mockGmapsResponseData := readMockFile("directions_response_multiple_routes_simplified_6.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
				}, nil
			} else {
				mockKitchenResponse := readMockFile("kitchen_response.json")
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       noopCloser{bytes.NewBuffer(mockKitchenResponse)},
				}, nil
			}
		},
	}

	closestClusterTruckInfo := findDriveTimeToClosestClusterTruckKitchen(client, "startingAddress")
	assertResult(t, "21 mins", closestClusterTruckInfo.DriveTime.Text)
	assertResult(t, 2001, closestClusterTruckInfo.DriveTime.Value)
	assertResult(t, "96.2 mi", closestClusterTruckInfo.DriveDistance.Text)
	assertResult(t, 154775, closestClusterTruckInfo.DriveDistance.Value)
	assertResult(t, "342 East Long Street, Columbus, OH, 43215", closestClusterTruckInfo.DestinationAddress)
}
