# ClusterTruck Take Home Problem #4
This is Arda Tugay's solution to one of ClusterTrucker's take home problems. The original problem can be found [here](https://github.com/ClusterTruck/take-home-problems/blob/master/engineering/problem-4.md).

## Usage
TODO

## Application Requirements
* Build an HTTP endpoint that will receive a street address and return the drive time to the closest ClusterTruck kitchen.
    * It should utilize the ClusterTruck Kitchen API to get the kitchen addresses.
    * It should utilize the Google Maps Directions API to get the directions from the address to the closest ClusterTruck kitchen.

## Assumptions
TODO

## Specifications and Design
### API
There will only be a single endpoint that will accept a `POST` request, with the request body being `Content-Type: application/json`. The JSON object in the request body will contain only one property, named `address`. For example:

```json
{
    "address": "123 Main St, Anywhere, OH"
}
```

**It is expected that the user will input a valid address that can be found by Google.**

Users can make a request using the endpoint and any software that lets them make HTTP requests. If using `cURL`, here's an example:

```bash
curl -X "POST" "http://www.example.com/url" \
     -H "Content-Type: application/json" \
     -d $'{"address": "123 Main St, Anywhere, OH"}'
```

The users can expect to receive a `200` response with header `Content-Type: application/json` and a body containing something like the following:

```json
{
    "drive_time": {
        "value": 74829 // This will be in seconds
        "text": "20 hours 47 mins"
    },
    "errors": []
}
```

If there is a client-related error, they will receive a `400` response, with content like the following:

```json
{
    "drive_time": null
    "errors": [
        {
            "code": 1000,
            "message": "The provided address was not valid, please check the address and try again."
            "address": "3400 Invalid Street, Unknown, UGR, 00000"
        }
    ]
}
```

The response will be similar if there was a server-related error, but the return code will be `500` instead.

### Backend

#### ClusterTruck Kitchen Information
This information will be retrieved from `https://api.staging.clustertruck.com/api/kitchens`, using the request header `Accept: application/vnd.api.clustertruck.com; version=2`.

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

* If no `Access-Key` is provided, the user will receive a `401` error.
* If an invalid `Access-Key` is provided, the user will receive a `403` error.

## Rationale Behind Technology Used
### Go (Programming Language)
Go is a great language to setup HTTP endpoints with the comprehensive http package it provides.

### Docker
Docker makes it very easy to create a container that can be easily deployed on any machine, without having to clutter the file system of the host system.

## Deploying and Running Locally
TODO

## Future Improvements
TODO
