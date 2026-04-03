
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/clients` Documentation

The `clients` SDK allows for interaction with Azure Resource Manager `eventgrid` (API Version `2025-02-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/clients"
```


### Client Initialization

```go
client := clients.NewClientsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClientsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := clients.NewClientID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "clientName")

payload := clients.Client{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClientsClient.Delete`

```go
ctx := context.TODO()
id := clients.NewClientID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "clientName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClientsClient.Get`

```go
ctx := context.TODO()
id := clients.NewClientID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "clientName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ClientsClient.ListByNamespace`

```go
ctx := context.TODO()
id := clients.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, clients.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, clients.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
