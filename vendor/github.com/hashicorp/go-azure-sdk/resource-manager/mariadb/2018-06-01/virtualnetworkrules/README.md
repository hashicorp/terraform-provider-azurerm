
## `github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/virtualnetworkrules` Documentation

The `virtualnetworkrules` SDK allows for interaction with the Azure Resource Manager Service `mariadb` (API Version `2018-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/virtualnetworkrules"
```


### Client Initialization

```go
client := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "virtualNetworkRuleValue")

payload := virtualnetworkrules.VirtualNetworkRule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkRulesClient.Delete`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "virtualNetworkRuleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkRulesClient.Get`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "virtualNetworkRuleValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VirtualNetworkRulesClient.ListByServer`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
