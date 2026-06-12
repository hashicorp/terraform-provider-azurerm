
## `github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/connectors` Documentation

The `connectors` SDK allows for interaction with Azure Resource Manager `storage` (API Version `2025-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/connectors"
```


### Client Initialization

```go
client := connectors.NewConnectorsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ConnectorsClient.Create`

```go
ctx := context.TODO()
id := connectors.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "connectorName")

payload := connectors.Connector{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectorsClient.Delete`

```go
ctx := context.TODO()
id := connectors.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "connectorName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectorsClient.Get`

```go
ctx := context.TODO()
id := connectors.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "connectorName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ConnectorsClient.ListByStorageAccount`

```go
ctx := context.TODO()
id := commonids.NewStorageAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName")

// alternatively `client.ListByStorageAccount(ctx, id)` can be used to do batched pagination
items, err := client.ListByStorageAccountComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ConnectorsClient.TestExistingConnection`

```go
ctx := context.TODO()
id := connectors.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "connectorName")

payload := connectors.TestExistingConnectionRequest{
	// ...
}


if err := client.TestExistingConnectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ConnectorsClient.Update`

```go
ctx := context.TODO()
id := connectors.NewConnectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "storageAccountName", "connectorName")

payload := connectors.ConnectorUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
