package clustertruck

import (
	"testing"
	"net/http/httptest"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

func TestAPIWithInvalidAccessKey(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_multiple_routes.json")
	mockKitchenResponse := readMockFile("kitchen_response.json")

	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "kitchens") {
				return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockKitchenResponse)), nil
			} else {
				return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockGmapsResponseData)), nil
			}
		},
	}

	api := SetupAPI(client)
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("POST", "/api/drive-time",
		noopCloser{bytes.NewBufferString(`{"address": "Martinsville, IN"}`)})
	request.Header.Add("Access-Key", "12345")

	api.ServeHTTP(recorder, request)

	assertResult(t, http.StatusUnauthorized, recorder.Code)

	result, _ := ioutil.ReadAll(recorder.Result().Body)
	var response HTTPError
	json.Unmarshal(result, &response)
	assertResult(t, "Your Access could not be verified. Please check your Access Key and try again.",
		response.Message)
	assertResult(t, "12345", response.Parameters["access_key"].(string))
}

func TestAPIWithValidAccessKey(t *testing.T) {
	mockGmapsResponseData := readMockFile("directions_response_multiple_routes.json")
	mockKitchenResponse := readMockFile("kitchen_response.json")

	client := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "kitchens") {
				return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockKitchenResponse)), nil
			} else {
				return createHttpResponseForTest(http.StatusOK, bytes.NewBuffer(mockGmapsResponseData)), nil
			}
		},
	}

	api := SetupAPI(client)
	recorder := httptest.NewRecorder()

	request := httptest.NewRequest("POST", "/api/drive-time",
		noopCloser{bytes.NewBufferString(`{"address": "Martinsville, IN"}`)})
	request.Header.Add("Access-Key", "JVvlYlqTBwhs2yu8")

	api.ServeHTTP(recorder, request)

	assertResult(t, http.StatusOK, recorder.Code)
}
