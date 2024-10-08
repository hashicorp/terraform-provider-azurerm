
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis` Documentation

The `managedapis` SDK allows for interaction with Azure Resource Manager `web` (API Version `2016-06-01`).

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


### Example Usage: `ManagedAPIsClient.Get`

```go
ctx := context.TODO()
id := managedapis.NewManagedApiID("12345678-1234-9876-4563-123456789012", "location", "apiName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ManagedAPIsClient.List`

```go
ctx := context.TODO()
id := managedapis.NewLocationID("12345678-1234-9876-4563-123456789012", "location")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
