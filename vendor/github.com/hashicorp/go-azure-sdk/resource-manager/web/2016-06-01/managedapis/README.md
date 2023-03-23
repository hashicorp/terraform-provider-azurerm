
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis` Documentation

The `managedapis` SDK allows for interaction with the Azure Resource Manager Service `web` (API Version `2016-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
```


### Client Initialization

```go
client := managedapis.NewManagedAPIsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ManagedAPIsClient.ManagedApisGet`

```go
ctx := context.TODO()
id := managedapis.NewManagedApiID("12345678-1234-9876-4563-123456789012", "locationValue", "managedApiValue")

read, err := client.ManagedApisGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedAPIsClient.ManagedApisList`

```go
ctx := context.TODO()
id := managedapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.ManagedApisList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
