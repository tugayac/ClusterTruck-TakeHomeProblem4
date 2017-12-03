package clustertruck

import (
	"testing"
	"net/http"
	"bytes"
)

func TestGetClusterTruckKitchenInfoAddress(t *testing.T) {
	mockGmapsResponseData := readMockFile("kitchen_response.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
			}, nil
		},
	}

	kitchens := GetClusterTruckKitchenInfo(client)
	assertResult(t, 6, len(kitchens))
	assertResult(t, "729 N. Pennsylvania St., Indianapolis, IN, 46204", kitchens[0].Address)
}

func TestGetClusterTruckKitchenInfoHours(t *testing.T) {
	mockGmapsResponseData := readMockFile("kitchen_response.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       noopCloser{bytes.NewBuffer(mockGmapsResponseData)},
			}, nil
		},
	}

	kitchens := GetClusterTruckKitchenInfo(client)
	assertResult(t, 6, len(kitchens))
	assertResult(t, "23:00", kitchens[2].Hours.Friday.CloseTime)
}
