
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/firewallrules` Documentation

The `firewallrules` SDK allows for interaction with the Azure Resource Manager Service `mysql` (API Version `2017-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/firewallrules"
```


### Client Initialization

```go
client := firewallrules.NewFirewallRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "firewallRuleValue")

payload := firewallrules.FirewallRule{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallRulesClient.Delete`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "firewallRuleValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `FirewallRulesClient.Get`

```go
ctx := context.TODO()
id := firewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue", "firewallRuleValue")

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
id := firewallrules.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.ListByServer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
