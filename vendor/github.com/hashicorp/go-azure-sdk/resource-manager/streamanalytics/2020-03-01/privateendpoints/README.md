
## `github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints` Documentation

The `privateendpoints` SDK allows for interaction with the Azure Resource Manager Service `streamanalytics` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
```


### Client Initialization

```go
client := privateendpoints.NewPrivateEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateEndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "privateEndpointValue")

payload := privateendpoints.PrivateEndpoint{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, privateendpoints.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointsClient.Delete`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "privateEndpointValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `PrivateEndpointsClient.Get`

```go
ctx := context.TODO()
id := privateendpoints.NewPrivateEndpointID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "privateEndpointValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateEndpointsClient.ListByCluster`

```go
ctx := context.TODO()
id := privateendpoints.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

// alternatively `client.ListByCluster(ctx, id)` can be used to do batched pagination
items, err := client.ListByClusterComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
