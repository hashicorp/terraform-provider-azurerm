
## `github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections` Documentation

The `connections` SDK allows for interaction with Azure Resource Manager `web` (API Version `2016-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
```


### Client Initialization

```go
client := connections.NewConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectionsClient.ConfirmConsentCode`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := connections.ConfirmConsentCodeDefinition{
	// ...
}


read, err := client.ConfirmConsentCode(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := connections.ApiConnectionDefinition{
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


### Example Usage: `ConnectionsClient.Delete`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionsClient.Get`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionsClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.List(ctx, id, connections.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionsClient.ListConsentLinks`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := connections.ListConsentLinksDefinition{
	// ...
}


read, err := client.ListConsentLinks(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectionsClient.Update`

```go
ctx := context.TODO()
id := connections.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "connectionName")

payload := connections.ApiConnectionDefinition{
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
