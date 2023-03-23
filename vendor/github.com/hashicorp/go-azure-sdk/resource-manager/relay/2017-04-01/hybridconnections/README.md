
## `github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/hybridconnections` Documentation

The `hybridconnections` SDK allows for interaction with the Azure Resource Manager Service `relay` (API Version `2017-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/relay/2017-04-01/hybridconnections"
```


### Client Initialization

```go
client := hybridconnections.NewHybridConnectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HybridConnectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue")

payload := hybridconnections.HybridConnection{
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


### Example Usage: `HybridConnectionsClient.CreateOrUpdateAuthorizationRule`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

payload := hybridconnections.AuthorizationRule{
	// ...
}


read, err := client.CreateOrUpdateAuthorizationRule(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.Delete`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.DeleteAuthorizationRule`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

read, err := client.DeleteAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.Get`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.GetAuthorizationRule`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

read, err := client.GetAuthorizationRule(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.ListAuthorizationRules`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue")

// alternatively `client.ListAuthorizationRules(ctx, id)` can be used to do batched pagination
items, err := client.ListAuthorizationRulesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HybridConnectionsClient.ListByNamespace`

```go
ctx := context.TODO()
id := hybridconnections.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue")

// alternatively `client.ListByNamespace(ctx, id)` can be used to do batched pagination
items, err := client.ListByNamespaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HybridConnectionsClient.ListKeys`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

read, err := client.ListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HybridConnectionsClient.RegenerateKeys`

```go
ctx := context.TODO()
id := hybridconnections.NewHybridConnectionAuthorizationRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceValue", "hybridConnectionValue", "authorizationRuleValue")

payload := hybridconnections.RegenerateAccessKeyParameters{
	// ...
}


read, err := client.RegenerateKeys(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
