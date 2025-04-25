
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs` Documentation

The `eventhubs` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
```


### Client Initialization

```go
client := eventhubs.NewEventHubsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EventHubsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := eventhubs.NewEventhubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName")

payload := eventhubs.Eventhub{
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


### Example Usage: `EventHubsClient.Delete`

```go
ctx := context.TODO()
id := eventhubs.NewEventhubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventHubsClient.DeleteAuthorizationRule`

```go
ctx := context.TODO()
id := eventhubs.NewEventhubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "authorizationRuleName")

read, err := client.DeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventHubsClient.Get`

```go
ctx := context.TODO()
id := eventhubs.NewEventhubID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventHubsClient.GetAuthorizationRule`

```go
ctx := context.TODO()
id := eventhubs.NewEventhubAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName", "eventhubName", "authorizationRuleName")

read, err := client.GetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EventHubsClient.ListByNamespace`

```go
ctx := context.TODO()
id := eventhubs.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.ListByNamespace(ctx, id, eventhubs.DefaultListByNamespaceOperationOptions())` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id, eventhubs.DefaultListByNamespaceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
