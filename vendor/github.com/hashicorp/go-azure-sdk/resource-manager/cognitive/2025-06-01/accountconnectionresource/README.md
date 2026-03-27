
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource` Documentation

The `accountconnectionresource` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
```


### Client Initialization

```go
client := accountconnectionresource.NewAccountConnectionResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountConnectionResourceClient.AccountConnectionsCreate`

```go
ctx := context.TODO()
id := accountconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "connectionName")

payload := accountconnectionresource.ConnectionPropertiesV2BasicResource{
	// ...
}


read, err := client.AccountConnectionsCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountConnectionResourceClient.AccountConnectionsDelete`

```go
ctx := context.TODO()
id := accountconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "connectionName")

read, err := client.AccountConnectionsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountConnectionResourceClient.AccountConnectionsGet`

```go
ctx := context.TODO()
id := accountconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "connectionName")

read, err := client.AccountConnectionsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountConnectionResourceClient.AccountConnectionsList`

```go
ctx := context.TODO()
id := accountconnectionresource.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.AccountConnectionsList(ctx, id, accountconnectionresource.DefaultAccountConnectionsListOperationOptions())` can be used to do batched pagination
items, err := client.AccountConnectionsListComplete(ctx, id, accountconnectionresource.DefaultAccountConnectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `AccountConnectionResourceClient.AccountConnectionsUpdate`

```go
ctx := context.TODO()
id := accountconnectionresource.NewConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "connectionName")

payload := accountconnectionresource.ConnectionUpdateContent{
	// ...
}


read, err := client.AccountConnectionsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
