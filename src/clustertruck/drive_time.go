package clustertruck

import (
	"math"
	"sync"
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

type KitchenIDDirectionsPair struct {
	ID         string
	Directions *GMapsDirections
}

func findDriveTimeToClosestClusterTruckKitchen(httpClient HttpClient, startingAddress string) *ClosestClusterTruck {
	kitchens := getClusterTruckKitchenInfo(httpClient)
	drivingTimesToKitchen := make(map[string]*Route)
	allDirections := make(chan *KitchenIDDirectionsPair, 6)

	getDirectionsConcurrently(kitchens, httpClient, startingAddress, allDirections)

	closestKitchenData, directionsToClosestKitchen :=
		findClosestKitchenAndAssociatedDrivingData(allDirections, drivingTimesToKitchen, kitchens)

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
func findClosestKitchenAndAssociatedDrivingData(allDirections chan *KitchenIDDirectionsPair,
	drivingTimesToKitchen map[string]*Route, kitchens map[string]Kitchen) (*Kitchen, *Leg) {

	for kitchenDirectionsPair := range allDirections {
		if len(kitchenDirectionsPair.Directions.Routes) > 1 {
			shortestDriveTimeRoute := findShortestDriveTimeOfAllRoutes(kitchenDirectionsPair.Directions.Routes)
			drivingTimesToKitchen[kitchenDirectionsPair.ID] = shortestDriveTimeRoute
		} else {
			drivingTimesToKitchen[kitchenDirectionsPair.ID] = &kitchenDirectionsPair.Directions.Routes[0]
		}
	}

	closestKitchenId := findClosestClusterTruckByDriveTime(drivingTimesToKitchen)
	closestKitchenData := kitchens[closestKitchenId]
	directionsToClosestKitchen := drivingTimesToKitchen[closestKitchenId].Legs[0]

	return &closestKitchenData, &directionsToClosestKitchen
}

// This function makes concurrent calls to the GMaps Directions API,
// to avoid having to wait for the previous call to the GMaps
// Directions API.
//
// Without this optimization, subsequent calls take ~750ms to complete.
// With this optimization, subsequent calls take ~150ms to complete,
// which is about a 500% improvement (i.e. 5 times more calls can be
// processed in the same amount of time).
func getDirectionsConcurrently(kitchens map[string]Kitchen, httpClient HttpClient, startingAddress string,
	allDirections chan *KitchenIDDirectionsPair) {

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(kitchens))
	for _, kitchen := range kitchens {
		go getGoogleMapsDirections(httpClient, startingAddress, kitchen.Address,
			kitchen.ID, allDirections, &waitGroup)
	}
	waitGroup.Wait()
	close(allDirections)
}

func findClosestClusterTruckByDriveTime(drivingTimesToKitchen map[string]*Route) string {
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
