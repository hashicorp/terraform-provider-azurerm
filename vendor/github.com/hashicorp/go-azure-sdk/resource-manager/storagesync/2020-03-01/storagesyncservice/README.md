
## `github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservice` Documentation

The `storagesyncservice` SDK allows for interaction with the Azure Resource Manager Service `storagesync` (API Version `2020-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservice"
```


### Client Initialization

```go
client := storagesyncservice.NewStorageSyncServiceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `StorageSyncServiceClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := storagesyncservice.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := storagesyncservice.CheckNameAvailabilityParameters{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
