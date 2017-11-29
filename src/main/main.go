package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"net/url"
)

func main() {
	//GetClusterTruckKitchenInfo()
	GetGoogleMapsDirections("3400 S Sare Rd, Apt 1022, Bloomington, IN, 47401",
		"2630 E Tenth St, Bloomington, IN 47408")
}

func GetGoogleMapsDirections(origin string, destination string) {
	apiKey := "AIzaSyB50Zxex5E1MEA_E3F7M4BFYKdKkrFxPkE"
	requestUrl := url.URL{
		Scheme: "https",
		Path: "maps.googleapis.com/maps/api/directions/json",
	}
	parameters := url.Values{}
	parameters.Add("key", apiKey)
	parameters.Add("origin", origin)
	parameters.Add("destination", destination)
	requestUrl.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		panic(err)
	}

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	str := string(body[:])
	fmt.Println(str)
}

func GetClusterTruckKitchenInfo() {
	req, err := http.NewRequest("GET", "https://api.staging.clustertruck.com/api/kitchens", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/vnd.api.clustertruck.com; version=2")

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	str := string(body[:])
	fmt.Println(str)
}
