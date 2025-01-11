
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork` Documentation

The `mobilenetwork` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
```


### Client Initialization

```go
client := mobilenetwork.NewMobileNetworkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MobileNetworkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := mobilenetwork.NewMobileNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName")

payload := mobilenetwork.MobileNetwork{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MobileNetworkClient.Delete`

```go
ctx := context.TODO()
id := mobilenetwork.NewMobileNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MobileNetworkClient.Get`

```go
ctx := context.TODO()
id := mobilenetwork.NewMobileNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MobileNetworkClient.UpdateTags`

```go
ctx := context.TODO()
id := mobilenetwork.NewMobileNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName")

payload := mobilenetwork.TagsObject{
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
