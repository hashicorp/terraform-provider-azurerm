
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/openapis` Documentation

The `openapis` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2024-12-30`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2024-12-30/openapis"
```


### Client Initialization

```go
client := openapis.NewOpenapisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OpenapisClient.CheckNameAvailabilityExecute`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := openapis.NameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailabilityExecute(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CheckNameAvailabilityWithoutLocationExecute`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

payload := openapis.NameAvailabilityRequest{
	// ...
}


read, err := client.CheckNameAvailabilityWithoutLocationExecute(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CheckVirtualNetworkSubnetUsageExecute`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := openapis.VirtualNetworkSubnetUsageParameter{
	// ...
}


read, err := client.CheckVirtualNetworkSubnetUsageExecute(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GetPrivateDnsZoneSuffixExecute`

```go
ctx := context.TODO()


read, err := client.GetPrivateDnsZoneSuffixExecute(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.LocationBasedCapabilitiesList`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.LocationBasedCapabilitiesList(ctx, id)` can be used to do batched pagination
items, err := client.LocationBasedCapabilitiesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.OperationProgressGet`

```go
ctx := context.TODO()
id := openapis.NewOperationProgressID("12345678-1234-9876-4563-123456789012", "locationName", "operationId")

read, err := client.OperationProgressGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
