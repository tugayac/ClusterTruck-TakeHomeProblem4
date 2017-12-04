# ClusterTruck Take Home Problem #4
This is Arda Tugay's solution to one of ClusterTrucker's take home problems. The original problem can be found [here](https://github.com/ClusterTruck/take-home-problems/blob/master/engineering/problem-4.md).

## Usage
All instructions were tested on a Mac, but should work without a problem on Linux devices as well.

### Requirements
* Docker (latest community edition)

### Running Docker Container Locally (Optional)
1. Create an environment file (`.env` file) and add the following to it:

    ```
    CT_API_ACCESS_KEY=<access_key>
    CT_GMAPS_API_KEY=<gmaps_directions_api_key>
    ```

    You can set whatever you want for the access_key. See the "Making a Call" section below. You need to enable the [Google Maps Directions API](https://developers.google.com/maps/documentation/directions/intro) in order to get an API key.
1. Build the docker container using the `docker-build.sh` script (provided)
1. Run the docker container using the `docker-run.sh` script (provided). You may change the port from `8090` to anything you like.

Note: Use `Ctrl + P` then `Ctrl + Q` to detach from the docker container after it starts.

### Making a Call
You can download a tool called [_Postman_](https://www.getpostman.com/) and make a request using that or use `cURL`:

```bash
curl -X "POST" "http://localhost:8090/api/drive-time" \
     -H "Content-Type: application/json" \
     -H "Access-Key: <access_key>" \
     -d $'{"address": "123 Main St, Anywhere, OH"}'
```

Where `<access_key>` is the same access key mentioned above. If you are not running Docker locally, an access key should have already been provided for you, in which case you can use that.

**Note:** Users are assumed to always give well-formed addresses that include the number, street, city, and state, such as `123 Main St, Anywhere, OH` (zip code can be included as well). If they do not use this format, they may not get the best results.

## Application Requirements
* Build an HTTP endpoint that will receive a street address and return the drive time to the closest ClusterTruck kitchen.
    * It should utilize the ClusterTruck Kitchen API to get the kitchen addresses.
    * It should utilize the Google Maps Directions API to get the directions from the address to the closest ClusterTruck kitchen.

## Assumptions
* All ClusterTruck Kitchens are assumed to be "active" even if their `active` status is set to `false`, so that the values returned by the drive time endpoint have some variance.
* User's locale is `en_US` (American English, country USA).
* Users are located in the USA and expect distance values to be in _miles_.
* Users are assumed to always give well-formed addresses that include the number, street, city, and state, such as `123 Main St, Anywhere, OH` (zip code can be included as well). If they do not use this format, they may not get the best results
* Google Maps Directions API can return multiple routes to a destination. As such, "Drive time to closest ClusterTruck" implies shortest drive time, regardless of driving distance.
* It does not matter whether a user requests for the drive time to the nearest ClusterTruck kitchen inside or outside of a delivery area. They will always be given the drive time to the closest ClusterTruck kitchen.
* The hours of ClusterTruck kitchens are ignored. This means a user will always be given the drive time to the closest ClusterTruck, regardless of whether it's open or closed.

## Specifications and Design
### API
There will only be a single endpoint that will accept a `POST` request, with the request body being `Content-Type: application/json`. The JSON object in the request body will contain only one property, named `address`. For example:

```json
{
    "address": "123 Main St, Anywhere, OH"
}
```

**It is expected that the user will input a valid address that can be found by Google Maps.**

Users can make a request using the endpoint and any software that lets them make HTTP requests. If using `cURL`, here's an example:

```bash
curl -X "POST" "http://www.example.com/api/drive-time" \
     -H "Content-Type: application/json" \
     -H "Access-Key: <access_key>" \
     -d $'{"address": "123 Main St, Anywhere, OH"}'
```

The users can expect to receive a `200` response with header `Content-Type: application/json` and a body containing the following:

```json
{
    "drive_time": {
        "text": "29 mins",
        "value": 1715,
        "value_unit": "seconds"
    },
    "drive_distance": {
        "text": "20.5 mi",
        "value": 33043,
        "value_unit": "meters"
    },
    "location_name": "Bloomington",
    "start_address": "50 Bill's Blvd, Martinsville, IN",
    "destination_address": "2618 E 10th St, Bloomington, IN, 47408"
}
```

If there is a client-related error, they will receive a `400` response, with content like the following:

```json
{
    "error": {
        "message": "The provided address was not valid, please check the address and try again.",
        "parameters": {
            "address": "3400 Invalid Street, Unknown, UGR, 00000"
        }
    }
}
```

The response will be similar if there was a server-related error, but the return code will be `500` instead.

### Backend
#### ClusterTruck Kitchen Information
This information will be retrieved from `https://api.staging.clustertruck.com/api/kitchens`, using the request header `Accept: application/vnd.api.clustertruck.com; version=2`.

To avoid having to call the ClusterTruck Kitchen API too often, a cache will be used, with a TTL of 24 hours. A TTL of 24 hours is chosen because kitchens are not likely to change location, hours, etc frequently, and any new kitchens that are added will appear within 24 hours.

#### Calculating Drive Time
The Google Maps Directions API will be used to get the drive time from one address to the other. The server will need to use an API key. Examples of requests and responses can be found [here](https://developers.google.com/maps/documentation/directions/intro).

#### Security
To prevent unwanted users from making requests to this server, anyone who wants to access the endpoint above will need to use a key. This key will need to be passed in as part of the request header, with name `Access-Key`. For example, if using `cURL`:

```bash
curl -X "POST" "http://www.example.com/url" \
     -H "Content-Type: application/json" \
     -H "Access-Key: <key>" \
     -d $'{"address": "123 Main St, Anywhere, OH"}'
```

If no `Access-Key` is provided or it's invalid, the user will receive an HTTP `401` error.

## Rationale Behind Technology Used
### Go (Programming Language)
Go is a great language to setup HTTP endpoints fast, with the comprehensive http package it provides.

### Docker
Docker makes it very easy to create a container that can be deployed on any machine, without having to clutter the file system of the host system.

## Future Improvements
### Caching Requests
If we detect that some users frequently make requests to our API from the same starting address, we may want to cache the driving time information we get from the GMaps Directions API (since this part of our application takes the longest).

Same can be done for requests made to the ClusterTruck Kitchen API: Since kitchen information is unlikely to change often, we could cache the information we get from this API with a TTL of 24 hours.

### Driving Time Based on Kitchen Hours
Currently, kitchen hours are ignored when returning driving times. However, the results can be improved to route a user to the closest _open_ ClusterTruck location. Here's one way of doing that without changing the input format:

1. Extract the city and state from the address sent in the request. For example `Bloomington, IN`.
1. Use the [GMaps Geocoding API](https://developers.google.com/maps/documentation/geocoding/start) to get the coordinates of `Bloomington, IN`. In this case, that would be `(39.16648, -86.526918)` (latitude and longitude).
1. Send the coordinates received in the step above to the [Google Time Zone API](https://developers.google.com/maps/documentation/timezone/start) to get timezone information. The result would be something like the following for `(39.16648, -86.526918)`:
    ```json
    {
       "dstOffset" : 3600,
       "rawOffset" : -18000,
       "status" : "OK",
       "timeZoneId" : "America/New_York",
       "timeZoneName" : "Eastern Daylight Time"
    }
    ```
1. Convert all ClusterTruck kitchen hours to UTC.
1. Make a request to the GMaps Directions API as usual.
1. Calculate the estimated arrival time (in UTC) based on results from the GMaps Directions API (we could ask the GMaps Directions API to return the arrival time as well, but it's [only for public transportation](https://developers.google.com/maps/documentation/directions/intro)).
1. Return the driving time to the closest kitchen that will be _open_ when the user arrives there.

### Better Error Handling
Depending on who the end user will be, better error messages can be returned. For example, instead of telling the user that the "Google Maps Directions API returned a 400 error", we could tell them that "No results could be found due to a problem with external services. Please try again in a minute.".

Currently, the end user is assumed to be a developer, and as such, detailed error messages (directly from the failing parts of the application) are returned in case of an error.
