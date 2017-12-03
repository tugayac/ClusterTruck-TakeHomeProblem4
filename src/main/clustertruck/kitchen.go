package clustertruck

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Kitchens []Kitchen

type Kitchen struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Address  string       `json:"-"`
	Hours    KitchenHours `json:"hours"`
	Timezone string       `json:"timezone"`
}

// List of kitchen hours for every day.
// Based on the response from the ClusterTruck Kitchen API, some assumptions have to be made for simplification:
//
// 1. Each day only has a single pair of open/close hours.
// 		Even though the response returns an array of open/close hours
//		for each day, no day has more than a single pair of open/close
// 		hours.
// 2. Hours listed as 01:00 to 01:01 are assumed to imply that the
// 		kitchen is closed on those days.
// 3. If no hours are listed for the kitchen (i.e. hours: null), the
// 		kitchen is assumed to be open 24/7.
type KitchenHours struct {
	Sunday    OpenClosePair `json:"sunday"`
	Monday    OpenClosePair `json:"monday"`
	Tuesday   OpenClosePair `json:"tuesday"`
	Wednesday OpenClosePair `json:"wednesday"`
	Thursday  OpenClosePair `json:"thursday"`
	Friday    OpenClosePair `json:"friday"`
	Saturday  OpenClosePair `json:"saturday"`
}

// The open and close times for the kitchen, given in hh:mm in local time
type OpenClosePair struct {
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

func GetClusterTruckKitchenInfo(httpClient HttpClient) Kitchens {
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

	return kitchens
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

		// Condense address into one variable, for easy searching with GMaps Directions API
		//
		// We perform a check to see if the value exists (ok == false if type is not string)
		// and also check the length of the variable to make sure it's not empty
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

		// Condense hours of a kitchen into an easier to read format
		kitchenHoursRaw, ok := kitchen["hours"].(map[string]interface{})
		if ok {
			kitchenHours := KitchenHours{
				Sunday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "sunday", true),
					CloseTime: getHours(kitchenHoursRaw, "sunday", false),
				},
				Monday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "monday", true),
					CloseTime: getHours(kitchenHoursRaw, "monday", false),
				},
				Tuesday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "tuesday", true),
					CloseTime: getHours(kitchenHoursRaw, "tuesday", false),
				},
				Wednesday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "wednesday", true),
					CloseTime: getHours(kitchenHoursRaw, "wednesday", false),
				},
				Thursday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "thursday", true),
					CloseTime: getHours(kitchenHoursRaw, "thursday", false),
				},
				Friday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "friday", true),
					CloseTime: getHours(kitchenHoursRaw, "friday", false),
				},
				Saturday: OpenClosePair{
					OpenTime:  getHours(kitchenHoursRaw, "saturday", true),
					CloseTime: getHours(kitchenHoursRaw, "saturday", false),
				},
			}
			(*k)[i].Hours = kitchenHours
		}

		timezone, ok := kitchen["timezone"].(string)
		if ok && len(timezone) > 0 {
			(*k)[i].Timezone = timezone
		}
	}

	return nil
}

func getHours(listOfHours map[string]interface{}, dayOfWeek string, openTime bool) string {
	if openTime {
		return listOfHours[dayOfWeek].([]interface{})[0].([]interface{})[0].(string)
	} else {
		return listOfHours[dayOfWeek].([]interface{})[0].([]interface{})[1].(string)
	}
}
