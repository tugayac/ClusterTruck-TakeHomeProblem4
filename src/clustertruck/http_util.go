package clustertruck

import (
	"net/http"
	"io"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}

	return &http.Response{}, nil
}

type noopCloser struct {
	io.Reader
}

func (noopCloser) Close() error {
	return nil
}
