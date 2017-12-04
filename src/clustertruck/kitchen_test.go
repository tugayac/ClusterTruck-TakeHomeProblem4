package clustertruck

import (
	"testing"
	"net/http"
	"bytes"
)

func TestGetClusterTruckKitchenInfoAddress(t *testing.T) {
	mockKitchenResponse := readMockFile("kitchen_response.json")
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       noopCloser{bytes.NewBuffer(mockKitchenResponse)},
			}, nil
		},
	}

	kitchens := getClusterTruckKitchenInfo(client)
	assertResult(t, 6, len(kitchens))
	kitchenId := "00000000-0000-0000-0000-000000000000"
	assertResult(t, "729 N. Pennsylvania St., Indianapolis, IN, 46204", kitchens[kitchenId].Address)
}
