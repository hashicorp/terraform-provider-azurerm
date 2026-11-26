
## `github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers` Documentation

The `trafficmanagers` SDK allows for interaction with Azure Resource Manager `trafficmanager` (API Version `2022-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
```


### Client Initialization

```go
client := trafficmanagers.NewTrafficmanagersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TrafficmanagersClient.EndpointsCreateOrUpdate`

```go
ctx := context.TODO()
id := trafficmanagers.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileName", "AzureEndpoints", "endpointName")

payload := trafficmanagers.Endpoint{
	// ...
}


read, err := client.EndpointsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.EndpointsDelete`

```go
ctx := context.TODO()
id := trafficmanagers.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileName", "AzureEndpoints", "endpointName")

read, err := client.EndpointsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.EndpointsGet`

```go
ctx := context.TODO()
id := trafficmanagers.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileName", "AzureEndpoints", "endpointName")

read, err := client.EndpointsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.EndpointsUpdate`

```go
ctx := context.TODO()
id := trafficmanagers.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileName", "AzureEndpoints", "endpointName")

payload := trafficmanagers.Endpoint{
	// ...
}


read, err := client.EndpointsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.HeatMapGet`

```go
ctx := context.TODO()
id := trafficmanagers.NewTrafficManagerProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "trafficManagerProfileName")

read, err := client.HeatMapGet(ctx, id, trafficmanagers.DefaultHeatMapGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.ProfilesCheckTrafficManagerRelativeDnsNameAvailability`

```go
ctx := context.TODO()

payload := trafficmanagers.CheckTrafficManagerRelativeDnsNameAvailabilityParameters{
	// ...
}


read, err := client.ProfilesCheckTrafficManagerRelativeDnsNameAvailability(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TrafficmanagersClient.ProfilescheckTrafficManagerNameAvailabilityV2`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := trafficmanagers.CheckTrafficManagerRelativeDnsNameAvailabilityParameters{
	// ...
}


read, err := client.ProfilescheckTrafficManagerNameAvailabilityV2(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
