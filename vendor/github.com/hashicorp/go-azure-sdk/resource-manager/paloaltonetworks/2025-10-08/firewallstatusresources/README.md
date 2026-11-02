
## `github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallstatusresources` Documentation

The `firewallstatusresources` SDK allows for interaction with Azure Resource Manager `paloaltonetworks` (API Version `2025-10-08`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/firewallstatusresources"
```


### Client Initialization

```go
client := firewallstatusresources.NewFirewallStatusResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `FirewallStatusResourcesClient.FirewallStatusGet`

```go
ctx := context.TODO()
id := firewallstatusresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

read, err := client.FirewallStatusGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FirewallStatusResourcesClient.FirewallStatusListByFirewalls`

```go
ctx := context.TODO()
id := firewallstatusresources.NewFirewallID("12345678-1234-9876-4563-123456789012", "example-resource-group", "firewallName")

// alternatively `client.FirewallStatusListByFirewalls(ctx, id)` can be used to do batched pagination
items, err := client.FirewallStatusListByFirewallsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
