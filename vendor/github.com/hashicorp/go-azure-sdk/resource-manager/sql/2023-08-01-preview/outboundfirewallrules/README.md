
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/outboundfirewallrules` Documentation

The `outboundfirewallrules` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/outboundfirewallrules"
```


### Client Initialization

```go
client := outboundfirewallrules.NewOutboundFirewallRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OutboundFirewallRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := outboundfirewallrules.NewOutboundFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "outboundFirewallRuleName")

if err := client.CreateOrUpdateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OutboundFirewallRulesClient.Delete`

```go
ctx := context.TODO()
id := outboundfirewallrules.NewOutboundFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "outboundFirewallRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OutboundFirewallRulesClient.Get`

```go
ctx := context.TODO()
id := outboundfirewallrules.NewOutboundFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "outboundFirewallRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OutboundFirewallRulesClient.ListByServer`

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
