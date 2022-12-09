
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork` Documentation

The `datanetwork` SDK allows for interaction with the Azure Resource Manager Service `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/datanetwork"
```


### Client Initialization

```go
client := datanetwork.NewDataNetworkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DataNetworkClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := datanetwork.NewDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "dataNetworkValue")

payload := datanetwork.DataNetwork{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DataNetworkClient.Delete`

```go
ctx := context.TODO()
id := datanetwork.NewDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "dataNetworkValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DataNetworkClient.Get`

```go
ctx := context.TODO()
id := datanetwork.NewDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "dataNetworkValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DataNetworkClient.UpdateTags`

```go
ctx := context.TODO()
id := datanetwork.NewDataNetworkID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkValue", "dataNetworkValue")

payload := datanetwork.TagsObject{
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
