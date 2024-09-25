
## `github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/firewallrules` Documentation

The `firewallrules` SDK allows for interaction with Azure Resource Manager `sql` (API Version `2023-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/firewallrules"
```


### Client Initialization

```go
client := firewallrules.NewFirewallRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "firewallRuleName")

payload := firewallrules.FirewallRule{
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


### Example Usage: `FirewallRulesClient.Delete`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "firewallRuleName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallRulesClient.Get`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName", "firewallRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallRulesClient.ListByServer`

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


### Example Usage: `FirewallRulesClient.Replace`

```go
ctx := context.TODO()
id := commonids.NewSqlServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverName")

payload := firewallrules.FirewallRuleList{
	// ...
}


read, err := client.Replace(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
