package clustertruck

import (
	"io/ioutil"
	"fmt"
	"os"
	"testing"
	"reflect"
	"bytes"
	"net/http"
)

func readMockFile(filename string) []byte {
	raw, err := ioutil.ReadFile(fmt.Sprintf("resources/test-data/%s", filename))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return raw
}

func assertResult(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Fatal(fmt.Sprintf("Expected types to match, but type of expected was %T, and type of actual was %T",
			expected, actual))
	}

	if actual != expected {
		t.Fatal(fmt.Sprintf("Expected %v, but got %v", expected, actual))
	}
}

func createHttpResponseForTest(statusCode int, body *bytes.Buffer) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       noopCloser{body},
	}
}
