package clustertruck

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Kitchens []Kitchen

// Contains ClusterTruck Kitchen Information
type Kitchen struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Address  string       `json:"-"`
}

func getClusterTruckKitchenInfo(httpClient HttpClient) map[string]Kitchen {
	req, err := http.NewRequest("GET", "https://api.staging.clustertruck.com/api/kitchens", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/vnd.api.clustertruck.com; version=2")

	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var kitchens Kitchens
	err = json.Unmarshal(body, &kitchens)
	if err != nil {
		panic(err)
	}

	kitchenMap := make(map[string]Kitchen)
	for _, kitchen := range kitchens {
		kitchenMap[kitchen.ID] = kitchen
	}

	return kitchenMap
}

func (k *Kitchens) UnmarshalJSON(b []byte) error {
	var kitchens []map[string]interface{}
	err := json.Unmarshal(b, &kitchens)
	if err != nil {
		return err
	}

	*k = make([]Kitchen, len(kitchens))
	for i, kitchen := range kitchens {
		(*k)[i].ID = kitchen["id"].(string)
		(*k)[i].Name = kitchen["name"].(string)

		unmarshalAddress(kitchen, k, i)
	}

	return nil
}

// Condense address into one variable, for easy searching with GMaps Directions API
//
// We perform a check to see if the value exists (ok == false if type is not string)
// and also check the length of the variable to make sure it's not empty
func unmarshalAddress(kitchen map[string]interface{}, k *Kitchens, i int) {
	fullAddress := ""
	address1, ok := kitchen["address_1"].(string)
	if ok && len(address1) > 0 {
		fullAddress += address1
	}
	address2, ok := kitchen["address_2"].(string)
	if ok && len(address2) > 0 {
		fullAddress += " " + address2
	}
	city, ok := kitchen["city"].(string)
	if ok && len(city) > 0 {
		fullAddress += ", " + city
	}
	state, ok := kitchen["state"].(string)
	if ok && len(state) > 0 {
		fullAddress += ", " + state
	}
	zipCode, ok := kitchen["zip_code"].(string)
	if ok && len(zipCode) > 0 {
		fullAddress += ", " + zipCode
	}

	(*k)[i].Address = fullAddress
}
