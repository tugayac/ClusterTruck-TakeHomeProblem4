package clustertruck

import (
	"math"
)

type RequestPayload struct {
	// The address given by the user
	StartingAddress string `json:"address"`
}

type ClosestClusterTruck struct {
	// Drive time to the closest ClusterTruck Kitchen
	DriveTime ResponseMeasurementValues `json:"drive_time"`
	// Drive distance to the closest ClusterTruck Kitchen, based on the drive time given above
	DriveDistance ResponseMeasurementValues `json:"drive_distance"`
	// Name of the ClusterTruck Kitchen
	LocationName string `json:"location_name"`
	// Address input by the user
	StartAddress string `json:"start_address"`
	// Address of the ClusterTruck Kitchen
	DestinationAddress string `json:"destination_address"`
}

type ResponseMeasurementValues struct {
	// Display value of the measured value
	Text string `json:"text"`
	// Internal representation of the measure value
	Value int `json:"value"`
	// Unit used for the internal representation
	Unit string `json:"value_unit"`
}

func findDriveTimeToClosestClusterTruckKitchen(httpClient HttpClient, startingAddress string) *ClosestClusterTruck {
	kitchens := getClusterTruckKitchenInfo(httpClient)
	drivingTimesToKitchen := make(map[string]*Route)
	for _, kitchen := range kitchens {
		directions := getGoogleMapsDirections(httpClient, startingAddress, kitchen.Address)
		if len(directions.Routes) > 1 {
			shortestDriveTimeRoute := findShortestDriveTime(directions.Routes)
			drivingTimesToKitchen[kitchen.ID] = shortestDriveTimeRoute
		} else {
			drivingTimesToKitchen[kitchen.ID] = &directions.Routes[0]
		}
	}

	closestKitchenId := findMinimumDriveTime(drivingTimesToKitchen)
	closestKitchenData := kitchens[closestKitchenId]
	directionsToClosestKitchen := drivingTimesToKitchen[closestKitchenId].Legs[0]

	return &ClosestClusterTruck{
		DriveTime: ResponseMeasurementValues{
			Text:  directionsToClosestKitchen.Duration.Text,
			Value: directionsToClosestKitchen.Duration.Value,
			Unit:  "seconds",
		},
		DriveDistance: ResponseMeasurementValues{
			Text:  directionsToClosestKitchen.Distance.Text,
			Value: directionsToClosestKitchen.Distance.Value,
			Unit:  "meters",
		},
		LocationName:       closestKitchenData.Name,
		StartAddress:       startingAddress,
		DestinationAddress: closestKitchenData.Address,
	}
}

func findMinimumDriveTime(drivingTimesToKitchen map[string]*Route) string {
	shortestDriveTime := math.MaxInt32
	closestKitchenId := ""
	for kitchenId, direction := range drivingTimesToKitchen {
		driveTime := direction.Legs[0].Duration.Value
		if driveTime < shortestDriveTime {
			shortestDriveTime = driveTime
			closestKitchenId = kitchenId
		}
	}

	return closestKitchenId
}

func findShortestDriveTime(routes []Route) *Route {
	shortestDriveTimeRoute := Route{
		Legs: []Leg{
			{
				Duration: MeasurementValues{
					Value: math.MaxInt32,
				},
			},
		},
	}
	for _, route := range routes {
		driveTime := route.Legs[0].Duration.Value
		if driveTime < shortestDriveTimeRoute.Legs[0].Duration.Value {
			shortestDriveTimeRoute = route
		}
	}

	return &shortestDriveTimeRoute
}
