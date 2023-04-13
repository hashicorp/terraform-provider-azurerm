
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service` Documentation

The `service` SDK allows for interaction with the Azure Resource Manager Service `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/service"
```


### Client Initialization

```go
client := service.NewServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ServiceClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := service.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "serviceValue")

payload := service.Service{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceClient.Delete`

```go
ctx := context.TODO()
id := service.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "serviceValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ServiceClient.Get`

```go
ctx := context.TODO()
id := service.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "serviceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ServiceClient.UpdateTags`

```go
ctx := context.TODO()
id := service.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "serviceValue")

payload := service.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
