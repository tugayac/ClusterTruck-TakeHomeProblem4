package clustertruck

import (
	"math"
	"sync"
)

// Represents the request sent by the user
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

	kitchenIdToRouteMap := make(map[string]*Route)
	allPossibleDirections := make(chan *KitchenIDDirectionsPair, len(kitchens))

	getDirectionsConcurrently(kitchens, httpClient, startingAddress, allPossibleDirections)

	closestKitchenData, directionsToClosestKitchen :=
		findClosestKitchenAndRoute(allPossibleDirections, kitchenIdToRouteMap, kitchens)

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
// This function makes concurrent calls to the GMaps Directions API,
// to avoid having to wait for the previous call to the GMaps
// Directions API.
//
// Without this optimization, subsequent calls take ~1000ms to complete.
// With this optimization, subsequent calls take ~250ms to complete,
// which is about a 400% improvement (i.e. 4 times more calls can be
// processed in the same amount of time).
func getDirectionsConcurrently(kitchens map[string]Kitchen, httpClient HttpClient, startingAddress string,
	allPossibleDirections chan *KitchenIDDirectionsPair) {

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(kitchens))

	for _, kitchen := range kitchens {
		go getGoogleMapsDirections(httpClient, startingAddress, kitchen.Address,
			kitchen.ID, allPossibleDirections, &waitGroup)
	}
	waitGroup.Wait()
	close(allPossibleDirections)
}

func findClosestKitchenAndRoute(allPossibleDirections chan *KitchenIDDirectionsPair,
	kitchenIdToRouteMap map[string]*Route, kitchens map[string]Kitchen) (*Kitchen, *Leg) {

	for kitchenIdDirectionsPair := range allPossibleDirections {
		if len(kitchenIdDirectionsPair.Directions.Routes) > 1 {
			shortestDriveTimeRoute := findShortestRouteByDriveTime(kitchenIdDirectionsPair.Directions.Routes)
			kitchenIdToRouteMap[kitchenIdDirectionsPair.ID] = shortestDriveTimeRoute
		} else {
			kitchenIdToRouteMap[kitchenIdDirectionsPair.ID] = &kitchenIdDirectionsPair.Directions.Routes[0]
		}
	}

	closestKitchenId := findClosestClusterTruckByDriveTime(kitchenIdToRouteMap)
	closestKitchenData := kitchens[closestKitchenId]
	directionsToClosestKitchen := kitchenIdToRouteMap[closestKitchenId].Legs[0]

	return &closestKitchenData, &directionsToClosestKitchen
}

func findClosestClusterTruckByDriveTime(kitchenIdToRouteMap map[string]*Route) string {
	shortestDriveTime := math.MaxInt32
	closestKitchenId := ""
	for kitchenId, directions := range kitchenIdToRouteMap {
		driveTime := directions.Legs[0].Duration.Value
		if driveTime < shortestDriveTime {
			shortestDriveTime = driveTime
			closestKitchenId = kitchenId
		}
	}

	return closestKitchenId
}
