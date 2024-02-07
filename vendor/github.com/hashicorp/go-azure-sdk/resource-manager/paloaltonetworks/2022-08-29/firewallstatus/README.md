
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewallstatus` Documentation

The `firewallstatus` SDK allows for interaction with the Azure Resource Manager Service `paloaltonetworks` (API Version `2022-08-29`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewallstatus"
```


### Client Initialization

```go
client := firewallstatus.NewFirewallStatusClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallStatusClient.Get`

```go
ctx := context.TODO()
id := firewallstatus.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallStatusClient.ListByFirewalls`

```go
ctx := context.TODO()
id := firewallstatus.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallValue")

// alternatively `client.ListByFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.ListByFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
