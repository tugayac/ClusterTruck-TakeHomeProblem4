package clustertruck

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
)

func SetupAPI(httpClient HttpClient) *http.ServeMux {
	httpMux := http.NewServeMux()

	driveTimeEndpoint := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
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

			closestClusterTruckInfo, err :=
				findDriveTimeToClosestClusterTruckKitchen(httpClient, requestPayload.StartingAddress)
			if err != nil {
				errorWhileSearchingForDriveTime(response, err)
			}

			responseBody, err := json.Marshal(closestClusterTruckInfo)
			if err != nil {
				resultsCouldNotBeReturnedError(response, err, closestClusterTruckInfo)
				return
			}

			response.WriteHeader(http.StatusOK)
			response.Write(responseBody)
		}
	})

	httpMux.Handle("/api/drive-time", verifyAccessKeyMiddleware(driveTimeEndpoint))

	return httpMux
}

func verifyAccessKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		accessKey := os.Getenv("CT_API_ACCESS_KEY")
		userAccessKey := request.Header.Get("Access-Key")
		if userAccessKey != accessKey {
			unauthorizedError(response, request)
			return
		}

		next.ServeHTTP(response, request)
	})
}

func unauthorizedError(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusUnauthorized)
	response.Write(marshalError(&HTTPError{
		Message: "Your Access Key could not be verified. Please check your Access Key and try again.",
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

func errorWhileSearchingForDriveTime(response http.ResponseWriter, err error) {
	response.WriteHeader(http.StatusBadRequest)
	response.Write(marshalError(&HTTPError{
		Message: fmt.Sprintf("An error occurred while searching for drive time: %s",
			err.Error()),
	}))
}
