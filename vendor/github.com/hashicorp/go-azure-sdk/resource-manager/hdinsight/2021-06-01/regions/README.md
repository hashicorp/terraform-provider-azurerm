
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/regions` Documentation

The `regions` SDK allows for interaction with Azure Resource Manager `hdinsight` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/regions"
```


### Client Initialization

```go
client := regions.NewRegionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RegionsClient.LocationsCheckNameAvailability`

```go
ctx := context.TODO()
id := regions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := regions.NameAvailabilityCheckRequestParameters{
	// ...
}


read, err := client.LocationsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegionsClient.LocationsGetCapabilities`

```go
ctx := context.TODO()
id := regions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.LocationsGetCapabilities(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegionsClient.LocationsListBillingSpecs`

```go
ctx := context.TODO()
id := regions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.LocationsListBillingSpecs(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegionsClient.LocationsListUsages`

```go
ctx := context.TODO()
id := regions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.LocationsListUsages(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RegionsClient.LocationsValidateClusterCreateRequest`

```go
ctx := context.TODO()
id := regions.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := regions.ClusterCreateRequestValidationParameters{
	// ...
}


read, err := client.LocationsValidateClusterCreateRequest(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
