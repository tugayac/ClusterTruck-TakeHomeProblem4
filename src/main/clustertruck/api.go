package clustertruck

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
)

func SetupAPI() *http.ServeMux {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/api/drive-time", func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			if accessKeyIsValid(request.Header.Get("Access-Key")) {
				var requestPayload RequestPayload
				body, err := ioutil.ReadAll(request.Body)
				if err != nil {
					requestBodyCouldNotBeReadError(response, err, request)
					return
				}

				err = json.Unmarshal(body, &requestPayload)
				if err != nil {
					requestBodyCouldNotBeDeserializedError(response, err, request)
					return
				}

				httpClient := http.Client{}
				closestClusterTruckInfo :=
					findDriveTimeToClosestClusterTruckKitchen(&httpClient, requestPayload.StartingAddress)

				responseBody, err := json.Marshal(closestClusterTruckInfo)
				if err != nil {
					resultsCouldNotBeReturnedError(response, err, closestClusterTruckInfo)
					return
				}

				response.WriteHeader(http.StatusOK)
				response.Write(responseBody)
			} else {
				unauthorizedError(response, request)
			}
		}
	})

	return httpMux
}

func unauthorizedError(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusUnauthorized)
	response.Write(marshalError(&HTTPError{
		Message: "Your Access could not be verified. Please check your Access Key and try again.",
		Parameters: map[string]interface{}{
			"access_key": request.Header.Get("Access-Key"),
		},
	}))
}

func resultsCouldNotBeReturnedError(response http.ResponseWriter, err error, closestClusterTruckInfo *ClosestClusterTruck) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(marshalError(&HTTPError{
		Message: fmt.Sprintf("There was a problem with returning you the results: %s",
			err.Error()),
		Parameters: map[string]interface{}{
			"drive_time_info": fmt.Sprintf("%+v", closestClusterTruckInfo),
		},
	}))
}

func requestBodyCouldNotBeDeserializedError(response http.ResponseWriter, err error, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
	response.Write(marshalError(&HTTPError{
		Message: fmt.Sprintf("The request body you provided could not be deserialized: %s",
			err.Error()),
		Parameters: map[string]interface{}{
			"body": request.Body,
		},
	}))
}

func requestBodyCouldNotBeReadError(response http.ResponseWriter, err error, request *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
	response.Write(marshalError(&HTTPError{
		Message: fmt.Sprintf("The request body you provided could not be read: %s",
			err.Error()),
		Parameters: map[string]interface{}{
			"body": request.Body,
		},
	}))
}

func accessKeyIsValid(userAccessKey string) bool {
	accessKey := os.Getenv("CT_API_ACCESS_KEY")

	return userAccessKey == accessKey
}
