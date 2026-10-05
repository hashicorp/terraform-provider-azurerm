
## `github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules` Documentation

The `ipfirewallrules` SDK allows for interaction with Azure Resource Manager `synapse` (API Version `2021-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/ipfirewallrules"
```


### Client Initialization

```go
client := ipfirewallrules.NewIPFirewallRulesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `IPFirewallRulesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := ipfirewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "firewallRuleName")

payload := ipfirewallrules.IPFirewallRuleInfo{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `IPFirewallRulesClient.Delete`

```go
ctx := context.TODO()
id := ipfirewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "firewallRuleName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `IPFirewallRulesClient.Get`

```go
ctx := context.TODO()
id := ipfirewallrules.NewFirewallRuleID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName", "firewallRuleName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `IPFirewallRulesClient.ListByWorkspace`

```go
ctx := context.TODO()
id := ipfirewallrules.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

// alternatively `client.ListByWorkspace(ctx, id)` can be used to do batched pagination
items, err := client.ListByWorkspaceComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `IPFirewallRulesClient.ReplaceAll`

```go
ctx := context.TODO()
id := ipfirewallrules.NewWorkspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "workspaceName")

payload := ipfirewallrules.ReplaceAllIPFirewallRulesRequest{
	// ...
}


if err := client.ReplaceAllThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
