
## `github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/endpoints` Documentation

The `endpoints` SDK allows for interaction with the Azure Resource Manager Service `trafficmanager` (API Version `2018-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2018-08-01/endpoints"
```


### Client Initialization

```go
client := endpoints.NewEndpointsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EndpointsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := endpoints.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileValue", "AzureEndpoints", "endpointValue")

payload := endpoints.Endpoint{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.Delete`

```go
ctx := context.TODO()
id := endpoints.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileValue", "AzureEndpoints", "endpointValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.Get`

```go
ctx := context.TODO()
id := endpoints.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileValue", "AzureEndpoints", "endpointValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EndpointsClient.Update`

```go
ctx := context.TODO()
id := endpoints.NewEndpointTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "profileValue", "AzureEndpoints", "endpointValue")

payload := endpoints.Endpoint{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
