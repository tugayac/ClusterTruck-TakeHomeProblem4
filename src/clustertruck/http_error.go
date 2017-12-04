package clustertruck

import "encoding/json"

// Used to carry error responses sent to users
type HTTPError struct {
	Message string `json:"message"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

func marshalError(error *HTTPError) []byte {
	marshaledError, err := json.Marshal(error)
	if err != nil {
		panic(err)
	}

	return marshaledError
}
