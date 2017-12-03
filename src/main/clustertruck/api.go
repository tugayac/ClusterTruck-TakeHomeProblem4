package clustertruck

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func SetupAPI() *http.ServeMux {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/api/drive-time", func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			var requestPayload RequestPayload
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(body, &requestPayload)
			if err != nil {
				panic(err)
			}

			httpClient := http.Client{}
			closestClusterTruckInfo :=
				findDriveTimeToClosestClusterTruckKitchen(&httpClient, requestPayload.StartingAddress)

			responseBody, err := json.Marshal(closestClusterTruckInfo)
			if err != nil {
				panic(err)
			}

			response.Write(responseBody)
		}
	})

	return httpMux
}
