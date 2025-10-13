
## `github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets` Documentation

The `networkrulesets` SDK allows for interaction with Azure Resource Manager `eventhub` (API Version `2024-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
```


### Client Initialization

```go
client := networkrulesets.NewNetworkRuleSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NetworkRuleSetsClient.NamespacesCreateOrUpdateNetworkRuleSet`

```go
ctx := context.TODO()
id := networkrulesets.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

payload := networkrulesets.NetworkRuleSet{
	// ...
}


read, err := client.NamespacesCreateOrUpdateNetworkRuleSet(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkRuleSetsClient.NamespacesGetNetworkRuleSet`

```go
ctx := context.TODO()
id := networkrulesets.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

read, err := client.NamespacesGetNetworkRuleSet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NetworkRuleSetsClient.NamespacesListNetworkRuleSet`

```go
ctx := context.TODO()
id := networkrulesets.NewNamespaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "namespaceName")

// alternatively `client.NamespacesListNetworkRuleSet(ctx, id)` can be used to do batched pagination
items, err := client.NamespacesListNetworkRuleSetComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
