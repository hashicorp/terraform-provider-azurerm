
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/virtualnetworkrules` Documentation

The `virtualnetworkrules` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/virtualnetworkrules"
```


### Client Initialization

```go
client := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VirtualNetworkRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "virtualNetworkRuleName")

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
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "virtualNetworkRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `VirtualNetworkRulesClient.Get`

```go
ctx := context.TODO()
id := virtualnetworkrules.NewVirtualNetworkRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "virtualNetworkRuleName")

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
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

// alternatively `client.ListByServer(ctx, id)` can be used to do batched pagination
items, err := client.ListByServerComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
