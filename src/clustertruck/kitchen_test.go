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
			return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockKitchenResponse)), nil
		},
	}

	kitchens, _ := getClusterTruckKitchenInfo(client)
	assertResult(t, 6, len(kitchens))
	kitchenId := "00000000-0000-0000-0000-000000000000"
	assertResult(t, "729 N. Pennsylvania St., Indianapolis, IN, 46204", kitchens[kitchenId].Address)
}

func TestGetClusterTruckKitchenInfoReturnError(t *testing.T) {
	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return createHttpResponseForTest(http.StatusOK, bytes.NewBufferString("invalidJson")), nil
		},
	}

	_, err := getClusterTruckKitchenInfo(client)
	expected := "There was an error deserializing the response from the ClusterTruck Kitchens API: invalid " +
		"character 'i' looking for beginning of value"
	assertResult(t, expected, err.Error())
}
